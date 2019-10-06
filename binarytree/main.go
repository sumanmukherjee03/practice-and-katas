package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sumanmukherjee03/practice-and-katas/binarytree/tree"
)

func main() {
	numOfNodes, err := strconv.ParseInt(os.Args[1], 0, 0)
	if err != nil {
		panic(err)
	}
	t := tree.New(int(numOfNodes))
	t = t.Insert(21)
	fmt.Println(t.Inorder())
	fmt.Println(t.PreOrder())
	fmt.Println(t.PostOrder())
	fmt.Println("Searching for 21 : ", t.Search(21))
	fmt.Println("Searching for 5 : ", t.Search(5))
	fmt.Println("Minimum element is : ", t.Min())
	fmt.Println("Maximum element is : ", t.Max())
	fmt.Println(t.String())
}
