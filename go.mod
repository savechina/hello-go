module hello

go 1.23

toolchain go1.24.3

require (
	github.com/mattn/go-sqlite3 v1.14.17
	github.com/spf13/cobra v1.8.1
	github.com/spf13/pflag v1.0.6 // indirect
	// github.com/spf13/viper v1.16.0
	go.etcd.io/bbolt v1.4.2
	gorm.io/driver/sqlite v1.5.1
	gorm.io/gorm v1.25.1
)

// indirect
require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
)
