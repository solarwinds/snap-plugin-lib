using System;
using System.Collections.Generic;

namespace SnapPluginLib
{
    public interface ICollectContext : IContext
    {
        void AddMetric(string ns, double value, params IPublicModifier[] modifiers);
        void AddMetric(string ns, int value, params IPublicModifier[] modifiers);
        void AddMetric(string ns, uint value, params IPublicModifier[] modifiers);
        void AddMetric(string ns, bool value, params IPublicModifier[] modifiers);
        void AlwaysApply(string ns, params IPublicModifier[] modifiers);
        void DismissAllModifiers();
        bool ShouldProcess(string ns);
        IList<string> RequestedMetrics();
    }
}