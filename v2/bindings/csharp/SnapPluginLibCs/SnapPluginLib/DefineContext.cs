/*
 Copyright (c) 2021 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

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