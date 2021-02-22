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
    public static class Runner
    {
        private static CollectorPluginBase _collector;

        internal delegate void DefineHandler();

        internal delegate void CollectHandler(string taskId);

        internal delegate void LoadHandler(string taskId);

        internal delegate void UnloadHandler(string taskId);

        private static void DefineHandlerFn()
        {
            _collector.DefinePlugin(new DefineContext());
        }

        private static void CollectHandlerFn(string taskId)
        {
            _collector.Collect(ContextMemory.Get(taskId));
        }

        private static void LoadHandlerFn(string taskId)
        {
            _collector.Load(ContextMemory.Get(taskId));
        }

        private static void UnloadHandlerFn(string taskId)
        {
            _collector.Unload(ContextMemory.Get(taskId));
        }

        public static void StartCollector(CollectorPluginBase collector)
        {
            _collector = collector;

            CBridge.start_collector(
                CollectHandlerFn, LoadHandlerFn, UnloadHandlerFn, DefineHandlerFn,
                collector.Name, Convertions.ToSemanticVersion(collector.Version));
        }
    }
}