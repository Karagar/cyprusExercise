package main

import (
	"github.com/Karagar/cyprusExercise/pkg/utils"
)

func main() {
	server := utils.ServerStruct{}
	server.Declare()
	server.Serve()
}
