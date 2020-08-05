using System;
using System.Collections.Generic;
using System.Runtime.InteropServices;

namespace SnapPluginLib
{
    public class Context : IContext
    {
        private string _taskId;

        public Context(string taskId)
        {
            _taskId = taskId;
        }

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern void ctx_config(string taskId, string key);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern string ctx_config_keys(string taskId);
        
        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern string ctx_raw_config(string taskId);
        
        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern string ctx_add_warning(string taskId, string message);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern string ctx_log(string taskId, int level, string message); // todo: adamik: map

        public string Config(string key)
        {
            ctx_config(_taskId, key); // todo: adamik: fix return value
            return "";
        }

        public IList<string> ConfigKeys()
        {
            ctx_config_keys(_taskId);
            var val = new List<string>();
            return val;
        }

        public string RawConfig()
        {
            ctx_raw_config(_taskId);
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
            ctx_add_warning(_taskId, message);
        }

        public void Log(int level, string message, Dictionary<string, string> fields)
        {
            
        }
    }
}