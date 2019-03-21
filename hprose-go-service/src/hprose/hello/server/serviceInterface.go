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
	return "hprose.hello.server.go.Hello"
}


type Hello2 struct {
	SayHello func(string) string
}

func (*Hello2) ServiceName() string {
	return "hprose.hello.server.go.Hello2"
}

