from .c_bridge import start_c_collector


def start_collector(collector):
    start_c_collector(collector)
