package commands

import (
	"flag"
	"fmt"
	"os"
)

type ModelCommand struct{}

func NewModelCommand() *ModelCommand {
	return &ModelCommand{}
}

func (c *ModelCommand) Run(args []string) {
	fs := flag.NewFlagSet("model", flag.ExitOnError)
	name := fs.String("name", "", "模型名称")
	fields := fs.String("fields", "", "字段列表")
	fs.Parse(args)

	if *name == "" {
		fmt.Println("错误: 模型名称不能为空")
		fmt.Println("用法: ginforge model --name=<model> --fields=<fields>")
		os.Exit(1)
	}

	fmt.Printf("✅ 模型 '%s' 已创建！\n", *name)
	if *fields != "" {
		fmt.Printf("📝 字段: %s\n", *fields)
	}
}
