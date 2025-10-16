package main

import (
	"fmt"
	"os"
	"path/filepath"

	"goweb/pkg/generator"

	"github.com/spf13/cobra"
)

var (
	// 全局标志
	verbose bool
	dryRun  bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "generator",
		Short:   "GinForge code generator",
		Long:    "Auto-generate CRUD code from database tables for admin-api service",
		Version: "1.0.0",
	}

	// 全局标志
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "显示详细输出")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "预览生成结果，不实际创建文件")

	// 添加子命令
	rootCmd.AddCommand(genCrudCmd())
	rootCmd.AddCommand(genModelCmd())
	rootCmd.AddCommand(initConfigCmd())
	rootCmd.AddCommand(listTablesCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// genCrudCmd 生成完整的 CRUD 代码
func genCrudCmd() *cobra.Command {
	// 本地变量，避免全局变量冲突
	var (
		tableName    string
		configFile   string
		outputDir    string
		withFrontend bool
		force        bool
		autoRegister bool
	)

	cmd := &cobra.Command{
		Use:   "gen:crud",
		Short: "生成完整的 CRUD 代码（后端+前端）",
		Long: `从数据库表或配置文件生成完整的 CRUD 代码

这个命令会生成：
  • Model（数据模型）
  • Repository（数据访问层）
  • Service（业务逻辑层）
  • Handler（HTTP 处理层）
  • 前端 API 定义
  • 前端 Vue 页面（列表 + 表单）
  • 路由配置提示

使用方式：
  # 从数据库表生成（生成到 admin-api 服务）
  generator gen:crud --table=articles

  # 从配置文件生成
  generator gen:crud --config=generator/articles.yaml

  # 只生成后端代码
  generator gen:crud --table=articles --frontend=false

  # 强制覆盖已存在的文件
  generator gen:crud --table=articles --force
  
注意：所有代码将生成到 services/admin-api/ 和 web/admin/
`,
		RunE: runGenCrud,
	}

	cmd.Flags().StringVarP(&tableName, "table", "t", "", "数据库表名（必填，除非使用 --config）")
	cmd.Flags().StringVarP(&configFile, "config", "c", "", "配置文件路径（YAML 格式）")
	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "输出目录（默认为空，自动使用项目根目录）")
	cmd.Flags().BoolVar(&withFrontend, "frontend", true, "生成前端代码")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "强制覆盖已存在的文件")
	cmd.Flags().BoolVarP(&autoRegister, "auto-register", "a", false, "自动注册路由和菜单")

	// 注意：所有代码生成到 admin-api 服务和 admin 前端

	return cmd
}

// genModelCmd 只生成 Model
func genModelCmd() *cobra.Command {
	// 本地变量
	var (
		tableName string
		outputDir string
		force     bool
	)

	cmd := &cobra.Command{
		Use:   "gen:model",
		Short: "只生成 Model 数据模型",
		Long: `从数据库表生成 Model 数据模型

示例：
  generator gen:model --table=articles
  
注意：模型将生成到 services/admin-api/internal/model/
`,
		RunE: runGenModel,
	}

	cmd.Flags().StringVarP(&tableName, "table", "t", "", "数据库表名（必填）")
	cmd.Flags().StringVarP(&outputDir, "output", "o", ".", "输出目录")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "强制覆盖已存在的文件")

	cmd.MarkFlagRequired("table")

	return cmd
}

// initConfigCmd 初始化配置文件
func initConfigCmd() *cobra.Command {
	// 本地变量
	var (
		tableName string
		outputDir string
	)

	cmd := &cobra.Command{
		Use:   "init:config",
		Short: "创建生成器配置文件模板",
		Long: `创建一个配置文件模板，用于自定义代码生成

示例：
  generator init:config --table=articles
  # 将创建 generator/articles.yaml 配置文件
`,
		RunE: runInitConfig,
	}

	cmd.Flags().StringVarP(&tableName, "table", "t", "", "数据库表名（必填）")
	cmd.Flags().StringVarP(&outputDir, "output", "o", "generator", "配置文件输出目录")

	cmd.MarkFlagRequired("table")

	return cmd
}

// listTablesCmd 列出所有数据库表
func listTablesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list:tables",
		Short: "列出数据库中的所有表",
		Long: `连接数据库并列出所有表名，方便选择要生成的表

示例：
  generator list:tables
`,
		RunE: runListTables,
	}

	return cmd
}

