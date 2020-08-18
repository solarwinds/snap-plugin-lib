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
            var exampleConfigYaml = "value: 113\nother: 142";

            // def.DefineMetric("/example/group1/metric1", "o", false, "metric description");
            // def.DefineGroup("dynGroup", "dyn group description");
            def.DefineExampleConfig(exampleConfigYaml);
            // def.DefineInstancesLimit(5);
            // def.DefineTaskPerInstanceLimit(10);
        }

        public override void Collect(ICollectContext ctx)
        {
            
            // Load object from memory
            var obj = ctx.Load<Dictionary<string, int>>("stored_object");
            obj["counter"]++;
            ctx.Store("stored_object", obj);
            
            // Log messages with information from stored object
            ctx.Log(LogLevel.Info, "Log message from C#", new Dictionary<string, string>
            {
                {"language", "c#"}, 
                {"counter", $"{obj["counter"]}"}
            });
            ctx.AddWarning("Warning from C#");
            
            // List requested metrics
            var reqMts = ctx.RequestedMetrics();
            if (reqMts.Count > 0)
            {
                Console.WriteLine("Requested metrics: ");
                foreach (var mt in reqMts)
                {
                    Console.WriteLine($"- {mt}");
                }
            }

            // Add metrics
            ctx.AlwaysApply("/example/group1/*", Modifiers.Tags(new Dictionary<string, string>
            {
                {"virtualization", "VirtualBox"}
            }));

            ctx.AddMetric("/example/group1/metric1", 12.4,
                Modifiers.Tags(new Dictionary<string, string>
                {
                    {"origin", "C# lang"},
                    {"system", "Windows"}
                }),
                Modifiers.Description("new custom description")
            );

            ctx.AddMetric("/example/group1/metric2", 20);
            ctx.AddMetric("/example/group1/metric3", (uint) 30);

            if (ctx.ShouldProcess("/example/group2/metric4"))
            {
                ctx.AddMetric("/example/group2/metric4", true);
            }
        }

        public override void Load(IContext ctx)
        {
            // Raw Config
            var rCfg = ctx.RawConfig();
            Console.WriteLine($"Raw config: {rCfg}");

            // List config keys
            var keys = ctx.ConfigKeys();
            if (keys.Count > 0)
            {
                Console.WriteLine($"Config keys:");
                foreach (var k in keys)
                {
                    Console.WriteLine($"- {k}");
                }

                Console.WriteLine();
            }

            // Store object for a later use
            var obj = new Dictionary<string, int>
            {
                {"iteration", 20},
                {"counter", 0}
            };
            ctx.Store("stored_object", obj);
        }
    }
}