package main

import(
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"time"
	"os"
	"log"
	"bufio"
	"./connector"
)

var (
	action = kingpin.Arg("action","Specify an action to run init/new/up/down").Required().String()
	option = kingpin.Arg("modifier","Adds extra information to the command, specifies the migration name on new command").String()
)

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

func (this Migrin) create_file(timestamp,filename string) {
	this.create_migrations_table() // Run concurrently
	f,err := os.Create(timestamp+"_"+filename+".go")
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
	}
}