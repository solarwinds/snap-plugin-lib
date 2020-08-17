using System;
using System.Collections.Generic;
using System.Runtime.InteropServices;

namespace SnapPluginLib
{
    public class Context : IContext
    {
        protected string TaskId { get; }

        public Context(string taskId)
        {
            TaskId = taskId;
        }


        public string Config(string key)
        {
            var ptr = CBridge.ctx_config(TaskId, key);

            if (ptr != IntPtr.Zero)
                return Marshal.PtrToStringAnsi(ptr);

            return "";
        }

        public IList<string> ConfigKeys()
        {
            return Convertions.NativeStringArrayToList(CBridge.ctx_config_keys(TaskId));
        }

        public string RawConfig()
        {
            var ptr = CBridge.ctx_raw_config(TaskId);

            if (ptr != IntPtr.Zero)
                return Marshal.PtrToStringAnsi(ptr);

            return "";
        }

        public void Store(string key, object obj)
        {
            // todo: adamik: implement
        }

        public object Load(string key)
        {
            // todo: adamik: implement
            return null;
        }

        public void AddWarning(string message)
        {
            CBridge.ctx_add_warning(TaskId, message);
        }


        public void Log(LogLevel level, string message, Dictionary<string, string> fields)
        {
            CBridge.ctx_log(TaskId, (int) level, message, Convertions.DictionaryToNativeMapMem(fields));
        }
    }
}