package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

// Tree data type
type Tree struct {
	Left  *Tree
	Value int
	Right *Tree
}

// New : Constructor for Tree
func New(numOfNodes int) *Tree {
	var t *Tree
	for _, val := range rand.Perm(numOfNodes) {
		t = t.Insert(val + rand.Intn(100))
	}
	return t
}

// Insert : Insert nodes into the Tree
func (t *Tree) Insert(v int) *Tree {
	if t == nil {
		return &Tree{nil, v, nil}
	} else if v < t.Value {
		t.Left = t.Left.Insert(v)
	} else if v > t.Value {
		t.Right = t.Right.Insert(v)
	}
	return t
}

// Search : Search for an element in the Tree
func (t *Tree) Search(v int) bool {
	if t == nil {
		return false
	}
	if t.Value == v {
		return true
	}
	if t.Left != nil && v < t.Value {
		return t.Left.Search(v)
	}
	if t.Right != nil && v > t.Value {
		return t.Right.Search(v)
	}
	return false
}

// Inorder : Inorder traversal of the Tree
func (t *Tree) Inorder() []int {
	traversal := []int{}
	if t == nil {
		return traversal
	}
	if t.Left != nil {
		traversal = append(traversal, t.Left.Inorder()...)
	}
	traversal = append(traversal, t.Value)
	if t.Right != nil {
		traversal = append(traversal, t.Right.Inorder()...)
	}
	return traversal
}

// PreOrder : Preorder traversal of the Tree
func (t *Tree) PreOrder() []int {
	traversal := []int{}
	if t == nil {
		return traversal
	}
	traversal = append(traversal, t.Value)
	if t.Left != nil {
		traversal = append(traversal, t.Left.PreOrder()...)
	}
	if t.Right != nil {
		traversal = append(traversal, t.Right.PreOrder()...)
	}
	return traversal
}

// PostOrder : Postorder traversal of Tree
func (t *Tree) PostOrder() []int {
	traversal := []int{}
	if t == nil {
		return traversal
	}
	if t.Left != nil {
		traversal = append(traversal, t.Left.PostOrder()...)
	}
	if t.Right != nil {
		traversal = append(traversal, t.Right.PostOrder()...)
	}
	traversal = append(traversal, t.Value)
	return traversal
}

// Min : Finds the minimum value in the Tree
func (t *Tree) Min() int {
	if t == nil {
		return 0
	}
	if t.Left != nil {
		return t.Left.Min()
	}
	return t.Value
}

// Max : Finds the maximum value in the Tree
func (t *Tree) Max() int {
	if t == nil {
		return 0
	}
	if t.Right != nil {
		return t.Right.Max()
	}
	return t.Value
}

func main() {
	numOfNodes, err := strconv.ParseInt(os.Args[1], 0, 0)
	if err != nil {
		panic(err)
	}
	t := New(int(numOfNodes))
	t = t.Insert(21)
	fmt.Println(t.Inorder())
	fmt.Println(t.PreOrder())
	fmt.Println(t.PostOrder())
	fmt.Println("Searching for 21 : ", t.Search(21))
	fmt.Println("Searching for 5 : ", t.Search(5))
	fmt.Println("Minimum element is : ", t.Min())
	fmt.Println("Maximum element is : ", t.Max())
}
