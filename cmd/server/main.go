package main

import (
	"fmt"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/config"
	_ "github.com/lib/pq"
)

func main() {
	config := config.Load()

	fmt.Println("Hello World")
	fmt.Println(config)
}
