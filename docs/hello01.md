开始指南（get startd guide）

工程目录结构：

参考：https://github.com/golang-standards/project-layout

```bash
hello-go git:(master) ✗ ls
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