package packet_entities

type 
namespace D2NET.Core.PacketEntities
{
    using System;
    using System.Collections.Generic;

    public class PacketEntity : ICloneable
    {
        private const int SERIAL_NUM_BITS = 10;
        public PacketEntity()
        {
            this.Values = new Dictionary<string, object>();
        }
        public int Tick { get; set; }
        public int Index { get; set; }
        public int SerialNum { get; set; }
        public int ClassId { get; set; }
        public int Handle { get { return (this.Index | (this.SerialNum << (SERIAL_NUM_BITS + 1))); } }
        public string Name { get; set; }
        public UpdateType Type { get; set; }
        public Dictionary<string, object> Values { get; set; }
        public object Clone()
        {
            PacketEntity clone = new PacketEntity();
            clone.ClassId = this.ClassId;
            clone.Index = this.Index;
            clone.Name = this.Name;
            clone.SerialNum = this.SerialNum;
            clone.Tick = this.Tick;
            clone.Type = this.Type;
            foreach (string key in this.Values.Keys) clone.Values.Add(key, this.Values[key]);
            return clone;
        }
    }
}

