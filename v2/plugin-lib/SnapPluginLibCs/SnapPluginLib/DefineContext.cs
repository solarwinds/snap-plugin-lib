using System.Runtime.InteropServices;

namespace SnapPluginLib
{
    public class DefineContext : IDefineContext
    {
        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern void define_metric(string ns, string unit, int idDefault, string description);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern void define_group(string name, string description);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern void define_example_config(string config); // todo: adamik: handle return error

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern void define_tasks_per_instance_limit(int limit);

        [DllImport("plugin-lib.dll", CharSet = CharSet.Ansi, SetLastError = true)]
        private static extern void define_instances_limit(int limit);

        public void DefineMetric(string ns, string unit, bool isDefault, string description)
        {
            define_metric(ns, unit, 0, description);
        }

        public void DefineGroup(string name, string description)
        {
            define_group(name, description);
        }

        public void DefineExampleConfig(string config)
        {
            define_example_config(config);
        }

        public void DefineTaskPerInstanceLimit(int limit)
        {
            define_tasks_per_instance_limit(5);
        }

        public void DefineInstancesLimit(int limit)
        {
            define_instances_limit(5);
        }
    }
}