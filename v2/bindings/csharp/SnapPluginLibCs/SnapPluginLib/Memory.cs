using System;
using System.Runtime.InteropServices;

namespace SnapPluginLib
{
    internal static class Memory
    {
        public static void FreeNativeMap(IntPtr nativeMapPtr)
        {
            var nativeMap = Marshal.PtrToStructure<NativeMap>(nativeMapPtr);
            Marshal.FreeHGlobal(nativeMap.elements);
            Marshal.FreeHGlobal(nativeMapPtr);
        }

        public static void FreeNativeModifiers(NativeModifiers nativeModifiers)
        {
            if (nativeModifiers.tagsToAdd != IntPtr.Zero)
            {
                Marshal.FreeHGlobal(nativeModifiers.tagsToAdd);
            }
            
            if (nativeModifiers.tagsToRemove != IntPtr.Zero)
            {
                Marshal.FreeHGlobal(nativeModifiers.tagsToRemove);
            }
        }
    }
}