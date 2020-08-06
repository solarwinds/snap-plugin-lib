import time

from snap_plugin_lib_py import BaseCollector, start_collector, LOGLEVEL_INFO
from snap_plugin_lib_py.exceptions import PluginLibException


class ExampleCollectorPlugin(BaseCollector):
    def define_plugin(self, ctx):
        example_config = {"ip": "127.0.0.1", "port": 5434}
        ctx.define_example_config(str(example_config))
        ctx.define_tasks_per_instance_limit(4)
        ctx.define_instances_limit(3)

        ctx.define_group("dyn", "Dynamic element from python")
        ctx.define_metric("/python/group1/metric1", "C", True, "1st metric")
        ctx.define_metric("/python/group1/metric2", "C", True, "2nd metric")
        ctx.define_metric("/python/group1/metric3", "C", False, "3rd metric")
        ctx.define_metric(
            "/python/group2/[dyn]/metric4", "C", False, "1st dynamic metric"
        )
        ctx.define_metric(
            "/python/group2/[dyn]/metric5", "C", False, "2nd dynamic metric"
        )

    def collect(self, ctx):
        ctx.log(LOGLEVEL_INFO, "Collect requested", {"name": self._name})
        print("Requested metrics: ", ctx.requested_metrics())

        mts = (
            (
                "/python/group1/metric1",
                10,
                {"bastide": "80", "succubine": "31"},
                "knifesmith threadmaker",
                "Albanian cheirosophy",
            ),
            (
                "/python/group1/metric2",
                20,
                {"untransferred": "27", "cupolaman": "7"},
                "Chastacosta Entoprocta",
                "octocotyloid lexicon",
            ),
            (
                "/python/group1/metric3",
                40,
                {"gymnotid": "59", "ineffable": "76"},
                "reliantly relick",
                "Cerastium chilenite",
            ),
            (
                "/python/group2/dyn1/metric4",
                40,
                {"heelstrap": "48", "corbiculate": "84"},
                "bastide koechlinite",
                "counteravouch euphon",
            ),
            (
                "/python/group2/dyn15/metric4",
                11,
                {"newelty": "64", "lungmotor": "36"},
                "rootlessness unwise",
                "prepartake Tabebuia",
            ),
            (
                "/python/group2/dyn55/metric4",
                5.34,
                {"wresting": "71", "semidiaphaneity": "3"},
                "overinsolent desilicify",
                "unindulgently planipetalous",
            ),
            (
                "/python/group2/dyn56/metric4",
                -21,
                {"replenishingly": "83", "exululate": "50"},
                "ostraite proliferant",
                "nonhomogenous unreeling",
            ),
            (
                "/python/group2/dyn57/metric4",
                9223372036854775999,
                {"masu": "84", "dallier": "55"},
                "unrelevant corncake",
                "proconquest predeterministic",
            ),
            (
                "/python/group2/dyn58/metric4",
                True,
                {"postcordial": "52", "cupping": "68"},
                "liturgize buhr",
                "presentiveness sulfopurpurate",
            ),
            (
                "/python/group2/dyn3/metric5",
                10,
                {"foyaite": "64", "osoberry": "38"},
                "jowler serape",
                "overknee pinnaclet",
            ),
        )

        for (mt, val, tags, desc, unit) in mts:
            ctx.add_metric(mt, val, tags=tags, description=desc, unit=unit)

        try:
            ctx.add_metric("/^$^%", 20)
        except PluginLibException as e:
            print("Exception: ", e)

    def load(self, ctx):
        ctx.log(LOGLEVEL_INFO, "Plugin is being loaded", {"name": self._name})

        print("Config keys:", ctx.config_keys())
        print("Config: ", ctx.raw_config())
        print("Config a", ctx.config("a"))

    def unload(self, ctx):
        ctx.log(LOGLEVEL_INFO, "Plugin is being unloaded", {"name": self._name})


if __name__ == "__main__":
    collector = ExampleCollectorPlugin("collector-example", "0.0.1")
    start_collector(collector)
