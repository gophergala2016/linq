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
  if len(os.Args) != 2 {
    helper()
    os.Exit(1)
  }
  var command string
  command = os.Args[1]

  if validateCommand(command){
    initFunction()
  }else{
    fmt.Println("Lalo")
    helper()
  }
}
