package namespace

import (
	"fmt"
	"github.com/suprunchuksergey/dpl/val"
)

type namespace struct {
	m      map[string]val.Val
	parent Namespace
}

func (ns *namespace) Get(key string) (val.Val, error) {
	v, ok := ns.m[key]
	if !ok {
		if ns.parent == nil {
			return nil, fmt.Errorf("переменная %s не существует", key)
		}
		return ns.parent.Get(key)
	}
	return v, nil
}

func (ns *namespace) Set(k string, v val.Val) { ns.m[k] = v }

type Namespace interface {
	Get(key string) (val.Val, error)
	Set(k string, v val.Val)
}

func New(init map[string]val.Val) Namespace {
	return WithParent(nil, init)
}

func WithParent(parent Namespace, init map[string]val.Val) Namespace {
	if init == nil {
		init = make(map[string]val.Val)
	}
	return &namespace{m: init, parent: parent}
}
