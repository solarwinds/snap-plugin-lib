namespace SnapPluginLib
{
    public abstract class PluginBase
    {
        protected string Name;
        protected string Version;

        PluginBase(string name, string version)
        {
            Name = name;
            Version = version;
        }

        public void DefinePlugin()
        {
        }

        public void Collect()
        {
        }

        public void Load()
        {
        }

        public void Unload()
        {
        }
    }
}