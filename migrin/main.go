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
	"../lib"
	"io/ioutil"
	"strings"
	"bytes"
)

type columns []lib.ColumnBuilder

func (this *columns) Set(value string) error{
	if value == ""{
		return fmt.Errorf("'%s' is not a valid column name", value)
	}else{
		components := strings.Split(value,":")
		c_b := lib.ColumnBuilder{Name: components[0]}
		if len(components) > 1{
			c_b.Data_type = components[1]
		}
		*this = append(*this,c_b)
		return nil
	}

}

func (this *columns) String() string{
	return ""
}

func (this *columns) IsCumulative() bool {
  return true
}

func ColumnList(s kingpin.Settings) (target *[]lib.ColumnBuilder) {
  target = new([]lib.ColumnBuilder)
  s.SetValue((*columns)(target))
  return
}



var (
	action = kingpin.Arg("action","Specify an action to run init/new/up/down").Required().String()
	option = kingpin.Arg("modifier","Adds extra information to the command, specifies the migration name on new command").String()
	method = kingpin.Arg("method", "method").String()
	table = kingpin.Arg("table","method").String()
	p = kingpin.Flag("production","Runs command queries on production ").Short('p').Bool()
	column_args = ColumnList(kingpin.Arg("columns","N number of columns to add to your migration"))
)
const initFolderName = "./database/"
const initFileName = "./database/config.yml"
const initFolderNameMigration = "./database/migrations"
const initFolderNameMigrationDown = "./database/migrations/downs"


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
	create_file_migration(folder+"/"+timestamp+"_"+filename+".go")

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
func (this Migrin) remove_migration_from_db(timestamp string){
	connector.RemoveMigration(timestamp)
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
	this.create_migrations_table()
	files, _ := ioutil.ReadDir(initFolderNameMigration)
  for _, f := range files {
  	extension := strings.Split(f.Name(),".")
  	if len(extension) < 1 || extension[len(extension)-1] != "go"{
  		continue
  	}
  	name_components := strings.Split(f.Name(),"_")
  	if len(name_components) > 0 && !migration_executed(name_components[0]){
  		execute_migration(f,initFolderNameMigration)
  		connector.Query("INSERT INTO migrations(migration_id) VALUES('"+name_components[0]+"')")
  	}
  }
}

func (this Migrin) down() {
	files, _ := ioutil.ReadDir(initFolderNameMigrationDown)
  for _, f := range files {

  	name_components := strings.Split(f.Name(),"_")
  	if len(name_components) > 0 && migration_executed(name_components[0]){
  		execute_migration(f,initFolderNameMigrationDown)
  		connector.Query("DELETE FROM migrations WHERE migration_id = '"+name_components[0]+"'")
  		break
  	}
  }
}

func migration_executed(timestamp string) bool{
	rows := connector.GetQuery("SELECT migration_id FROM migrations WHERE migration_id = '"+timestamp+"'")
	return rows.Next()
}

func create_file_migration(file_path string){
	customeTable := *table
	customeMethod := *method

	f,err := os.Create(file_path)
	defer f.Close()
	if err != nil{
		log.Fatal(err)
	}

	option := GetOption(customeMethod)

	columnBuilderName := "column"
	w := bufio.NewWriter(f)
	imports := "\n\t\"../../lib\"\n\t \"os\"\n"
	main_body := "\n\t//Write here your migration sentences. Next line is necessary for configuration\n\tlib.Options(os.Args)\n"
	if len(*column_args) > 0{
		main_body += "\n\t" + columnBuilderName +":= []lib.ColumnBuilder{"
		for i,column := range *column_args{
			main_body += column.Go_code_string()
			if i < (len(*column_args)-1){
				main_body += ","
			}
		}
		main_body += "}\n"
	}
	line := "package main \n\nimport("+imports+")\n\nfunc main(){"+main_body
	line += "\n\tlib."+ option +"(" + "\"" +customeTable + "\"" + "," + columnBuilderName +")"
	line += "\n}"

	_,err = w.WriteString(line)
	if err != nil{
		log.Fatal(err)
	}
	f.Sync()
	w.Flush()
}

func GetOption(option string) string{
	switch option {
	case "create":
		return "CreateTable"
	case "add_column":
		return "AddColum"
	}
	return ""
}

func execute_migration(file os.FileInfo,file_path string) bool{
	if file != nil{
		cmd := exec.Command("go", "run",file_path+"/"+file.Name(),production_arg())
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

func production_arg() string{
	if *p{
		return "production"
	}
	return ""
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
