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
	Last   *ListItem
	First  *ListItem
	Length int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.Length
}

func (l *list) Front() *ListItem {
	return l.First
}

func (l *list) Back() *ListItem {
	return l.Last
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := ListItem{
		Next:  l.First,
		Value: v,
	}
	if l.Last == nil {
		l.Last = &item
	}
	if l.First != nil {
		l.First.Prev = &item
	}
	l.First = &item
	l.Length++
	return &item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := ListItem{
		Prev:  l.Last,
		Value: v,
	}
	if l.Last != nil {
		l.Last.Next = &item
	}
	if l.First == nil {
		l.First = &item
	}
	l.Last = &item
	l.Length++
	return &item
}

func (l *list) Remove(t *ListItem) {
	if l.First != t {
		t.Prev.Next = t.Next
	} else {
		l.First = t.Next
	}
	if l.Last != t {
		t.Next.Prev = t.Prev
	} else {
		l.Last = t.Prev
	}
	t.Next = nil
	t.Prev = nil
	l.Length--
}

func (l *list) MoveToFront(t *ListItem) {
	l.Remove(t)
	l.PushFront(t.Value)
}
