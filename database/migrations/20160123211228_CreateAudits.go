package main 

import(
	"github.com/gophergala2016/linq/lib"
)

func main(){
	lib.CreateTable("users",[]lib.ColumnBuilder{{Name:"audits"}})
	lib.AddColumn("users")	
}