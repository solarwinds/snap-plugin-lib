using System;
using System.Collections.Generic;
using System.Runtime.InteropServices;

namespace SnapPluginLib
{
    internal class CollectContext : Context, ICollectContext
    {
        public CollectContext(string taskId) : base(taskId)
        {
        }

        public void AddMetric(string ns, double value)
        {
        }

        public void AlwaysApply(string ns)
        {
        }

        public void DismissAllModifiers()
        {
        }

        public void ShouldProcess(string ns)
        {
        }

        public IList<string> RequestedMetrics()
        {
            return new List<string>();
        }

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern void ctx_add_metric(string taskId, string ns, Value value, Modifiers modifiers);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern void ctx_always_apply(string taskId, string ns, Modifiers modifiers);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern void ctx_dismiss_all_modifiers(string taskId);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern int ctx_should_process(string taskId, string ns);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern IList<string> ctx_requested_metrics(string taskId);
    }
}