using System.Collections.Generic;

namespace SnapPluginLib
{
    internal class CollectContext : Context, ICollectContext
    {
        private enum ValueType
        {
            TypeInt64 = 1,
            TypeUint64,
            TypeDouble,
            TypeBool,
        }

        public CollectContext(string taskId) : base(taskId)
        {
        }

        public void AddMetric(string ns, int value, params IPublicModifier[] modifiers)
        {
            var nativeValue = new NativeValue
            {
                v_int64 = value,
                vtype = (int) ValueType.TypeInt64
            };

            AddMetric(ns, nativeValue, modifiers);
        }

        public void AddMetric(string ns, uint value, params IPublicModifier[] modifiers)
        {
            var nativeValue = new NativeValue
            {
                v_uint64 = value,
                vtype = (int) ValueType.TypeUint64
            };

            AddMetric(ns, nativeValue, modifiers);
        }

        public void AddMetric(string ns, double value, params IPublicModifier[] modifiers)
        {
            var nativeValue = new NativeValue
            {
                v_double = value,
                vtype = (int) ValueType.TypeDouble
            };

            AddMetric(ns, nativeValue, modifiers);
        }

        public void AddMetric(string ns, bool value, params IPublicModifier[] modifiers)
        {
            var nativeValue = new NativeValue
            {
                v_bool = value ? 1 : 0,
                vtype = (int) ValueType.TypeBool
            };

            AddMetric(ns, nativeValue, modifiers);
        }

        private void AddMetric(string ns, NativeValue nativeValue, params IPublicModifier[] modifiers)
        {
            var nativeModifiers = new NativeModifiers();

            foreach (var m in modifiers)
            {
                ((IModifier) m).Apply(nativeModifiers);
            }

            CBridge.ctx_add_metric(TaskId, ns, nativeValue, nativeModifiers);
        }

        public void AlwaysApply(string ns)
        {
        }

        public void DismissAllModifiers()
        {
        }

        public void ShouldProcess(string ns)
        {
        }

        public IList<string> RequestedMetrics()
        {
            return new List<string>();
        }
    }
}