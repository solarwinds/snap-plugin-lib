using System;
using System.Runtime.InteropServices;

namespace SnapPluginLib
{
    [StructLayout(LayoutKind.Sequential)]
    internal class NativeMap
    {
        public IntPtr elements;
        public int length;
    }

    [StructLayout(LayoutKind.Sequential)]
    internal class NativeMapElements
    {
        public string key;
        public string value;
    }

    [StructLayout(LayoutKind.Sequential)]
    internal class NativeModifiers
    {
        public IntPtr tagsToAdd; // NativeMap
        public IntPtr tagsToRemove; //NativeMap
        public NativeTimeWithNs timestamp;
        
        [MarshalAs(UnmanagedType.LPStr)] public string description;
        [MarshalAs(UnmanagedType.LPStr)] public string unit;
    }

    [StructLayout(LayoutKind.Sequential)]
    internal class NativeTimeWithNs
    {
        public int sec;
        public int nsec;
    }
    
    internal enum ValueType
    {
        TypeInt64 = 1,
        TypeUint64,
        TypeInt32,
        TypeUint32,
        TypeFloat,
        TypeDouble,
        TypeBool,
        TypeCString,
    }

    [StructLayout(LayoutKind.Explicit)]
    internal class NativeValue
    {
        // Union
        [FieldOffset(0)] public Int64 v_int64;
        [FieldOffset(0)] public UInt64 v_uint64;
        [FieldOffset(0)] public Int32 v_int32;
        [FieldOffset(0)] public UInt32 v_uint32;
        [FieldOffset(0)] public float v_float;
        [FieldOffset(0)] public Double v_double;
        [FieldOffset(0)] public int v_bool;
        [FieldOffset(0)] public IntPtr v_cstring;

        // Std values
        [FieldOffset(8)] public int vtype;
    }

    [StructLayout(LayoutKind.Sequential)]
    internal class NativeError
    {
        public string errorMessage;
    }
}