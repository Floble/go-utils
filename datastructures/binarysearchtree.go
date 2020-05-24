package datastructures

import (
	"fmt"
)

type BinarySearchTree struct {
	Root *Node
}

type Node struct {
	Key int
	Data interface{}
	Left *Node
	Right *Node
	Parent *Node
}

func NewNode(key int, data interface{}) *Node {
	node := new(Node)
	node.Key = key
	node.Data = data

	return node
}

func NewBinarySearchTree() *BinarySearchTree {
	bst := new(BinarySearchTree)
	bst.Root = nil

	return bst
}

func (bst *BinarySearchTree) Print(node *Node) {
	if node != nil {
		bst.Print(node.Left)
		fmt.Println(node.Key)
		bst.Print(node.Right)
	}
}

func (bst *BinarySearchTree) Search(node *Node, key int) *Node {
	for node != nil && node.Key != key {
		if node.Key < key {
			node = node.Left
		} else {
			node = node.Right
		}
	}

	return node
}

func (bst *BinarySearchTree) Minimum(node *Node) *Node {
	for node.Left != nil {
		node = node.Left
	}

	return node
}

func (bst *BinarySearchTree) Maximum(node *Node) *Node {
	for node.Right != nil {
		node = node.Right
	}

	return node
}

func (bst *BinarySearchTree) Successor(node *Node) *Node {
	if node.Right != nil {
		return bst.Minimum(node.Right)
	}

	parent := node.Parent

	for parent != nil && node == parent.Right {
		node = parent
		parent = parent.Parent
	}

	return parent
}

func (bst *BinarySearchTree) Predecessor(node *Node) *Node {
	if node.Left != nil {
		return bst.Maximum(node.Left)
	}

	parent := node.Parent

	for parent != nil && node == parent.Left {
		node = parent
		parent = parent.Parent
	}

	return parent
}

func (bst *BinarySearchTree) Insert(node *Node) {
	var y *Node
	x := bst.Root

	for x != nil {
		y = x
		if node.Key < x.Key {
			x = y.Left
		} else {
			x = y.Right
		}
	}

	if y == nil {
		bst.Root = node
	} else if node.Key < y.Key {
		y.Left = node
	} else {
		y.Right = node
	}
}

func (bst *BinarySearchTree) transplant(u, v *Node) {
	if u.Parent == nil {
		bst.Root = v
	} else if u == u.Parent.Left {
		u.Parent.Left = v
	} else {
		u.Parent.Right = v
	}

	if v != nil {
		v.Parent = u.Parent
	}
}

func (bst *BinarySearchTree) Delete(node *Node) {
	if node.Left == nil {
		bst.transplant(node, node.Right)
	} else if node.Right == nil {
		bst.transplant(node, node.Left)
	} else {
		y := bst.Successor(node)

		if y != node.Right {
			bst.transplant(y, y.Right)
			y.Right = node.Right
			y.Right.Parent = y
		}
		
		bst.transplant(node, y)
		y.Left = node.Left
		y.Left.Parent = y
	}
}