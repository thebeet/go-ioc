package ioc

type IocPostConstruct interface {
	IocPostConstruct()
}

type IocInstanceNameAware interface {
	SetIocInstanceName(name string)
}

type IocContainerAware interface {
	SetIocContainer(container Ioc)
}
