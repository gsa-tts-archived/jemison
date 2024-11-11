package common

type ErrorType int

const (
	NonIndexableContentType ErrorType = iota
)

func (et ErrorType) String() string {
	return [...]string{"NONDXTYPE"}[et]
}
