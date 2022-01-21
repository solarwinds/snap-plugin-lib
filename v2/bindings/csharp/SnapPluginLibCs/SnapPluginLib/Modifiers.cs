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
using System.Collections.Generic;

namespace SnapPluginLib
{
    public abstract class Modifier : IModifier
    {
        internal Modifier()
        {
        }

        void IModifier.Apply(NativeModifiers nModifier)
        {
            Apply(nModifier);
        }

        internal abstract void Apply(NativeModifiers nModifier);
    }

    public static class Modifiers
    {
        public static Modifier Tags(Dictionary<string, string> tags)
        {
            return new MetricTags(tags);
        }

        public static Modifier TagsToRemove(Dictionary<string, string> tags)
        {
            return new MetricRemoveTags(tags);
        }

        public static Modifier Timestamp(DateTime dt)
        {
            return new MetricTimestamp(dt);
        }

        public static Modifier Description(string description)
        {
            return new MetricDescription(description);
        }

        public static Modifier Unit(string unit)
        {
            return new MetricUnit(unit);
        }
    }

    internal class MetricTags : Modifier
    {
        public MetricTags(Dictionary<string, string> tags)
        {
            _tagsToAdd = tags;
        }

        internal override void Apply(NativeModifiers nModifier)
        {
            nModifier.tagsToAdd = Convertions.DictionaryToNativeMapMem(_tagsToAdd);
        }

        private readonly Dictionary<string, string> _tagsToAdd;
    }

    internal class MetricRemoveTags : Modifier
    {
        public MetricRemoveTags(Dictionary<string, string> tags)
        {
            _tags = tags;
        }

        internal override void Apply(NativeModifiers nModifier)
        {
            nModifier.tagsToRemove = Convertions.DictionaryToNativeMapMem(_tags);
        }

        private readonly Dictionary<string, string> _tags;
    }

    internal class MetricTimestamp : Modifier
    {
        const int MilliToNanoFactor = (int) 1e6;

        public MetricTimestamp(DateTime timestamp)
        {
            _timestamp = timestamp;
        }

        internal override void Apply(NativeModifiers nModifier)
        {
            nModifier.timestamp = new NativeTimeWithNs
            {
                sec = _timestamp.Second,
                nsec = _timestamp.Millisecond * MilliToNanoFactor
            };
        }

        private readonly DateTime _timestamp;
    }

    internal class MetricDescription : Modifier
    {
        public MetricDescription(string description)
        {
            _description = description;
        }

        internal override void Apply(NativeModifiers nModifier)
        {
            nModifier.description = _description;
        }

        private readonly string _description;
    }

    internal class MetricUnit : Modifier
    {
        public MetricUnit(string unit)
        {
            _unit = unit;
        }

        internal override void Apply(NativeModifiers nModifier)
        {
            nModifier.unit = _unit;
        }

        private readonly string _unit;
    }

    internal class MetricType : Modifier
    {
        public MetricType(NativeMetricType type)
        {
            _type = type;
        }

        internal override void Apply(NativeModifiers nModifier)
        {
            nModifier.metric_type = (int)_type;
        }

        private readonly NativeMetricType _type;
    }
}