package connectedapps

import (
	"encoding/xml"

	"github.com/ForceCLI/force-md/internal"
)

const NAME = "ConnectedApp"

func init() {
	internal.TypeRegistry.Register(NAME, func(path string) (internal.RegisterableMetadata, error) { return Open(path) })
}

type ConnectedApp struct {
	internal.MetadataInfo
	XMLName      xml.Name `xml:"ConnectedApp"`
	Xmlns        string   `xml:"xmlns,attr"`
	ContactEmail struct {
		Text string `xml:",chardata"`
	} `xml:"contactEmail"`
	Label struct {
		Text string `xml:",chardata"`
	} `xml:"label"`
	OauthConfig struct {
		CallbackUrl struct {
			Text string `xml:",chardata"`
		} `xml:"callbackUrl"`
		Certificate struct {
			Text string `xml:",chardata"`
		} `xml:"certificate"`
		ConsumerKey struct {
			Text string `xml:",chardata"`
		} `xml:"consumerKey"`
		IsAdminApproved struct {
			Text string `xml:",chardata"`
		} `xml:"isAdminApproved"`
		IsConsumerSecretOptional struct {
			Text string `xml:",chardata"`
		} `xml:"isConsumerSecretOptional"`
		IsIntrospectAllTokens struct {
			Text string `xml:",chardata"`
		} `xml:"isIntrospectAllTokens"`
		IsSecretRequiredForRefreshToken struct {
			Text string `xml:",chardata"`
		} `xml:"isSecretRequiredForRefreshToken"`
		Scopes []struct {
			Text string `xml:",chardata"`
		} `xml:"scopes"`
	} `xml:"oauthConfig"`
	OauthPolicy struct {
		IpRelaxation struct {
			Text string `xml:",chardata"`
		} `xml:"ipRelaxation"`
		RefreshTokenPolicy struct {
			Text string `xml:",chardata"`
		} `xml:"refreshTokenPolicy"`
	} `xml:"oauthPolicy"`
	Description struct {
		Text string `xml:",chardata"`
	} `xml:"description"`
	PermissionSetName struct {
		Text string `xml:",chardata"`
	} `xml:"permissionSetName"`
}

func (c *ConnectedApp) SetMetadata(m internal.MetadataInfo) {
	c.MetadataInfo = m
}

func Open(path string) (*ConnectedApp, error) {
	p := &ConnectedApp{}
	return p, internal.ParseMetadataXml(p, path)
}

func (c *ConnectedApp) Type() internal.MetadataType {
	return NAME
}
