package objects

import (
	"github.com/pkg/errors"
)

func (o *CustomObject) GetFieldSets() []FieldSet {
	return o.FieldSets
}

func (p *CustomObject) DeleteFieldSet(fieldSetName string) error {
	found := false
	newFieldSets := p.FieldSets[:0]
	for _, f := range p.FieldSets {
		if f.FullName == fieldSetName {
			found = true
		} else {
			newFieldSets = append(newFieldSets, f)
		}
	}
	if !found {
		return errors.New("field set not found")
	}
	p.FieldSets = newFieldSets
	return nil
}
