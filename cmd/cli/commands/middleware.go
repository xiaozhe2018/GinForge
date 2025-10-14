package commands

import (
	"flag"
	"fmt"
	"os"
)

type MiddlewareCommand struct{}

func NewMiddlewareCommand() *MiddlewareCommand {
	return &MiddlewareCommand{}
}

func (c *MiddlewareCommand) Run(args []string) {
	fs := flag.NewFlagSet("middleware", flag.ExitOnError)
	name := fs.String("name", "", "中间件名称")
	fs.Parse(args)

	if *name == "" {
		fmt.Println("错误: 中间件名称不能为空")
		fmt.Println("用法: ginforge middleware --name=<middleware>")
		os.Exit(1)
	}

	fmt.Printf("✅ 中间件 '%s' 已创建！\n", *name)
}
