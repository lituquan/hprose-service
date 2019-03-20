package server

type IEcho struct {
	SayHello func(string) string
}
func (*IEcho)ServiceName() string{
	return "hprose.hello.server.IEcho"
}

type Hello struct {
	SayHello func(string) string
}

func (*Hello) ServiceName() string {
	return "hprose.hello.server.go.IEcho"
}

