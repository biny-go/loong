package main

import (
	"log"

	"github.com/biny-go/loong/cmd/loong/internal/project"
	"github.com/biny-go/loong/cmd/loong/internal/proto"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "loong",
	Short:   "loong CLI脚手架工具",
	Long:    `loong CLI脚手架工具，协助快速创建项目。创建新proto：loong proto add api/test/hello.proto; 创建新项目：loong new hello; 运行项目：loong run hello`,
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
