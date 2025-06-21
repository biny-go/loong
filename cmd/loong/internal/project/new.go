package project

import "github.com/spf13/cobra"

var CmdNew = &cobra.Command{
	Use:   "new",
	Short: "创建新项目模板",
	Long:  "创建新项目模板. Example: loong new helloworld",
}
