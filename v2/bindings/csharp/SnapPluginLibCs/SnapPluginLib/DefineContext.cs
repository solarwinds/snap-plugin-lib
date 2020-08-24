using System;
using System.Runtime.InteropServices;

namespace SnapPluginLib
{
    internal class DefineContext : IDefineContext
    {
        public void DefineMetric(string ns, string unit, bool isDefault, string description)
        {
            CBridge.define_metric(ns, unit, isDefault ? 1 : 0, description);
        }

        public void DefineGroup(string name, string description)
        {
            CBridge.define_group(name, description);
        }

        public void DefineExampleConfig(string config)
        {
            var errPtr = CBridge.define_example_config(config);
            Exceptions.ThrowExceptionIfError(errPtr);
        }

        public void DefineTaskPerInstanceLimit(int limit)
        {
            CBridge.define_tasks_per_instance_limit(limit);
        }

        public void DefineInstancesLimit(int limit)
        {
            CBridge.define_instances_limit(limit);
        }
    }
}