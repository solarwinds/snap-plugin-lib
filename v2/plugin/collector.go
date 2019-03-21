package plugin

type Collector interface {
	Collect(ctx Context) error
}
