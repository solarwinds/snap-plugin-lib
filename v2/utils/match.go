package utils

import (
	"github.com/solarwinds/snap-plugin-lib/v2/internal/util/metrictree"
)

func MatchNsToFilter(ns string, filter string) (bool, error) {
	return metrictree.MatchNsToFilter(ns, filter)
}
