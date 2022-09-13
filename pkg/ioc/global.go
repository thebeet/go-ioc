package ioc

var (
	global = New()
)

func Register(resolver any) {
	global.Register(resolver)
}

func RegisterWithName(resolver any, name string) {
	global.RegisterWithName(resolver, name)
}

func RegisterInstance(instance any) {
	global.RegisterInstance(instance)
}

func RegisterInstanceWithName(instance any, name string) {
	global.RegisterInstanceWithName(instance, name)
}

func Call(f any) error {
	return global.Call(f)
}

func Fill(f any) error {
	return global.Fill(f)
}
