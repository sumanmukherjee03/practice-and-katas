package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
)

func main() {
	path, err := filepath.Abs("./wind-direction/metarfiles")
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		data, err := ioutil.ReadFile(filepath.Join(path, f.Name()))
		if err != nil {
			log.Fatal(err)
		}
		text := string(data)
		text += ";"
		// fmt.Println(text)
	}
}
