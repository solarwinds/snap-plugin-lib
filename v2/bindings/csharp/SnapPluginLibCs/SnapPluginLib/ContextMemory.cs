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

using System.Collections.Generic;
using System.Threading;

namespace SnapPluginLib
{
    public static class ContextMemory
    {
        private static Dictionary<string, ICollectContext> _contexts = new Dictionary<string, ICollectContext>();
        private static Mutex _mutex = new Mutex();

        public static ICollectContext Get(string id)
        {
            _mutex.WaitOne();

            if (!_contexts.ContainsKey(id))
            {
                _contexts.Add(id, new CollectContext(id));
            }

            var retValue = _contexts[id];

            _mutex.ReleaseMutex();

            return retValue;
        }

        public static void Clear(string id)
        {
            _mutex.WaitOne();

            if (_contexts.ContainsKey(id))
            {
                _contexts.Remove(id);
            }

            _mutex.ReleaseMutex();
        }
    }
}