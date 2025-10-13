package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func main() {
	var command = flag.String("command", "", "命令类型: service, api, model")
	var name = flag.String("name", "", "名称")
	flag.Parse()

	if *command == "" || *name == "" {
		fmt.Println("使用方法:")
		fmt.Println("  go run ./cmd/generator -command=service -name=my-service")
		fmt.Println("  go run ./cmd/generator -command=api -name=my-api")
		fmt.Println("  go run ./cmd/generator -command=model -name=my-model")
		return
	}

	switch *command {
	case "service":
		createService(*name)
	case "api":
		createAPI(*name)
	case "model":
		createModel(*name)
	default:
		fmt.Printf("未知命令: %s\n", *command)
	}
}

func createService(name string) {
	fmt.Printf("创建服务: %s\n", name)

	// 创建服务目录
	serviceDir := filepath.Join("services", name)
	if err := os.MkdirAll(serviceDir, 0755); err != nil {
		fmt.Printf("创建目录失败: %v\n", err)
		return
	}

	// 创建子目录
	subdirs := []string{"cmd/server", "internal/handler", "internal/service", "internal/router", "api"}
	for _, subdir := range subdirs {
		if err := os.MkdirAll(filepath.Join(serviceDir, subdir), 0755); err != nil {
			fmt.Printf("创建子目录失败: %v\n", err)
			return
		}
	}

	// 生成 main.go
	mainTemplate := `package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goweb/pkg/config"
	"goweb/pkg/logger"
	"goweb/services/{{.Name}}/internal/handler"
	"goweb/services/{{.Name}}/internal/router"
	"goweb/services/{{.Name}}/internal/service"
)

func main() {
	// 加载配置（新版）
	cfg := config.New()
	serviceName := "{{.Name}}"
	log := logger.New(serviceName, cfg.GetString("log.level"))

	// 初始化服务
	myService := service.New{{.Name}}Service()
	myHandler := handler.New{{.Name}}Handler(myService)

	// 初始化路由
	r := router.NewRouter(cfg, log, myHandler)

	// 启动HTTP服务
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.GetInt("app.port")),
		Handler:      r,
		ReadTimeout:  cfg.GetDuration("app.read_timeout"),
		WriteTimeout: cfg.GetDuration("app.write_timeout"),
		IdleTimeout:  cfg.GetDuration("app.idle_timeout"),
	}

	go func() {
		log.Info("{{.Name}} service starting", "port", cfg.GetInt("app.port"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("{{.Name}} service start error", err)
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("{{.Name}} service shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("{{.Name}} service shutdown error", err)
	}
}`

	tmpl, _ := template.New("main").Parse(mainTemplate)
	file, _ := os.Create(filepath.Join(serviceDir, "cmd/server/main.go"))
	defer file.Close()
	tmpl.Execute(file, map[string]string{"Name": name})

	fmt.Printf("服务 %s 创建成功！\n", name)
	fmt.Printf("目录: %s\n", serviceDir)
	fmt.Println("请根据需要修改配置和服务类型")
}

func createAPI(name string) {
	fmt.Printf("创建API: %s\n", name)
	// TODO: 实现API生成逻辑
}

func createModel(name string) {
	fmt.Printf("创建模型: %s\n", name)
	// TODO: 实现模型生成逻辑
}
