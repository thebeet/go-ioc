package ioc

import (
	"log"
	"testing"
)

type Foo interface {
	Hello() string
}

type foo struct {
}

func NewFoo() Foo {
	log.Printf("foo init")
	return &foo{}
}

func (c *foo) Hello() string {
	return "Hello"
}

type Bar interface {
	World() string
}

type bar struct {
	foo Foo
	s   string
}

func NewBar(foo Foo) Bar {
	log.Printf("foo init")
	return &bar{foo: foo}
}

func (b *bar) World() string {
	return b.foo.Hello() + ",Wolrd:" + b.s
}

func (b *bar) PostConstruct() {
	b.s = "PostConstruct"
}

func TestIoc(t *testing.T) {
	ioc := New()

	ioc.Register(NewFoo)
	ioc.Register(NewBar)

	ioc.Call(func(bar Bar) {
		s := bar.World()
		if s != "Hello,Wolrd:PostConstruct" {
			t.Errorf("ioc fail, expect Hello,Wolrd:PostConstruct but %s", s)
		}
	})
}

func TestAutoWire(t *testing.T) {
	type Foo struct {
		S string
	}
	type Bar struct {
		FooImpl *Foo `autowire:""`
	}
	NewFoo := func() *Foo {
		f := Foo{}
		f.S = "foo"
		return &f
	}
	NewBar := func() *Bar {
		return &Bar{}
	}

	ioc := New()

	ioc.Register(NewFoo)
	ioc.Register(NewBar)

	ioc.Call(func(bar *Bar) {
		if bar.FooImpl.S != "foo" {
			t.Errorf("autowire fail")
		}
	})
}
