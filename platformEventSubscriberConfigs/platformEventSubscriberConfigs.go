package platformEventSubscriberConfig

import (
	"encoding/xml"

	"github.com/ForceCLI/force-md/internal"
)

const NAME = "PlatformEventSubscriberConfig"

func init() {
	internal.TypeRegistry.Register(NAME, func(path string) (internal.RegisterableMetadata, error) { return Open(path) })
}

type PlatformEventSubscriberConfig struct {
	internal.MetadataInfo
	XMLName   xml.Name `xml:"PlatformEventSubscriberConfig"`
	Xmlns     string   `xml:"xmlns,attr"`
	BatchSize struct {
		Text string `xml:",chardata"`
	} `xml:"batchSize"`
	MasterLabel struct {
		Text string `xml:",chardata"`
	} `xml:"masterLabel"`
	PlatformEventConsumer struct {
		Text string `xml:",chardata"`
	} `xml:"platformEventConsumer"`
	User struct {
		Text string `xml:",chardata"`
	} `xml:"user"`
}

func (c *PlatformEventSubscriberConfig) SetMetadata(m internal.MetadataInfo) {
	c.MetadataInfo = m
}

func (c *PlatformEventSubscriberConfig) Type() internal.MetadataType {
	return NAME
}

func Open(path string) (*PlatformEventSubscriberConfig, error) {
	p := &PlatformEventSubscriberConfig{}
	return p, internal.ParseMetadataXml(p, path)
}
