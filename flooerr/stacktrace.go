package flooerr

import "fmt"

type stacktrace struct {
	Function string
	File     string
	Line     int
}

func (s *stacktrace) String() string {
	return fmt.Sprintf("%s:%s:%d", s.Function, s.File, s.Line)
}
