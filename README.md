## Migrin

A migration toolkit writted in Golang that allows you to create the SQL of the migrations using a DSL in Go.

## Installation

Execute the go get command to the toolkit repository
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
		lib.DropTable("courses"})	
	}
```

## Contribuitors

* [Eduardo](https://github.com/eduardo78d)
* [Uriel](https://github.com/urielhdz)

