# 配置管理 (Configuration)

## 开篇故事

想象你要开一家连锁餐厅。每家店都需要一套配置：

- **默认配置**：菜单、餐具、基本流程（每家店都一样）
- **本地配置**：装修风格、当地特色菜（每家店不同）
- **环境配置**：开业时间、员工数量（根据商圈调整）

如果你把这些信息硬编码（hard-coded）在员工手册里，每次开新店都要重写整本手册——很快就会乱套。

Go 程序的配置管理也是这个道理。初学者常把端口、数据库连接串写死在代码里，项目一多、环境一复杂（开发、测试、生产），维护成本就会爆炸。这一章教你如何设计灵活、可维护的配置系统。

## 本章适合谁

- ✅ 写过"把数据库连接串硬编码在代码里"的程序，现在想改进
- ✅ 需要区分开发/测试/生产环境配置，但不知道如何组织
- ✅ 用过 Viper 等配置库，但想理解底层原理
- ✅ 想学习用反射（reflection）和结构体标签（struct tags）实现配置绑定

如果你曾经为"为什么测试环境连到生产数据库"而恐慌过，本章就是为你准备的。

## 你会学到什么

完成本章后，你将能够：

1. **设计配置结构体**：用 Go 结构体和标签组织配置项，支持 JSON/YAML/环境变量
2. **实现分层加载**：默认值 → 配置文件 → 环境变量，理解优先级顺序
3. **使用反射读取标签**：用 `reflect` 包自动绑定配置值到结构体字段
4. **处理配置错误**：给出清晰的错误信息，包含字段名和期望值
5. **实现环境隔离**：用环境变量覆盖配置，支持不同部署场景

## 前置要求

在开始之前，请确保你已掌握：

- Go 结构体（struct）和标签（tags）语法
- JSON 基本格式（key-value、嵌套对象）
- 环境变量概念（`os.Getenv`）
- 错误处理模式（`error` 返回值、`fmt.Errorf`）

了解反射（reflection）有帮助，但本章会有详细解释。

## 第一个例子

让我们从一个最简单的问题开始：如何从环境变量读取一个端口号？

```go
package main

import (
    "fmt"
    "os"
    "strconv"
)

func main() {
    // 默认值
    port := 8080
    
    // 环境变量覆盖
    if portStr := os.Getenv("APP_PORT"); portStr != "" {
        if p, err := strconv.Atoi(portStr); err == nil {
            port = p
        }
    }
    
    fmt.Printf("服务端口：%d\n", port)
}
```

**运行**：
```bash
$ go run main.go
服务端口：8080

$ APP_PORT=9090 go run main.go
服务端口：9090
```

**核心思想**：
1. 代码中设置**默认值**（安全起点）
2. 环境变量**可选覆盖**（部署时灵活）
3. 解析失败时**保留默认值**（安全降级）

这个简单模式是配置管理的基石。接下来我们会扩展它，支持更多配置项和文件格式。

## 原理解析

### 1. 配置的三个来源

一个完整的配置系统通常有三个层次：

```
┌─────────────────────────────────────┐
│   环境变量 (Environment Variables)   │  ← 最高优先级
│   - 部署时覆盖                         │    (容器、CI/CD)
│   - 格式：HELLO_SERVER_PORT=9090      │
└─────────────────────────────────────┘
              ↓ 覆盖
┌─────────────────────────────────────┐
│   配置文件 (Config File)             │  ← 中等优先级
│   - 项目级设置                         │    (JSON/YAML)
│   - 格式：{"server": {"port": 8080}} │
└─────────────────────────────────────┘
              ↓ 覆盖
┌─────────────────────────────────────┐
│   默认值 (Default Values)            │  ← 最低优先级
│   - 安全兜底                           │    (代码中定义)
│   - 格式：Port: 8080                 │
└─────────────────────────────────────┘
```

**为什么这个顺序重要**？
- **默认值**保证程序"开箱即用"，无需任何配置
- **配置文件**允许项目定制化（如数据库地址）
- **环境变量**允许部署时动态调整（如 Kubernetes ConfigMap）

**代码中的体现**：
```go
func resolveConfig(paths []string, prefix string, lookup func(string) (string, bool)) (appConfig, error) {
    // 1. 从默认值开始
    cfg := defaultConfig()
    
    // 2. 逐个加载配置文件（后加载的覆盖先加载的）
    for _, path := range paths {
        next, err := loadConfigFile(cfg, path)
        if err != nil {
            return appConfig{}, err
        }
        cfg = next
    }
    
    // 3. 最后用环境变量覆盖
    return loadEnvConfig(cfg, prefix, lookup)
}
```

