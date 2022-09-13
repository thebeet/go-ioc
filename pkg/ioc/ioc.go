package ioc

import (
	"errors"
	"log"
	"reflect"
	"sync"
)

type definition struct {
	name      string
	kind      reflect.Type
	construct any
	instance  *reflect.Value
	opts      struct {
		lazy  bool
		scope string
	}
}

type Ioc interface {
	Register(resolver any)
	Call(f any) error
}

func New() Ioc {
	return &ioc{}
}

type ioc struct {
	container map[reflect.Type]map[string]*definition
	mutex     sync.Mutex
}

type opts struct {
}

func (i *ioc) Register(resolver any) {
	i.bind(resolver, "")
}

func (i *ioc) bind(resolver any, name string) {
	if i.container == nil {
		i.container = make(map[reflect.Type]map[string]*definition, 0)
	}

	t := reflect.TypeOf(resolver)
	if t.Kind() != reflect.Func {
		panic("")
	}
	if t.NumOut() == 1 {
		if _, exist := i.container[t.Out(0)]; !exist {
			i.container[t.Out(0)] = make(map[string]*definition)
		}
	} else {
		panic("resolver should return only one ")
	}
	outType := t.Out(0)
	i.container[outType][name] = &definition{
		name:      name,
		kind:      outType,
		construct: resolver,
		instance:  nil,
	}
}

func (i *ioc) resolve(r reflect.Type, name string) reflect.Value {
	d := i.container[r][name]
	return i.instance(d)
}

func (i *ioc) instance(d *definition) reflect.Value {
	if d.instance == nil {
		//i.mutex.Lock()
		//defer i.mutex.Unlock()
		log.Printf("init %+v", *d)

		in, _ := i.arguments(d.construct)
		f := reflect.ValueOf(d.construct)
		out := f.Call(in)
		bean := out[0].Interface()
		d.instance = &out[0]

		i, ok := bean.(PostConstruct)
		if ok {
			i.PostConstruct()
		}
	}
	return *d.instance
}

func (i *ioc) arguments(f any) ([]reflect.Value, error) {
	t := reflect.TypeOf(f)
	in := []reflect.Value{}
	for j := 0; j < t.NumIn(); j++ {
		log.Printf("%+v\n", t.In(j).String())
		in = append(in, i.resolve(t.In(j), ""))
	}
	return in, nil
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
