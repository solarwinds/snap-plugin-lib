namespace SnapPluginLib
{
    public interface IStreamingCollectorPlugin : IPlugin
    {
        void StreamingCollect(ICollectContext ctx);
        
    }
}