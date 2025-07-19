# Getting Started

## Install

1. 下载源代码

`git clone git@github.com:savechina/hello-go.git`

2. 编译构建
`make clean compile`

Make 全部任务
```zsh
➜  hello-go git:(main) make

 Choose a command run in hello-go:

  install   Install missing dependencies. Runs `go get` internally. e.g; make install get=github.com/foo/bar
  start     Start in development mode. Auto-starts when code changes.
  stop      Stop development mode.
  watch     Run given command when code changes. e.g; make watch run="echo 'hey'"
  build     Build the binary
  compile   Compile the binary.
  exec      Run given command, wrapped with custom GOPATH. e.g; make exec run="go test ./..."
  clean     Clean build files. Runs `go clean` internally.
```

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
