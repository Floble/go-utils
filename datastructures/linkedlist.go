package datastructures

import (
	"fmt"
)

type LinkedList struct {
	Head *ListElement
	Tail *ListElement
}

type ListElement struct {
	Prev *ListElement
	Next *ListElement
	Key int
}

func NewLinkedList() *LinkedList {
	list := new(LinkedList)
	list.Head = nil
	list.Tail = nil

	return list
}

func NewListElement(key int) *ListElement {
	e := new(ListElement)
	e.Key = key

	return e
}

func (list *LinkedList) Search(key int) *ListElement {
	e := list.Head

	for e != nil {
		if e.Key == key {
			return e
		}
		e = e.Next
	}

	return nil
}

func (list *LinkedList) Insert(key int) {
	e := NewListElement(key)
	e.Next = list.Head
	if list.Head != nil {
		list.Head.Prev = e
	}
	e.Prev = nil
	list.Head = e
}

func (list *LinkedList) DeleteElement(e *ListElement) {
	if e.Prev != nil {
		e.Prev.Next = e.Next
	} else {
		list.Head = e.Next
	}
	
	if e.Next != nil {
		e.Next.Prev = e.Prev
	} else {
		list.Tail = e.Prev
	}

	e.Next = nil
	e.Prev = nil
}

func (list *LinkedList) Delete(key int) {
	e := list.Search(key)
	list.DeleteElement(e)
}

func (list *LinkedList) Print() {
	e := list.Head

	for e != nil {
		fmt.Println(e.Key)
		e = e.Next
	}
}