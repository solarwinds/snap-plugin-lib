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