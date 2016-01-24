package main 

import(
	"../../lib"
	 "os"
)

func main(){
	//Write here your migration sentences. Next line is necessary for configuration
	lib.Options(os.Args)

	column:= []lib.ColumnBuilder{{Name:"email",Data_type:"nvarchar{255}"}}

	lib.AddColum("table",column)
}