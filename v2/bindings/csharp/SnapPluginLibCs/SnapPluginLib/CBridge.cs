using System;

namespace SnapPluginLib
{
    /*
     * Responsible for calling all exported (native) C functions 
     */
    internal static class CBridge
    {
        // Runner

        internal static void start_collector(
            Runner.CollectHandler collectHandler, Runner.LoadHandler loadHandler, Runner.UnloadHandler unloadHandler,
            Runner.DefineHandler defineHandler, string name, string version
        )
        {
            if (IsWindows())
            {
                CBridgeWin.start_collector(collectHandler, loadHandler, unloadHandler, defineHandler, name, version);
            }
            else if (IsLinux())
            {
                CBridgeLinux.start_collector(collectHandler, loadHandler, unloadHandler, defineHandler, name, version);
            }
            else
            {
                throw new NotImplementedException(NoImplementedError);
            }
        }

        // Collect context related functions

        internal static IntPtr /* NativeError */
            ctx_add_metric(string taskId, string ns, NativeValue nativeValue, NativeModifiers nativeModifiers)
        {
            if (IsWindows())
            {
                return CBridgeWin.ctx_add_metric(taskId, ns, nativeValue, nativeModifiers);
            }

            if (IsLinux())
            {
                return CBridgeLinux.ctx_add_metric(taskId, ns, nativeValue, nativeModifiers);
            }

            throw new NotImplementedException(NoImplementedError);
        }

        internal static IntPtr /* NativeError */
            ctx_always_apply(string taskId, string ns, NativeModifiers nativeModifiers)
        {
            if (IsWindows())
            {
                return CBridgeWin.ctx_always_apply(taskId, ns, nativeModifiers);
            }

            if (IsLinux())
            {
                return CBridgeLinux.ctx_always_apply(taskId, ns, nativeModifiers);
            }

            throw new NotImplementedException(NoImplementedError);
        }

        internal static void ctx_dismiss_all_modifiers(string taskId)
        {
            if (IsWindows())
            {
                CBridgeWin.ctx_dismiss_all_modifiers(taskId);
            }
            else if (IsLinux())
            {
                CBridgeLinux.ctx_dismiss_all_modifiers(taskId);
            }
            else
            {
                throw new NotImplementedException(NoImplementedError);
            }
        }

        internal static int ctx_should_process(string taskId, string ns)
        {
            if (IsWindows())
            {
                return CBridgeWin.ctx_should_process(taskId, ns);
            }

            if (IsLinux())
            {
                return CBridgeLinux.ctx_should_process(taskId, ns);
            }

            throw new NotImplementedException(NoImplementedError);
        }

        internal static IntPtr ctx_requested_metrics(string taskId)
        {
            if (IsWindows())
            {
                return CBridgeWin.ctx_requested_metrics(taskId);
            }

            if (IsLinux())
            {
                return CBridgeLinux.ctx_requested_metrics(taskId);
            }

            throw new NotImplementedException(NoImplementedError);
        }

        // Context related functions

        internal static IntPtr ctx_config_value(string taskId, string key)
        {
            if (IsWindows())
            {
                return CBridgeWin.ctx_config_value(taskId, key);
            }

            if (IsLinux())
            {
                return CBridgeLinux.ctx_config_value(taskId, key);
            }

            throw new NotImplementedException(NoImplementedError);
        }

        internal static IntPtr ctx_config_keys(string taskId)
        {
            if (IsWindows())
            {
                return CBridgeWin.ctx_config_keys(taskId);
            }

            if (IsLinux())
            {
                return CBridgeLinux.ctx_config_keys(taskId);
            }

            throw new NotImplementedException(NoImplementedError);
        }

        internal static IntPtr ctx_raw_config(string taskId)
        {
            if (IsWindows())
            {
                return CBridgeWin.ctx_raw_config(taskId);
            }

            if (IsLinux())
            {
                return CBridgeLinux.ctx_raw_config(taskId);
            }

            throw new NotImplementedException(NoImplementedError);
        }

        internal static string ctx_add_warning(string taskId, string message)
        {
            if (IsWindows())
            {
                return CBridgeWin.ctx_add_warning(taskId, message);
            }

            if (IsLinux())
            {
                return CBridgeLinux.ctx_add_warning(taskId, message);
            }

            throw new NotImplementedException(NoImplementedError);
        }

        internal static string ctx_log(string taskId, int level, string message, IntPtr /* NativeMap */ fields)
        {
            if (IsWindows())
            {
                return CBridgeWin.ctx_log(taskId, level, message, fields);
            }

            if (IsLinux())
            {
                return CBridgeLinux.ctx_log(taskId, level, message, fields);
            }

            throw new NotImplementedException(NoImplementedError);
        }

        // DefinePlugin related functions 

        internal static void define_metric(string ns, string unit, int isDefault, string description)
        {
            if (IsWindows())
            {
                CBridgeWin.define_metric(ns, unit, isDefault, description);
            }
            else if (IsLinux())
            {
                CBridgeLinux.define_metric(ns, unit, isDefault, description);
            }
            else
            {
                throw new NotImplementedException(NoImplementedError);
            }
        }

        internal static void define_group(string name, string description)
        {
            if (IsWindows())
            {
                CBridgeWin.define_group(name, description);
            }
            else if (IsLinux())
            {
                CBridgeLinux.define_group(name, description);
            }
            else
            {
                throw new NotImplementedException(NoImplementedError);
            }
        }

        internal static IntPtr /* NativeError */ define_example_config(string config)
        {
            if (IsWindows())
            {
                return CBridgeWin.define_example_config(config);
            }

            if (IsLinux())
            {
                return CBridgeLinux.define_example_config(config);
            }

            throw new NotImplementedException(NoImplementedError);
        }

        internal static void define_tasks_per_instance_limit(int limit)
        {
            if (IsWindows())
            {
                CBridgeWin.define_tasks_per_instance_limit(limit);
            }
            else if (IsLinux())
            {
                CBridgeLinux.define_tasks_per_instance_limit(limit);
            }
            else
            {
                throw new NotImplementedException(NoImplementedError);
            }
        }

        internal static void define_instances_limit(int limit)
        {
            if (IsWindows())
            {
                CBridgeWin.define_instances_limit(limit);
            }
            else if (IsLinux())
            {
                CBridgeLinux.define_instances_limit(limit);
            }
            else
            {
                throw new NotImplementedException(NoImplementedError);
            }
        }

        private static bool IsWindows()
        {
            return Environment.OSVersion.Platform == PlatformID.Win32Windows
                   || Environment.OSVersion.Platform == PlatformID.Win32NT;
        }

        private static bool IsLinux()
        {
            return Environment.OSVersion.Platform == PlatformID.Unix;
        }

        private static string NoImplementedError = $"no implementation for {Environment.OSVersion.Platform.ToString()}";
    }
}