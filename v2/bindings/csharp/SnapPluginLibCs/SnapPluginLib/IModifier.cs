using System.Collections.Generic;

namespace SnapPluginLib
{
    internal interface IModifier
    {
        void Apply(NativeModifiers nModifier);
    }
}