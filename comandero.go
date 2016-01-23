package main

import(
  "os"
  "fmt"
)
const initFolderName = "./Init"
const initFileName = "./Init/confing.xaml"

func helper(){
  fmt.Println("Aqui ira el archvio helper")
}

func existFolder(folderName string) bool {
    _, err := os.Stat(folderName)
    return !os.IsNotExist(err)
}

func validateCommand(command string)bool  {
  if command != "Init"{ //Se deben de validar por una lista de palabras no por if
    return false
  }
  return true
}

func createInitFolder(){
  folder := initFolderName //Refactor en caso de verse lento
  localPathFile :=  initFileName

  if !existFolder(folder){
    os.Mkdir(folder,0777)
    file, _ :=  os.Create(localPathFile)
    file.WriteString("test\nhello")
  }
}

func initFunction(){
  createInitFolder()
}

func main(){
  if len(os.Args) != 2 {
    helper()
    os.Exit(1)
  }
  var command string
  command = os.Args[1]

  if validateCommand(command){
    initFunction()
  }else{
    helper()
  }
}
