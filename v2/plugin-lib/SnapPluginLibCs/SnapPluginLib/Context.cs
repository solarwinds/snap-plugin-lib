using System;
using System.Collections.Generic;
using System.Runtime.CompilerServices;
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
        private static extern IntPtr ctx_config(string taskId, string key);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern IntPtr ctx_config_keys(string taskId);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern IntPtr ctx_raw_config(string taskId);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern string ctx_add_warning(string taskId, string message);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern string
            ctx_log(string taskId, int level, string message, LogMap fields); // todo: adamik: map

        public string Config(string key)
        {
            var ptr = ctx_config(_taskId, key);

            if (ptr != IntPtr.Zero)
                return Marshal.PtrToStringAnsi(ptr);

            return "";
        }

        public IList<string> ConfigKeys()
        {
            var ptr = ctx_config_keys(_taskId);
            return new List<string>(); // todo: how to process char **
        }

        public string RawConfig()
        {
            var ptr = ctx_raw_config(_taskId);

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
            ctx_add_warning(_taskId, message);
        }

        public void Log(int level, string message, Dictionary<string, string> fields)
        {
            var m = new LogMap();
            m.length = 2;
            m.elements = new LogMapElements();
            m.elements.key = "field1";
            m.elements.value = "value1";

            ctx_log(_taskId, 1, message, m);
        }
    }

    [StructLayout(LayoutKind.Sequential)]
    public class LogMap
    {
        public LogMapElements elements;
        public int length;
    }

    [StructLayout(LayoutKind.Sequential)]
    public class LogMapElements
    {
        public string key;
        public string value;
    }
}