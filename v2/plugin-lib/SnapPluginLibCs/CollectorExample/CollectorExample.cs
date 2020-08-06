using System;
using System.Collections.Generic;
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

        public override void Collect(IContext ctx)
        {
            Console.WriteLine("C# Collect executed");

            var cTest = ctx.Config("test");
            Console.WriteLine($"USERCOLLECT:ctxConfig : {cTest}");

            var rCfg = ctx.RawConfig();
            Console.WriteLine($"Raw confgi: {rCfg}");

            var keys = ctx.ConfigKeys();
            Console.WriteLine($"Config keys: {keys}");

            ctx.AddWarning("Warning from C#");

            ctx.Log(0, "Log message from C#", new Dictionary<string, string>()
                {{"language", "c#"}, {"function", "Collect"}});
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