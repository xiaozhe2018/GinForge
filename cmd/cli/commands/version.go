package commands

import (
	"fmt"
	"runtime"
)

type VersionCommand struct{}

func NewVersionCommand() *VersionCommand {
	return &VersionCommand{}
}

func (c *VersionCommand) Run(args []string) {
	fmt.Println("GinForge CLI v1.0.0")
	fmt.Printf("Go版本: %s\n", runtime.Version())
	fmt.Printf("操作系统: %s\n", runtime.GOOS)
	fmt.Printf("架构: %s\n", runtime.GOARCH)
}
