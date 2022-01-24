package utils

type State struct {
	Code    int32
	Message string
}

type Status struct {
	state *State
}

func (s *Status) Code() int32 {
	if s.state == nil {
		return 0
	}
	return s.state.Code
}

func (s *Status) Error() string {
	if s.state == nil {
		return "OK"
	}
	return s.state.Message
}

func (s *Status) SetError(code int32, msg string) {
	s.state.Code = code
	s.state.Message = msg
}

func (s *Status) OK() bool {
	return s.state == nil
}
