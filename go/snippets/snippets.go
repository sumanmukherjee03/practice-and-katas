package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	fmt.Println(base64.StdEncoding.EncodeToString([]byte("username:password")))
}
