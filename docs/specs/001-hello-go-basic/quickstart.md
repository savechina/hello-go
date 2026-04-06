# Quickstart: hello-go 学习样例

## 前置条件

```bash
# Go 1.24+
go version

# 依赖安装
make install
```

## 快速开始

### 1. 构建项目

```bash
make build
```

### 2. 运行基础章节

```bash
# 变量与表达式
go run ./cmd/hello basic variables

# 数据类型
go run ./cmd/hello basic datatype

# 并发
go run ./cmd/hello basic concurrency
```

### 3. 运行高级章节

```bash
# 数据库示例
go run ./cmd/hello advance database

# Web 开发
go run ./cmd/hello advance web
```

### 4. 运行实战项目

```bash
# Web 服务
go run ./cmd/hello awesome webservice
```

### 5. 查看文档

```bash
# 本地预览
mdbook serve docs/

# 浏览器打开 http://localhost:3000
```

### 6. 运行测试

```bash
make test
```

### 7. 代码质量检查

```bash
make fmt    # 格式化
make vet    # 静态分析
make lint   # 完整 lint
make verify # 完整质量门禁
```

## 项目结构速览

```
hello-go/
├── cmd/hello/          # 统一入口
├── internal/basic/     # 基础入门 (≥12 章)
├── internal/advance/   # 高级进阶 (≥8 章)
├── internal/awesome/   # 精选实战 (≥3 项目)
└── docs/src/           # mdBook 文档
```

## 学习路径

1. **初学者**: 从 `basic/variables` 开始，按顺序学习每个章节
2. **进阶者**: 直接跳转到 `advance/` 对应主题
3. **实战**: 完成基础+高级后，参考 `awesome/` 项目
