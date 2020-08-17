using System;
using System.Collections.Generic;
using System.Runtime.InteropServices;
using System.Threading.Tasks;

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
            CBridge.ctx_add_metric(TaskId, ns, nativeValue, ToNativeModifiers(modifiers));
        }

        public void AlwaysApply(string ns, params IPublicModifier[] modifiers)
        {
            CBridge.ctx_always_apply(TaskId, ns, ToNativeModifiers(modifiers));
        }

        public void DismissAllModifiers()
        {
            CBridge.ctx_dismiss_all_modifiers(TaskId);
        }

        public bool ShouldProcess(string ns)
        {
            return CBridge.ctx_should_process(TaskId, ns) > 0;
        }

        public IList<string> RequestedMetrics()
        {
            var ptr = CBridge.ctx_requested_metrics(TaskId);
            
            var ptr2 = Marshal.ReadIntPtr(ptr);

            var f1 = Marshal.PtrToStringAnsi(Marshal.ReadIntPtr(ptr));
            var f2 = Marshal.PtrToStringAnsi(Marshal.ReadIntPtr(ptr+8));
            var f3 = Marshal.PtrToStringAnsi(Marshal.ReadIntPtr(ptr+16));

            Console.WriteLine(f1);
            Console.WriteLine(f2);
            Console.WriteLine(f3);
            
            return new List<string>();
        }
        
        private NativeModifiers ToNativeModifiers(params IPublicModifier[] modifiers)
        {
            var nativeModifiers = new NativeModifiers();

            foreach (var m in modifiers)
            {
                ((IModifier) m).Apply(nativeModifiers);
            }

            return nativeModifiers;
        }
    }
}