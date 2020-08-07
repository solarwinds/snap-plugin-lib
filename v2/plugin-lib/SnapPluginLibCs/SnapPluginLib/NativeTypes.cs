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
        public IntPtr tagsToAdd;
        public NativeMap tagsToRemove;
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

    [StructLayout(LayoutKind.Explicit)]
    internal class NativeValue
    {
        // Union
        [FieldOffset(0)] public Int64 v_int64;
        [FieldOffset(0)] public UInt64 v_uint64;
        [FieldOffset(0)] public Double v_double;
        [FieldOffset(0)] public int v_bool;

        // Std values
        [FieldOffset(8)] public int vtype;
    }
}