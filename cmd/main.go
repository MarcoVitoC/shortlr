package main

import "github.com/MarcoVitoC/shortlr/internal"

func main() {
	server := internal.NewServer("localhost:8080")
	server.Run()
}