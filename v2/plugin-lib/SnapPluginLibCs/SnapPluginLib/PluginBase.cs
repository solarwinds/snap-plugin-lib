namespace SnapPluginLib
{
    public abstract class PluginBase
    {
        public string Name { get; }
        public string Version { get; }

        protected PluginBase(string name, string version)
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

        public virtual void Load()
        {
        }

        public virtual void Unload()
        {
        }
    }
}