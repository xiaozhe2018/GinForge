package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// AutoRegisterOptions 自动注册选项
type AutoRegisterOptions struct {
	RegisterBackend  bool // 注册后端路由
	RegisterFrontend bool // 注册前端路由
	RegisterMenu     bool // 注册菜单
	DryRun           bool // 预览模式
	Verbose          bool // 详细输出
}

// AutoRegister 自动注册路由和菜单
func (g *Generator) AutoRegister(config *CRUDConfig, opts *AutoRegisterOptions) error {
	results := []string{}

	// 1. 注册后端路由
	if opts.RegisterBackend {
		if err := g.registerBackendRouter(config, opts); err != nil {
			results = append(results, fmt.Sprintf("❌ 后端路由注册失败: %v", err))
		} else {
			results = append(results, "✅ 后端路由注册成功")
		}
	}

	// 2. 注册前端路由
	if opts.RegisterFrontend {
		if err := g.registerFrontendRouter(config, opts); err != nil {
			results = append(results, fmt.Sprintf("❌ 前端路由注册失败: %v", err))
		} else {
			results = append(results, "✅ 前端路由注册成功")
		}
	}

	// 3. 注册菜单
	if opts.RegisterMenu {
		if err := g.registerMenu(config, opts); err != nil {
			results = append(results, fmt.Sprintf("❌ 菜单注册失败: %v", err))
		} else {
			results = append(results, "✅ 菜单注册成功")
		}
	}

	// 输出结果
	if opts.Verbose {
		for _, result := range results {
			fmt.Println(result)
		}
	}

	return nil
}

// registerBackendRouter 注册后端路由
func (g *Generator) registerBackendRouter(config *CRUDConfig, opts *AutoRegisterOptions) error {
	routerFile := filepath.Join("services", config.Module+"-api", "internal", "router", "router.go")

	// 读取文件
	content, err := os.ReadFile(routerFile)
	if err != nil {
		return fmt.Errorf("读取路由文件失败: %w", err)
	}

	fileContent := string(content)

	// 检查是否已经注册
	if strings.Contains(fileContent, config.ModelName+"Handler") {
		return fmt.Errorf("路由已经注册，跳过")
	}

	// 1. 添加 import（如果需要）
	importCode := fmt.Sprintf(`	"goweb/services/%s-api/internal/handler"
	"goweb/services/%s-api/internal/repository"
	"goweb/services/%s-api/internal/service"`, config.Module, config.Module, config.Module)

	if !strings.Contains(fileContent, `"goweb/services/`+config.Module+`-api/internal/handler"`) {
		// 找到 import 块并添加
		importPattern := regexp.MustCompile(`import\s*\(([\s\S]*?)\)`)
		fileContent = importPattern.ReplaceAllStringFunc(fileContent, func(match string) string {
			return strings.Replace(match, ")", importCode+"\n)", 1)
		})
	}

	// 2. 生成初始化代码
	initCode := fmt.Sprintf(`
	// 初始化 %s
	%sRepo := repository.New%sRepository(database)
	%sService := service.New%sService(%sRepo, log)
	%sHandler := handler.New%sHandler(%sService, log)`,
		config.ModelName,
		config.ModelNameCamel, config.ModelName,
		config.ModelNameCamel, config.ModelName, config.ModelNameCamel,
		config.ModelNameCamel, config.ModelName, config.ModelNameCamel,
	)

	// 3. 生成路由注册代码
	routeCode := fmt.Sprintf(`
		// %s 路由
		auth.GET("/%s", %sHandler.List)
		auth.GET("/%s/:id", %sHandler.Get)
		auth.POST("/%s", %sHandler.Create)
		auth.PUT("/%s/:id", %sHandler.Update)
		auth.DELETE("/%s/:id", %sHandler.Delete)`,
		config.Frontend.Title,
		config.ResourceName, config.ModelNameCamel,
		config.ResourceName, config.ModelNameCamel,
		config.ResourceName, config.ModelNameCamel,
		config.ResourceName, config.ModelNameCamel,
		config.ResourceName, config.ModelNameCamel,
	)

	// 找到合适的位置插入初始化代码
	// 在 NewRouter 函数中，database 初始化之后
	handlerPattern := regexp.MustCompile(`(.*Handler\s*:=\s*handler\.New.*Handler.*\n)`)
	matches := handlerPattern.FindAllStringIndex(fileContent, -1)
	if len(matches) > 0 {
		// 在最后一个 Handler 初始化之后插入
		lastMatch := matches[len(matches)-1]
		insertPos := lastMatch[1]
		fileContent = fileContent[:insertPos] + initCode + "\n" + fileContent[insertPos:]
	}

	// 找到合适的位置插入路由注册代码
	// 在 auth 路由组中
	authPattern := regexp.MustCompile(`auth\.DELETE\([^)]+\)[^\n]*\n`)
	authMatches := authPattern.FindAllStringIndex(fileContent, -1)
	if len(authMatches) > 0 {
		// 在最后一个路由之后插入
		lastAuthMatch := authMatches[len(authMatches)-1]
		insertPos := lastAuthMatch[1]
		fileContent = fileContent[:insertPos] + routeCode + "\n" + fileContent[insertPos:]
	}

	// 写入文件
	if !opts.DryRun {
		if err := os.WriteFile(routerFile, []byte(fileContent), 0644); err != nil {
			return fmt.Errorf("写入路由文件失败: %w", err)
		}
	}

	if opts.Verbose {
		fmt.Printf("后端路由注册位置: %s\n", routerFile)
	}

	return nil
}

