using System;
using System.Collections.Generic;
using System.Runtime.InteropServices;

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
            var requestedMetrics = new List<string>();
            IntPtr arrPtr = CBridge.ctx_requested_metrics(TaskId);

            for (int offset = 0;; offset += 8)
            {
                var charPtr = Marshal.ReadIntPtr(arrPtr + offset);
                if (charPtr == IntPtr.Zero)
                    break;

                requestedMetrics.Add(Marshal.PtrToStringAnsi(charPtr));
            }

            return requestedMetrics;
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