### 2. 结构体标签的作用

结构体标签（struct tags）是配置绑定的"元数据"：

```go
type appConfig struct {
    AppName  string `json:"app_name" config:"app_name" env:"APP_NAME"`
    LogLevel string `json:"log_level" config:"log_level" env:"LOG_LEVEL"`
    Server   serverConfig `json:"server" config:"server"`
}
```

**三种标签**：
- `json:"app_name"`：JSON 解析时用（`encoding/json` 包）
- `config:"app_name"`：YAML/自定义解析时用
- `env:"APP_NAME"`：环境变量绑定时用

**用反射读取标签**：
```go
fieldType := targetType.Field(index)
key := fieldType.Tag.Get("env") // 读取 env 标签
if key == "" {
    continue // 没有标签的字段跳过
}

// 用 key 查找环境变量
rawValue, ok := lookup(key)
if ok {
    setValueFromString(fieldValue, rawValue)
}
```

**好处**：
- 配置映射关系**声明式**，清晰可见
- 不依赖框架，纯 Go 标准库
- 易于测试和扩展

### 3. 反射绑定的核心逻辑

`bindMap()` 和 `bindEnv()` 是配置系统的核心函数。它们的工作流程类似：

```go
func bindMap(target any, values map[string]any) error {
    root := reflect.ValueOf(target)
    
    // 遍历结构体所有字段
    for index := range target.NumField() {
        fieldValue := target.Field(index)
        fieldType := targetType.Field(index)
        
        // 读取标签
        key := fieldType.Tag.Get("config")
        if key == "" {
            continue
        }
        
        // 处理嵌套结构体
        if fieldValue.Kind() == reflect.Struct {
            nestedValues := values[key].(map[string]any)
            bindMapValue(fieldValue, nestedValues)
            continue
        }
        
        // 设置值（类型转换）
        rawValue := values[key]
        setValueFromAny(fieldValue, rawValue)
    }
}
```

**关键点**：
1. **递归处理嵌套**：`serverConfig` 这样的嵌套结构体需要递归绑定
2. **类型转换**：配置文件中的数字是 `float64`，需要转成 `int`
3. **错误处理**：类型不匹配时返回清晰错误

### 4. TypeScript 类型转换的细节

配置文件中的值到 Go 类型需要转换：

```go
func setValueFromAny(target reflect.Value, rawValue any) error {
    switch value := rawValue.(type) {
    case string:
        return setValueFromString(target, value)
    case float64: // JSON 数字默认是 float64
        if target.Kind() != reflect.Int {
            return fmt.Errorf("expected %s, got number", target.Kind())
        }
        target.SetInt(int64(value))
        return nil
    case bool:
        if target.Kind() != reflect.Bool {
            return fmt.Errorf("expected %s, got bool", target.Kind())
        }
        target.SetBool(value)
        return nil
    }
}
```

**常见陷阱**：
- JSON 解析数字 → `float64`，需要 `int64(value)` 转换
- YAML 解析数字 → 可能是 `int` 或 `string`，需要判断
- 类型不匹配时**立即报错**，不要静默失败

### 5. 简易 YAML 解析器

为了演示原理，代码中实现了一个极简 YAML 解析器：

```go
func parseSimpleYAML(content string) (map[string]any, error) {
    result := map[string]any{}
    currentSection := ""
    
    for _, line := range strings.Split(content, "\n") {
        trimmed := strings.TrimSpace(line)
        
        // 跳过空行和注释
        if trimmed == "" || strings.HasPrefix(trimmed, "#") {
            continue
        }
        
        parts := strings.SplitN(trimmed, ":", 2)
        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])
        indent := len(line) - len(strings.TrimLeft(line, " "))
        
        // 顶层 key
        if indent == 0 {
            if value == "" {
                result[key] = map[string]any{} // 嵌套对象
                currentSection = key
            } else {
                result[key] = parseScalar(value) // 标量值
            }
            continue
        }
        
        // 嵌套 key
        sectionValues := result[currentSection].(map[string]any)
        sectionValues[key] = parseScalar(value)
    }
    
    return result, nil
}
```

**支持格式**：
```yaml
app_name: hello-go
log_level: info
server:
  host: 127.0.0.1
  port: 8080
```

