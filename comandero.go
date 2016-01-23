package main

import(
  "os"
  "fmt"
  "gopkg.in/alecthomas/kingpin.v2"
)
var (
	action = kingpin.Arg("action","Init").Required().String()
)


const initFolderName = "./Init"
const initFileName = "./Init/confing.xaml"

func helper(){
  fmt.Println("Aqui ira el archvio helper")
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func existFolder(folderName string) bool {
    _, err := os.Stat(folderName)
    return !os.IsNotExist(err)
}

func validateCommand(command string)bool  {
  listCommands := []string {"Init", "init"}
  return contains( listCommands , command)
}

func createInitFolder(){
  folder := initFolderName //Refactor en caso de verse lento
  localPathFile :=  initFileName

  if !existFolder(folder){
    os.Mkdir(folder,0777)
    file, _ :=  os.Create(localPathFile)
    file.WriteString("path:\nusername:\npassword:\nport:\n")
  }
}

func initFunction(){
  createInitFolder()
}

func main(){
  kingpin.Parse()
	switch *action{
  case "Init":
    initFunction()
  case "init":
			initFunction()
	}
}
