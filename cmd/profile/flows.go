package profile

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/ForceCLI/force-md/internal"
	"github.com/ForceCLI/force-md/metadata/profile"
)

var flowName string

func init() {
	deleteFlowCmd.Flags().StringVarP(&flowName, "flow", "f", "", "flow name")
	deleteFlowCmd.MarkFlagRequired("flow")

	FlowCmd.AddCommand(deleteFlowCmd)
	FlowCmd.AddCommand(listFlowCmd)
}

var FlowCmd = &cobra.Command{
	Use:   "flow",
	Short: "Manage flow visibility",
}

var deleteFlowCmd = &cobra.Command{
	Use:   "delete -f FlowName [flags] [filename]...",
	Short: "Delete flow access",
	Long:  "Delete flow access in profiles",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			deleteFlowVisibility(file, flowName)
		}
	},
}

var listFlowCmd = &cobra.Command{
	Use:                   "list [filename]...",
	Short:                 "List flows",
	Long:                  "List flows in profiles",
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			listFlowVisibility(file)
		}
	},
}

func deleteFlowVisibility(file string, flowName string) {
	p, err := profile.Open(file)
	if err != nil {
		log.Warn("parsing profile failed: " + err.Error())
		return
	}
	err = p.DeleteFlowAccess(flowName)
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

func listFlowVisibility(file string) {
	p, err := profile.Open(file)
	if err != nil {
		log.Warn("parsing profile failed: " + err.Error())
		return
	}
	tabs := p.GetFlows()
	for _, t := range tabs {
		fmt.Printf("%s: %v\n", t.Flow, t.Enabled.ToBool())
	}
}
