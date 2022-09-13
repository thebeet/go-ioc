package ioc

var (
	global = New()
)

func Register(resolver any) {
	global.Register(resolver)
}

func Call(f any) error {
	return global.Call(f)
}
