package main

import (
"os"
"runtime"
	"fmt"
	"path"
)

func main() {
	_, filename, _, _ := runtime.Caller(1)
	f, err := os.Open(path.Join(path.Dir(filename), "test.txt"))

	fmt.Println(f)
}
