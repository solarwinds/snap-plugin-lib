namespace SnapPluginLib
{
    public interface ICollectorPlugin : IPlugin

    {
        void Collect(ICollectContext ctx);
    }
}