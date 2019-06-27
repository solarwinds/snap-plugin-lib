package plugin

type Meta struct {
	GRPCVersion string

	Plugin struct {
		Name    string
		Version string
	}

	GRPC struct {
		IP   string
		Port int
	}

	PProf struct {
		Enabled bool
		IP      string
		Port    int
	}

	Stats struct {
		Enabled bool
		IP      string
		Port    int
	}
}
