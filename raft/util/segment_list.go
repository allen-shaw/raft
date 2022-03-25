package util

const (
	SegmentShift = 7
	SegmentSize  = 2 << (SegmentShift - 1)
)

// SegmentList A list implementation based on segments. Only supports removing elements from start or end.
// The list keep the elements in a segment list, every segment contains at most 128 elements.
//
//               [segment, segment, segment ...]
//            /                 |                 \
//        segment             segment              segment
//     [0, 1 ... 127]    [128, 129 ... 255]    [256, 257 ... 383]
type SegmentList[T any] struct {
	segments ArrayDeque[T]
}

type Segment[T any] struct {
	elements []T
	pos      int // end offset(exclusive)
	offset   int // start offset(inclusive)
}

func NewSegment[T any]() *Segment[T] {
	s := &Segment[T]{}
	s.elements = make([]T, 0, SegmentSize)
	s.pos = 0
	s.offset = 0
	return s
}

func (s *Segment[T]) Clear() {
	s.pos = 0
	s.offset = 0
	s.elements = make([]T, 0, SegmentSize)
}

func (s *Segment[T]) Cap() int {
	return SegmentSize - s.pos
}

func (s *Segment[T]) IsReachEnd() bool {
	return s.pos == SegmentSize
}

func (s *Segment[T]) IsEmpty() bool {
	return s.pos == 0
}
