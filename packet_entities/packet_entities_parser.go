namespace D2NET.Core.PacketEntities {
    using System;
    using System.Collections.Generic;
    using System.Linq;

    using D2NET.Core.Parser;
    using D2NET.Core.Protobuf;
    using D2NET.Core.SendTables;
    using D2NET.Core.Utils;

    public class PacketEntitiesParser {
        private Dictionary<int, Dictionary<string, object>> baseline = new Dictionary<int, Dictionary<string, object>>();
        private int classIdNumBits = 0;
        private Dictionary<string, int> classInfosIdMapping = new Dictionary<string, int>();
        private Dictionary<int, string> classInfosNameMapping = new Dictionary<int, string>();
        private Dictionary<string, Handler> createClassHandlers = new Dictionary<string, Handler>();
        private Dictionary<string, Handler> deleteClassHandlers = new Dictionary<string, Handler>();
        private Dictionary<int, SendProp[]> mapping = new Dictionary<int, SendProp[]>();
        private Dictionary<int, Dictionary<string, int>> multiples = new Dictionary<int, Dictionary<string, int>>();
        private List<ParserBaseItem> packets;
        private Dictionary<string, HandlerPreserve> preserveClassHandlers = new Dictionary<string, HandlerPreserve>();
        private SendTablesHelper sendTablesHelper;

        public PacketEntitiesParser(List<ParserBaseItem> items) {
            this.Entities = new PacketEntity[2048];
            this.packets = items.Where(o => o.ItemType == typeof(CSVCMsg_PacketEntities) && o.From == ParserBaseEvent.DEM_Packet)
                                .OrderBy(o => o.Tick)
                                .ToList();
            CSVCMsg_ServerInfo server_info = (CSVCMsg_ServerInfo)items.First(o => o.ItemType == typeof(CSVCMsg_ServerInfo)).Value;
            this.classIdNumBits = (int)(Math.Log(server_info.max_classes) / Math.Log(2)) + 1;
            CDemoClassInfo class_infos = (CDemoClassInfo)items.First(o => o.ItemType == typeof(CDemoClassInfo)).Value;
            this.classInfosNameMapping = class_infos.classes.ToDictionary(o => o.class_id, o => o.table_name);
            this.classInfosIdMapping = class_infos.classes.ToDictionary(o => o.table_name, o => o.class_id);
            this.sendTablesHelper = new SendTablesHelper(items.Where(o => o.ItemType == typeof(CSVCMsg_SendTable))
                                                              .ToDictionary(o => ((CSVCMsg_SendTable)o.Value).net_table_name, o => (CSVCMsg_SendTable)o.Value));
            foreach (string key in this.classInfosIdMapping.Keys) {
                SendProp[] props = this.sendTablesHelper.LoadSendTable(key);
                this.mapping.Add(this.classInfosIdMapping[key], props);
                this.multiples.Add(this.classInfosIdMapping[key], (from o in props group o by new { o.dt_name, o.var_name } into g select new { Key = g.Key.dt_name + "." + g.Key.var_name, Count = g.Count() }).ToDictionary(o => o.Key, o => o.Count));
            }

            CDemoStringTables.table_t instanceBaseline = items.Where(o => o.ItemType == typeof(CDemoStringTables))
                                                              .Select(o => (CDemoStringTables)o.Value)
                                                              .Last(o => o.tables.Any(p => p.table_name == "instancebaseline"))
                                                              .tables
                                                              .First(o => o.table_name == "instancebaseline");
            foreach (CDemoStringTables.items_t item in instanceBaseline.items) {
                int class_id = int.Parse(item.str);
                string class_name = this.classInfosNameMapping[class_id];
                BitReader br = new BitReader(item.data);
                List<int> indexes = br.ReadPropertiesIndex();
                try {
                    this.baseline.Add(int.Parse(item.str), br.ReadPropertiesValues(this.mapping[class_id], this.multiples[class_id], indexes));
                }
                catch (Exception ex) {
                    var toto = true;
                }
            }
        }

        private PacketEntitiesParser() {
        }

        public delegate void Procedure(PacketEntity pe);

        public delegate void ProcedurePreserve(PacketEntity pe, Dictionary<string, object> values);

        public PacketEntity[] Entities {
            get; private set;
        }

        public void AddCreateHandler(string class_name, Procedure callback) {
            Handler h = new Handler();
            h.ClassName = class_name;
            h.Callback = callback;
            this.createClassHandlers.Add(class_name, h);
        }

        public void AddDeleteHandler(string class_name, Procedure callback) {
            Handler h = new Handler();
            h.ClassName = class_name;
            h.Callback = callback;
            this.deleteClassHandlers.Add(class_name, h);
        }

        public void AddPreserveHandler(string class_name, ProcedurePreserve callback) {
            HandlerPreserve h = new HandlerPreserve();
            h.ClassName = class_name;
            h.Callback = callback;
            this.preserveClassHandlers.Add(class_name, h);
        }

        public void Parse() {
            for (int i = 0; i < this.packets.Count; i++) {
                this.ParsePacket(i);
            }
        }

        private void EntityCreate(BitReader br, int current_index, int tick) {
            PacketEntity pe = new PacketEntity();
            pe.Tick = tick;
            pe.ClassId = (int)br.ReadUBits(this.classIdNumBits);
            pe.SerialNum = (int)br.ReadUBits(10);
            pe.Index = current_index;
            pe.Name = this.classInfosNameMapping[pe.ClassId];
            pe.Type = UpdateType.Create;
            pe.Values = new Dictionary<string, object>();
            List<int> indexes = br.ReadPropertiesIndex();
            Dictionary<string, object> values = br.ReadPropertiesValues(this.mapping[pe.ClassId], this.multiples[pe.ClassId], indexes);
            if (this.baseline.ContainsKey(pe.ClassId)) {
                foreach (string basekey in this.baseline[pe.ClassId].Keys) {
                    object object_value;
                    if (values.TryGetValue(basekey, out object_value)) {
                        pe.Values.Add(basekey, object_value);
                    } else {
                        pe.Values.Add(basekey, this.baseline[pe.ClassId][basekey]);
                    }
                }
            } else {
                foreach (string key in values.Keys) {
                    pe.Values.Add(key, values[key]);
                }
            }

            this.Entities[pe.Index] = pe;
            Handler handler;
            if (this.createClassHandlers.TryGetValue(pe.Name, out handler)) {
                handler.Callback(pe);
            }
        }

        private void EntityDelete(BitReader br, int current_index, int tick) {
            PacketEntity pe = (PacketEntity)this.Entities[current_index].Clone();
            pe.Tick = tick;
            pe.Type = UpdateType.Delete;
            this.Entities[current_index] = null;
            Handler handler;
            if (this.deleteClassHandlers.TryGetValue(pe.Name, out handler)) {
                handler.Callback(pe);
            }
        }

        private void EntityPreserve(BitReader br, int current_index, int tick) {
            PacketEntity pe = this.Entities[current_index];
            pe.Tick = tick;
            pe.Type = UpdateType.Preserve;
            List<int> indexes = br.ReadPropertiesIndex();
            Dictionary<string, object> values = br.ReadPropertiesValues(this.mapping[this.classInfosIdMapping[pe.Name]], this.multiples[this.classInfosIdMapping[pe.Name]], indexes);
            foreach (string key in values.Keys) {
                object object_value;
                if (pe.Values.TryGetValue(key, out object_value)) {
                    pe.Values[key] = values[key];
                } else {
                    pe.Values.Add(key, values[key]);
                }
            }

            HandlerPreserve handler;
            if (this.preserveClassHandlers.TryGetValue(pe.Name, out handler)) {
                handler.Callback(pe, values);
            }
        }

        private void ParsePacket(int index) {
            CSVCMsg_PacketEntities pe = (CSVCMsg_PacketEntities)this.packets[index].Value;
            BitReader br = new BitReader(pe.entity_data);
            int current_index = -1;
            for (int i = 0; i < pe.updated_entries; i++) {
                current_index = br.ReadNextEntityIndex(current_index);
                UpdateType type = br.ReadUpdateType();
                if (type == UpdateType.Preserve) {
                    this.EntityPreserve(br, current_index, this.packets[index].Tick);
                } else if (type == UpdateType.Create) {
                    this.EntityCreate(br, current_index, this.packets[index].Tick);
                } else if (type == UpdateType.Delete) {
                    this.EntityDelete(br, current_index, this.packets[index].Tick);
                }
            }
        }

        internal class Handler {
            public Procedure Callback {
                get; set;
            }

            public string ClassName {
                get; set;
            }

            public override int GetHashCode() {
                return this.ClassName.GetHashCode();
            }
        }

        internal class HandlerPreserve {
            public ProcedurePreserve Callback {
                get; set;
            }

            public string ClassName {
                get; set;
            }

            public override int GetHashCode() {
                return this.ClassName.GetHashCode();
            }
        }
    }
}

