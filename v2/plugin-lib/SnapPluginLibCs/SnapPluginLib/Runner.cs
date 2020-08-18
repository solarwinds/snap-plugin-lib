using System.Runtime.InteropServices;

namespace SnapPluginLib
{
    public static class Runner
    {
        private static PluginBase _collector;

        internal delegate void DefineHandler();

        internal delegate void CollectHandler(string taskId);

        internal delegate void LoadHandler(string taskId);

        internal delegate void UnloadHandler(string taskId);

        private static void DefineHandlerFn()
        {
            _collector.DefinePlugin(new DefineContext());
        }

        private static void CollectHandlerFn(string taskId)
        {
            _collector.Collect(ContextMemory.Get(taskId));
        }

        private static void LoadHandlerFn(string taskId)
        {
            _collector.Load(ContextMemory.Get(taskId));
        }

        private static void UnloadHandlerFn(string taskId)
        {
            _collector.Unload(ContextMemory.Get(taskId));
        }

        public static void StartCollector(PluginBase collector)
        {
            _collector = collector;

            CBridge.start_collector(
                CollectHandlerFn, LoadHandlerFn, UnloadHandlerFn, DefineHandlerFn,
                collector.Name, collector.Version);
        }
    }
}