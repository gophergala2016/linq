//package lib
package main

import(
	//"github.com/gophergala2016/linq/connector"
	"strconv"
	"fmt"
	"../connector/"
	"log"
	"os"
	"strings"
)

func get_default_length(data_type string) int{
	switch data_type{
		case "varchar":
			return 250
	}
	return 10
}

type ColumnBuilder struct{
	name string
	data_type string
	length int
	null bool
	primary_key bool
	index bool
	auto_increment bool
	default_value string
	new_name string
	table string
	foreignKey string
}

func (this ColumnBuilder) creation_string() string{
	return this.name+this.new_name_get()+this.data_type_get()+this.null_get()+this.primary_key_get()+this.default_value_get()+this.auto_increment_get()
}

func (this ColumnBuilder) null_get() string{
	if this.null{
		return ""
	}
	return " NOT NULL "
}
func (this ColumnBuilder) primary_key_get() string{
	if this.primary_key{
		return " PRIMARY KEY "
	}
	return ""
}

func (this ColumnBuilder) data_type_get() string{
	if this.data_type != "" && this.data_type != "varchar"{
		return " "+this.data_type
	}else if(this.data_type == "varchar" || this.data_type == "nvarchar"){
		return " "+this.data_type+"("+this.length_get()+")"
	}
	return " varchar("+this.length_get()+") "
}

func (this ColumnBuilder) length_get() string{
	var str string
	if this.length == 0{
		str = strconv.Itoa(get_default_length(this.data_type))
		return str
	}
	str = strconv.Itoa(this.length)
	return str
}

func (this ColumnBuilder) default_value_get() string{
	if this.default_value != ""{
		return " DEFAULT '"+this.default_value+"' "
	}
	return ""
}

func (this ColumnBuilder) auto_increment_get() string{
	if this.auto_increment{
		return " AUTO_INCREMENT "
	}
	return ""
}

func (this ColumnBuilder) new_name_get() string{
	return " "+this.new_name+" "
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
	fmt.Println("el query es ", query)
	connector.Query(query)
}

func RemoveColumn(table,column string){
	query := "ALTER TABLE "+table+" DROP COLUMN "+ column
	connector.Query(query)
}

func AddColum(table string, this ColumnBuilder){
	acceptValues := []string {"nvarchar","varchar"}
	query := "ALTER TABLE "+table+" ADD COLUMN "+ this.name + " " + this.data_type + ""
	if(contains(acceptValues, this.data_type)){
		if(this.length <= 0){
			 this.length = 255
		}
		query += "(" + 	strconv.Itoa(this.length)  + ")"
	}
	connector.Query(query)
}

func ChangeColumn(table string,column ColumnBuilder){
	var modifier string
	if column.new_name != ""{
		modifier = "CHANGE"
	}else{
		modifier = "MODIFY"
	}
	query := "ALTER TABLE "+table+" "+ modifier +" "+ column.creation_string()

	connector.Query(query)
}

func AddIndex(table,index_name,column string) {
	query := "CREATE INDEX "+index_name+" ON "+ table + "("+column+")"
	connector.Query(query)
}

func RemoveIndex(table,index_name string){
	query := "DROP INDEX "+index_name+" ON "+ table
	fmt.Println(query)
	connector.Query(query)
}

func AddForeignKey(col1 ColumnBuilder, col2 ColumnBuilder ){
	query := "ALTER TABLE " + col1.table + " ADD FOREIGN KEY (" + col1.foreignKey + ") "
	query += "REFERENCES " + col2.table +  "(" +  col2.foreignKey  + ")"
	connector.Query(query)
}

func RemoveForeigKey(this ColumnBuilder){
	query := "ALTER TABLE " + this.table + " DROP FOREIGN KEY "  + this.foreignKey
	connector.Query(query)
}

func DropTable(table string){
	query := "DROP TABLE " + table
	connector.Query(query)
}

/*Region Internal*/
func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func main() {
	GenerateShema()
}

func ExistFolder(folderName string) bool {
    _, err := os.Stat(folderName)
    return !os.IsNotExist(err)
}

func GenerateShema(){
	filePath := "./database/shema.yml"
	if ExistFolder(filePath){
		os.Remove(filePath)
	}
	file, _ :=  os.Create(filePath)

	rows := connector.GetQuery("SHOW TABLES")
	for rows.Next(){
		var table string
		err := rows.Scan(&table)
		if err != nil{
			log.Fatal(err)
		}
		finalQuery :="CREATE TABLE " + table + "(\n"
		columnQuery := GetValuesForTable(table)
		finalQuery += columnQuery
		finalQuery += "\n)\n"

		file.WriteString(finalQuery)
	}
}

func GetValuesForTable(table string) string {
	var result string
	query := "SHOW COLUMNS FROM "+ table
	rows := connector.GetQuery(query)

	for rows.Next(){
		var field string
		var data_type string
		var null string
		var key string
		var data_default *string //Regresa null
		var extra string

		err := rows.Scan(&field, &data_type, &null, &key, &data_default, &extra)
		if err != nil{
			log.Fatal(err)
		}

		query := "\t" + GetFormatLine(field, data_type, null, key, extra) + "-"
		result += query
	}
	return SetFinalFormat(result)
}
func GetFormatLine(field string, data_type string, null string, key string, extra string)string{
	valueNull := ""
	valueKey := ""
	valueExtra := strings.ToUpper(extra)

	if key == "PRI"{
		valueKey = "PRIMARY KEY"
	}

	if null == "NO"{
		valueNull = "NOT NULL"
	}
	return fmt.Sprintf("%s %s %s %s %s", field, data_type, valueNull, valueExtra, valueKey)
}

func SetFinalFormat(querySentence string)string{
	var finalString string
	splitLine := strings.Split(querySentence, "-")
	length := cap(splitLine) - 1

	for i := 0;  i < length  ; i++{
		value := splitLine[i]
		if(i != length - 1){
				value = value + "," + "\n"
		}
		finalString+= value
	}
	return finalString
}
