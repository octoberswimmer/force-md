package permissionset

import (
	"encoding/xml"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/ForceCLI/force-md/internal"
	"github.com/ForceCLI/force-md/metadata/permissionset"
)

var (
	applicationName string
	visible         bool
)

func init() {
	deleteApplicationCmd.Flags().StringVarP(&applicationName, "application", "a", "", "application name")
	deleteApplicationCmd.MarkFlagRequired("application")

	addApplicationCmd.Flags().StringVarP(&applicationName, "application", "a", "", "application name")
	addApplicationCmd.MarkFlagRequired("application")

	listApplicationsCmd.Flags().BoolVarP(&visible, "visible", "v", false, "visible applications only")

	showApplicationCmd.Flags().StringVarP(&applicationName, "application", "a", "", "application name")
	showApplicationCmd.MarkFlagRequired("application")

	editApplicationCmd.Flags().StringVarP(&applicationName, "application", "a", "", "application name")
	editApplicationCmd.Flags().BoolP("default", "d", false, "is default application")
	editApplicationCmd.Flags().BoolP("visible", "v", false, "is visible")
	editApplicationCmd.Flags().BoolP("no-default", "D", false, "is not default")
	editApplicationCmd.Flags().BoolP("no-visible", "V", false, "is not visible")
	editApplicationCmd.MarkFlagRequired("application")

	ApplicationCmd.AddCommand(addApplicationCmd)
	ApplicationCmd.AddCommand(deleteApplicationCmd)
	ApplicationCmd.AddCommand(editApplicationCmd)
	ApplicationCmd.AddCommand(listApplicationsCmd)
	ApplicationCmd.AddCommand(showApplicationCmd)
}

var ApplicationCmd = &cobra.Command{
	Use:   "application",
	Short: "Manage application visibility",
}

var deleteApplicationCmd = &cobra.Command{
	Use:   "delete -a ApplicationName [flags] [filename]...",
	Short: "Delete application visibility",
	Long:  "Delete application visibility in profiles",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			deleteApplicationVisibility(file, applicationName)
		}
	},
}

var showApplicationCmd = &cobra.Command{
	Use:   "show -a ApplicationName [filename]...",
	Short: "Show application visibility",
	Long:  "Show application visibility in profiles",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			showApplicationVisibility(file, applicationName)
		}
	},
}

var editApplicationCmd = &cobra.Command{
	Use:   "edit -a ApplicationName [flags] [filename]...",
	Short: "Update application visibility",
	Long:  "Update application visibility in profiles",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		perms := applicationVisibilityToUpdate(cmd)
		for _, file := range args {
			updateApplicationVisibility(file, perms)
		}
	},
}

var addApplicationCmd = &cobra.Command{
	Use:   "add -a ApplicationName [flags] [filename]...",
	Short: "Add application visibility",
	Long:  "Add application visibility in profiles",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			addApplicationVisibility(file, applicationName)
		}
	},
}

var listApplicationsCmd = &cobra.Command{
	Use:   "list [flags] [filename]...",
	Short: "List applications assigned",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			listApplications(file)
		}
	},
}

func deleteApplicationVisibility(file string, applicationName string) {
	p, err := permissionset.Open(file)
	if err != nil {
		log.Warn("parsing permissionset failed: " + err.Error())
		return
	}
	err = p.DeleteApplicationVisibility(applicationName)
	if err != nil {
		log.Warn(fmt.Sprintf("update failed for %s: %s", file, err.Error()))
		return
	}
	err = internal.WriteToFile(p, file)
	if err != nil {
		log.Warn("update failed: " + err.Error())
		return
	}
}

func addApplicationVisibility(file string, applicationName string) {
	p, err := permissionset.Open(file)
	if err != nil {
		log.Warn("parsing permissionset failed: " + err.Error())
		return
	}
	err = p.AddApplicationVisibility(applicationName)
	if err != nil {
		log.Warn(fmt.Sprintf("update failed for %s: %s", file, err.Error()))
		return
	}
	err = internal.WriteToFile(p, file)
	if err != nil {
		log.Warn("update failed: " + err.Error())
		return
	}
}

func listApplications(file string) {
	p, err := permissionset.Open(file)
	if err != nil {
		log.Warn("parsing permissionset failed: " + err.Error())
		return
	}
	var filters []permissionset.ApplicationFilter
	if visible {
		filters = append(filters, func(v permissionset.ApplicationVisibility) bool {
			return v.Visible.Text == "true"
		})
	}
	apps := p.GetApplications(filters...)
	for _, a := range apps {
		defaultString := ""
		visibleString := "not visible"
		if a.Visible.Text == "true" {
			visibleString = "visible"
		}
		fmt.Printf("%s%s: %s\n", a.Application, defaultString, visibleString)
	}
}

func showApplicationVisibility(file string, application string) {
	p, err := permissionset.Open(file)
	if err != nil {
		log.Warn("parsing permissionset set failed: " + err.Error())
		return
	}
	filter := func(v permissionset.ApplicationVisibility) bool {
		return strings.ToLower(v.Application) == strings.ToLower(application)
	}
	apps := p.GetApplications(filter)
	b, err := xml.MarshalIndent(apps[0], "", "    ")
	if err != nil {
		log.Warn("marshal failed: " + err.Error())
		return
	}
	fmt.Println(string(b))
}

func updateApplicationVisibility(file string, perms permissionset.ApplicationVisibility) {
	p, err := permissionset.Open(file)
	if err != nil {
		log.Warn("parsing permissionset failed: " + err.Error())
		return
	}
	err = p.SetApplicationVisibility(applicationName, perms)
	if err != nil {
		log.Warn(fmt.Sprintf("update failed for %s: %s", file, err.Error()))
		return
	}
	err = internal.WriteToFile(p, file)
	if err != nil {
		log.Warn("update failed: " + err.Error())
		return
	}
}

func applicationVisibilityToUpdate(cmd *cobra.Command) permissionset.ApplicationVisibility {
	perms := permissionset.ApplicationVisibility{}
	perms.Visible = textValue(cmd, "visible")
	return perms
}
