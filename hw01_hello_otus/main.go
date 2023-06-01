package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	out := stringutil.Reverse("Hello, OTUS!")
	fmt.Println(out)
}