**局限性**（生产环境请用 `gopkg.in/yaml.v3`）：
- 不支持数组
- 不支持多行字符串
- 不支持复杂嵌套

## 常见错误

### 错误 1：环境变量 key 大小写错误

```go
// ❌ 错误代码
type serverConfig struct {
    Port int `env:"server_port"` // 小写
}

// 环境变量通常是 HELLO_SERVER_PORT，匹配不上
```

**如何修复**：
```go
// ✅ 修复：用大写，与环境变量一致
type serverConfig struct {
    Port int `env:"SERVER_PORT"`
}

// 然后用前缀拼接
lookup("HELLO_SERVER_PORT") // "HELLO_" + "SERVER_PORT"
```

**最佳实践**：环境变量用 `PREFIX_SECTION_KEY` 格式，如 `HELLO_SERVER_PORT`。

### 错误 2：类型转换失败不报错

```go
// ❌ 错误代码
func setValueFromString(target reflect.Value, value string) {
    parsed, _ := strconv.Atoi(value) // 忽略错误！
    target.SetInt(int64(parsed))     // 解析失败时设为 0
}

// 后果：SERVER_PORT=abc 被设为 0，难以排查
```

**修复**：
```go
// ✅ 修复：返回错误
func setValueFromString(target reflect.Value, value string) error {
    switch target.Kind() {
    case reflect.Int:
        parsed, err := strconv.Atoi(value)
        if err != nil {
            return fmt.Errorf("parse int: %w", err)
        }
        target.SetInt(int64(parsed))
    }
    return nil
}
```

### 错误 3：嵌套结构体标签不完整

```go
// ❌ 错误代码
type appConfig struct {
    Server serverConfig `json:"server"` // 缺少 config 标签
}

// loadConfigFile 无法识别嵌套字段
```

**修复**：
```go
// ✅ 修复：所有需要的标签都要写
type appConfig struct {
    Server serverConfig `json:"server" config:"server"`
}

type serverConfig struct {
    Port int `json:"port" config:"port" env:"SERVER_PORT"`
}
```

## 动手练习

### 练习 1：预测输出

阅读以下配置代码，预测输出（先自己想，再看答案）：

```go
// 默认配置
func defaultConfig() appConfig {
    return appConfig{
        AppName:  "hello-go",
        LogLevel: "info",
        Server: serverConfig{
            Port: 8080,
        },
    }
}

// YAML 文件
// app_name: hello-go-yaml
// server:
//   port: 8081

// 环境变量
// HELLO_LOG_LEVEL=warn
// HELLO_SERVER_PORT=9090

cfg, _ := resolveConfig([]string{"config.yaml"}, "HELLO", lookup)
fmt.Println(cfg.AppName, cfg.LogLevel, cfg.Server.Port)
```

<details>
<summary>点击查看答案</summary>

**输出**：`hello-go-yaml warn 9090`

**解析**：
1. 默认值：`hello-go`, `info`, `8080`
2. YAML 覆盖：`hello-go-yaml`, `8081`（LogLevel 不变）
3. 环境变量覆盖：`warn`, `9090`（AppName 不变）

**优先级**：默认值 < 文件 < 环境变量
</details>

### 练习 2：添加新的配置项

在现有配置结构中添加一个新字段 `Database.MaxIdleConns`，支持从 JSON、YAML、环境变量读取。

<details>
<summary>点击查看答案</summary>

```go
// 1. 修改结构体
type databaseConfig struct {
    Driver       string `json:"driver" config:"driver" env:"DATABASE_DRIVER"`
    DSN          string `json:"dsn" config:"dsn" env:"DATABASE_DSN"`
    MaxOpenConns int    `json:"max_open_conns" config:"max_open_conns" env:"DATABASE_MAX_OPEN_CONNS"`
    MaxIdleConns int    `json:"max_idle_conns" config:"max_idle_conns" env:"DATABASE_MAX_IDLE_CONNS"`
}

// 2. 修改默认值
func defaultConfig() appConfig {
    return appConfig{
        // ...
        Database: databaseConfig{
            MaxOpenConns: 2,
            MaxIdleConns: 1, // 新增
        },
    }
}

// 3. 无需修改绑定逻辑（反射自动处理）
```

**关键**：只要标签完整，`bindMap` 和 `bindEnv` 会自动处理新字段。
</details>

### 练习 3：实现配置验证

添加一个 `Validate()` 方法，检查配置是否合法（如端口范围 1-65535）。

