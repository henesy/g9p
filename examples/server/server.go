package main

import (
	"fmt"
	"time"
p9	"github.com/henesy/g9p"
)


// Demonstrate server functionality
func main() {
	srv, err := p9.MkSrv("tcp", "8080")
	fmt.Println(err)
	srv.Init()
	for {
		time.Sleep(5 * time.Millisecond)
	}
}
