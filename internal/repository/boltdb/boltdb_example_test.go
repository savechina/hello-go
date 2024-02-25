package boltdb

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// Ensure that the page type can be returned in human readable format.
func TestBoltDB_sample(t *testing.T) {

	// Bolt_demo()
}

func TestSomething(t *testing.T) {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("获取项目根目录失败: %v", err)
		return
	}

	// 使用项目根目录进行操作
	fmt.Println("项目根目录:", wd)
}

func TestSomething2(t *testing.T) {
	// 获取当前文件路径
	path, err := os.Executable()
	if err != nil {
		t.Errorf("获取当前文件路径失败: %v", err)
		return
	}

	// 获取项目根目录
	dir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		t.Errorf("获取项目根目录失败: %v", err)
		return
	}

	// 使用项目根目录进行操作
	fmt.Println("项目根目录:", dir)
}

func TestSomethingTemp(t *testing.T) {
	// 创建临时目录
	dir := t.TempDir()

	// 使用临时目录进行操作
	fmt.Println("临时目录:", dir)

	// ...

	// 删除临时目录
	defer os.RemoveAll(dir)
}
