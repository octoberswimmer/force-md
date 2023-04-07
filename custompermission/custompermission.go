package custompermission

import (
	"encoding/xml"

	"github.com/ForceCLI/force-md/internal"
)

type CustomPermission struct {
	XMLName     xml.Name `xml:"CustomPermission"`
	Xmlns       string   `xml:"xmlns,attr"`
	Label       string   `xml:"label"`
	Description *string  `xml:"description"`
}

func (p *CustomPermission) MetaCheck() {}

func Open(path string) (*CustomPermission, error) {
	p := &CustomPermission{}
	return p, internal.ParseMetadataXml(p, path)
}

func New(label string) CustomPermission {
	p := CustomPermission{
		Label: label,
		Xmlns: "http://soap.sforce.com/2006/04/metadata",
	}
	return p
}

func (p *CustomPermission) EditDescription(description string) {
	p.Description = &description
}
