package profile

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	. "github.com/octoberswimmer/force-md/general"
	"github.com/octoberswimmer/force-md/internal"
	"github.com/octoberswimmer/force-md/profile"
)

var (
	objectName string
)

func init() {
	editObjectCmd.Flags().StringVarP(&objectName, "object", "o", "", "object name")
	editObjectCmd.Flags().BoolP("create", "c", false, "allow create")
	editObjectCmd.Flags().BoolP("delete", "d", false, "allow delete")
	editObjectCmd.Flags().BoolP("edit", "e", false, "allow edit")
	editObjectCmd.Flags().BoolP("read", "r", false, "allow read")
	editObjectCmd.Flags().BoolP("modify-all", "m", false, "allow modify all")
	editObjectCmd.Flags().BoolP("view-all", "v", false, "allow view all")
	editObjectCmd.Flags().BoolP("no-create", "C", false, "disallow create")
	editObjectCmd.Flags().BoolP("no-delete", "D", false, "disallow delete")
	editObjectCmd.Flags().BoolP("no-edit", "E", false, "disallow edit")
	editObjectCmd.Flags().BoolP("no-read", "R", false, "disallow read")
	editObjectCmd.Flags().BoolP("no-modify-all", "M", false, "disallow modify all")
	editObjectCmd.Flags().BoolP("no-view-all", "V", false, "disallow view all")
	editObjectCmd.Flags().SortFlags = false
	editObjectCmd.MarkFlagRequired("object")

	addObjectCmd.Flags().StringVarP(&objectName, "object", "o", "", "object name")
	addObjectCmd.MarkFlagRequired("object")

	deleteObjectCmd.Flags().StringVarP(&objectName, "object", "o", "", "object name")
	deleteObjectCmd.MarkFlagRequired("object")

	showObjectCmd.Flags().StringVarP(&objectName, "object", "o", "", "object name")
	showObjectCmd.MarkFlagRequired("object")

	listObjectCmd.Flags().BoolP("create", "c", false, "has create")
	listObjectCmd.Flags().BoolP("delete", "d", false, "has delete")
	listObjectCmd.Flags().BoolP("edit", "e", false, "has edit")
	listObjectCmd.Flags().BoolP("read", "r", false, "has read")
	listObjectCmd.Flags().BoolP("modify-all", "m", false, "has modify all")
	listObjectCmd.Flags().BoolP("view-all", "v", false, "has view all")
	listObjectCmd.Flags().BoolP("no-create", "C", false, "does not have create")
	listObjectCmd.Flags().BoolP("no-delete", "D", false, "does not have delete")
	listObjectCmd.Flags().BoolP("no-edit", "E", false, "does not have edit")
	listObjectCmd.Flags().BoolP("no-read", "R", false, "does not have read")
	listObjectCmd.Flags().BoolP("no-modify-all", "M", false, "does not have modify all")
	listObjectCmd.Flags().BoolP("no-view-all", "V", false, "does not have view all")

	ObjectPermissionsCmd.AddCommand(editObjectCmd)
	ObjectPermissionsCmd.AddCommand(addObjectCmd)
	ObjectPermissionsCmd.AddCommand(showObjectCmd)
	ObjectPermissionsCmd.AddCommand(deleteObjectCmd)
	ObjectPermissionsCmd.AddCommand(listObjectCmd)
}

var ObjectPermissionsCmd = &cobra.Command{
	Use:   "object-permissions",
	Short: "Update object permissions",
}

var editObjectCmd = &cobra.Command{
	Use:   "edit -o SObject [flags] [filename]...",
	Short: "Update object permissions",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		perms := objectPermissionsFromFlags(cmd)
		for _, file := range args {
			updateObjectPermissions(file, perms)
		}
	},
}

var addObjectCmd = &cobra.Command{
	Use:   "add -o SObject [flags] [filename]...",
	Short: "Add object permissions",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			addObjectPermissions(file)
		}
	},
}

var deleteObjectCmd = &cobra.Command{
	Use:   "delete -o SObject [flags] [filename]...",
	Short: "Delete object permissions",
	Long:  "Delete object permissions and related field permissions in profiles",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			deleteObjectPermissions(file, objectName)
		}
	},
}

var showObjectCmd = &cobra.Command{
	Use:                   "show -f Object [filename]...",
	Short:                 "Show object permissions",
	Args:                  cobra.MinimumNArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			showObjectPermissions(file, objectName)
		}
	},
}

