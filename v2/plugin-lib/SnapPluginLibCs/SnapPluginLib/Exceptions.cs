using System;
using System.Runtime.InteropServices;

namespace SnapPluginLib
{
    internal static class Exceptions
    {
        public static void ThrowExceptionIfError(IntPtr errPtr)
        {
            var err = Marshal.PtrToStructure<NativeError>(errPtr);
            if (err != null && err.errorMessage != null)
            {
                throw new PluginLibException(err.errorMessage);
            }
        }
    } 
    
    public class PluginLibException : Exception
    {
        public PluginLibException(string nativeMsg) : base(nativeMsg)
        {
        }
    }
}