// registerFrontendRouter 注册前端路由
func (g *Generator) registerFrontendRouter(config *CRUDConfig, opts *AutoRegisterOptions) error {
	routerFile := filepath.Join("web", "admin", "src", "router", "index.ts")

	// 读取文件
	content, err := os.ReadFile(routerFile)
	if err != nil {
		return fmt.Errorf("读取前端路由文件失败: %w", err)
	}

	fileContent := string(content)

	// 检查是否已经注册
	if strings.Contains(fileContent, config.ModelName+"List") {
		return fmt.Errorf("前端路由已经注册，跳过")
	}

	// 生成路由代码
	routeCode := fmt.Sprintf(`      // %s
      {
        path: '%s',
        name: '%sList',
        component: () => import('@/views/%s/index.vue'),
        meta: { title: '%s', requiresAuth: true }
      },`,
		config.Frontend.Title,
		config.ResourceName,
		config.ModelName,
		config.ModelName,
		config.Frontend.Title,
	)

	// 找到 dashboard 的 children 数组，在最后插入
	childrenPattern := regexp.MustCompile(`(path:\s*'dashboard'[\s\S]*?children:\s*\[[\s\S]*?)(\s*\]\s*\})`)
	fileContent = childrenPattern.ReplaceAllString(fileContent, "${1}"+routeCode+"\n${2}")

	// 写入文件
	if !opts.DryRun {
		if err := os.WriteFile(routerFile, []byte(fileContent), 0644); err != nil {
			return fmt.Errorf("写入前端路由文件失败: %w", err)
		}
	}

	if opts.Verbose {
		fmt.Printf("前端路由注册位置: %s\n", routerFile)
	}

	return nil
}

// registerMenu 注册菜单
func (g *Generator) registerMenu(config *CRUDConfig, opts *AutoRegisterOptions) error {
	menuFile := filepath.Join("web", "admin", "src", "layout", "index.vue")

	// 读取文件
	content, err := os.ReadFile(menuFile)
	if err != nil {
		return fmt.Errorf("读取菜单文件失败: %w", err)
	}

	fileContent := string(content)

	// 检查是否已经注册
	if strings.Contains(fileContent, "/dashboard/"+config.ResourceName) {
		return fmt.Errorf("菜单已经注册，跳过")
	}

	// 生成菜单代码
	menuCode := fmt.Sprintf(`
          <!-- %s -->
          <el-menu-item index="/dashboard/%s">
            <el-icon><%s /></el-icon>
            <span>%s</span>
          </el-menu-item>`,
		config.Frontend.Title,
		config.ResourceName,
		config.Frontend.Icon,
		config.Frontend.Title,
	)

	// 找到 el-menu 结束标签之前插入
	menuEndPattern := regexp.MustCompile(`(\s*</el-menu>)`)
	fileContent = menuEndPattern.ReplaceAllString(fileContent, menuCode+"\n${1}")

	// 检查是否需要添加 icon import
	iconImport := config.Frontend.Icon
	if !strings.Contains(fileContent, iconImport) {
		// 找到 import 的 icons 行并添加
		iconPattern := regexp.MustCompile(`(import\s*\{[^}]*)\}\s*from\s*'@element-plus/icons-vue'`)
		fileContent = iconPattern.ReplaceAllString(fileContent, "${1}, "+iconImport+"} from '@element-plus/icons-vue'")
	}

	// 写入文件
	if !opts.DryRun {
		if err := os.WriteFile(menuFile, []byte(fileContent), 0644); err != nil {
			return fmt.Errorf("写入菜单文件失败: %w", err)
		}
	}

	if opts.Verbose {
		fmt.Printf("菜单注册位置: %s\n", menuFile)
	}

	return nil
}
