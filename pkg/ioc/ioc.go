package ioc

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type definition struct {
	name      string
	kind      reflect.Type
	realKind  reflect.Type
	construct any
	instance  *reflect.Value
	opts      definitionOpt
	mutex     sync.Mutex
}

type definitionOpt struct {
	lazy  bool
	scope string
}

type definitionGroup struct {
	group  map[string]*definition
	primer *definition
}

type Ioc interface {
	Register(resolver any)
	RegisterWithName(resolver any, name string)

	RegisterInstance(instance any)
	RegisterInstanceWithName(instance any, name string)
	Call(f any) error
	Fill(f any) error
}

func New(opts ...Opts) Ioc {
	if len(opts) == 0 {
		return &ioc{}
	}
	return &ioc{opts: opts[0]}
}

type ioc struct {
	container map[reflect.Type]*definitionGroup
	opts      Opts
	mutex     sync.Mutex
}

type Opts struct {
	EnablePostConstruct bool
}

func (i *ioc) Register(resolver any) {
	i.RegisterWithName(resolver, "")
}

func (i *ioc) RegisterWithName(resolver any, name string) {
	i.bind(resolver, name)
}

func (i *ioc) RegisterInstance(instance any) {
	i.RegisterInstanceWithName(instance, "")
}

func (i *ioc) RegisterInstanceWithName(instance any, name string) {
	i.bindInstance(instance, name)
}

func (i *ioc) Call(f any) error {
	t := reflect.ValueOf(f)
	if t.Kind() != reflect.Func {
		return errors.New("func only")
	}
	arguments, err := i.arguments(f)
	if err != nil {
		return err
	}
	_ = t.Call(arguments)
	return nil
}

func (i *ioc) Fill(f any) error {
	if reflect.TypeOf(f).Kind() == reflect.Pointer {
		t := reflect.TypeOf(f).Elem()
		v := reflect.ValueOf(f).Elem()
		for j := 0; j < t.NumField(); j++ {
			field := t.Field(j)
			if value, ok := field.Tag.Lookup("autowire"); ok {
				valueField := v.Field(j)
				valueField.Set(i.resolve(valueField.Type(), value))
			}
		}
	}
	return nil
}

func (i *ioc) fill(v reflect.Value) error {

	return nil
}

func (i *ioc) bind(resolver any, name string, opts ...definitionOpt) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if i.container == nil {
		i.container = make(map[reflect.Type]*definitionGroup, 0)
	}

	t := reflect.TypeOf(resolver)
	if t.Kind() != reflect.Func {
		panic("resolver should be a construct func")
	}
	if t.NumOut() != 1 {
		panic("resolver should be a construct func return only one")
	}
	outType := t.Out(0)
	if outType.Kind() == reflect.Pointer {
		outType = outType.Elem()
	}
	if _, exist := i.container[outType]; !exist {
		i.container[outType] = &definitionGroup{group: make(map[string]*definition), primer: nil}
	}

	i.container[outType].group[name] = &definition{
		name:      name,
		kind:      outType,
		realKind:  t.Out(0),
		construct: resolver,
		instance:  nil,
	}
	if len(i.container[outType].group) == 1 {
		i.container[outType].primer = i.container[outType].group[name]
	}
}

func (i *ioc) bindInstance(instance any, name string, opts ...definitionOpt) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if i.container == nil {
		i.container = make(map[reflect.Type]*definitionGroup, 0)
	}

	t := reflect.TypeOf(instance)
	v := reflect.ValueOf(instance)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if _, exist := i.container[t]; !exist {
		i.container[t] = &definitionGroup{group: make(map[string]*definition), primer: nil}
	}

	i.container[t].group[name] = &definition{
		name:      name,
		kind:      t,
		realKind:  reflect.TypeOf(instance),
		construct: nil,
		instance:  &v,
	}
	if len(i.container[t].group) == 1 {
		i.container[t].primer = i.container[t].group[name]
	}
}

func (i *ioc) resolve(r reflect.Type, name string) reflect.Value {
	tr := r
	if r.Kind() == reflect.Pointer {
		tr = r.Elem()
	}

	if _, exist := i.container[tr]; !exist {
		panic(fmt.Sprintf("type [%+v] not exits", r))
	}
	var d *definition
	if len(i.container[tr].group) == 1 {
		d = i.container[tr].primer
	} else {
		var find bool
		d, find = i.container[tr].group[name]
		if !find {
			panic(fmt.Sprintf("type [%+v<%s>] not exits", r, name))
		}
	}
	instance := i.instance(d)

	if r.Kind() == reflect.Pointer && d.realKind.Kind() == reflect.Struct {
		if instance.CanAddr() {
			return instance.Addr()
		}
	} else if r.Kind() == reflect.Struct && d.realKind.Kind() == reflect.Pointer {
		return instance.Elem()
	}
	return instance
}

func (i *ioc) instance(d *definition) reflect.Value {
	if d.instance == nil {
		d.mutex.Lock()
		defer d.mutex.Unlock()

		in, _ := i.arguments(d.construct)
		f := reflect.ValueOf(d.construct)
		out := f.Call(in)
		beanReflect := out[0]
		bean := beanReflect.Interface()
		d.instance = &beanReflect

		if beanReflect.Kind() == reflect.Pointer && beanReflect.Elem().Kind() == reflect.Struct {
			for j := 0; j < beanReflect.Elem().NumField(); j++ {
				t := d.kind.Field(j)
				if value, ok := t.Tag.Lookup("autowire"); ok {
					f := beanReflect.Elem().Field(j)
					// log.Printf("autowire: %+v<name:%s> for field %s", f.Type(), value, t.Name)
					f.Set(i.resolve(f.Type(), value))
				}
			}
		}

		if t, ok := bean.(IocPostConstruct); ok {
			t.IocPostConstruct()
		}

		if t, ok := bean.(IocContainerAware); ok {
			t.SetIocContainer(i)
		}

		if t, ok := bean.(IocInstanceNameAware); ok {
			t.SetIocInstanceName(d.name)
		}
	}
	return *d.instance
}

func (i *ioc) arguments(f any) ([]reflect.Value, error) {
	t := reflect.TypeOf(f)
	in := []reflect.Value{}
	for j := 0; j < t.NumIn(); j++ {
		in = append(in, i.resolve(t.In(j), ""))
	}
	return in, nil
}
