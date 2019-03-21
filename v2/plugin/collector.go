/*
The package "plugin" provides interfaces to define custom plugins and Context interface
which allows to perform any collection-related operation.
*/
package plugin

type Collector interface {
	Collect(ctx Context) error
}
