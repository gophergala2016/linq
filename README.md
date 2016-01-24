## Migrin

A migration toolkit writted in Golang that allows you to create the SQL of the migrations using a DSL in Go.

## Installation

Execute the following go get command to install the toolkit from the github repository
```
go get github.com/gophergala2016/linq
```

This gives you access to the migrin command to execute different actions

##Getting started

1. Install the toolkit
2. Execute the initializer
	```
		migrin init
	```
3. Modify database/config.yml with the credentials for your database
4. Create your first migration
	```
		migrin new <MigrationName>
	```
5. Execute your migration
	```
		migrate up
	```
6. In case something went wrong, you can reverse your migrations
	```
		migrate down
	```

##Features

* DSL to modify a database without writting SQL in a more expresive way <3
* Generates a schema file with all the instructions to regenerate the DB
* Creates two files per migration (one to execute the migration, and another to undo it)
* Accepts config data for both production and development environments
* Independent from your project, can be used with other languages projects as long as you have Go installed

##API

The beauty of migrin is that you don't need to write SQL to define what your migration should do, you use a simple API to modify your database

###CreateTable(table_name,[]ColumnBuilder)
Creates a new table with the specified name (1st argument) and the specified columns, defined by an slice of ColumnBuilder's

Example
```go
package main 
import(
	"github.com/gophergala2016/linq/lib"
)
func main(){
	lib.CreateTable("courses",[]lib.ColumnBuilder{{Name:"title"},{Name:"description"}})	
}
```

###DropTable(table_name)

Drops the specified table from the database

Example
```go
package main 
import(
	"github.com/gophergala2016/linq/lib"
)
func main(){
	lib.DropTable("courses")	
}
```

###AddColumn(table_name,ColumnBuilder{})

Adds a column to an already created table

Example
```go
package main 
import(
	"github.com/gophergala2016/linq/lib"
)
func main(){
	lib.AddColumn("courses",ColumnBuilder{Name:'status',Data_type:'int'})	
}
```

###RemoveColumn(table_name,column_name)

Removes a column from the specified table

Example
```go
package main 
import(
	"github.com/gophergala2016/linq/lib"
)
func main(){
	lib.RemoveColumn("courses","status")	
}
```

###ChangeColumn(table_name,ColumnBuilder{})

Changes the column structure of an existing column from the specified table, the column name is obtained from the ColumBuilder struct, if you want to change the name of the column you need to fill the new_name attribute from the ColumnBuilder struct

Example
```go
package main 
import(
	"github.com/gophergala2016/linq/lib"
)
func main(){
	//Changes the column name status to state
	lib.ChangeColumn("courses",ColumnBuilder{Name:"status",New_name:"state"})	
}
```

###AddIndex(table,index_name,column)

Adds an index to the specified table and column

Example
```go
package main 
import(
	"github.com/gophergala2016/linq/lib"
)
func main(){
	lib.AddIndex("courses","status_index","status")	
}
```

###RemoveIndex(table,index_name)

Removes the specified index from the specified talbe

Example
```go
package main 
import(
	"github.com/gophergala2016/linq/lib"
)
func main(){
	lib.RemoveIndex("courses","status_index")	
}
```

###AddForeignKey(ColumnBuilder{}, ColumnBuilder{})

Creates a foreign key between two columns, the tables of each column are specified in a table attribute in the ColumnBuilder struct

Example
```go
package main 
import(
	"github.com/gophergala2016/linq/lib"
)
func main(){
	lib.AddForeignKey(ColumnBuilder{Name:"id",Table:"courses"},ColumnBuilder{Name:"course_id",Table:"videos"})	
}
```

###RemoveForeigKey(ColumnBuilder{})

Removes a foreign key, the foreign key to eliminate is specified in the ForeignKey attribute of a ColumnBuilder

Example
```go
package main 
import(
	"github.com/gophergala2016/linq/lib"
)
func main(){
	lib.RemoveForeignKey(ColumnBuilder{ForeignKey:"foreign_key"})	
}
```

###ColumnBuilder

A struct that defines the atributtes for a column, it's used for multiple methods of the DSL, it accepts the following attributes:

```go
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
```

## Contribuitors

* [Eduardo](https://github.com/eduardo78d)
* [Uriel](https://github.com/urielhdz)

