using System;
using System.Collections.Generic;
using System.Runtime.InteropServices;

namespace SnapPluginLib
{
    public class Context : IContext
    {
        protected string TaskId { get; }
        private Dictionary<string, object> _storedObjects;

        public Context(string taskId)
        {
            TaskId = taskId;
            _storedObjects = new Dictionary<string, object>();
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
            if (_storedObjects.ContainsKey(key))
            {
                _storedObjects.Remove(key);
            }
            
            _storedObjects.Add(key, obj);
        }

        public T Load<T>(string key)
        {
            try
            {
                return (T) _storedObjects[key];
            }
            catch (Exception e)
            {
                throw new PluginLibException($"Can't load object ({key}): " + e.Message);
            }
        }

        public void AddWarning(string message)
        {
            CBridge.ctx_add_warning(TaskId, message);
        }


        public void Log(LogLevel level, string message, Dictionary<string, string> fields)
        {
            var nativeMapPtr = Convertions.DictionaryToNativeMapMem(fields);
            CBridge.ctx_log(TaskId, (int) level, message, nativeMapPtr);
            Memory.FreeNativeMap(nativeMapPtr);
        }
    }
}