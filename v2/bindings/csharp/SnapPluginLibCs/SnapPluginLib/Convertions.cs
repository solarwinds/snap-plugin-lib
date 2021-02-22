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
                elements = Marshal.AllocHGlobal(Marshal.SizeOf(new NativeMapElements()) * dictionary.Count)
            };

            var i = 0;
            foreach (KeyValuePair<string, string> entry in dictionary)
            {
                var nativeMapElem = new NativeMapElements();
                nativeMapElem.key = entry.Key;
                nativeMapElem.value = entry.Value;

                Marshal.StructureToPtr(nativeMapElem,
                    (IntPtr) nativeMap.elements.ToInt64() + i * Marshal.SizeOf(new NativeMapElements()), false);

                i++;
            }

            var nativeMapAsMemBlock = Marshal.AllocHGlobal(Marshal.SizeOf(new NativeMap()));
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

        public static string ToSemanticVersion(Version version)
        {
            var buildVer = version.Build != -1 ? version.Build : 0;
            var semanticVersion = $"{version.Major}.{version.Minor}.{buildVer}";

            return semanticVersion;
        }
    }
}