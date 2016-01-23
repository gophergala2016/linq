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
const pathConfig = "database/config.yml"
var (
	username string
	password string
	database string
	format string
	connector string
)

func Run(){
	Initialize()
	//db,err := sql.Open("mysql","root:@/spam_db")
	db,err := sql.Open(connector,format)
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

		if(splitLine[0] == "username"){
				username = splitLine[1]
		}

		elif(splitLine[0] == "password"){
				password = splitLine[1]
		}
		elif(splitLine[0] == "database"){
				database = splitLine[1]
		}
	}
	format := getFormat(username, password, database)

	if err := scanner.Err(); err != nil {
		log.Fatal(scanner.Err())
	}

}

func getFormat(username string, password string, database string)string{
	return fmt.Sprintf("%s:%s@/%s", username, password, database)
}
