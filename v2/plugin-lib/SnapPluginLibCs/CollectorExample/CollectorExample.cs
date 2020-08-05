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
            Console.WriteLine("C# Define plugin");
        }

        public override void Collect()
        {
            Console.WriteLine("C# Collect executed");
        }

        public override void Load()
        {
            Console.WriteLine("C# Load executed");
        }

        public override void Unload()
        {
            Console.WriteLine("C# Unload executed");
        }
    }
}