<details>
<summary>点击查看答案</summary>

```go
func (c *appConfig) Validate() error {
    if c.Server.Port < 1 || c.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", c.Server.Port)
    }
    if c.LogLevel != "debug" && c.LogLevel != "info" && 
       c.LogLevel != "warn" && c.LogLevel != "error" {
        return fmt.Errorf("invalid log level: %s", c.LogLevel)
    }
    if c.Database.Driver == "" {
        return errors.New("database driver is required")
    }
    return nil
}

// 使用
cfg, err := resolveConfig(paths, prefix, lookup)
if err != nil {
    return err
}
if err := cfg.Validate(); err != nil {
    return fmt.Errorf("invalid config: %w", err)
}
```
</details>

## 故障排查 (FAQ)

### Q1: 为什么环境变量没有生效？

**排查步骤**：

1. **检查 key 是否匹配**：
   ```bash
   # 打印所有 HELLO_ 开头的环境变量
   env | grep ^HELLO_
   ```

2. **检查标签是否正确**：
   ```go
   type serverConfig struct {
       Port int `env:"SERVER_PORT"` // 必须和实际环境变量一致
   }
   ```

3. **检查前缀拼接**：
   ```go
   // 如果 prefix="HELLO"，实际查找 "HELLO_SERVER_PORT"
   loadEnvConfig(cfg, "HELLO", lookup)
   ```

**常见原因**：大小写不一致、前缀错误、标签缺失。

### Q2: JSON 和 YAML 应该如何选择？

**JSON 的优势**：
- 严格的语法，解析库成熟
- 适合机器生成（如脚本、工具）
- 支持复杂类型（数组、嵌套）

**YAML 的优势**：
- 可读性更好，适合手编辑
- 支持注释
- 支持多行字符串

**建议**：
- **开发环境**：用 YAML，方便手动调整
- **CI/CD、生产**：用 JSON，减少解析错误
- **混合使用**：默认 YAML，CI 生成 JSON

### Q3: 如何在测试中注入配置？

**方法 1：用 map 模拟环境变量**：
```go
func TestConfig(t *testing.T) {
    lookup := mapLookup(map[string]string{
        "HELLO_SERVER_PORT": "9999",
    })
    
    cfg, err := resolveConfig([]string{}, "HELLO", lookup)
    if err != nil {
        t.Fatal(err)
    }
    
    if cfg.Server.Port != 9999 {
        t.Errorf("expected 9999, got %d", cfg.Server.Port)
    }
}

// 辅助函数
func mapLookup(values map[string]string) func(string) (string, bool) {
    return func(key string) (string, bool) {
        value, ok := values[key]
        return value, ok
    }
}
```

**方法 2：用临时文件**：
```go
func TestConfigFile(t *testing.T) {
    file, _ := os.CreateTemp("", "config-*.json")
    defer os.Remove(file.Name())
    
    content := `{"server": {"port": 8888}}`
    file.WriteString(content)
    file.Close()
    
    cfg, _ := loadConfigFile(defaultConfig(), file.Name())
    // 断言...
}
```

## 知识扩展 (选学)

### Viper 库的设计思路

