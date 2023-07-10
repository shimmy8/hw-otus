package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	Length    int
	FirstItem *ListItem
	LastItem  *ListItem
}

func (lst list) Len() int {
	return lst.Length
}

func (lst list) Front() *ListItem {
	return lst.FirstItem
}

func (lst list) Back() *ListItem {
	return lst.LastItem
}

func (lst *list) PushFront(v interface{}) *ListItem {
	listItem := ListItem{Value: v}
	if lst.FirstItem == nil {
		lst.FirstItem = &listItem
		lst.LastItem = &listItem
	} else {
		listItem.Next = lst.FirstItem
		lst.FirstItem.Prev = &listItem
		lst.FirstItem = &listItem
	}
	lst.Length++
	return &listItem
}

func (lst *list) PushBack(v interface{}) *ListItem {
	listItem := ListItem{Value: v}
	if lst.LastItem == nil { // lst is empty
		lst.LastItem = &listItem
		lst.FirstItem = &listItem
	} else {
		listItem.Prev = lst.LastItem
		lst.LastItem.Next = &listItem
		lst.LastItem = &listItem
	}
	lst.Length++
	return &listItem
}

func (lst *list) Remove(i *ListItem) {
	if i.Prev == nil {
		lst.FirstItem = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		lst.LastItem = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	lst.Length--
	i = nil
}

func (lst *list) MoveToFront(i *ListItem) {
	if lst.FirstItem == i {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	lst.FirstItem.Prev = i
	i.Next = lst.FirstItem
	i.Prev = nil
	lst.FirstItem = i
}

func NewList() List {
	return new(list)
}
