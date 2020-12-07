using System;

namespace SnapPluginLib
{
    public interface IPlugin
    {
        void DefinePlugin(IDefineContext ctx);
        void Load(IContext ctx);
        void Unload(IContext ctx);
        string Name { get; }
        Version Version { get; }
    }
}