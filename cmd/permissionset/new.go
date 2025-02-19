package permissionset

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	. "github.com/ForceCLI/force-md/general"
	"github.com/ForceCLI/force-md/internal"
	"github.com/ForceCLI/force-md/metadata/permissionset"
)

var (
	label       string
	description string
	license     string
)

func init() {
	NewCmd.Flags().StringVarP(&label, "label", "l", "", "label")
	NewCmd.Flags().StringVarP(&description, "description", "d", "", "description")
	NewCmd.MarkFlagRequired("label")

	EditCmd.Flags().StringVarP(&license, "license", "i", "", "license")
}

var NewCmd = &cobra.Command{
	Use:   "new [flags] [filename]...",
	Short: "Create new permission set",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			addPermissionSet(file)
		}
	},
}

var EditCmd = &cobra.Command{
	Use:   "edit [flags] [filename]...",
	Short: "Edit permission set",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			updatePermissionSet(file)
		}
	},
}

func addPermissionSet(file string) {
	p := &permissionset.PermissionSet{
		Xmlns:                 "http://soap.sforce.com/2006/04/metadata",
		HasActivationRequired: FalseText,
		Label:                 label,
	}
	if description != "" {
		p.Description = &permissionset.Description{
			Text: description,
		}
	}
	err := internal.WriteToFile(p, file)
	if err != nil {
		log.Warn("write failed: " + err.Error())
		return
	}
}

func updatePermissionSet(file string) {
	p, err := permissionset.Open(file)
	if err != nil {
		log.Warn("parsing permission set failed: " + err.Error())
		return
	}
	if license != "" {
		p.License = &permissionset.License{
			Text: license,
		}
	}
	err = internal.WriteToFile(p, file)
	if err != nil {
		log.Warn("update failed: " + err.Error())
		return
	}
}
