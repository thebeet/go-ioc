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

func (b *bar) IocPostConstruct() {
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

func TestFill(t *testing.T) {

	type DB struct {
		dsn string
	}

	type Bar struct {
		Db2 *DB `autowire:"db2"`
	}

	type Foo struct {
		Db1 *DB  `autowire:"db1"`
		Bar *Bar `autowire:""`
	}

	NewDb := func(dsn string) *DB {
		return &DB{dsn}
	}

	NewBar := func() *Bar {
		return &Bar{}
	}
	ioc := New()

	ioc.RegisterInstanceWithName(NewDb("dsn1"), "db1")
	ioc.RegisterInstanceWithName(NewDb("dsn2"), "db2")
	ioc.Register(NewBar)

	var foo Foo

	ioc.Fill(&foo)
	if foo.Db1.dsn != "dsn1" {
		t.Errorf("fill fail")
	}
	if foo.Bar.Db2.dsn != "dsn2" {
		t.Errorf("fill fail")
	}
}

func TestInstanceRegister(t *testing.T) {
	type DB struct {
		dsn string
	}

	type Foo struct {
		Db1 *DB `autowire:"db1"`
	}

	type Bar struct {
		Db2 *DB `autowire:"db2"`
	}

	NewDb := func(dsn string) *DB {
		return &DB{dsn}
	}

	NewFoo := func() *Foo {
		return &Foo{}
	}

	NewBar := func() *Bar {
		return &Bar{}
	}

	ioc := New()

	ioc.RegisterInstanceWithName(NewDb("dsn1"), "db1")
	ioc.RegisterInstanceWithName(NewDb("dsn2"), "db2")
	ioc.Register(NewFoo)
	ioc.Register(NewBar)

	ioc.Call(func(foo *Foo, bar *Bar) {
		if foo.Db1.dsn != "dsn1" {
			t.Errorf("autowire fail")
		}
		if bar.Db2.dsn != "dsn2" {
			t.Errorf("autowire fail")
		}
	})
}
