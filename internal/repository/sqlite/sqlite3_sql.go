package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func Sqlite3_demo() {
	var file = "data/foo.db"

	os.Remove(file)

	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table foo (id integer not null primary key, name text,age integer,sex integer);
	delete from foo;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	// begin transcation and insert into data
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into foo(id, name,age) values(?, ?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := 0; i < 100; i++ {
		//随机生成年龄
		var age = rand.Intn(30-20) + 20

		_, err = stmt.Exec(i, fmt.Sprintf("你好，世界.%03d", i), age)
		if err != nil {
			log.Fatal(err)
		}
	}
	// commit transcation
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	//query data by select all data
	rows, err := db.Query("select id, name, age from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var age int
		err = rows.Scan(&id, &name, &age)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name, age)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	//query one record by select id
	stmt, err = db.Prepare("select name from foo where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow("3").Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)

	// delete table
	_, err = db.Exec("delete from foo")
	if err != nil {
		log.Fatal(err)
	}

	//again insert record
	_, err = db.Exec("insert into foo(id, name,age) values(1, 'foo',21), (2, 'bar',22), (3, 'baz',21)")
	if err != nil {
		log.Fatal(err)
	}

	//select record
	query(db, "select id, name, age from foo")

}

func query(db *sql.DB, sql string) {
	var rows, err = db.Query(sql)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var age int
		err = rows.Scan(&id, &name, &age)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, name, age)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
