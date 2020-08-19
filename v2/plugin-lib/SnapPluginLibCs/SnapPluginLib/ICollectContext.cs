using System;
using System.Collections.Generic;

namespace SnapPluginLib
{
    public interface ICollectContext : IContext
    {
        void AddMetric(string ns, double value, params Modifier[] modifiers);
        void AddMetric(string ns, int value, params Modifier[] modifiers);
        void AddMetric(string ns, uint value, params Modifier[] modifiers);
        void AddMetric(string ns, bool value, params Modifier[] modifiers);
        void AlwaysApply(string ns, params Modifier[] modifiers);
        void DismissAllModifiers();
        bool ShouldProcess(string ns);
        IList<string> RequestedMetrics();
    }
}