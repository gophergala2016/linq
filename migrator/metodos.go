package migrator

import(
	"github.com/gophergala2016/linq/connector"
	"strconv"
	"fmt"
)

func init(){
	connector.SetEnv("development")
}

func get_default_length(data_type string) int{
	switch data_type{
		case "varchar":
			return 250
	}
	return 10
}

type ColumnBuilder struct{
	Name string
	Data_type string
	Length int
	Null bool
	Primary_key bool
	Index bool
	Auto_increment bool
	Default_value string
	New_name string
	Table string
	ForeignKey string
}

func (this ColumnBuilder) Go_code_string() string{
	return "{Name:\""+this.Name+"\",Data_type:\""+this.Data_type+"\"}"
}

func (this ColumnBuilder) creation_string() string{
	return this.Name+this.new_name_get()+this.data_type_get()+this.null_get()+this.primary_key_get()+this.default_value_get()+this.auto_increment_get()
}

func (this ColumnBuilder) null_get() string{
	if this.Null{
		return ""
	}
	return " NOT NULL "
}
func (this ColumnBuilder) primary_key_get() string{
	if this.Primary_key{
		return " PRIMARY KEY "
	}
	return ""
}

func (this ColumnBuilder) data_type_get() string{
	if this.Data_type != "" && this.Data_type != "varchar"{
		return " "+this.Data_type
	}else if(this.Data_type == "varchar" || this.Data_type == "nvarchar"){
		return " "+this.Data_type+"("+this.length_get()+")"
	}
	return " varchar("+this.length_get()+") "
}



func (this ColumnBuilder) length_get() string{
	var str string
	if this.Length == 0{
		str = strconv.Itoa(get_default_length(this.Data_type))
		return str
	}
	str = strconv.Itoa(this.Length)
	return str
}

func (this ColumnBuilder) default_value_get() string{
	if this.Default_value != ""{
		return " DEFAULT '"+this.Default_value+"' "
	}
	return ""
}

func (this ColumnBuilder) auto_increment_get() string{
	if this.Auto_increment{
		return " AUTO_INCREMENT "
	}
	return ""
}
func (this ColumnBuilder) new_name_get() string{
	return " "+this.New_name+" "
}

func CreateTable(table_name string, columns []ColumnBuilder){
	query := "CREATE TABLE "+table_name+"("
	for index,column := range columns{
		query += column.creation_string()
		if index < (len(columns) -1){
			query += ","
		}
	}
	query +=")"
	connector.Query(query)
}

func RemoveColumn(table,column string){
	query := "ALTER TABLE "+table+" DROP COLUMN "+ column
	connector.Query(query)
}

func ChangeColumn(table string,column ColumnBuilder){
	var modifier string
	if column.New_name != ""{
		modifier = "CHANGE"
	}else{
		modifier = "MODIFY"
	}
	query := "ALTER TABLE "+table+" "+ modifier +" "+ column.creation_string()

	connector.Query(query)
}


func AddColum(table string, this ColumnBuilder){
	acceptValues := []string {"nvarchar","varchar"}
	query := "ALTER TABLE "+table+" ADD COLUMN "+ this.Name + " " + this.Data_type + ""
	if(Contains(acceptValues, this.Data_type)){
		if(this.Length <= 0){
			 this.Length = 255
		}
		query += "(" + 	strconv.Itoa(this.Length)  + ")"
	}
	connector.Query(query)
}

func AddIndex(table,index_name,column string) {
	query := "CREATE INDEX "+index_name+" ON "+ table + "("+column+")"
	connector.Query(query)
}

func AddForeignKey(col1 ColumnBuilder, col2 ColumnBuilder ){
	query := "ALTER " + col1.Table + "ADD FOREIGN KEY (" + col1.ForeignKey + ")"
	query += "RERERENCES " + col2.Table +  "(" +  col2.ForeignKey  + ")"
	connector.Query(query)
}

func RemoveForeigKey(this ColumnBuilder){
	query := "ALTER TABLE" + this.Name + "DROP FOREIGN KEY"  + this.ForeignKey
	connector.Query(query)
}

func RemoveIndex(table,index_name string){
	query := "DROP INDEX "+index_name+" ON "+ table
	fmt.Println(query)
	connector.Query(query)
}

func Options(args []string){
	if len(args) > 1{
		if args[1] == "production"{
			connector.SetEnv("production")
		}
	}
}

func DropTable(table string){
	query := "DROP TABLE " + table
	connector.Query(query)
}

/*Region Internal*/
func Contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
