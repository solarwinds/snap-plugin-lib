using System;
using System.Collections.Generic;

namespace SnapPluginLib
{
    public interface IContext
    {
        string GetConfigValue(string configKey);
        IList<string> ConfigKeys();
        string RawConfig();

        void Store(string objectKey, Object obj);
        T Load<T>(string objectKey);

        void AddWarning(string message);
        void Log(LogLevel level, string message, Dictionary<string, string> fields);
    }
}