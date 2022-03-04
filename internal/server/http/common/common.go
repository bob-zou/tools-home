package common

type ResponseCode int

const (
	ReplyCodeOK = iota
	ReplyCodeErr
)

type Reply struct {
	Code ResponseCode `json:"code"`
	Msg  string       `json:"msg"`
	Data interface{}  `json:"data"`
}
