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

        public virtual void DefinePlugin()
        {
        }

        public virtual void Collect()
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