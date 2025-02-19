package pkg

import (
	"bytes"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/ForceCLI/force-md/general"
	"github.com/ForceCLI/force-md/internal"
	"github.com/ForceCLI/force-md/metadata"
	"github.com/ForceCLI/force-md/metadata/pkg"
)

var (
	metadataType string
	name         string
	version      string
)

var defaultVersion = "51.0"

func init() {
	AddCmd.Flags().StringVarP(&metadataType, "type", "t", "", "metadata type")
	AddCmd.Flags().StringVarP(&name, "name", "n", "", "metadata item name")

	DeleteCmd.Flags().StringVarP(&metadataType, "type", "t", "", "metadata type")
	DeleteCmd.Flags().StringVarP(&name, "name", "n", "", "metadata item name")

	TidyCmd.Flags().BoolP("list", "l", false, "list files that need tidying")

	NewCmd.Flags().StringVarP(&version, "version", "v", defaultVersion, "API version")

	AddCmd.MarkFlagRequired("type")
	AddCmd.MarkFlagRequired("name")
	DeleteCmd.MarkFlagRequired("type")
	DeleteCmd.MarkFlagRequired("name")
}

var AddCmd = &cobra.Command{
	Use:                   "add -t Type -n Name [filename]...",
	Short:                 "Add metadata item to package.xml",
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			add(file, metadataType, name)
		}
	},
}

var DeleteCmd = &cobra.Command{
	Use:                   "delete -t Type -n Name [filename]...",
	Short:                 "Remove metadata item from package.xml",
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			deleteMember(file, metadataType, name)
		}
	},
}

var NewCmd = &cobra.Command{
	Use:                   "new [filename]...",
	Short:                 "Create new package.xml file",
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			createFile(file)
		}
	},
}

var TidyCmd = &cobra.Command{
	Use:                   "tidy [filename]...",
	Short:                 "Tidy package.xml",
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		changes := false
		for _, file := range args {
			list, _ := cmd.Flags().GetBool("list")
			if list {
				needsTidying := checkIfChanged(file)
				changes = needsTidying || changes
			} else {
				tidy(file)
			}
		}
		if changes {
			os.Exit(1)
		}
	},
}

var ListCmd = &cobra.Command{
	Use:                   "list [filename]...",
	Short:                 "list items in package.xml",
	DisableFlagsInUseLine: true,
	Args:                  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			listMembers(file)
		}
	},
}

func add(file string, metadataType string, member string) {
	p, err := pkg.Open(file)
	if err != nil {
		log.Warn("parsing package.xml failed: " + err.Error())
		return
	}
	err = p.Add(metadataType, member)
	if err != nil {
		log.Warn(fmt.Sprintf("update failed for %s: %s", file, err.Error()))
		return
	}
	if err := general.Tidy(p, metadata.MetadataFilePath(file)); err != nil {
		log.Warn("tidying failed: " + err.Error())
	}
	err = internal.WriteToFile(p, file)
	if err != nil {
		log.Warn("update failed: " + err.Error())
		return
	}
}

func deleteMember(file string, metadataType string, member string) {
	p, err := pkg.Open(file)
	if err != nil {
		log.Warn("parsing package.xml failed: " + err.Error())
		return
	}
	err = p.Delete(metadataType, member)
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

func createFile(file string) {
	p := pkg.NewPackage(version)
	err := internal.WriteToFile(p, file)
	if err != nil {
		log.Warn("create failed: " + err.Error())
		return
	}
}

func listMembers(file string) {
	p, err := pkg.Open(file)
	if err != nil {
		log.Warn("parsing package.xml failed: " + err.Error())
		return
	}
	for _, t := range p.Types {
		for _, m := range t.Members {
			fmt.Printf("%s: %s\n", t.Name, m)
		}
	}
}

func checkIfChanged(file string) (changed bool) {
	o := &pkg.Package{}
	contents, err := metadata.ParseMetadataXmlIfPossible(o, file)
	if err != nil {
		log.Warn("parse failure:" + err.Error())
		return
	}
	o.Tidy()
	newContents, err := internal.Marshal(o)
	if err != nil {
		log.Warn("serializing failed: " + err.Error())
		return
	}
	if !bytes.Equal(contents, newContents) {
		fmt.Println(file)
		return true
	}
	return false
}

func tidy(file string) {
	p, err := pkg.Open(file)
	if err != nil {
		log.Warn("parsing package.xml failed: " + err.Error())
		return
	}
	if err := general.Tidy(p, metadata.MetadataFilePath(file)); err != nil {
		log.Warn("tidying failed: " + err.Error())
	}
}
