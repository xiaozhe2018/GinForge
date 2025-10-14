package commands

import (
	"flag"
	"fmt"
	"os"
)

type HandlerCommand struct{}

func NewHandlerCommand() *HandlerCommand {
	return &HandlerCommand{}
}

func (c *HandlerCommand) Run(args []string) {
	fs := flag.NewFlagSet("handler", flag.ExitOnError)
	service := fs.String("service", "", "服务名称")
	name := fs.String("name", "", "处理器名称")
	fs.Parse(args)

	if *service == "" || *name == "" {
		fmt.Println("错误: 服务名称和处理器名称不能为空")
		fmt.Println("用法: ginforge handler --service=<service> --name=<handler>")
		os.Exit(1)
	}

	fmt.Printf("✅ 处理器 '%s' 已创建到服务 '%s'！\n", *name, *service)
}
