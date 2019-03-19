package server

type IEcho struct {
	SayHello func(string) string
}
func (*IEcho)ServiceName() string{
	return "hprose.hello.server.IEcho"
}
