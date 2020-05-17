package datastructures

import (
	"fmt"
)

type BinarySearchTree struct {
	Nodes []*Node
}

type Node struct {
	Key int
	Data interface{}
	Left *Node
	Right *Node
	Parent *Node
}

func NewBinarySearchTree() *BinarySearchTree {
	bst := new(BinarySearchTree)
	bst.Nodes = make([]*Node, 0)

	return bst
}

func (bst *BinarySearchTree) Print(node *Node) {
	for node != nil {
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