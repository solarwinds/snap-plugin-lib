using System;
using System.Collections.Generic;
using System.Runtime.InteropServices;

namespace SnapPluginLib
{
    internal static class Convertions
    {
        // Conversion: Dictionary<string, string> -> map_t* (NativeMap)
        public static IntPtr DictionaryToNativeMapMem(Dictionary<string, string> dictionary)
        {
            var nativeMap = new NativeMap
            {
                length = dictionary.Count,
                elements = Marshal.AllocHGlobal(Marshal.SizeOf(typeof(NativeMapElements)) * dictionary.Count)
            };

            var i = 0;
            foreach (KeyValuePair<string, string> entry in dictionary)
            {
                var nativeMapElem = new NativeMapElements();
                nativeMapElem.key = entry.Key;
                nativeMapElem.value = entry.Value;

                Marshal.StructureToPtr(nativeMapElem,
                    (IntPtr) nativeMap.elements.ToInt64() + i * Marshal.SizeOf(typeof(NativeMapElements)), false);

                i++;
            }

            var nativeMapAsMemBlock = Marshal.AllocHGlobal(Marshal.SizeOf(typeof(NativeMap)));
            Marshal.StructureToPtr(nativeMap, nativeMapAsMemBlock, false);

            return nativeMapAsMemBlock;
        }
        
        // Conversion: char** -> List<string>
        public static List<string> NativeStringArrayToList(IntPtr arrPtr)
        {
            var requestedMetrics = new List<string>();

            for (int offset = 0;; offset += 8)
            {
                var charPtr = Marshal.ReadIntPtr(arrPtr + offset);
                if (charPtr == IntPtr.Zero)
                    break;

                requestedMetrics.Add(Marshal.PtrToStringAnsi(charPtr));
            }
            
            return requestedMetrics;
        }
    }
}