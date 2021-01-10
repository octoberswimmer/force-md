package profile

import (
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
)

func (p *Profile) SetFieldPermissions(fieldName string, updates FieldPermissions) error {
	found := false
	for i, f := range p.FieldPermissions {
		if f.Field.Text == fieldName {
			found = true
			if err := mergo.Merge(&updates, f); err != nil {
				return errors.Wrap(err, "merging permissions")
			}
			p.FieldPermissions[i] = updates
		}
	}
	if !found {
		return errors.New("field not found")
	}
	return nil
}
