package entity

const (
	EVT_CODE_MIN = 656
	EVT_CODE_MAX = 685
)

type EVFileType int

const (
	EV_FILE_TYPE_SEARCH EVFileType = iota
	EV_FILE_TYPE_EXEC
)
