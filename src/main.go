package main

import "fmt"

func main() {
	done := make(chan struct{})
	go func() {
		defer close(done)
		fmt.Println("hi")
		fmt.Println("bye")
	}()
	<-done
}
