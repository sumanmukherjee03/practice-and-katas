package tree

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Tree data type
type Tree struct {
	Left  *Tree
	Value int
	Right *Tree
}

// GenDataForTree : Generates an array of random numbers to insert in a tree
func GenDataForTree(size int) []int {
	arr := []int{}
	rand.Seed(time.Now().UnixNano())
	for _, val := range rand.Perm(size) {
		arr = append(arr, val+rand.Intn(1000))
	}
	return arr
}

// New : Constructor for Tree
func New(arr []int) *Tree {
	var t *Tree
	for _, val := range arr {
		t = t.Insert(val)
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

// String : String representation of the Tree
func (t *Tree) String() string {
	if t == nil {
		return "()"
	}
	str := ""
	str += fmt.Sprint(t.Value) + " "
	if t.Left != nil {
		str += t.Left.String() + " "
	} else {
		str += "()" + " "
	}
	if t.Right != nil {
		str += t.Right.String() + " "
	} else {
		str += "()" + " "
	}
	return "(" + str + ")"
}

// Same : compares 2 binary trees
func (t *Tree) Same(x *Tree) bool {
	if t == nil && x == nil {
		return true
	}
	if (t == nil && x != nil) || (t != nil && x == nil) {
		return false
	}
	if t != nil && x != nil {
		if t.Value == x.Value {
			return t.Left.Same(x.Left) && t.Right.Same(x.Right)
		}
	}
	return false
}

// MaxDepth : Finds the maximum depth of the tree
func (t *Tree) MaxDepth() int {
	leftDepth := 0
	rightDepth := 0
	if t == nil {
		return 0
	}
	if t.Left != nil {
		leftDepth += 1 + t.Left.MaxDepth()
	}
	if t.Right != nil {
		rightDepth += 1 + t.Right.MaxDepth()
	}
	return int(math.Max(float64(leftDepth), float64(rightDepth)))
}

// Mirror : Mirrors a given tree. Resulting tree wont be a bst any more
func (t *Tree) Mirror() *Tree {
	if t == nil {
		return nil
	}
	x := &Tree{nil, t.Value, nil}
	mirroredLeft := t.Left.Mirror()
	mirroredRight := t.Right.Mirror()
	x.Left = mirroredRight
	x.Right = mirroredLeft
	return x
}

// RootToLeafPaths : Returns an array of root -> node -> leaf paths for the Tree
func (t *Tree) RootToLeafPaths() [][]int {
	paths := [][]int{}
	path := []int{}
	if t == nil {
		return paths
	}

	var dfs func(*Tree, int)
	dfs = func(n *Tree, level int) {
		if level >= len(path) {
			path = append(path, n.Value)
		} else {
			path[level] = n.Value
		}
		if n.Left == nil && n.Right == nil {
			temp := make([]int, len(path))
			copy(temp, path)
			paths = append(paths, temp)
		}
		if n.Left != nil {
			dfs(n.Left, level+1)
		}
		if n.Right != nil {
			dfs(n.Right, level+1)
		}
	}

	dfs(t, 0)
	return paths
}
