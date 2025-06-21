package main

import (
	"log"

	"github.com/biny-go/cLoong/cmd/cloong/internal/project"
	"github.com/biny-go/cLoong/cmd/cloong/internal/proto"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "cloong",
	Short:   "cloong CLI脚手架工具",
	Long:    `cloong CLI脚手架工具，协助快速创建项目。创建新proto：cloong proto add api/test/hello.proto; 创建新项目：cloong new hello; 运行项目：cloong run hello`,
	Version: "v1.0.0",
}

func init() {
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(proto.CmdProto)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