// runGenCrud 执行 CRUD 生成
func runGenCrud(cmd *cobra.Command, args []string) error {
	// 从 flags 获取参数值
	tableName, _ := cmd.Flags().GetString("table")
	configFile, _ := cmd.Flags().GetString("config")
	outputDir, _ := cmd.Flags().GetString("output")
	withFrontend, _ := cmd.Flags().GetBool("frontend")
	force, _ := cmd.Flags().GetBool("force")
	autoRegister, _ := cmd.Flags().GetBool("auto-register")

	fmt.Println("🚀 GinForge CRUD 代码生成器")
	fmt.Println("================================")
	fmt.Println()

	// 创建生成器实例
	gen, err := generator.New()
	if err != nil {
		return fmt.Errorf("初始化生成器失败: %w", err)
	}

	var config *generator.CRUDConfig

	// 从配置文件或命令行参数读取配置
	if configFile != "" {
		fmt.Printf("📖 读取配置文件: %s\n", configFile)
		config, err = generator.LoadConfigFromFile(configFile)
		if err != nil {
			return fmt.Errorf("读取配置文件失败: %w", err)
		}
	} else {
		if tableName == "" {
			return fmt.Errorf("必须指定 --table 或 --config 参数")
		}

		fmt.Printf("📊 读取数据库表: %s\n", tableName)
		// 固定使用 admin 模块（后台管理）
		config, err = gen.GenerateConfigFromTable(tableName, "admin")
		if err != nil {
			return fmt.Errorf("读取表结构失败: %w", err)
		}
	}

	// 设置输出选项
	// 如果outputDir为空，自动检测项目根目录
	if outputDir == "" || outputDir == "." {
		// 查找项目根目录（包含go.mod的目录）
		if wd, err := os.Getwd(); err == nil {
			// 向上查找go.mod
			for dir := wd; dir != "/"; dir = filepath.Dir(dir) {
				if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
					outputDir = dir
					break
				}
			}
		}
		if outputDir == "" || outputDir == "." {
			outputDir, _ = os.Getwd() // 降级为当前目录
		}
	}
	opts := &generator.GenerateOptions{
		OutputDir:    outputDir,
		WithFrontend: withFrontend,
		Force:        force,
		DryRun:       dryRun,
		Verbose:      verbose,
	}

	fmt.Println()
	fmt.Println("📝 生成配置:")
	fmt.Printf("  • 表名: %s\n", config.Table)
	fmt.Printf("  • 目标服务: admin-api (后台管理)\n")
	fmt.Printf("  • 模型名: %s\n", config.ModelName)
	fmt.Printf("  • 字段数: %d\n", len(config.Fields))
	fmt.Printf("  • 生成前端: %v\n", opts.WithFrontend)
	if verbose {
		fmt.Printf("  • 输出目录: %s\n", opts.OutputDir)
		if opts.OutputDir == "" {
			fmt.Printf("  • 实际路径: 项目根目录（services/admin-api/ 和 web/admin/）\n")
		}
	}
	fmt.Println()

	if dryRun {
		fmt.Println("🔍 预览模式（不会实际创建文件）")
		fmt.Println()
	}

	// 生成代码
	result, err := gen.GenerateCRUD(config, opts)
	if err != nil {
		return fmt.Errorf("生成代码失败: %w", err)
	}

	// 显示生成结果
	fmt.Println()
	fmt.Println("✅ 代码生成完成！")
	fmt.Println()
	fmt.Println("📁 生成的文件:")
	for _, file := range result.Files {
		if file.Created {
			fmt.Printf("  ✅ %s\n", file.Path)
		} else if file.Skipped {
			fmt.Printf("  ⏭️  %s (已存在，跳过)\n", file.Path)
		}
	}

	if len(result.Errors) > 0 {
		fmt.Println()
		fmt.Println("⚠️  警告/错误:")
		for _, err := range result.Errors {
			fmt.Printf("  ❌ %s\n", err)
		}
	}

	// 自动注册路由和菜单
	if autoRegister {
		fmt.Println()
		fmt.Println("🔧 自动注册路由和菜单...")

		autoOpts := &generator.AutoRegisterOptions{
			RegisterBackend:  true,
			RegisterFrontend: opts.WithFrontend,
			RegisterMenu:     opts.WithFrontend,
			DryRun:           opts.DryRun,
			Verbose:          opts.Verbose,
		}

		if err := gen.AutoRegister(config, autoOpts); err != nil {
			fmt.Printf("⚠️  自动注册部分失败: %v\n", err)
			fmt.Println("💡 提示: 您可以手动完成剩余步骤")
		} else {
			fmt.Println("✅ 路由和菜单注册完成！")
		}
	}

	// 显示后续步骤提示
	fmt.Println()
	if autoRegister && !dryRun {
		fmt.Println("📌 后续步骤:")
		fmt.Println("  ✅ 后端路由已自动注册")
		if opts.WithFrontend {
			fmt.Println("  ✅ 前端路由已自动注册")
			fmt.Println("  ✅ 菜单已自动注册")
		}
		fmt.Println()
		fmt.Println("  🚀 现在只需重启服务即可使用！")
		fmt.Println("     后端: cd services/admin-api && go run cmd/server/main.go")
		fmt.Println("     前端: 刷新浏览器 (http://localhost:3000)")
	} else {
		fmt.Println("📌 后续步骤:")
		fmt.Println("  1. 在路由文件中注册新的 Handler")
		fmt.Println("     在 services/admin-api/internal/router/router.go 中添加:")
		fmt.Printf("     %sHandler := handler.New%sHandler(%sService, log)\n", config.ModelNameCamel, config.ModelName, config.ModelNameCamel)
		fmt.Printf("     auth.GET(\"/%s\", %sHandler.List)\n", config.ResourceName, config.ModelNameCamel)
		fmt.Printf("     auth.POST(\"/%s\", %sHandler.Create)\n", config.ResourceName, config.ModelNameCamel)
		fmt.Printf("     auth.PUT(\"/%s/:id\", %sHandler.Update)\n", config.ResourceName, config.ModelNameCamel)
		fmt.Printf("     auth.DELETE(\"/%s/:id\", %sHandler.Delete)\n", config.ResourceName, config.ModelNameCamel)
		fmt.Println()

		if opts.WithFrontend {
			fmt.Println("  2. 在前端路由中添加新页面")
			fmt.Println("     在 web/admin/src/router/index.ts 中添加:")
			fmt.Printf("     { path: '/%s', name: '%sList', component: () => import('@/views/%s/List.vue') }\n",
				config.ResourceName, config.ModelName, config.ModelName)
			fmt.Println()
			fmt.Println("  3. 在菜单中添加入口（通过管理后台菜单管理功能）")
			fmt.Println()
		}

		fmt.Println("  4. 重启 admin-api 服务并测试功能")
		fmt.Println("     cd services/admin-api && go run cmd/server/main.go")
		fmt.Println()
		fmt.Println("💡 提示: 使用 --auto-register 或 -a 选项可以自动完成上述步骤")
	}
	fmt.Println()
	fmt.Println("🎉 完成！祝您开发愉快！")

	return nil
}

