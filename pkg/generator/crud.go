package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GenerateCRUD 生成完整的 CRUD 代码
func (g *Generator) GenerateCRUD(config *CRUDConfig, opts *GenerateOptions) (*GenerateResult, error) {
	result := &GenerateResult{
		Files:  []FileResult{},
		Errors: []string{},
	}

	// 准备模板数据
	data := g.prepareTemplateData(config)

	// 1. 生成 Model
	if fileResult := g.generateModel(data, config, opts); fileResult != nil {
		result.Files = append(result.Files, *fileResult)
	}

	// 2. 生成 Repository
	if fileResult := g.generateRepository(data, config, opts); fileResult != nil {
		result.Files = append(result.Files, *fileResult)
	}

	// 3. 生成 Service
	if fileResult := g.generateService(data, config, opts); fileResult != nil {
		result.Files = append(result.Files, *fileResult)
	}

	// 4. 生成 Handler
	if fileResult := g.generateHandler(data, config, opts); fileResult != nil {
		result.Files = append(result.Files, *fileResult)
	}

	// 5. 生成前端代码
	if opts.WithFrontend {
		// API 定义
		if fileResult := g.generateFrontendAPI(data, config, opts); fileResult != nil {
			result.Files = append(result.Files, *fileResult)
		}

		// 列表页面
		if fileResult := g.generateFrontendListView(data, config, opts); fileResult != nil {
			result.Files = append(result.Files, *fileResult)
		}

		// 表单页面
		if fileResult := g.generateFrontendFormView(data, config, opts); fileResult != nil {
			result.Files = append(result.Files, *fileResult)
		}
	}

	return result, nil
}

// GenerateModel 生成 Model
func (g *Generator) GenerateModel(config *CRUDConfig, opts *GenerateOptions) (*GenerateResult, error) {
	result := &GenerateResult{
		Files:  []FileResult{},
		Errors: []string{},
	}

	data := g.prepareTemplateData(config)

	if fileResult := g.generateModel(data, config, opts); fileResult != nil {
		result.Files = append(result.Files, *fileResult)
	}

	return result, nil
}

// prepareTemplateData 准备模板数据
func (g *Generator) prepareTemplateData(config *CRUDConfig) *TemplateData {
	data := &TemplateData{
		Table:          config.Table,
		Module:         config.Module,
		ModelName:      config.ModelName,
		ModelNameCamel: config.ModelNameCamel,
		ResourceName:   config.ResourceName,
		PackageName:    "model", // 默认包名
		Fields:         config.Fields,
		PrimaryKey:     getPrimaryKeyField(config.Fields),
		HasSoftDelete:  config.Features.SoftDelete,
		HasTimestamps:  config.Features.Timestamps,
		HasPagination:  config.Features.Pagination,
		HasSearch:      config.Features.Search,
		HasSort:        config.Features.Sort,
		Title:          config.Frontend.Title,
		Icon:           config.Frontend.Icon,
		Imports:        generateImports(config),
		GeneratedAt:    time.Now().Format("2006-01-02 15:04:05"),
	}

	return data
}

// generateModel 生成 Model 文件
func (g *Generator) generateModel(data *TemplateData, config *CRUDConfig, opts *GenerateOptions) *FileResult {
	// 生成文件路径
	outputPath := filepath.Join(opts.OutputDir, "services", config.Module+"-api", "internal", "model")
	filename := toSnakeCase(config.ModelName) + ".go"
	fullPath := filepath.Join(outputPath, filename)

	// 检查文件是否存在
	if !opts.Force && !opts.DryRun {
		if _, err := os.Stat(fullPath); err == nil {
			return &FileResult{
				Path:    fullPath,
				Skipped: true,
			}
		}
	}

	// 渲染模板
	content, err := renderModelTemplate(data)
	if err != nil {
		return &FileResult{
			Path:  fullPath,
			Error: err.Error(),
		}
	}

	// 写入文件
	if !opts.DryRun {
		if err := os.MkdirAll(outputPath, 0755); err != nil {
			return &FileResult{
				Path:  fullPath,
				Error: err.Error(),
			}
		}

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return &FileResult{
				Path:  fullPath,
				Error: err.Error(),
			}
		}
	}

	if opts.Verbose {
		fmt.Printf("生成 Model: %s\n", fullPath)
	}

	return &FileResult{
		Path:    fullPath,
		Created: true,
	}
}

// generateRepository 生成 Repository 文件
func (g *Generator) generateRepository(data *TemplateData, config *CRUDConfig, opts *GenerateOptions) *FileResult {
	outputPath := filepath.Join(opts.OutputDir, "services", config.Module+"-api", "internal", "repository")
	filename := toSnakeCase(config.ModelName) + "_repository.go"
	fullPath := filepath.Join(outputPath, filename)

	if !opts.Force && !opts.DryRun {
		if _, err := os.Stat(fullPath); err == nil {
			return &FileResult{Path: fullPath, Skipped: true}
		}
	}

	content, err := renderRepositoryTemplate(data)
	if err != nil {
		return &FileResult{Path: fullPath, Error: err.Error()}
	}

	if !opts.DryRun {
		if err := os.MkdirAll(outputPath, 0755); err != nil {
			return &FileResult{Path: fullPath, Error: err.Error()}
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return &FileResult{Path: fullPath, Error: err.Error()}
		}
	}

	if opts.Verbose {
		fmt.Printf("生成 Repository: %s\n", fullPath)
	}

	return &FileResult{Path: fullPath, Created: true}
}

