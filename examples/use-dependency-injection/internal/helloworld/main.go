package helloworld

import (
	"fmt"

	"go.uber.org/dig"
)

// Main function to run hello-world
func Main(text string) {
	fmt.Println(text)
}

type greeter struct {
	dig.In
	Text string `name:"typical"`
}

// Main2 function to run hello-world with name-constructor
func Main2(greeter greeter) {
	fmt.Println(greeter.Text)
}

// HelloWorld text
// @ctor
func HelloWorld() string {
	return "Hello World"
}

// HelloTypical text
// @ctor (name:"typical")
func HelloTypical() string {
	return "Hello Typical"
}

// Close the application
// @dtor
func Close() {
	fmt.Println("close the app")
}