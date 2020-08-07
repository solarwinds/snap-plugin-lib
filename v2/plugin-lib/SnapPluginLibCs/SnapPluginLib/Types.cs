using System;
using System.Runtime.InteropServices;

namespace SnapPluginLib
{
    [StructLayout(LayoutKind.Sequential)]
    internal class LogMap
    {
        public IntPtr elements;
        public int length;
    }

    [StructLayout(LayoutKind.Sequential)]
    internal class LogMapElements
    {
        public string key;
        public string value;
    }
    
    [StructLayout(LayoutKind.Sequential)]
    internal class Modifiers
    {
        
    }
    
    [StructLayout(LayoutKind.Sequential)]
    internal class Value
    {
        
    }
}