package hw04_lru_cache // nolint:golint,stylecheck

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(elem *ListItem)
	MoveToFront(elem *ListItem)
	insert(elem, at *ListItem) *ListItem
	move(elem, at *ListItem) *ListItem
	initElem(v interface{}) *ListItem
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len  int
	root ListItem
}

// lazyInit lazily initializes a zero List value.

func (l *list) lazyInit() {
	if l.root.Next == nil {
		l.Init()
	}
}

// Init initializes or clears list l.
func (l *list) Init() *list {
	l.root.Next = &l.root
	l.root.Prev = &l.root
	l.len = 0

	return l
}

func (l *list) insertValue(v interface{}, at *ListItem) *ListItem {
	return l.insert(&ListItem{Value: v}, at)
}

func (l *list) insert(elem, at *ListItem) *ListItem {
	elem.Prev = at
	elem.Next = at.Next
	elem.Prev.Next = elem
	elem.Next.Prev = elem

	l.len++
	return elem
}

func (l *list) initElem(v interface{}) *ListItem {
	return &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}
}

func (l *list) move(elem, at *ListItem) *ListItem {
	if elem == at {
		return elem
	}
	elem.Prev.Next = elem.Next
	elem.Next.Prev = elem.Prev
	elem.Prev = at
	elem.Next = at.Next
	elem.Prev.Next = elem
	elem.Next.Prev = elem

	return elem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	if l.len == 0 {
		return nil
	}

	return l.root.Next
}

func (l *list) Back() *ListItem {
	if l.len == 0 {
		return nil
	}

	return l.root.Prev
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.lazyInit()
	return l.insertValue(v, l.root.Prev)
}

func (l *list) Remove(elem *ListItem) {
	elem.Prev.Next = elem.Next
	elem.Next.Prev = elem.Prev
	elem.Next = nil
	elem.Prev = nil

	l.len--
}

func (l *list) MoveToFront(elem *ListItem) {
	if l.root.Next == elem {
		return
	}

	l.move(elem, &l.root)
}

func NewList() List {
	return &list{}
}