// runGenModel 执行 Model 生成
func runGenModel(cmd *cobra.Command, args []string) error {
	// 从 flags 获取参数
	tableName, _ := cmd.Flags().GetString("table")
	outputDir, _ := cmd.Flags().GetString("output")
	force, _ := cmd.Flags().GetBool("force")

	fmt.Println("🚀 GinForge Model 生成器")
	fmt.Println("================================")
	fmt.Println()

	gen, err := generator.New()
	if err != nil {
		return fmt.Errorf("初始化生成器失败: %w", err)
	}

	// 固定使用 admin 模块
	config, err := gen.GenerateConfigFromTable(tableName, "admin")
	if err != nil {
		return fmt.Errorf("读取表结构失败: %w", err)
	}

	opts := &generator.GenerateOptions{
		OutputDir: outputDir,
		Force:     force,
		DryRun:    dryRun,
		Verbose:   verbose,
	}

	result, err := gen.GenerateModel(config, opts)
	if err != nil {
		return fmt.Errorf("生成 Model 失败: %w", err)
	}

	fmt.Println("✅ Model 生成完成！")
	fmt.Println("📁 生成位置: services/admin-api/internal/model/")
	if len(result.Files) > 0 {
		fmt.Printf("📄 文件名: %s\n", result.Files[0].Path)
	}

	return nil
}

// runInitConfig 执行配置文件初始化
func runInitConfig(cmd *cobra.Command, args []string) error {
	// 从 flags 获取参数
	tableName, _ := cmd.Flags().GetString("table")
	outputDir, _ := cmd.Flags().GetString("output")

	fmt.Println("🚀 GinForge 配置文件生成器")
	fmt.Println("================================")
	fmt.Println()

	gen, err := generator.New()
	if err != nil {
		return fmt.Errorf("初始化生成器失败: %w", err)
	}

	// 固定使用 admin 模块
	config, err := gen.GenerateConfigFromTable(tableName, "admin")
	if err != nil {
		return fmt.Errorf("读取表结构失败: %w", err)
	}

	configPath, err := gen.SaveConfigToFile(config, outputDir)
	if err != nil {
		return fmt.Errorf("保存配置文件失败: %w", err)
	}

	fmt.Println("✅ 配置文件已创建！")
	fmt.Printf("📁 文件: %s\n", configPath)
	fmt.Println()
	fmt.Println("💡 提示:")
	fmt.Println("  • 您可以编辑这个配置文件来自定义生成规则")
	fmt.Println("  • 然后使用 generator gen:crud --config=" + configPath + " 来生成代码")

	return nil
}

// runListTables 执行列出表
func runListTables(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 GinForge 数据库表列表")
	fmt.Println("================================")
	fmt.Println()

	gen, err := generator.New()
	if err != nil {
		return fmt.Errorf("初始化生成器失败: %w", err)
	}

	tables, err := gen.ListTables()
	if err != nil {
		return fmt.Errorf("获取表列表失败: %w", err)
	}

	fmt.Printf("找到 %d 个表:\n\n", len(tables))
	for i, table := range tables {
		fmt.Printf("  %d. %s\n", i+1, table)
	}

	fmt.Println()
	fmt.Println("💡 使用示例:")
	fmt.Println("  generator gen:crud --table=<表名>")
	fmt.Println("  # 所有代码将生成到 admin-api 服务和 admin 前端")

	return nil
}
