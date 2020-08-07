using System;
using System.Collections.Generic;

namespace SnapPluginLib
{
    public interface ICollectContext : IContext
    {
        void AddMetric(string ns, double value);
        void AlwaysApply(string ns);
        void DismissAllModifiers();
        void ShouldProcess(string ns);
        IList<string> RequestedMetrics();
    }
}