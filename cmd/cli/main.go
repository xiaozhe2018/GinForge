package main

import (
	"fmt"
	"os"

	"goweb/cmd/cli/commands"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "service":
		commands.NewServiceCommand().Run(args)
	case "handler":
		commands.NewHandlerCommand().Run(args)
	case "model":
		commands.NewModelCommand().Run(args)
	case "middleware":
		commands.NewMiddlewareCommand().Run(args)
	case "config":
		commands.NewConfigCommand().Run(args)
	case "test":
		commands.NewTestCommand().Run(args)
	case "deploy":
		commands.NewDeployCommand().Run(args)
	case "init":
		commands.NewInitCommand().Run(args)
	case "version":
		commands.NewVersionCommand().Run(args)
	case "help", "-h", "--help":
		showHelp()
	default:
		fmt.Printf("未知命令: %s\n", command)
		showHelp()
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("GinForge CLI - 微服务开发脚手架")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  ginforge <command> [flags]")
	fmt.Println()
	fmt.Println("可用命令:")
	fmt.Println("  service    创建新的微服务")
	fmt.Println("  handler    创建新的处理器")
	fmt.Println("  model      创建新的数据模型")
	fmt.Println("  middleware 创建新的中间件")
	fmt.Println("  config     管理配置文件")
	fmt.Println("  test       运行测试")
	fmt.Println("  deploy     部署服务")
	fmt.Println("  init       初始化项目")
	fmt.Println("  version    显示版本信息")
	fmt.Println("  help       显示帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  ginforge service --name=payment --port=8086")
	fmt.Println("  ginforge handler --service=user --name=profile")
	fmt.Println("  ginforge model --name=user --fields=name,email,age")
	fmt.Println("  ginforge test --service=user --coverage")
	fmt.Println("  ginforge deploy --env=production")
}
