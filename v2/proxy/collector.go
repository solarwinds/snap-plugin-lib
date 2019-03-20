package proxy

type Collector interface {
	RequestCollect(id int) ([]Metric, error)
	LoadTask(id int, config string, selectors []string) error
	UnloadTask(id int) error
	RequestInfo() info
}
