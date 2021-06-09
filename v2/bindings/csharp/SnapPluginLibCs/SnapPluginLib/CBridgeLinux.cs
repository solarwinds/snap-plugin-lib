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
    /*
     * Responsible for calling all exported (native) C functions 
     */
    internal static class CBridgeLinux
    {
        private const string PluginLibDllName = "swisnap-plugin-lib.so";

        // Runner

        [DllImport(PluginLibDllName, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern void start_collector(
            Runner.CollectHandler collectHandler,
            Runner.LoadHandler loadHandler,
            Runner.UnloadHandler unloadHandler,
            Runner.DefineHandler defineHandler,
            string name,
            string version
        );
        
        [DllImport(PluginLibDllName, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern void start_streaming_collector(
            Runner.CollectHandler collectHandler,
            Runner.LoadHandler loadHandler,
            Runner.UnloadHandler unloadHandler,
            Runner.DefineHandler defineHandler,
            string name,
            string version
        );

        // Collect context related functions

        [DllImport(PluginLibDllName, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern IntPtr /* NativeError */
            ctx_add_metric(string taskId, string ns, NativeValue nativeValue, NativeModifiers nativeModifiers);

        [DllImport(PluginLibDllName, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern IntPtr /* NativeError */
            ctx_always_apply(string taskId, string ns, NativeModifiers nativeModifiers);

        [DllImport(PluginLibDllName, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern void ctx_dismiss_all_modifiers(string taskId);

        [DllImport(PluginLibDllName, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern int ctx_should_process(string taskId, string ns);

        [DllImport(PluginLibDllName, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern IntPtr ctx_requested_metrics(string taskId);

        // Context related functions

        [DllImport(PluginLibDllName, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern IntPtr ctx_config_value(string taskId, string key);

        [DllImport(PluginLibDllName, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern IntPtr ctx_config_keys(string taskId);

        [DllImport(PluginLibDllName, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern IntPtr ctx_raw_config(string taskId);

        [DllImport(PluginLibDllName, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern string ctx_add_warning(string taskId, string message);

        [DllImport(PluginLibDllName, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern string ctx_log(string taskId, int level, string message, IntPtr /* NativeMap */ fields);

        // DefinePlugin related functions 

        [DllImport(PluginLibDllName, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern void define_metric(string ns, string unit, int idDefault, string description);

        [DllImport(PluginLibDllName, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern void define_group(string name, string description);

        [DllImport(PluginLibDllName, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern IntPtr /* NativeError */ define_example_config(string config);

        [DllImport(PluginLibDllName, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern void define_tasks_per_instance_limit(int limit);

        [DllImport(PluginLibDllName, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern void define_instances_limit(int limit);
    }
}