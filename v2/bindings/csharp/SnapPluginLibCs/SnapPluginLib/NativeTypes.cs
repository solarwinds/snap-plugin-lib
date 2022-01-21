/*
 Copyright (c) 2021 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

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

        public int metric_type;
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

    internal enum NativeMetricType {
        Unknown,
        Gauge,
        Sym,
        Summary,
        Histogram,
    }
}