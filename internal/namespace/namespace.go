package namespace

import (
	"fmt"
	"github.com/suprunchuksergey/dpl/internal/value"
)

type Namespace interface {
	//создать переменную в текущем пространстве
	Create(name string, v value.Value) error
	//изменить переменную если есть или создать переменную в текущем пространстве
	Set(name string, v value.Value)
	//получить переменную
	Get(name string) (value.Value, error)
	//создать дочернее пространство (возвращает дочернее пространство)
	New(init map[string]value.Value) Namespace
}

type namespace struct {
	value  map[string]value.Value
	parent *namespace
}

func (n *namespace) New(init map[string]value.Value) Namespace {
	if init == nil {
		init = make(map[string]value.Value)
	}

	return &namespace{
		value:  init,
		parent: n,
	}
}

func VarAlreadyExists(name string) error {
	return fmt.Errorf("переменная с именем %s уже существует", name)
}

func (n *namespace) Create(name string, v value.Value) error {
	_, ok := n.value[name]
	if ok {
		return VarAlreadyExists(name)
	}
	n.value[name] = v
	return nil
}

func (n *namespace) Set(name string, v value.Value) {
	if n.set(name, v) {
		return
	}
	n.Create(name, v)
}

func (n *namespace) set(name string, v value.Value) bool {
	_, ok := n.value[name]
	if ok {
		n.value[name] = v
		return true
	}

	if n.parent != nil {
		return n.parent.set(name, v)
	}

	return false
}

func VarDoesNotExist(name string) error {
	return fmt.Errorf("переменной с именем %s не существует", name)
}

func (n *namespace) Get(name string) (value.Value, error) {
	v, ok := n.value[name]
	if ok {
		return v, nil
	}

	if n.parent != nil {
		return n.parent.Get(name)
	}

	return nil, VarDoesNotExist(name)
}

func New(init map[string]value.Value) Namespace {
	if init == nil {
		init = make(map[string]value.Value)
	}

	return &namespace{value: init}
}
