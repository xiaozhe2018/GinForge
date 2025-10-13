# 数据库示例（db）

```go
package demo

import (
    "goweb/pkg/config"
    "goweb/pkg/logger"
    "goweb/pkg/db"
    "goweb/pkg/model"
    "gorm.io/gorm"
)

func ExampleDB() (*gorm.DB, error) {
    cfg := config.New()
    log := logger.New("demo", cfg.GetString("log.level"))

    manager := db.NewManager(cfg, log)
    gormDB, err := manager.Open()
    if err != nil { return nil, err }

    // 迁移
    if err := gormDB.AutoMigrate(&model.User{}); err != nil { return nil, err }

    // 事务
    err = gormDB.Transaction(func(tx *gorm.DB) error {
        // 示例：插入数据
        // return tx.Create(&model.User{Username:"u1"}).Error
        return nil
    })
    return gormDB, err
}
```
