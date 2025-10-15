# 工具函数

GinForge 提供了丰富的工具函数，简化常见操作。

## 🔐 加密工具 (`pkg/utils/crypto.go`)

### 哈希函数

```go
import "goweb/pkg/utils"

// MD5 哈希
hash := utils.MD5("hello world")
// 输出：5eb63bbbe01eeed093cb22bb8f5acdc3

// SHA1 哈希
hash := utils.SHA1("hello world")

// SHA256 哈希
hash := utils.SHA256("hello world")

// SHA512 哈希
hash := utils.SHA512("hello world")
```

### 密码加密

```go
import "golang.org/x/crypto/bcrypt"

// 加密密码
hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

// 验证密码
err := bcrypt.CompareHashAndPassword(hashedPassword, []byte("password123"))
if err == nil {
    fmt.Println("密码正确")
}
```

## 📝 字符串工具 (`pkg/utils/string.go`)

### 常用方法

```go
import "goweb/pkg/utils"

// 生成随机字符串
randomStr := utils.RandomString(8)
// 输出：a1B2c3D4

// 生成 UUID
uuid := utils.GenerateUUID()
// 输出：550e8400-e29b-41d4-a716-446655440000

// 驼峰转下划线
snakeCase := utils.CamelToSnake("userName")
// 输出：user_name

// 下划线转驼峰
camelCase := utils.SnakeToCamel("user_name")
// 输出：userName

// 截取字符串
truncated := utils.Truncate("这是一个很长的字符串", 10)
// 输出：这是一个很长的字...

// 判断是否为空
isEmpty := utils.IsEmpty("")  // true
isEmpty := utils.IsEmpty("  ")  // true（空格也算空）
```

## ⏰ 时间工具 (`pkg/utils/time.go`)

### 时间格式化

```go
import "goweb/pkg/utils"

// 格式化时间
formatted := utils.FormatTime(time.Now())
// 输出：2025-10-15 14:30:00

// 格式化日期
formatted := utils.FormatDate(time.Now())
// 输出：2025-10-15

// 时间戳转时间
t := utils.TimestampToTime(1697356800)

// 时间转时间戳
timestamp := utils.TimeToTimestamp(time.Now())
```

### 时间计算

```go
// 获取今天开始时间
startOfDay := utils.BeginOfDay(time.Now())

// 获取今天结束时间
endOfDay := utils.EndOfDay(time.Now())

// 获取本周开始时间
startOfWeek := utils.BeginOfWeek(time.Now())

// 获取本月开始时间
startOfMonth := utils.BeginOfMonth(time.Now())

// 计算时间差
duration := utils.DiffDays(time.Now(), yesterday)
```

## 🎲 随机工具

### 生成随机数

```go
// 生成随机整数 (0-99)
num := rand.Intn(100)

// 生成随机浮点数 (0.0-1.0)
f := rand.Float64()

// 生成指定范围的随机数 (10-20)
num := rand.Intn(11) + 10
```

### 生成随机字符串

```go
// 数字和字母
str := utils.RandomString(10)

// 只有数字
str := utils.RandomNumericString(6)

// 只有字母
str := utils.RandomAlphaString(8)
```

## 🔢 数字工具

### 数字转换

```go
import "strconv"

// 字符串转整数
num, err := strconv.Atoi("123")

// 整数转字符串
str := strconv.Itoa(123)

// 字符串转 int64
num, err := strconv.ParseInt("123", 10, 64)

// float64 转字符串
str := strconv.FormatFloat(123.45, 'f', 2, 64)
```

### 数字格式化

```go
import "fmt"

// 格式化金额（保留2位小数）
amount := fmt.Sprintf("%.2f", 123.456)
// 输出：123.46

// 千分位格式化
formatted := utils.FormatNumber(1234567.89)
// 输出：1,234,567.89
```

## 📋 数组和切片

### 切片操作

```go
// 判断元素是否在切片中
func Contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

// 去重
func Unique(slice []string) []string {
    keys := make(map[string]bool)
    list := []string{}
    
    for _, entry := range slice {
        if _, ok := keys[entry]; !ok {
            keys[entry] = true
            list = append(list, entry)
        }
    }
    return list
}

// 过滤
func Filter(slice []int, predicate func(int) bool) []int {
    result := []int{}
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}
```

## 🗺️ Map 操作

### 常用方法

```go
// 获取 Map 的所有键
func Keys(m map[string]interface{}) []string {
    keys := make([]string, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}

// 合并 Map
func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
    result := make(map[string]interface{})
    for _, m := range maps {
        for k, v := range m {
            result[k] = v
        }
    }
    return result
}
```

## 📁 文件操作

### 文件工具

```go
import "os"

// 检查文件是否存在
func FileExists(path string) bool {
    _, err := os.Stat(path)
    return !os.IsNotExist(err)
}

// 读取文件
content, err := os.ReadFile("config.yaml")

// 写入文件
err := os.WriteFile("output.txt", []byte("content"), 0644)

// 创建目录
err := os.MkdirAll("path/to/dir", 0755)
```

### 路径操作

```go
import "path/filepath"

// 获取文件扩展名
ext := filepath.Ext("document.pdf")  // .pdf

// 获取文件名
name := filepath.Base("/path/to/file.txt")  // file.txt

// 获取目录
dir := filepath.Dir("/path/to/file.txt")  // /path/to

// 拼接路径
path := filepath.Join("uploads", "images", "photo.jpg")
// uploads/images/photo.jpg
```

## 🌐 HTTP 工具

### HTTP 客户端

```go
import "net/http"

// GET 请求
resp, err := http.Get("https://api.example.com/data")
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)

// POST 请求
data := map[string]interface{}{"key": "value"}
jsonData, _ := json.Marshal(data)
resp, err := http.Post("https://api.example.com/data", "application/json", bytes.NewBuffer(jsonData))
```

## 📚 JSON 操作

### JSON 处理

```go
// 结构体转 JSON
user := model.User{Username: "john"}
jsonData, err := json.Marshal(user)

// JSON 转结构体
var user model.User
err := json.Unmarshal(jsonData, &user)

// 美化 JSON 输出
jsonData, err := json.MarshalIndent(user, "", "  ")
```

## 🎯 实用示例

### 生成文件名

```go
func generateFileName(originalName string) string {
    ext := filepath.Ext(originalName)
    baseName := strings.TrimSuffix(originalName, ext)
    timestamp := time.Now().UnixNano()
    randomStr := utils.RandomString(8)
    
    return fmt.Sprintf("%s_%d_%s%s", baseName, timestamp, randomStr, ext)
}
```

### 分页计算

```go
func calculatePagination(page, pageSize int, total int64) map[string]interface{} {
    totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
    
    return map[string]interface{}{
        "page":        page,
        "page_size":   pageSize,
        "total":       total,
        "total_pages": totalPages,
        "has_next":    page < int(totalPages),
        "has_prev":    page > 1,
    }
}
```

## 📚 更多工具

查看完整工具函数：

- **加密工具**: `pkg/utils/crypto.go`
- **字符串工具**: `pkg/utils/string.go`
- **时间工具**: `pkg/utils/time.go`
- **验证工具**: `pkg/validator/validator.go`

## 🎯 下一步

- [配置选项](./config-options) - 完整的配置项说明
- [最佳实践](../best-practices/code-style) - 代码规范

---

**提示**: 善用工具函数可以让代码更简洁、更易维护！

