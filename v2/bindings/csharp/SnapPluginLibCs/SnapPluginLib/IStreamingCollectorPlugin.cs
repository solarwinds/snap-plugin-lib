namespace SnapPluginLib
{
    public interface IStreamingCollectorPlugin : ICollectorPlugin
    {
        void StreamingCollect(ICollectContext ctx);
    }
}