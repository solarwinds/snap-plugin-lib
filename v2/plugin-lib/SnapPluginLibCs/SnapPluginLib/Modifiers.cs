using System.Collections.Generic;

namespace SnapPluginLib
{
    public interface IPublicModifier
    {
    }

    public class Modifiers
    {
        public static IPublicModifier Tags(Dictionary<string, string> tags)
        {
            return null;
        }

        public static IPublicModifier TagsToRemove(Dictionary<string, string> tags)
        {
            return null;
        }

        public static IPublicModifier Timestamp()
        {
            return null;
        }

        public static IPublicModifier Description(string description)
        {
            return null;
        }

        public static IPublicModifier Unit(string unit)
        {
            return null;
        }
    }


    internal class MetricTags : IModifier, IPublicModifier
    {
        private Dictionary<string, string> _tagsToAdd;

        public void Apply(NativeModifiers nModifier)
        {
        }
    }

    internal class MetricRemoveTags : IModifier, IPublicModifier
    {
        private Dictionary<string, string> _tagsToRemove;

        public void Apply(NativeModifiers nModifier)
        {
        }
    }

    internal class MetricTimestamp : IModifier, IPublicModifier
    {
        public void Apply(NativeModifiers nModifier)
        {
        }
    }

    internal class MetricDescription : IModifier, IPublicModifier
    {
        private string _description;

        public void Apply(NativeModifiers nModifier)
        {
        }
    }

    internal class MetricUnit : IModifier, IPublicModifier
    {
        private string _unit;

        public void Apply(NativeModifiers nModifier)
        {
        }
    }
}