/*
  trader API Engine
*/

package config

import (
	"fmt"
	"os"
)

// Envs
const (
	Dev  string = "development"
	Prod string = "production"
	Test string = "test"
)

// Env is the environment that Martini is executing in. The MARTINI_ENV is read on initialization to set this variable.
var Env = Dev
var Root string

func setENV(e string) {
	if len(e) > 0 {
		Env = e
	}
}

func init() {
	var err error
	Root, err = os.Getwd()
	fmt.Println("Root:", Root)
	if err != nil {
		panic(err)
	}
}
