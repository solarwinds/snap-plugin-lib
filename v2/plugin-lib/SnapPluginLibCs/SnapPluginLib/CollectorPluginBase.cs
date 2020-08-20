using System;

namespace SnapPluginLib
{
    public abstract class CollectorPluginBase
    {
        public string Name { get; }
        public Version Version { get; }

        protected CollectorPluginBase(string name, Version version)
        {
            Name = name;
            Version = version;
        }

        public virtual void DefinePlugin(IDefineContext ctx)
        {
        }

        public virtual void Collect(ICollectContext ctx)
        {
        }

        public virtual void Load(IContext ctx)
        {
        }

        public virtual void Unload(IContext ctx)
        {
        }
    }
}