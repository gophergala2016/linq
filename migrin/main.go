package main

import(
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"time"
	"os"
	"log"
	"bufio"
	"github.com/gophergala/linq/migrin"
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
	_,err = w.WriteString("package main \n\n import(\n\t 'fmt' \n)\n\n func main(){}")

	if err != nil{
		log.Fatal(err)
	}

	f.Sync()
	w.Flush()

	this.save_migration_in_db(timestamp)
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
func main() {
	kingpin.Parse()
	m := Migrin{}
	switch *action{
		case "new":
			m.new()
		case "init":
			m.init()
	}
}