var listObjectCmd = &cobra.Command{
	Use:   "list [flags] [filename]...",
	Short: "List object permissions",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		perms := objectPermissionsFromFlags(cmd)
		for _, file := range args {
			listObjectPermissions(file, perms)
		}
	},
}

func textValue(cmd *cobra.Command, flag string) (t BooleanText) {
	if cmd.Flags().Changed(flag) {
		val, _ := cmd.Flags().GetBool(flag)
		t = BooleanText{
			Text: strconv.FormatBool(val),
		}
	}
	antiFlag := "no-" + flag
	if cmd.Flags().Changed(antiFlag) {
		val, _ := cmd.Flags().GetBool(antiFlag)
		t = BooleanText{
			Text: strconv.FormatBool(!val),
		}
	}
	return t
}

func objectPermissionsFromFlags(cmd *cobra.Command) profile.ObjectPermissions {
	perms := profile.ObjectPermissions{}
	perms.AllowCreate = textValue(cmd, "create")
	perms.AllowDelete = textValue(cmd, "delete")
	perms.AllowEdit = textValue(cmd, "edit")
	perms.AllowRead = textValue(cmd, "read")
	perms.ModifyAllRecords = textValue(cmd, "modify-all")
	perms.ViewAllRecords = textValue(cmd, "view-all")
	return perms
}

func updateObjectPermissions(file string, perms profile.ObjectPermissions) {
	p, err := profile.Open(file)
	if err != nil {
		log.Warn("parsing profile failed: " + err.Error())
		return
	}
	err = p.SetObjectPermissions(objectName, perms)
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

func addObjectPermissions(file string) {
	p, err := profile.Open(file)
	if err != nil {
		log.Warn("parsing profile failed: " + err.Error())
		return
	}
	err = p.AddObjectPermissions(objectName)
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

func deleteObjectPermissions(file string, objectName string) {
	p, err := profile.Open(file)
	if err != nil {
		log.Warn("parsing profile failed: " + err.Error())
		return
	}
	p.DeleteObjectPermissions(objectName)
	err = internal.WriteToFile(p, file)
	if err != nil {
		log.Warn("update failed: " + err.Error())
		return
	}
}

func showObjectPermissions(file string, objectName string) {
	p, err := profile.Open(file)
	if err != nil {
		log.Warn("parsing profile failed: " + err.Error())
		return
	}
	objects := p.GetObjectPermissions(func(o profile.ObjectPermissions) bool {
		return strings.ToLower(o.Object.Text) == strings.ToLower(objectName)
	})
	if len(objects) == 0 {
		log.Warn(fmt.Sprintf("object not found in %s", file))
		return
	}
	b, err := xml.MarshalIndent(objects[0], "", "    ")
	if err != nil {
		log.Warn("marshal failed: " + err.Error())
		return
	}
	fmt.Println(string(b))
}

func listObjectPermissions(file string, filter profile.ObjectPermissions) {
	p, err := profile.Open(file)
	if err != nil {
		log.Warn("parsing profile failed: " + err.Error())
		return
	}
	flagFilter := func(o profile.ObjectPermissions) bool {
		if filter.AllowCreate.Text != "" && filter.AllowCreate.ToBool() != o.AllowCreate.ToBool() {
			return false
		}
		if filter.AllowRead.Text != "" && filter.AllowRead.ToBool() != o.AllowRead.ToBool() {
			return false
		}
		if filter.AllowEdit.Text != "" && filter.AllowEdit.ToBool() != o.AllowEdit.ToBool() {
			return false
		}
		if filter.AllowDelete.Text != "" && filter.AllowDelete.ToBool() != o.AllowDelete.ToBool() {
			return false
		}
		if filter.ViewAllRecords.Text != "" && filter.ViewAllRecords.ToBool() != o.ViewAllRecords.ToBool() {
			return false
		}
		if filter.ModifyAllRecords.Text != "" && filter.ModifyAllRecords.ToBool() != o.ModifyAllRecords.ToBool() {
			return false
		}
		return true
	}
	objects := p.GetObjectPermissions(flagFilter)
	for _, o := range objects {
		perms := ""
		if o.AllowCreate.ToBool() {
			perms += "c"
		}
		if o.ViewAllRecords.ToBool() {
			perms += "R"
		} else if o.AllowRead.ToBool() {
			perms += "r"
		}
		if o.ModifyAllRecords.ToBool() {
			perms += "U"
		} else if o.AllowEdit.ToBool() {
			perms += "u"
		}
		if o.AllowDelete.ToBool() {
			perms += "d"
		}

		fmt.Printf("%s: %s\n", o.Object.Text, perms)
	}
}
