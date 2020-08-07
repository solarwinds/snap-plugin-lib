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
            var ptr = ctx_config(TaskId, key);

            if (ptr != IntPtr.Zero)
                return Marshal.PtrToStringAnsi(ptr);

            return "";
        }

        public IList<string> ConfigKeys()
        {
            var ptr = ctx_config_keys(TaskId);
            return new List<string>(); // todo: how to process char **
        }

        public string RawConfig()
        {
            var ptr = ctx_raw_config(TaskId);

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
            ctx_add_warning(TaskId, message);
        }


        public void Log(int level, string message, Dictionary<string, string> fields)
        {
            var m = new LogMap();
            m.length = fields.Count;
            m.elements = Marshal.AllocHGlobal(Marshal.SizeOf(typeof(LogMapElements)) * fields.Count);

            var i = 0;
            foreach (KeyValuePair<string, string> entry in fields)
            {
                var nativeMapElem = new LogMapElements();
                nativeMapElem.key = entry.Key;
                nativeMapElem.value = entry.Value;

                Marshal.StructureToPtr(nativeMapElem,
                    (IntPtr) m.elements.ToInt64() + i * Marshal.SizeOf(typeof(LogMapElements)), false);

                i++;
            }

            ctx_log(TaskId, 1, message, m);

            Marshal.FreeHGlobal(m.elements);
        }

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern IntPtr ctx_config(string taskId, string key);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern IntPtr ctx_config_keys(string taskId);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern IntPtr ctx_raw_config(string taskId);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern string ctx_add_warning(string taskId, string message);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern string
            ctx_log(string taskId, int level, string message, LogMap fields);
    }

    
}