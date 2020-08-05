using System;
using SnapPluginLib;

namespace CollectorExample
{
    public class CollectorExample : PluginBase
    {
        public CollectorExample(string name, string version) : base(name, version)
        {
        }

        public override void DefinePlugin(IDefineContext def)
        {
            Console.WriteLine("C# Define plugin");

            def.DefineMetric("/example/group1/metric1", "o", false, "metric description");
            def.DefineGroup("dynGroup", "dyn group description");
            def.DefineExampleConfig("{\"test\": \"a\"}");
            def.DefineInstancesLimit(5);
            def.DefineTaskPerInstanceLimit(10);
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