using System;
using System.Collections.Generic;
using System.Runtime.InteropServices;

namespace SnapPluginLib
{
    public interface IPublicModifier
    {
    }

    public class Modifiers
    {
        public static IPublicModifier Tags(Dictionary<string, string> tags)
        {
            return new MetricTags(tags);
        }

        public static IPublicModifier TagsToRemove(Dictionary<string, string> tags)
        {
            return new MetricRemoveTags(tags);
        }

        public static IPublicModifier Timestamp()
        {
            return new MetricTimestamp();
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
            Console.WriteLine("MODIFIER APPLY");

            var m = new NativeMap();
            m.length = _tagsToAdd.Count;
            m.elements = Marshal.AllocHGlobal(Marshal.SizeOf(typeof(NativeMapElements)) * _tagsToAdd.Count);

            var i = 0;
            foreach (KeyValuePair<string, string> entry in _tagsToAdd)
            {
                var nativeMapElem = new NativeMapElements();
                nativeMapElem.key = entry.Key;
                nativeMapElem.value = entry.Value;

                Marshal.StructureToPtr(nativeMapElem,
                    (IntPtr) m.elements.ToInt64() + i * Marshal.SizeOf(typeof(NativeMapElements)), false);

                i++;
            }

            var all = Marshal.AllocHGlobal(Marshal.SizeOf(typeof(NativeMap)));
            Marshal.StructureToPtr(m, all, false);
            nModifier.tagsToAdd = all;
        }

        private Dictionary<string, string> _tagsToAdd;
    }

    internal class MetricRemoveTags : IModifier, IPublicModifier
    {
        public MetricRemoveTags(Dictionary<string, string> tags)
        {
            _tags = tags;
        }

        public void Apply(NativeModifiers nModifier)
        {
        }

        private Dictionary<string, string> _tags;
    }

    internal class MetricTimestamp : IModifier, IPublicModifier
    {
        public void Apply(NativeModifiers nModifier)
        {
        }
    }

    internal class MetricDescription : IModifier, IPublicModifier
    {
        public MetricDescription(string description)
        {
            _description = description;
        }

        public void Apply(NativeModifiers nModifier)
        {
            Console.WriteLine("MODIFIER DESCRIPTIOTN");

            nModifier.description = _description;
        }

        private string _description;
    }

    internal class MetricUnit : IModifier, IPublicModifier
    {
        public MetricUnit(string unit)
        {
            _unit = unit;
        }

        public void Apply(NativeModifiers nModifier)
        {
        }

        private string _unit;
    }
}