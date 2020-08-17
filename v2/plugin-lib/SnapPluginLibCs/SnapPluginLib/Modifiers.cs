using System;
using System.Collections.Generic;

namespace SnapPluginLib
{
    public interface IPublicModifier
    {
    }

    public static class Modifiers
    {
        public static IPublicModifier Tags(Dictionary<string, string> tags)
        {
            return new MetricTags(tags);
        }

        public static IPublicModifier TagsToRemove(Dictionary<string, string> tags)
        {
            return new MetricRemoveTags(tags);
        }

        public static IPublicModifier Timestamp(DateTime dt)
        {
            return new MetricTimestamp(dt);
        }

        public static IPublicModifier Description(string description)
        {
            return new MetricDescription(description);
        }

        public static IPublicModifier Unit(string unit)
        {
            return new MetricUnit(unit);
        }
    }

    internal class MetricTags : IModifier, IPublicModifier
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

    internal class MetricRemoveTags : IModifier, IPublicModifier
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

    internal class MetricTimestamp : IModifier, IPublicModifier
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

    internal class MetricDescription : IModifier, IPublicModifier
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

    internal class MetricUnit : IModifier, IPublicModifier
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