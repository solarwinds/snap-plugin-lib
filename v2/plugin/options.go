package plugin

type Options struct {
	GrpcIp   string
	GrpcPort int

	LogLevel     string
	EnablePprof  bool
	EnableStats  bool
	DebugMode    bool
	PluginConfig string
}
