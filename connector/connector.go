package connector

import(
	"database/sql"
	"log"
	"fmt"
	"bufio"
	"strings"
	"os"
)
import _ "github.com/go-sql-driver/mysql"

const table_name = "migrations"

var (
	username string
	password string
	database string
	format string
	connector string
)
const pathConfig = "database/config.yml"

func connect_db() *sql.DB{
	// Change to config
	Initialize()
	db,err := sql.Open(connector,format)

	if(err != nil){
		log.Fatal(err)
		return nil
	}

	return db
}

func Run(){
	db := connect_db()
	rows,err := db.Query("SHOW TABLES LIKE 'migrations'")
	if(err != nil){
		log.Fatal(err)
	}

	if !rows.Next(){
		_,err = db.Exec("CREATE TABLE migrations(id int NOT NULL AUTO_INCREMENT PRIMARY KEY,migration_id varchar(11) NOT NULL,status int DEFAULT 0)")
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

func Initialize(){
	connector = "mysql"
	setValuesConfig()
}

func setValuesConfig(){
	inputFile, err := os.Open(pathConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer inputFile.Close()
	scanner := bufio.NewScanner(inputFile)

	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, ":")

		switch splitLine[0]{
			case "username":
				username = splitLine[1]
			case "password":
				password = splitLine[1]
			case "database":
				database = splitLine[1]
		}
	}
	format = getFormat(username, password, database)

	if err := scanner.Err(); err != nil {
		log.Fatal(scanner.Err())
	}
}

func getFormat(username string, password string, database string)string{
	return fmt.Sprintf("%s:%s@/%s", username, password, database)
}
