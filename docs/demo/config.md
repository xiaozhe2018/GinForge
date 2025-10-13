# 配置示例（config）

示例代码：

```go
package demo

import (
    "fmt"
    "goweb/pkg/config"
    "goweb/pkg/logger"
)

func ExampleConfig() {
    cfg := config.New()
    log := logger.New("demo", cfg.GetString("log.level"))

    fmt.Println("app.port:", cfg.GetInt("app.port"))
    fmt.Println("db.driver:", cfg.GetString("database.driver"))
    fmt.Println("jwt.secret:", cfg.GetString("jwt.secret"))

    if cfg.IsProduction() {
        log.Info("running in production")
    }
}
```

YAML 覆盖示例：

```yaml
app:
  port: 8081
jwt:
  secret: dev-secret
```

.env 覆盖示例：

```bash
APP_PORT=8088
JWT_SECRET=override-secret
```