// generateService 生成 Service 文件
func (g *Generator) generateService(data *TemplateData, config *CRUDConfig, opts *GenerateOptions) *FileResult {
	outputPath := filepath.Join(opts.OutputDir, "services", config.Module+"-api", "internal", "service")
	filename := toSnakeCase(config.ModelName) + "_service.go"
	fullPath := filepath.Join(outputPath, filename)

	if !opts.Force && !opts.DryRun {
		if _, err := os.Stat(fullPath); err == nil {
			return &FileResult{Path: fullPath, Skipped: true}
		}
	}

	content, err := renderServiceTemplate(data)
	if err != nil {
		return &FileResult{Path: fullPath, Error: err.Error()}
	}

	if !opts.DryRun {
		if err := os.MkdirAll(outputPath, 0755); err != nil {
			return &FileResult{Path: fullPath, Error: err.Error()}
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return &FileResult{Path: fullPath, Error: err.Error()}
		}
	}

	if opts.Verbose {
		fmt.Printf("生成 Service: %s\n", fullPath)
	}

	return &FileResult{Path: fullPath, Created: true}
}

// generateHandler 生成 Handler 文件
func (g *Generator) generateHandler(data *TemplateData, config *CRUDConfig, opts *GenerateOptions) *FileResult {
	outputPath := filepath.Join(opts.OutputDir, "services", config.Module+"-api", "internal", "handler")
	filename := toSnakeCase(config.ModelName) + "_handler.go"
	fullPath := filepath.Join(outputPath, filename)

	if !opts.Force && !opts.DryRun {
		if _, err := os.Stat(fullPath); err == nil {
			return &FileResult{Path: fullPath, Skipped: true}
		}
	}

	content, err := renderHandlerTemplate(data)
	if err != nil {
		return &FileResult{Path: fullPath, Error: err.Error()}
	}

	if !opts.DryRun {
		if err := os.MkdirAll(outputPath, 0755); err != nil {
			return &FileResult{Path: fullPath, Error: err.Error()}
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return &FileResult{Path: fullPath, Error: err.Error()}
		}
	}

	if opts.Verbose {
		fmt.Printf("生成 Handler: %s\n", fullPath)
	}

	return &FileResult{Path: fullPath, Created: true}
}

// generateFrontendAPI 生成前端 API 文件
func (g *Generator) generateFrontendAPI(data *TemplateData, config *CRUDConfig, opts *GenerateOptions) *FileResult {
	outputPath := filepath.Join(opts.OutputDir, "web", "admin", "src", "api")
	filename := toSnakeCase(config.ModelName) + ".ts"
	fullPath := filepath.Join(outputPath, filename)

	if !opts.Force && !opts.DryRun {
		if _, err := os.Stat(fullPath); err == nil {
			return &FileResult{Path: fullPath, Skipped: true}
		}
	}

	content, err := renderFrontendAPITemplate(data)
	if err != nil {
		return &FileResult{Path: fullPath, Error: err.Error()}
	}

	if !opts.DryRun {
		if err := os.MkdirAll(outputPath, 0755); err != nil {
			return &FileResult{Path: fullPath, Error: err.Error()}
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return &FileResult{Path: fullPath, Error: err.Error()}
		}
	}

	if opts.Verbose {
		fmt.Printf("生成前端 API: %s\n", fullPath)
	}

	return &FileResult{Path: fullPath, Created: true}
}

// generateFrontendListView 生成前端列表页面
func (g *Generator) generateFrontendListView(data *TemplateData, config *CRUDConfig, opts *GenerateOptions) *FileResult {
	outputPath := filepath.Join(opts.OutputDir, "web", "admin", "src", "views", config.ModelName)
	filename := "index.vue"
	fullPath := filepath.Join(outputPath, filename)

	if !opts.Force && !opts.DryRun {
		if _, err := os.Stat(fullPath); err == nil {
			return &FileResult{Path: fullPath, Skipped: true}
		}
	}

	content, err := renderFrontendListTemplate(data)
	if err != nil {
		return &FileResult{Path: fullPath, Error: err.Error()}
	}

	if !opts.DryRun {
		if err := os.MkdirAll(outputPath, 0755); err != nil {
			return &FileResult{Path: fullPath, Error: err.Error()}
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return &FileResult{Path: fullPath, Error: err.Error()}
		}
	}

	if opts.Verbose {
		fmt.Printf("生成前端列表页面: %s\n", fullPath)
	}

	return &FileResult{Path: fullPath, Created: true}
}

// generateFrontendFormView 生成前端表单页面
func (g *Generator) generateFrontendFormView(data *TemplateData, config *CRUDConfig, opts *GenerateOptions) *FileResult {
	outputPath := filepath.Join(opts.OutputDir, "web", "admin", "src", "views", config.ModelName)
	filename := "Form.vue"
	fullPath := filepath.Join(outputPath, filename)

	if !opts.Force && !opts.DryRun {
		if _, err := os.Stat(fullPath); err == nil {
			return &FileResult{Path: fullPath, Skipped: true}
		}
	}

	content, err := renderFrontendFormTemplate(data)
	if err != nil {
		return &FileResult{Path: fullPath, Error: err.Error()}
	}

	if !opts.DryRun {
		if err := os.MkdirAll(outputPath, 0755); err != nil {
			return &FileResult{Path: fullPath, Error: err.Error()}
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return &FileResult{Path: fullPath, Error: err.Error()}
		}
	}

	if opts.Verbose {
		fmt.Printf("生成前端表单页面: %s\n", fullPath)
	}

	return &FileResult{Path: fullPath, Created: true}
}
