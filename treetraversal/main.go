package main

import (
	"fmt"
	"math/rand"
)

// Tree data type
type Tree struct {
	Left  *Tree
	Value int
	Right *Tree
}

// New : Constructor for Tree
func New(numOfLeaves int) *Tree {
	var t *Tree
	for _, val := range rand.Perm(numOfLeaves) {
		t = insert(t, val)
	}
	return t
}

func insert(t *Tree, v int) *Tree {
	if t == nil {
		return &Tree{nil, v, nil}
	} else if v < t.Value {
		t.Left = insert(t.Left, v)
	} else if v > t.Value {
		t.Right = insert(t.Right, v)
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
	t := New(5)
	fmt.Println(t.Inorder())
	fmt.Println(t.PreOrder())
	fmt.Println(t.PostOrder())
}
