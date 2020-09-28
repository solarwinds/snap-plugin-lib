using System.Collections.Generic;
using System.Threading;

namespace SnapPluginLib
{
    public static class ContextMemory
    {
        private static Dictionary<string, ICollectContext> _contexts = new Dictionary<string, ICollectContext>();
        private static Mutex _mutex = new Mutex();

        public static ICollectContext Get(string id)
        {
            _mutex.WaitOne();

            if (!_contexts.ContainsKey(id))
            {
                _contexts.Add(id, new CollectContext(id));
            }

            var retValue = _contexts[id];

            _mutex.ReleaseMutex();

            return retValue;
        }

        public static void Clear(string id)
        {
            _mutex.WaitOne();

            if (_contexts.ContainsKey(id))
            {
                _contexts.Remove(id);
            }

            _mutex.ReleaseMutex();
        }
    }
}