package main

import(
	"fmt"
	"../connector"
	"strconv"
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
}

func (this ColumnBuilder) creation_string() string{
	return this.name+this.data_type_get()+this.null_get()+this.primary_key_get()+this.default_value_get()+this.auto_increment_get()
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
	}else if(this.data_type == "varchar"){
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

func main() {
	columns := []ColumnBuilder{{name:"name",null:true},{name:"email",data_type:"varchar"}}
	CreateTable("audits",columns)
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

func RemoveColumn(table,column){

}



