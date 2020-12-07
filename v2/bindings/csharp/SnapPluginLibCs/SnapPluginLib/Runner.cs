using System;

namespace SnapPluginLib
{
    public class Runner : IRunner
    {
        private IPlugin _plugin;

        public void StartCollector(ICollectorPlugin collector)
        {
            _plugin = collector ?? throw new ArgumentNullException(nameof(collector));

            CBridge.start_collector(
                InvokeCollect, InvokeLoad, InvokeUnload, InvokeDefine,
                collector.Name, Convertions.ToSemanticVersion(collector.Version));
        }

        public void StartStreamingCollector(IStreamingCollectorPlugin collector)
        {
            _plugin = collector ?? throw new ArgumentNullException(nameof(collector));

            CBridge.start_streaming_collector(
                InvokeStreamingCollect, InvokeLoad, InvokeUnload, InvokeDefine,
                collector.Name, Convertions.ToSemanticVersion(collector.Version));
        }


        private void InvokeDefine()
        {
            _plugin.DefinePlugin(new DefineContext());
        }

        private void InvokeCollect(string taskId)
        {
            (_plugin as ICollectorPlugin)?.Collect(ContextMemory.Get(taskId));
        }

        private void InvokeStreamingCollect(string taskId)
        {
            (_plugin as IStreamingCollectorPlugin)?.StreamingCollect(ContextMemory.Get(taskId));
        }

        private void InvokeLoad(string taskId)
        {
            _plugin.Load(ContextMemory.Get(taskId));
        }

        private void InvokeUnload(string taskId)
        {
            _plugin.Unload(ContextMemory.Get(taskId));
        }

    }
}