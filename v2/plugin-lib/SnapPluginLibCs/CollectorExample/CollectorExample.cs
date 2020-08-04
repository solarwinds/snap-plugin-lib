using System;
using SnapPluginLib;

namespace CollectorExample
{
    public class CollectorExample : PluginBase
    {
        public CollectorExample(string name, string version) : base(name, version)
        {
        }

        public override void DefinePlugin()
        {
            Console.WriteLine("Define plugin");
        }

        public override void Collect()
        {
            Console.WriteLine("Collect executed");
        }

        public override void Load()
        {
            Console.WriteLine("Load executed");
        }

        public override void Unload()
        {
            Console.WriteLine("Unload executed");
        }
    }
}