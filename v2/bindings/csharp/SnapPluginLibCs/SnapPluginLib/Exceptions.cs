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