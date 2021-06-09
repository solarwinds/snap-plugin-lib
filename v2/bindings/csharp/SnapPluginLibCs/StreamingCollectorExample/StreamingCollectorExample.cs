using System;
using System.Threading;
using SnapPluginLib;

namespace StreamingCollectorExample
{
    public class StreamingCollectorExample : StreamingCollectorPluginBase
    {
        public StreamingCollectorExample(string name, Version version) : base(name, version)
        {
        }

        public override void StreamingCollect(ICollectContext ctx)
        {
            var metricValue = 0;
            for (;;)
            {
                metricValue += 2;
                ctx.AddMetric("/example/group1/metric1", metricValue);
                
                Thread.Sleep(new TimeSpan(0, 0, 3));
            }
        }
    }
}