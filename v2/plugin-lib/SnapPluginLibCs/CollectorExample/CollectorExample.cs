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
            // def.DefineMetric("/example/group1/metric1", "o", false, "metric description");
            // def.DefineGroup("dynGroup", "dyn group description");
            // def.DefineExampleConfig("{\"test\": \"a\"}");
            // def.DefineInstancesLimit(5);
            // def.DefineTaskPerInstanceLimit(10);
        }

        public override void Collect(ICollectContext ctx)
        {
            ctx.AddWarning("Warning from C#");

            ctx.Log(LogLevel.Info, "Log message from C#", new Dictionary<string, string>
            {
                {"language", "c#"}, {"function", "Collect"}
            });

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
                Modifiers.TagsToRemove(new Dictionary<string, string>
                {
                    {"origin", "C# lang"}
                }), // todo: adamik: why lib can't remove it
                Modifiers.Description("new custom description"),
                Modifiers.Unit("new custom unit")
                // Modifiers.Timestamp(DateTime.Now) // todo: adamik: doesn't work - why
            );

            ctx.AddMetric("/example/group1/m2", 20);
            ctx.AddMetric("/example/group1/m3", (uint) 30);

            if (ctx.ShouldProcess("/example/group2/m4"))
            {
                ctx.AddMetric("/example/group2/m4", true);
            }

            var reqMts = ctx.RequestedMetrics();
            if (reqMts.Count > 0)
            {
                Console.WriteLine("Requested metrics: ");
                foreach (var mt in reqMts)
                {
                    Console.WriteLine($"- {mt}");
                }
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
        }
    }
}