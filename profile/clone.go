package profile

import (
	"sort"

	"github.com/pkg/errors"
)

func (p *Profile) CloneFieldPermissions(src, dest string) error {
	found := false
	for _, f := range p.FieldPermissions {
		if f.Field.Text == src {
			found = true
			clone := FieldPermission{}
			clone.Editable.Text = f.Editable.Text
			clone.Readable.Text = f.Readable.Text
			clone.Field.Text = dest
			p.FieldPermissions = append(p.FieldPermissions, clone)
		}
	}
	if !found {
		return errors.New("source field not found")
	}
	sort.Slice(p.FieldPermissions, func(i, j int) bool {
		return p.FieldPermissions[i].Field.Text < p.FieldPermissions[j].Field.Text
	})
	return nil
}
