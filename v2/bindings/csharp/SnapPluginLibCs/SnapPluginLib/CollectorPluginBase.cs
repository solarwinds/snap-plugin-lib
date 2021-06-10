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
    public class PluginBase
    {
        public string Name { get; }
        public Version Version { get; }

        protected PluginBase(string name, Version version)
        {
            Name = name;
            Version = version;
        }

        public virtual void DefinePlugin(IDefineContext ctx)
        {
        }

        public virtual void Load(IContext ctx)
        {
        }

        public virtual void Unload(IContext ctx)
        {
        }
    }

    public abstract class CollectorPluginBase : PluginBase
    {
        public CollectorPluginBase(string name, Version version) : base(name, version)
        {
        }

        public abstract void Collect(ICollectContext ctx);
    }

    public abstract class StreamingCollectorPluginBase : PluginBase
    {
        public StreamingCollectorPluginBase(string name, Version version) : base(name, version)
        {
        }

        public abstract void StreamingCollect(ICollectContext ctx);
    }
}