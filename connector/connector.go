package connector

import(
	"database/sql"
	"log"
	"fmt"
)

import _ "github.com/go-sql-driver/mysql"

func connect_db() *sql.DB{
	// Change to config
	db,err := sql.Open("mysql","root:@/spam_db")
	if(err != nil){
		log.Fatal(err)
		return nil
	}
	return db
}
const table_name = "migrations"

func Run(){
	db := connect_db()
	rows,err := db.Query("SHOW TABLES LIKE 'migrations'")
	if(err != nil){
		log.Fatal(err)
	}
	if rows.Next(){
		fmt.Println("Hola DB")
	}else{
		_,err = db.Exec("CREATE TABLE "+table_name+"(id int NOT NULL AUTO_INCREMENT PRIMARY KEY,migration_id varchar(11) NOT NULL,status int DEFAULT 0)")
	}
}

func InsertMigration(timestamp string){
	db := connect_db()
	_,err := db.Exec("INSERT INTO "+table_name+" (migration_id) VALUES("+timestamp+") ")
	if err != nil{
		log.Fatal(err)
	}
}

func Query(query string){
	db := connect_db()
	_,err := db.Exec(query)
	if err != nil{
		log.Fatal(err)
	}
}
