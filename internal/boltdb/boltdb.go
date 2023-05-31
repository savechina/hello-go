package boltdb

import (
	"fmt"
	"log"

	"go.etcd.io/bbolt"
)

func Bolt_test() {
	// 打开数据库文件，如果不存在则创建
	db, err := bbolt.Open("data/hello.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 开始一个读写事务
	err = db.Update(func(tx *bbolt.Tx) error {
		// 创建一个名为 "MyBucket" 的 bucket
		b, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			return err
		}
		// 向 bucket 中写入一些键值对
		err = b.Put([]byte("foo"), []byte("bar"))
		if err != nil {
			return err
		}
		err = b.Put([]byte("hello"), []byte("world"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// 开始一个只读事务
	err = db.View(func(tx *bbolt.Tx) error {
		// 获取 "MyBucket" 的引用
		b := tx.Bucket([]byte("MyBucket"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}
		// 从 bucket 中读取键值对并打印出来
		v := b.Get([]byte("foo"))
		fmt.Printf("foo: %s\n", v)
		v = b.Get([]byte("hello"))
		fmt.Printf("hello: %s\n", v)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
