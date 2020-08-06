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
            Console.WriteLine("############## C# Log");
            
            var m = new LogMap();
            m.length = fields.Count;
            
            LogMapElements[] nativeElements = new LogMapElements[fields.Count];
            var i = 0;
            foreach (KeyValuePair<string,string> entry in fields)
            {
                nativeElements[i] = new LogMapElements();
                nativeElements[i].key = entry.Key;
                nativeElements[i].value = entry.Value;
                i++;
            }

            Console.WriteLine($"C# ADDR {m.elements}");
            
            ctx_log(_taskId, 1, message, m);
        }
    }

    [StructLayout(LayoutKind.Sequential)]
    public class LogMap
    {
        public IntPtr elements;
        public int length;
    }

    [StructLayout(LayoutKind.Sequential)]
    public class LogMapElements
    {
        public string key;
        public string value;
    }
}