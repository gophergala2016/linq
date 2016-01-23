package connector

import(
	"database/sql"
	"log"
	"fmt"
)

var (
	manage string
	username string
	password string
	datebase string
	openSql string
)

import "github.com/go-sql-driver/mysql"

func Run(){
	//db,err := sql.Open("mysql","root:@/spam_db")
	db,err := sql.Open("mysql","root:xxpesar1020@/spam_db")

	if(err != nil){
		log.Fatal(err)
	}

	rows,err := db.Query("SHOW TABLES LIKE 'migrations'")
	if(err != nil){
		log.Fatal(err)
	}
	if rows.Next(){
		fmt.Println("Hola DB")
	}else{
		_,err = db.Exec("CREATE TABLE migrations(id int NOT NULL AUTO_INCREMENT PRIMARY KEY,migration_id varchar(11) NOT NULL,status int DEFAULT 0)")
	}
}
