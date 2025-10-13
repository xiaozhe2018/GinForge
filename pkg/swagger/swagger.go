package swagger

import (
	"os"
	"os/exec"
	"path/filepath"
)

// GenerateDocs 生成 Swagger 文档
func GenerateDocs(servicePath string) error {
	// 切换到服务目录
	originalDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(servicePath); err != nil {
		return err
	}

	// 运行 swag init
	cmd := exec.Command("swag", "init", "-g", "cmd/server/main.go", "-o", "docs")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// GenerateAllDocs 为所有服务生成 Swagger 文档
func GenerateAllDocs(projectRoot string) error {
	services := []string{
		"services/user-api",
		"services/merchant-api",
		"services/admin-api",
		"services/gateway",
	}

	for _, service := range services {
		servicePath := filepath.Join(projectRoot, service)
		if err := GenerateDocs(servicePath); err != nil {
			return err
		}
	}

	return nil
}
