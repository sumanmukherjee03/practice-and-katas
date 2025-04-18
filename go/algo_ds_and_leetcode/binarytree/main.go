package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/sumanmukherjee03/practice-and-katas/go/binarytree/tree"
)

func main() {
	numOfNodes, err := strconv.ParseInt(os.Args[1], 0, 0)
	if err != nil {
		panic(err)
	}
	arr := tree.GenDataForTree(int(numOfNodes))
	rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	fmt.Println("Original array of numbers inserting into the tree : ", arr)
	t := tree.New(arr)
	t = t.Insert(21)
	fmt.Println("Inorder traversal : ", t.Inorder())
	fmt.Println("Preorder traversal : ", t.PreOrder())
	fmt.Println("Postorder traversal : ", t.PostOrder())
	fmt.Println("Searching for 21 : ", t.Search(21))
	fmt.Println("Searching for 5 : ", t.Search(5))
	fmt.Println("Minimum element is : ", t.Min())
	fmt.Println("Maximum element is : ", t.Max())
	fmt.Println("Maximum depth is : ", t.MaxDepth())
	fmt.Println("Paths to leaf nodes : ", t.RootToLeafPaths())
	fmt.Println(">>>>>>>> Comparing binary trees")
	fmt.Println("Tree 1 : ", t.String())
	t1 := tree.New(arr)
	t1 = t1.Insert(21)
	rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	t2 := tree.New(arr)
	t2 = t2.Insert(21)
	fmt.Println("Tree 2 : ", t1.String())
	fmt.Println("Tree 3 : ", t2.String())
	fmt.Println(">>>>>>>> Comparing if tree 1 and 2 are equal :", t.Equal(t1))
	fmt.Println(">>>>>>>> Comparing if tree 2 and 3 are equal :", t1.Equal(t2))
	fmt.Println(">>>>>>>> Comparing if tree 2 and 3 are same :", t1.Same(t2))
	fmt.Println(">>>>>>>> Mirrored binary trees")
	fmt.Println("Tree 1 : ", t.String())
	fmt.Println("Tree 2 : ", t.Mirror().String())
}
