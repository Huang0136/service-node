package service

type Service struct {
	ServiceId   string
	ServiceName string
	URL         string
	Method      string
	Remark      string
	InParams    []InParam
	OutParams   []OutParam
}

// 入参结构体
type InParam struct {
	ParamCode string
	ParamName string
	ParamType int
	Length    int
	Remark    string
}

// 出参结构体
type OutParam struct {
	ParamCode string
	ParamName string
	ParamType int
	Remark    string
}
