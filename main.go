package main 

import (
	"github.com/PyMarcus/gobank/api"
)

func main(){
	a := api.NewAPIServer("127.0.0.1:3000")
	a.Run()
}
