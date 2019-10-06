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

func (t *Tree) Inorder() []int {
	traversal := []int{}
	if t == nil {
		return traversal
	}
	if t.Left != nil {
		for _, v := range t.Left.Inorder() {
			traversal = append(traversal, v)
		}
	}
	traversal = append(traversal, t.Value)
	if t.Right != nil {
		for _, v := range t.Right.Inorder() {
			traversal = append(traversal, v)
		}
	}
	return traversal
}

func (t *Tree) PreOrder() []int {
	traversal := []int{}
	if t == nil {
		return traversal
	}
	traversal = append(traversal, t.Value)
	if t.Left != nil {
		for _, v := range t.Left.PreOrder() {
			traversal = append(traversal, v)
		}
	}
	if t.Right != nil {
		for _, v := range t.Right.PreOrder() {
			traversal = append(traversal, v)
		}
	}
	return traversal
}

func (t *Tree) PostOrder() []int {
	traversal := []int{}
	if t == nil {
		return traversal
	}
	if t.Left != nil {
		for _, v := range t.Left.PostOrder() {
			traversal = append(traversal, v)
		}
	}
	if t.Right != nil {
		for _, v := range t.Right.PostOrder() {
			traversal = append(traversal, v)
		}
	}
	traversal = append(traversal, t.Value)
	return traversal
}

func main() {
	numOfNodes, err := strconv.ParseInt(os.Args[1], 0, 0)
	if err != nil {
		panic(err)
	}
	t := New(int(numOfNodes))
	fmt.Println(t.Inorder())
	fmt.Println(t.PreOrder())
	fmt.Println(t.PostOrder())
}
