// +build !api

package main

import "fmt"

func listen(context map[string]interface{}) {
	fmt.Println("To use this feature you must compile the builder with the -api flag")
}
