package main

import(
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"time"
	"os"
	"os/exec"
	"log"
	"bufio"
	"../connector"
	"io/ioutil"
	"strings"
	"bytes"
	"strconv"
)

var (
	action = kingpin.Arg("action","Specify an action to run init/new/up/down").Required().String()
	option = kingpin.Arg("modifier","Adds extra information to the command, specifies the migration name on new command").String()
	p = kingpin.Flag("production","Runs command queries on production ").Short('p').Bool()
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
	create_file_migration(folder+"/"+timestamp+"_"+filename+".go")
	this.save_migration_in_db(timestamp)
}

func (this Migrin) create_down_file(timestamp,filename string) {
	folder := initFolderNameMigration+"/downs"
	if !existFolder(folder){
    os.Mkdir(folder,0777)
 	}
 	create_file_migration(folder+"/"+timestamp+"_"+filename+".go")
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
	fields := "\n  username:\n  password:\n  port:\n  database:"
	file.WriteString("development:"+fields+"\nproduction:"+fields)
}

func (this Migrin) save_migration_in_db(timestamp string){
	connector.InsertMigration(timestamp)
}

func (this Migrin) create_migrations_table() {
	waiting_channel := make(chan bool)
	go func(){
		connector.Run()
		waiting_channel <- true	
	}()
	b := <-waiting_channel	
	if !b{
		fmt.Println("Error creating migrations table")
	}
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

func create_file_migration(file_path string){
	f,err := os.Create(file_path)
	defer f.Close()
	if err != nil{
		log.Fatal(err)
	}
	w := bufio.NewWriter(f)
	_,err = w.WriteString("package main \n\nimport(\n\t\"github.com/gophergala2016/linq/lib\"\n)\n\nfunc main(){\n\t//Write here your migration sentences\n}")
	if err != nil{
		log.Fatal(err)
	}
	f.Sync()
	w.Flush()

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
	if *p{
		connector.SetEnv("production")	
	}else{
		connector.SetEnv("development")
	}
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
