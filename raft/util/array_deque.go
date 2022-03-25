package util

type ArrayDeque[E any] struct {
	list []E
}

// PeekFirst Get the first element of list.
func (d *ArrayDeque[E]) PeekFirst() E {
	return ListPeekFirst(d.list)
}

// PeekLast Get the last element of list.
func (d *ArrayDeque[E]) PeekLast() E {
	return ListPeekLast(d.list)
}

// PollFirst Remove the first element from list and return it.
func (d *ArrayDeque[E]) PollFirst() E {
	return ListPollFirst(d.list)
}

// PollLast Remove the last element from list and return it.
func (d *ArrayDeque[E]) PollLast() E {
	return ListPollLast(d.list)
}

// ListPeekFirst Get the first element of list.
func ListPeekFirst[E any](list []E) E {
	return list[0]
}

// ListPollFirst Remove the first element from list and return it.
func ListPollFirst[E any](list []E) E {
	e := ListPeekFirst(list)
	list = list[1:]
	return e
}

// ListPeekLast Get the last element of list.
func ListPeekLast[E any](list []E) E {
	return list[len(list)-1]
}

// ListPollLast Remove the first element from list and return it.
func ListPollLast[E any](list []E) E {
	e := ListPeekLast(list)
	list = list[:len(list)-1]
	return e
}
