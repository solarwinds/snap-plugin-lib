from snap_plugin_lib_py import *


class ExamplePlugin(BasePlugin):
    def define_plugin(self, ctx):
        ctx.define_example_config('{"ip": "127.0.0.1", "port": 5434}')
        ctx.define_tasks_per_instance_limit(4)
        ctx.define_instances_limit(3)

        ctx.define_group("dyn", "Dynamic element from python")
        ctx.define_metric("/python/group1/metric1", "C", True, "1st metric")
        ctx.define_metric("/python/group1/metric2", "C", True, "2nd metric")
        ctx.define_metric("/python/group1/metric3", "C", False, "3rd metric")
        ctx.define_metric("/python/group2/[dyn]/metric4", "C", False, "1st dynamic metric")
        ctx.define_metric("/python/group2/[dyn]/metric5", "C", False, "2nd dynamic metric")

    def collect(self, ctx):
        ctx.add_metric("/python/group1/metric1", 10)
        ctx.add_metric("/python/group1/metric2", 20)
        ctx.add_metric("/python/group1/metric3", 40)
        ctx.add_metric("/python/group2/dyn1/metric4", 40)
        ctx.add_metric("/python/group2/dyn15/metric4", 11)
        ctx.add_metric_with_tags("/python/group2/dyn21/metric5", 35, {"tag1": "v1", "tags56": "v2"})
        ctx.add_metric_ex("/python/group2/dyn55/metric4", 5.34)
        ctx.add_metric_ex("/python/group2/dyn56/metric4", -21)
        ctx.add_metric_ex("/python/group2/dyn57/metric4", 9223372036854775999)
        ctx.add_metric_ex("/python/group2/dyn58/metric4", True)

    def load(self, ctx):
        pass

    def unload(self, ctx):
        pass


if __name__ == '__main__':
    collector = ExamplePlugin("example", "0.0.1")
    start_collector(collector)
