package yuma

import (
)

type Node struct {
	subprocess string
	nodes []*Node 
}

func NewNode(subprocess string) *Node {
	node := new(Node)
	node.subprocess = subprocess

	return node
}

func (node *Node) GetSubprocess() string {
	return node.subprocess
}

func (node *Node) GetNodes() []*Node {
	return node.nodes
}

func (node *Node) AddNode(child *Node) {
	node.nodes = append(node.nodes, child)
}

func DetermineServiceTree(node *Node, target string, solutions map[string][]string) *Node {
	for _, action := range solutions[target] {
		if (action != target) && (len(solutions[action]) > 1) {
			branch := NewNode(action)
			if checkValidChild(target, branch, solutions[target], solutions) {
				node.AddNode(branch)
				DetermineServiceTree(branch, action, solutions)
			}
		} else if (action != target) {
			child := NewNode(action)
			if checkValidChild(target, child, solutions[target], solutions) {
				node.AddNode(child)
			}
		}
	}

	return node
}

func PrintServiceTree(node *Node, space string, servicetree string) string {
	if space == "" {
		servicetree += space + "──" + node.GetSubprocess() + "\n"
	} else {
		servicetree += space + "└──" + node.GetSubprocess() + "\n"
	}
	for _, child := range node.GetNodes() {
		servicetree = PrintServiceTree(child, space + " ", servicetree)
	}

	return servicetree
}

/* func addNodes(node *Node, target string, solutions map[string][]string ) {
	for _, action := range solutions[target] {
		if (action != target) && (len(solutions[action]) > 1) {
			branch := NewNode(action)
			if checkValidChild(target, branch, solutions[target], solutions) {
				node.AddNode(branch)
				addNodes(branch, action, solutions)
			}
		} else if (action != target) {
			child := NewNode(action)
			if checkValidChild(target, child, solutions[target], solutions) {
				node.AddNode(child)
			}
		}
	}
} */

func checkValidChild(target string, node *Node, solution []string, solutions map[string][]string) bool {
	for _, subprocess := range solution {
		if subprocess == target || subprocess == node.GetSubprocess() {
			continue
		} else {
			for _, s := range solutions[subprocess] {
				if s == node.GetSubprocess() {
					return false
				}
			}
		}
	}
	
	return true
}