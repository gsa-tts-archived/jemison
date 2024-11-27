package common

type ErrorType int

const (
	NonIndexableContentType ErrorType = iota
	FileTooLargeToFetch
	FileTooSmallToProcess
)

func (et ErrorType) String() string {
	return [...]string{"NONDXTYPE", "TOOLARGETOFETCH", "TOOSMALLTOPROCESS"}[et]
}
