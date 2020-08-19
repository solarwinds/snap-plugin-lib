using System;
using System.Collections.Generic;

namespace SnapPluginLib
{
    public abstract class Modifier
    {
        internal Modifier() {}
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

    internal class MetricTags : Modifier, IModifier 
    {
        public MetricTags(Dictionary<string, string> tags)
        {
            _tagsToAdd = tags;
        }

        public void Apply(NativeModifiers nModifier)
        {
            nModifier.tagsToAdd = Convertions.DictionaryToNativeMapMem(_tagsToAdd);
        }

        private readonly Dictionary<string, string> _tagsToAdd;
    }

    internal class MetricRemoveTags : Modifier, IModifier 
    {
        public MetricRemoveTags(Dictionary<string, string> tags)
        {
            _tags = tags;
        }

        public void Apply(NativeModifiers nModifier)
        {
            nModifier.tagsToRemove = Convertions.DictionaryToNativeMapMem(_tags);
        }

        private readonly Dictionary<string, string> _tags;
    }

    internal class MetricTimestamp : Modifier, IModifier 
    {
        const int MilliToNanoFactor = (int) 1e6;

        public MetricTimestamp(DateTime timestamp)
        {
            _timestamp = timestamp;
        }

        public void Apply(NativeModifiers nModifier)
        {
            nModifier.timestamp = new NativeTimeWithNs
            {
                sec = _timestamp.Second,
                nsec = _timestamp.Millisecond * MilliToNanoFactor
            };
        }

        private readonly DateTime _timestamp;
    }

    internal class MetricDescription : Modifier, IModifier 
    {
        public MetricDescription(string description)
        {
            _description = description;
        }

        public void Apply(NativeModifiers nModifier)
        {
            nModifier.description = _description;
        }

        private readonly string _description;
    }

    internal class MetricUnit : Modifier, IModifier 
    {
        public MetricUnit(string unit)
        {
            _unit = unit;
        }

        public void Apply(NativeModifiers nModifier)
        {
            nModifier.unit = _unit;
        }

        private readonly string _unit;
    }
}