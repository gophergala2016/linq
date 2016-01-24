package main

import(
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"time"
	"os"
	"os/exec"
	"log"
	"bufio"
	"github.com/gophergala2016/linq/connector"
	"io/ioutil"
	"strings"
	"bytes"
	"strconv"
)

var (
	action = kingpin.Arg("action","Specify an action to run init/new/up/down").Required().String()
	option = kingpin.Arg("modifier","Adds extra information to the command, specifies the migration name on new command").String()
)
const initFolderName = "./database/"
const initFileName = "./database/config.yml"
const initFolderNameMigration = "./database/migrations"

type Migrin struct{

}

func (this Migrin) new() {
	if *option == ""{
		fmt.Println("Missing migration name.")
		return
	}else{
		t := time.Now()
		timestamp := t.Format("20060102150405")
		this.create_file(timestamp,*option)
		this.create_down_file(timestamp,*option)
	}
}

func existFolder(folderName string) bool {
    _, err := os.Stat(folderName)
    return !os.IsNotExist(err)
}

func (this Migrin) create_file(timestamp,filename string) {
	folder := initFolderNameMigration
	if !existFolder(folder){
    os.Mkdir(folder,0777)
 	}
	this.create_migrations_table() // Run concurrently
	f,err := os.Create(folder+"/"+timestamp+"_"+filename+".go")
	defer f.Close()
	if err != nil{
		log.Fatal(err)
	}
	w := bufio.NewWriter(f)
	_,err = w.WriteString("package main \n\nimport(\n\t\"github.com/gophergala2016/linq/lib\"\n)\n\nfunc main(){}")

	if err != nil{
		log.Fatal(err)
	}

	f.Sync()
	w.Flush()

	this.save_migration_in_db(timestamp)
}

func (this Migrin) create_down_file(timestamp,filename string) {
	folder := initFolderNameMigration+"/downs"
	if !existFolder(folder){
    os.Mkdir(folder,0777)
 	}
	f,err := os.Create(folder+"/"+timestamp+"_"+filename+".go")
	defer f.Close()
	if err != nil{
		log.Fatal(err)
	}
	w := bufio.NewWriter(f)
	_,err = w.WriteString("package main \n\nimport(\n\t\"github.com/gophergala2016/linq/lib\" \n)\n\nfunc main(){}")
	if err != nil{
		log.Fatal(err)
	}
	f.Sync()
	w.Flush()
}

func (this Migrin) init(){
  folder := initFolderName //Obtencion de variables globales para realizar  la operación más rápido
  localPathFile :=  initFileName

  if !existFolder(folder){
    err := os.Mkdir(folder,0777)
    if err != nil{
    	log.Fatal(err)
    }
  }

	file, _ :=  os.Create(localPathFile)
	file.WriteString("path:\nusername:\npassword:\nport:\ndatabase:\n")
}

func (this Migrin) save_migration_in_db(timestamp string){
	connector.InsertMigration(timestamp)
}

func (this Migrin) create_migrations_table() {
	connector.Run()
}

func (this Migrin) up() {
	rows := connector.GetQuery("SELECT id,migration_id FROM migrations WHERE status = 0")
	for rows.Next(){
		var id int
		var timestamp string 
		err := rows.Scan(&id,&timestamp)
		if err != nil{
			log.Fatal(err)
		}
		if execute_migration(timestamp,"./database/migrations/"){
			fmt.Println("Migration "+timestamp+" was executed")
			connector.Query("UPDATE migrations SET status = 1 WHERE id = "+strconv.Itoa(id))
		}
	}
}

func (this Migrin) down() {
	rows := connector.GetQuery("SELECT id,migration_id FROM migrations WHERE status = 1 ORDER BY id DESC LIMIT 1 ")
	for rows.Next(){
		var id int
		var timestamp string 
		err := rows.Scan(&id,&timestamp)
		if err != nil{
			log.Fatal(err)
		}
		if execute_migration(timestamp,"./database/migrations/downs/"){
			fmt.Println("Migration "+timestamp+" was reversed")
			connector.Query("UPDATE migrations SET status = 0 WHERE id = "+strconv.Itoa(id))
		}
	}
}

func execute_migration(timestamp,file_path string) bool{
	fmt.Println(timestamp)
	file := find_file(timestamp)
	if file != nil{
		cmd := exec.Command("go", "run",file_path+file.Name())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
		    fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		    return false
		}else{
			return true
		}
	}
	return false
}

func find_file(timestamp string) os.FileInfo{
	files, _ := ioutil.ReadDir("./database/migrations")
  for _, f := range files {
		if strings.Contains(f.Name(),timestamp) {
			return f
		}
  }
  return nil
}

func main() {
	kingpin.Parse()
	m := Migrin{}
	switch *action{
		case "new":
			m.new()
		case "init":
			m.init()
		case "up":
			m.up()
		case "down":
			m.down()
	}
}
