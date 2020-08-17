using System;
using System.Collections.Generic;
using System.Runtime.InteropServices;

namespace SnapPluginLib
{
    /*
     * Responsible for calling all exported (native) C functions 
     */
    internal static class CBridge
    {
        private const string PluginLibDllPath = "plugin-lib.dll";

        // Collect context related functions

        [DllImport(PluginLibDllPath, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern void ctx_add_metric(string taskId, string ns, NativeValue nativeValue,
            NativeModifiers nativeModifiers);

        [DllImport(PluginLibDllPath, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern void ctx_always_apply(string taskId, string ns, NativeModifiers nativeModifiers);

        [DllImport(PluginLibDllPath, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern void ctx_dismiss_all_modifiers(string taskId);

        [DllImport(PluginLibDllPath, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern int ctx_should_process(string taskId, string ns);

        [DllImport(PluginLibDllPath, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern IntPtr ctx_requested_metrics(string taskId);

        // Context related functions

        [DllImport(PluginLibDllPath, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern IntPtr ctx_config(string taskId, string key);

        [DllImport(PluginLibDllPath, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern IntPtr ctx_config_keys(string taskId);

        [DllImport(PluginLibDllPath, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern IntPtr ctx_raw_config(string taskId);

        [DllImport(PluginLibDllPath, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern string ctx_add_warning(string taskId, string message);

        [DllImport(PluginLibDllPath, CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern string ctx_log(string taskId, int level, string message, IntPtr /* NativeMap */ fields);

        // DefinePlugin related functions 

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern void define_metric(string ns, string unit, int idDefault, string description);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern void define_group(string name, string description);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern void define_example_config(string config); // todo: adamik: handle return error

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern void define_tasks_per_instance_limit(int limit);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        internal static extern void define_instances_limit(int limit);
    }
}