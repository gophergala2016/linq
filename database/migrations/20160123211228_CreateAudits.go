package main 

import(
	"github.com/gophergala2016/linq/lib"
)

func main(){
	lib.CreateTable("brothers",[]lib.ColumnBuilder{{Name:"audits"}})	
}