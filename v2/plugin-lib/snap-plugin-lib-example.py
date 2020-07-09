import time

from snap_plugin_lib_py import BasePlugin, start_collector, LOGLEVEL_INFO


class ExamplePlugin(BasePlugin):
    def define_plugin(self, ctx):
        example_config = {
            "ip": "127.0.0.1",
            "port": 5434
        }
        ctx.define_example_config(str(example_config))
        ctx.define_tasks_per_instance_limit(4)
        ctx.define_instances_limit(3)

        ctx.define_group("dyn", "Dynamic element from python")
        ctx.define_metric("/python/group1/metric1", "C", True, "1st metric")
        ctx.define_metric("/python/group1/metric2", "C", True, "2nd metric")
        ctx.define_metric("/python/group1/metric3", "C", False, "3rd metric")
        ctx.define_metric("/python/group2/[dyn]/metric4", "C", False, "1st dynamic metric")
        ctx.define_metric("/python/group2/[dyn]/metric5", "C", False, "2nd dynamic metric")

    def collect(self, ctx):
        ctx.log(LOGLEVEL_INFO, "Plugin is being loaded", {
            "name": "py-example"
        })
        print("Requested metrics: ", ctx.requested_metrics())

        ctx.add_metric("/python/group1/metric1", 10,
                       tags={"a": "10", "b": "20"},
                       timestamp=time.time(),
                       description="Custom metric description",
                       unit="Custom metric unit")

        ctx.add_metric("/python/group1/metric2", 20)
        ctx.add_metric("/python/group1/metric3", 40)
        ctx.add_metric("/python/group2/dyn1/metric4", 40)
        ctx.add_metric("/python/group2/dyn15/metric4", 11)
        ctx.add_metric("/python/group2/dyn55/metric4", 5.34)
        ctx.add_metric("/python/group2/dyn56/metric4", -21)
        ctx.add_metric("/python/group2/dyn57/metric4", 9223372036854775999)
        ctx.add_metric("/python/group2/dyn58/metric4", True)

    def load(self, ctx):
        ctx.log(LOGLEVEL_INFO, "Plugin is being loaded", {
            "name": "py-example"
        })

        print("Config keys:", ctx.config_keys())
        print("Config: ", ctx.raw_config())

    def unload(self, ctx):
        ctx.log(LOGLEVEL_INFO, "Plugin is being unloaded", {
            "name": "py-example"
        })


if __name__ == '__main__':
    collector = ExamplePlugin("example", "0.0.1")
    start_collector(collector)
