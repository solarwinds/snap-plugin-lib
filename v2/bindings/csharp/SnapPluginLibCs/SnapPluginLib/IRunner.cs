namespace SnapPluginLib
{
    public interface IRunner
    {
        void StartCollector(ICollectorPlugin collector);
        void StartStreamingCollector(IStreamingCollectorPlugin streamingCollector);
    }
}