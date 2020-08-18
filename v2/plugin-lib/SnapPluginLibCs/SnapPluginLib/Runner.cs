using System.Runtime.InteropServices;

namespace SnapPluginLib
{
    public static class Runner
    {
        private static PluginBase _collector;

        private delegate void DefineHandler();

        private delegate void CollectHandler(string taskId);

        private delegate void LoadHandler(string taskId);

        private delegate void UnloadHandler(string taskId);

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

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern void start_collector(
            CollectHandler collectHandler,
            LoadHandler loadHandler,
            UnloadHandler unloadHandler,
            DefineHandler defineHandler,
            string name,
            string version
        );

        public static void StartCollector(PluginBase collector)
        {
            _collector = collector;
            
            start_collector(
                CollectHandlerFn, LoadHandlerFn, UnloadHandlerFn, DefineHandlerFn,
                collector.Name, collector.Version);
        }
    }
}