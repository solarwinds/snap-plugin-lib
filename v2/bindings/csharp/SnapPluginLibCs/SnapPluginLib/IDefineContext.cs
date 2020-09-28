namespace SnapPluginLib
{
    public interface IDefineContext
    {
        void DefineMetric(string ns, string unit, bool isDefault, string description);
        void DefineGroup(string name, string description);
        void DefineExampleConfig(string config);
        void DefineTaskPerInstanceLimit(int limit);
        void DefineInstancesLimit(int limit);
    }
}