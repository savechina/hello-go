开始指南（getstartd guide）

## Install

1. 下载源代码

`git clone git@github.com:savechina/hello-go.git`


2. 编译构建
`make clean compile`

3. 输出：
```
➜  hello-go git:(master) ✗ bin/hello-go
1.0.0
hello world!
```

## 工程目录结构：

参考：[Project Structure](https://github.com/golang-standards/project-layout)

```bash
hello-go git:(master) ✗ tree
├── Makefile
├── README.md
├── bin
├── build
├── cmd
│   └── hello
│       └── main.go
├── configs
├── docs
│   └── rfc0001.md
├── examples
├── go.mod
├── go.sum
├── internal
│   ├── boltdb
│   │   └── boltdb.go
│   ├── first
│   │   └── first_example.go
│   └── version
│       └── version.go
├── pkg
├── test
└── tools
```

