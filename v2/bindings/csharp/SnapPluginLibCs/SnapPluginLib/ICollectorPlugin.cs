using System;

namespace SnapPluginLib
{
    public interface ICollectorPlugin
    {
        void DefinePlugin(IDefineContext ctx);
        void Collect(ICollectContext ctx);
        void Load(IContext ctx);
        void Unload(IContext ctx);
        string Name { get; }
        Version Version { get; }
    }
}