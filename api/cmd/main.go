package main

import (
	"github.com/Karagar/cyprusExercise/pkg/server"
)

func main() {
	server := server.New()
	server.Serve()
}
