using System;
using System.Collections.Generic;

namespace SnapPluginLib
{
    public interface ICollectContext
    {
        void AddMetric(string ns, double value); // todo: adamik: modifiers
        void AlwaysApply(string ns);
        void DismissAllModifiers();
        void ShouldProcess(string ns);
        IList<string> RequestedMetrics();
    }
}