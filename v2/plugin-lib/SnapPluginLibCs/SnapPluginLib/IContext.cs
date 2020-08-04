using System;
using System.Collections.Generic;

namespace SnapPluginLib
{
    public interface IContext
    {
        string Config(string key);
        IList<string> ConfigKeys();
        string RawConfig();

        void Store(string key, Object obj);
        Object Load(string key);

        void AddWarning(string message);
        void Log(int level, string message, Dictionary<string, string> fields);
    }
}