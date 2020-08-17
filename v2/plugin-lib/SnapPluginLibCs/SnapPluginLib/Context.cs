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


        public void Log(int level, string message, Dictionary<string, string> fields)
        {
            // todo: adamik: make more common with conversion function
            var m = new NativeMap();
            m.length = fields.Count;
            m.elements = Marshal.AllocHGlobal(Marshal.SizeOf(typeof(NativeMapElements)) * fields.Count);

            var i = 0;
            foreach (KeyValuePair<string, string> entry in fields)
            {
                var nativeMapElem = new NativeMapElements();
                nativeMapElem.key = entry.Key;
                nativeMapElem.value = entry.Value;

                Marshal.StructureToPtr(nativeMapElem,
                    (IntPtr) m.elements.ToInt64() + i * Marshal.SizeOf(typeof(NativeMapElements)), false);

                i++;
            }

            CBridge.ctx_log(TaskId, 1, message, m);

            Marshal.FreeHGlobal(m.elements);
        }


    }

    
}