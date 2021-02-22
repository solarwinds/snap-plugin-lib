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

using System;

namespace SnapPluginLib
{
    public abstract class CollectorPluginBase
    {
        public string Name { get; }
        public Version Version { get; }

        protected CollectorPluginBase(string name, Version version)
        {
            Name = name;
            Version = version;
        }

        public virtual void DefinePlugin(IDefineContext ctx)
        {
        }

        public virtual void Collect(ICollectContext ctx)
        {
        }

        public virtual void Load(IContext ctx)
        {
        }

        public virtual void Unload(IContext ctx)
        {
        }
    }
}