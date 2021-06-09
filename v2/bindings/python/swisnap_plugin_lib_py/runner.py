from .c_bridge import start_c_collector, start_c_publisher, start_c_streaming_collector


# Starting collector
def start_collector(collector):
    start_c_collector(collector)


def start_streaming_collector(collector):
    start_c_streaming_collector(collector)


# Starting publisher
def start_publisher(publisher):
    start_c_publisher(publisher)
