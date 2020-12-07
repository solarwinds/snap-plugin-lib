using System;

namespace SnapPluginLib
{
    public class Runner : IRunner
    {
        private ICollectorPlugin _collector;

        public void StartCollector(ICollectorPlugin collector)
        {
            _collector = collector ?? throw new ArgumentNullException(nameof(collector));

            CBridge.start_collector(
                InvokeCollect, InvokeLoad, InvokeUnload, InvokeDefine,
                collector.Name, Convertions.ToSemanticVersion(collector.Version));
        }

        public void StartStreamingCollector(IStreamingCollectorPlugin collector)
        {
            _collector = collector ?? throw new ArgumentNullException(nameof(collector));

            CBridge.start_streaming_collector(
                InvokeStreamingCollect, InvokeLoad, InvokeUnload, InvokeDefine,
                collector.Name, Convertions.ToSemanticVersion(collector.Version));
        }


        private void InvokeDefine()
        {
            _collector.DefinePlugin(new DefineContext());
        }

        private void InvokeCollect(string taskId)
        {
            _collector.Collect(ContextMemory.Get(taskId));
        }

        private void InvokeStreamingCollect(string taskId)
        {
            (_collector as IStreamingCollectorPlugin)?.StreamingCollect(ContextMemory.Get(taskId));
        }

        private void InvokeLoad(string taskId)
        {
            _collector.Load(ContextMemory.Get(taskId));
        }

        private void InvokeUnload(string taskId)
        {
            _collector.Unload(ContextMemory.Get(taskId));
        }

    }
}