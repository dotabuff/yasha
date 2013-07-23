namespace D2NET.Core.PacketEntities {
    using System;
    using System.Collections.Generic;
    using System.IO;
    using System.Text;
    using D2NET.Core.SendTables;
    using D2NET.Core.Utils;

    public static class BitReaderExtensions {
        public static int ReadNextEntityIndex(this BitReader br, int old_entity) {
            uint ret = br.ReadUBits(4);
            bool more1 = br.ReadBoolean();
            bool more2 = br.ReadBoolean();
            if (more1) ret += (br.ReadUBits(4) << 4);
            if (more2) ret += (br.ReadUBits(8) << 4);
            return old_entity + 1 + (int)ret;
        }

        public static UpdateType ReadUpdateType(this BitReader br) {
            UpdateType result = UpdateType.Preserve;
            if (!br.ReadBoolean()) {
                if (br.ReadBoolean()) result = UpdateType.Create;
            } else {
                result = UpdateType.Leave;
                if (br.ReadBoolean()) result = UpdateType.Delete;
            }
            return result;
        }

        public static List<int> ReadPropertiesIndex(this BitReader br) {
            List<int> props = new List<int>();
            int prop = -1;
            while (true) {
                if (br.ReadBoolean()) {
                    prop += 1;
                    props.Add(prop);
                } else {
                    uint value = br.ReadVarInt();
                    if (value == 16383) break;
                    prop += 1;
                    prop += (int)value;
                    props.Add(prop);
                }
            }
            return props;
        }

        public static Dictionary<string, object> ReadPropertiesValues(this BitReader br, SendProp[] mapping, Dictionary<string, int> multiples, List<int> indexes) {
            Dictionary<string, object> values = new Dictionary<string, object>();
            for (int j = 0; j < indexes.Count; j++) {
                SendProp prop = mapping[indexes[j]];
                bool multiple = multiples[prop.dt_name + "." + prop.var_name] > 1;
                int elements = 1;
                if ((prop.flags & Core.SendTables.FlagType.SPROP_INSIDEARRAY) != 0) elements = (int)br.ReadUBits(6);
                for (int k = 0; k < elements; k++) {
                    string key = string.Format("{0}.{1}", prop.dt_name, prop.var_name);
                    if (multiple) key += "-" + indexes[j];
                    if (elements > 1) key += "-" + k.ToString();
                    switch (prop.type) {
                        case D2NET.Core.SendTables.DPTType.DPT_Int:
                            if ((prop.flags & Core.SendTables.FlagType.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0) {
                                values.Add(key, br.ReadVarInt());
                            } else {
                                values.Add(key, br.ReadInt(prop));
                            }
                            break;
                        case D2NET.Core.SendTables.DPTType.DPT_Float:
                            if ((prop.flags & Core.SendTables.FlagType.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0) {
                                throw new InvalidOperationException("SPROP_ENCODED_AGAINST_TICKCOUNT");
                            } else {
                                values.Add(key, br.ReadFloat(prop));
                            }
                            break;
                        case D2NET.Core.SendTables.DPTType.DPT_Vector:
                            if ((prop.flags & Core.SendTables.FlagType.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0) {
                                throw new InvalidOperationException("SPROP_ENCODED_AGAINST_TICKCOUNT");
                            } else {
                                values.Add(key, br.ReadVector(prop));
                            }
                            break;
                        case D2NET.Core.SendTables.DPTType.DPT_VectorXY:
                            if ((prop.flags & Core.SendTables.FlagType.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0) {
                                throw new InvalidOperationException("SPROP_ENCODED_AGAINST_TICKCOUNT");
                            } else {
                                values.Add(key, br.ReadVectorXY(prop));
                            }
                            break;
                        case D2NET.Core.SendTables.DPTType.DPT_String:
                            if ((prop.flags & Core.SendTables.FlagType.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0) {
                                throw new InvalidOperationException("SPROP_ENCODED_AGAINST_TICKCOUNT");
                            } else {
                                values.Add(key, br.ReadLengthPrefixedString());
                            }
                            break;
                        case D2NET.Core.SendTables.DPTType.DPT_Int64:
                            if ((prop.flags & Core.SendTables.FlagType.SPROP_ENCODED_AGAINST_TICKCOUNT) != 0) {
                                throw new InvalidOperationException("SPROP_ENCODED_AGAINST_TICKCOUNT");
                            } else {
                                values.Add(key, br.ReadInt64(prop));
                            }
                            break;
                    }
                }
            }
            return values;
        }
        
        public static Vector ReadVector(this BitReader br, SendProp prop) {
            Vector result = new Vector();
            result.X = br.ReadFloat(prop);
            result.Y = br.ReadFloat(prop);
            if ((prop.flags & FlagType.SPROP_NORMAL) == 0) result.Z = br.ReadFloat(prop);
            else {
                bool signbit = br.ReadBoolean();
                float v0v0v1v1 = result.X * result.X + result.Y * result.Y;
                if (v0v0v1v1 < 1.0f) result.Z = (float)Math.Sqrt(1.0f - v0v0v1v1);
                else result.Z = 0.0f;
                if (signbit) result.Z *= -1.0f;
            }

            return result;
        }
        
        public static Vector ReadVectorXY(this BitReader br, SendProp prop) {
            Vector result = new Vector();
            result.X = br.ReadFloat(prop);
            result.Y = br.ReadFloat(prop);
            return result;
        }
        
        public static int ReadInt(this BitReader br, SendProp prop) {
            if ((prop.flags & FlagType.SPROP_UNSIGNED) != 0) return (int)br.ReadUBits(prop.num_bits);
            return br.ReadBits(prop.num_bits);
        }
        
        public static bool ReadSpecialFloat(this BitReader br, SendProp prop, ref float value) {
            if ((prop.flags & FlagType.SPROP_COORD) != 0) {
                value = br.ReadBitCoord();
                return true;
            } else if ((prop.flags & FlagType.SPROP_COORD_MP) != 0) {
                throw new InvalidOperationException("BitReader.ReadSpecialFloat");
            } else if ((prop.flags & FlagType.SPROP_COORD_MP_INTEGRAL) != 0) {
                throw new InvalidOperationException("BitReader.ReadSpecialFloat");
            } else if ((prop.flags & FlagType.SPROP_COORD_MP_LOWPRECISION) != 0) {
                throw new InvalidOperationException("BitReader.ReadSpecialFloat");
            } else if ((prop.flags & FlagType.SPROP_CELL_COORD) != 0) {
                value = br.ReadBitCellCoord(prop.num_bits, false, false);
                return true;
            } else if ((prop.flags & FlagType.SPROP_CELL_COORD_INTEGRAL) != 0) {
                value = br.ReadBitCellCoord(prop.num_bits, true, false);
                return true;
            } else if ((prop.flags & FlagType.SPROP_CELL_COORD_LOWPRECISION) != 0) {
                value = br.ReadBitCellCoord(prop.num_bits, false, true);
                return true;
            } else if ((prop.flags & FlagType.SPROP_NOSCALE) != 0) {
                value = br.ReadBitFloat();
                return true;
            } else if ((prop.flags & FlagType.SPROP_NORMAL) != 0) {
                value = br.ReadBitNormal();
                return true;
            }
            return false;
        }
        
        public static float ReadFloat(this BitReader br, SendProp prop) {
            float result = float.MinValue;
            if (br.ReadSpecialFloat(prop, ref result)) return result;
            uint dwInterp = br.ReadUBits(prop.num_bits);
            result = (float)dwInterp / ((1 << prop.num_bits) - 1);
            result = prop.low_value + (prop.high_value - prop.low_value) * result;
            return result;
        }

        public static string ReadLengthPrefixedString(this BitReader br) {
            uint stringLength = br.ReadUBits(9);
            if (stringLength > 0) return Encoding.UTF8.GetString(br.ReadBytes(stringLength));
            return string.Empty;
        }
        
        public static ulong ReadInt64(this BitReader br, SendProp prop) {
            uint low, high;
            if ((prop.flags & FlagType.SPROP_UNSIGNED) != 0) {
                low = br.ReadUBits(32);
                high = br.ReadUBits(32);
            } else {
                br.SeekBits(1, SeekOrigin.Current);
                low = br.ReadUBits(32);
                high = br.ReadUBits(31);
            }
            ulong res = high;
            res = (res << 32);
            res = res | (ulong)low;
            return res;
        }
    }
}