[Viper](https://github.com/spf13/viper) 是最流行的 Go 配置库。它的核心思想和本章示例一致：

```go
viper.SetDefault("port", 8080)      // 默认值
viper.SetConfigFile("config.yaml")  // 配置文件
viper.AutomaticEnv()                // 环境变量
viper.ReadInConfig()                // 加载
```

**区别**：
- Viper 支持更多格式（TOML、HCL、.env）
- Viper 支持配置热重载（watch）
- Viper 支持远程配置中心（etcd、Consul）

**建议**：小项目用本章方法（轻量），大项目用 Viper（功能全）。

### 配置中心的原理

在微服务架构中，配置通常集中存储在 etcd、Consul 或 Apollo 等配置中心：

```
┌─────────────┐      ┌─────────────┐
│   App       │ ←─── │   etcd      │
│             │ poll │  /config    │
└─────────────┘      └─────────────┘
```

**工作流程**：
1. 应用启动时从配置中心拉取配置
2. 定时轮询或 watch 配置变化
3. 配置更新后热重载（无需重启）

**Go 实现思路**：
```go
func watchConfig(etcdClient *clientv3.Client) {
    ch := etcdClient.Watch(context.Background(), "/config")
    for resp := range ch {
        cfg := parseConfig(resp.Kvs[0].Value)
        hotReload(cfg)
    }
}
```

### 配置加密

敏感配置（如数据库密码）通常需要加密存储：

```go
// 环境变量存储加密值
// HELLO_DATABASE_PASSWORD=enc:AES256(xxx)

func decryptIfNeeded(value string) (string, error) {
    if strings.HasPrefix(value, "enc:") {
        return decrypt(value[4:])
    }
    return value, nil
}
```

**最佳实践**：
- 开发环境：明文
- 生产环境：用 KMS（如 AWS Secrets Manager）
- 代码中不存储密钥

## 工业界应用

### 场景 1：Kubernetes ConfigMap

Kubernetes 用 ConfigMap 管理配置，通过环境变量或文件挂载注入：

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
      - name: app
        env:
        - name: HELLO_SERVER_PORT
          value: "9090"
        - name: HELLO_DATABASE_DSN
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: dsn
```

**Go 代码无需修改**，环境变量自动覆盖：
```go
cfg, _ := resolveConfig([]string{}, "HELLO", os.LookupEnv)
```

### 场景 2：多环境配置

```bash
# 目录结构
configs/
├── default.yaml     # 默认配置
├── development.yaml # 开发环境覆盖
├── test.yaml        # 测试环境覆盖
└── production.yaml  # 生产环境覆盖
```

**启动命令**：
```bash
# 开发环境
./app --config=configs/default.yaml --config=configs/development.yaml

# 生产环境
./app --config=configs/default.yaml --config=configs/production.yaml
```

**代码支持**：
```go
cfg, err := resolveConfig(
    []string{"configs/default.yaml", envFile},
    "HELLO",
    os.LookupEnv,
)
```

### 场景 3：特性开关 (Feature Flags)

```go
type featureFlags struct {
    EnableNewUI    bool `env:"ENABLE_NEW_UI"`
    EnableBetaAPI  bool `env:"ENABLE_BETA_API"`
    RolloutPercent int  `env:"ROLLOUT_PERCENT"`
}

func (f *featureFlags) isEnabled(user string) bool {
    if !f.EnableNewUI {
        return false
    }
    // 灰度发布：按用户 ID 哈希决定
    return hash(user)%100 < f.RolloutPercent
}
```

**价值**：配置决定功能开关，无需重新部署。

## 小结

### 核心要点

1. **分层配置**：默认值 → 配置文件 → 环境变量，优先级递增
2. **结构体标签**：用 `json`、`config`、`env` 标签声明映射关系
3. **反射绑定**：用 `reflect` 包自动将配置值赋给结构体字段
4. **类型转换**：处理 JSON/YAML 到 Go 类型的转换（如 float64→int）
5. **错误处理**：解析失败时返回清晰错误，包含字段名

### 关键术语

| 英文 | 中文 | 说明 |
|------|------|------|
| configuration | 配置 | 程序运行参数 |
| default value | 默认值 | 代码中预设的安全值 |
| environment variable | 环境变量 | 操作系统级别配置 |
| struct tag | 结构体标签 | 字段的元数据 |
| reflection | 反射 | 运行时检查类型信息 |
| hot reload | 热重载 | 不重启程序更新配置 |

### 下一步建议

1. 为你的项目添加配置结构体和默认值
2. 实现 JSON/YAML 文件加载支持
3. 添加环境变量覆盖功能
4. 用 `go test` 编写配置加载测试
5. 考虑是否需要引入 Viper 等成熟框架

## 术语表

| 术语 | 英文 | 说明 |
|------|------|------|
| 配置管理 | Configuration Management | 管理程序运行参数的系统 |
| 环境变量 | Environment Variable | 操作系统级别的键值对，用于部署时覆盖配置 |
| 配置文件 | Config File | 存储配置信息的 JSON/YAML 等格式文件 |
| 默认值 | Default Value | 代码中定义的兜底配置值 |
| 结构体标签 | Struct Tag | Go 结构体字段的元数据，用于映射配置键 |
| 反射 | Reflection | 运行时检查类型信息的能力，用于自动绑定配置 |
| 分层配置 | Layered Configuration | 多个配置来源按优先级合并的模式 |
| 优先级 | Precedence | 配置来源的覆盖顺序（环境变量 > 文件 > 默认值） |
| 哨兵错误 | Sentinel Error | 预定义的错误值，用于标识特定配置错误类型 |

## 源码

完整示例代码位于：[internal/advance/config/config.go](../../internal/advance/config/config.go)
