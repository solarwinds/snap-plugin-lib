package metrictree

import (
	"errors"
	"strings"
)

const (
	_ = iota
	metricDefinitionStrategy
	metricFilteringStrategy
)

type TreeValidator struct {
	strategy       int
	validationList []*Namespace
}

func NewMetricDefinition() *TreeValidator {
	return &TreeValidator{
		strategy: metricDefinitionStrategy,
	}
}

func NewMetricFilter() *TreeValidator {
	return &TreeValidator{
		strategy: metricFilteringStrategy,
	}
}

func (tv *TreeValidator) AddRule(ns string) error {
	parsedNs, err := ParseNamespace(ns)
	if err != nil {
		return err
	}

	return tv.add(parsedNs)
}

func (tv *TreeValidator) IsValid(ns string) bool {
	splitNs := strings.Split(ns, nsSeparator)[1:]

validatorsIt:
	for _, nsValidator := range tv.validationList {
		for i, nsElemValidator := range nsValidator.elements {
			if !nsElemValidator.Match(splitNs[i]) {
				continue validatorsIt
			}
		}
		return true
	}
	return false
}

func (tv *TreeValidator) add(parsedNs *Namespace) error {
	switch tv.strategy {
	case metricDefinitionStrategy:
		if !parsedNs.isUsableForDefinition() {
			return errors.New("can't add rule")
		}
	case metricFilteringStrategy:
		if !parsedNs.isUsableForSelection() {
			return errors.New("can't add rule")
		}
	default:
		panic("invalid strategy")
	}

	tv.validationList = append(tv.validationList, parsedNs)
	return nil
}
