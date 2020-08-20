﻿using System.Collections.Generic;

namespace SnapPluginLib
{
    internal class CollectContext : Context, ICollectContext
    {
        private enum ValueType
        {
            TypeInt64 = 1,
            TypeUint64,
            TypeFloat,
            TypeDouble,
            TypeBool,
        }
        
        public CollectContext(string taskId) : base(taskId)
        {
        }

        public void AddMetric(string ns, int value, params Modifier[] modifiers)
        {
            var nativeValue = new NativeValue
            {
                v_int64 = value,
                vtype = (int) ValueType.TypeInt64
            };

            AddMetricWithNativeValue(ns, nativeValue, modifiers);
        }

        public void AddMetric(string ns, uint value, params Modifier[] modifiers)
        {
            var nativeValue = new NativeValue
            {
                v_uint64 = value,
                vtype = (int) ValueType.TypeUint64
            };

            AddMetricWithNativeValue(ns, nativeValue, modifiers);
        }

        public void AddMetric(string ns, double value, params Modifier[] modifiers)
        {
            var nativeValue = new NativeValue
            {
                v_double = value,
                vtype = (int) ValueType.TypeDouble
            };

            AddMetricWithNativeValue(ns, nativeValue, modifiers);
        }

        public void AddMetric(string ns, bool value, params Modifier[] modifiers)
        {
            var nativeValue = new NativeValue
            {
                v_bool = value ? 1 : 0,
                vtype = (int) ValueType.TypeBool
            };

            AddMetricWithNativeValue(ns, nativeValue, modifiers);
        }

        public void AlwaysApply(string ns, params Modifier[] modifiers)
        {
            var nativeModifiers = ToNativeModifiers(modifiers);
            var errPtr = CBridge.ctx_always_apply(TaskId, ns, nativeModifiers);
            Memory.FreeNativeModifiers(nativeModifiers);

            Exceptions.ThrowExceptionIfError(errPtr);
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
            return Convertions.NativeStringArrayToList(CBridge.ctx_requested_metrics(TaskId));
        }
        
        private void AddMetricWithNativeValue(string ns, NativeValue nativeValue, params Modifier[] modifiers)
        {
            var nativeModifiers = ToNativeModifiers(modifiers);
            var errPtr = CBridge.ctx_add_metric(TaskId, ns, nativeValue, nativeModifiers);
            Memory.FreeNativeModifiers(nativeModifiers);
            
            Exceptions.ThrowExceptionIfError(errPtr);
        }

        private NativeModifiers ToNativeModifiers(params Modifier[] modifiers)
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