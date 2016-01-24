## Migrin

A migration toolkit writted in Golang and inspired on Rails' ActiveRecord::Migration, that allows you to create migration files to make change on your DB using a DSL.

## Installation

Execute the following go get command to install the toolkit from the github repository
```
go get github.com/gophergala2016/linq/migrin
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
		migrin up
	```
6. In case something went wrong, you can reverse your migrations
	```
		migrin down
	```

##Features

* DSL to modify a database without writting SQL in a more expresive way <3
* Generates a schema file with all the instructions to regenerate the DB
* Creates two files per migration (one to execute the migration, and another to undo it)
* Accepts config data for both production and development environments
* Independent from your project, can be used with other languages projects as long as you have Go installed
* Build migrations from the command line, without touching a file or thinking on SQL

##Production

As you may see, there are two environments on config.yml production|development this allows you to separate the credentials on each database, and only run changs on your production DB when explicitly specified.

To run commands in production append the `-production` flag, or it shorthand `-p`, to migrin commands

The following example executes migrations with production credentials
```
migrin up -p
```

##Auto generate migrations

Migrin commands are pretty awesome, you can't even tell it what changes dou you want to make tou your database through the terminal, no need to open migrations or remembering SQL. Here's an example

Given that I want to create a users table with email,password and age fields:

```
migrin new CreateUsersTable create_table users email:varchar password:varchar age:int
```

Migrin will generate the following migration:

```go
package main 

import(
	"github.com/gophergala2016/linq/migrator"
	 "os"
)

func main(){
	migrator.Options(os.Args)
	columns = []migrator.ColumnBuilder{{Name:"email",Data_type:"varchar"},{Name:"password",Data_type:"varchar"},{Name:"age",Data_type:"int"}
	migrator.CreateTable("users",columns)
}
```

And there you have, your table ready to go... this is the sintax to generate such awesome migrations:

```
migrin new migration_name command table_name column_name:data_type...
```

There's no limit on how much columns you can add, every argument after the table_name is considered a new column for the table.

Here are the available commands:

| command       | action                                                |
|---------------|-------------------------------------------------------|
| create_table  | To create a new table                                 |
| add_columns   | Add specified columns to the indicated table          |
| remove_column | Remove the specified columns from the indicated table |

##API

The beauty of migrin is that you don't need to write SQL to define what your migration should do, you use a simple API to modify your database

###CreateTable(table_name,[]ColumnBuilder)
Creates a new table with the specified name (1st argument) and the specified columns, defined by an slice of ColumnBuilder's

Example
```go
package main 
import(
	"github.com/gophergala2016/linq/migrator"
)
func main(){
	migrator.Options(os.Args)
	migrator.CreateTable("courses",[]migrator.ColumnBuilder{{Name:"title"},{Name:"description"}})	
}
```

###DropTable(table_name)

Drops the specified table from the database

Example
```go
package main 
import(
	"github.com/gophergala2016/linq/migrator"
)
func main(){
	migrator.Options(os.Args)
	migrator.DropTable("courses")	
}
```

###AddColumn(table_name,ColumnBuilder{})

Adds a column to an already created table

Example
```go
package main 
import(
	"github.com/gophergala2016/linq/migrator"
)
func main(){
	migrator.Options(os.Args)
	migrator.AddColumn("courses",ColumnBuilder{Name:'status',Data_type:'int'})	
}
```

###RemoveColumn(table_name,column_name)

Removes a column from the specified table

Example
```go
package main 
import(
	"github.com/gophergala2016/linq/migrator"
)
func main(){
	migrator.Options(os.Args)
	migrator.RemoveColumn("courses","status")	
}
```

###ChangeColumn(table_name,ColumnBuilder{})

Changes the column structure of an existing column from the specified table, the column name is obtained from the ColumBuilder struct, if you want to change the name of the column you need to fill the new_name attribute from the ColumnBuilder struct

Example
```go
package main 
import(
	"github.com/gophergala2016/linq/migrator"
)
func main(){
	migrator.Options(os.Args)
	//Changes the column name status to state
	migrator.ChangeColumn("courses",migrator.ColumnBuilder{Name:"status",New_name:"state"})	
}
```

###AddIndex(table,index_name,column)

Adds an index to the specified table and column

Example
```go
package main 
import(
	"github.com/gophergala2016/linq/migrator"
)
func main(){
	migrator.Options(os.Args)
	migrator.AddIndex("courses","status_index","status")	
}
```

###RemoveIndex(table,index_name)

Removes the specified index from the specified talbe

Example
```go
package main 
import(
	"github.com/gophergala2016/linq/migrator"
)
func main(){
	migrator.Options(os.Args)
	migrator.RemoveIndex("courses","status_index")	
}
```

###AddForeignKey(ColumnBuilder{}, ColumnBuilder{})

Creates a foreign key between two columns, the tables of each column are specified in a table attribute in the ColumnBuilder struct

Example
```go
package main 
import(
	"github.com/gophergala2016/linq/migrator"
)
func main(){
	migrator.Options(os.Args)
	migrator.AddForeignKey(migrator.ColumnBuilder{Name:"id",Table:"courses"},ColumnBuilder{Name:"course_id",Table:"videos"})	
}
```

###RemoveForeigKey(ColumnBuilder{})

Removes a foreign key, the foreign key to eliminate is specified in the ForeignKey attribute of a ColumnBuilder

Example
```go
package main 
import(
	"github.com/gophergala2016/linq/migrator"
)
func main(){
	migrator.Options(os.Args)
	migrator.RemoveForeignKey(migrator.ColumnBuilder{ForeignKey:"foreign_key"})	
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

##GopherGala 2016

Created for the GopherGala2016 hackathon.

## RoadMap

* Add other drives (for now it only supports mySQL)
* Generate complete migrations (including methods) using the command line as seen in Rails ActiveRecord::Migration
* Use Go concurrency efficiency for better performance

## Contribuitors

* [Eduardo](https://github.com/eduardo78d)
* [Uriel](https://github.com/urielhdz)

