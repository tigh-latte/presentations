package store

import (
	"context"

	"github.com/tigh-latte/heap"
)

type personstore struct {
	store map[int]heap.Person
}

func NewPersonStore() heap.Store[heap.Person] {
	return &personstore{
		store: map[int]heap.Person{
			1: {ID: 1, Name: "Alex", Age: 45},
		},
	}
}

func (p *personstore) Get(ctx context.Context, key int, out *heap.Person) error {
	*out = p.store[key]

	return nil
}

func (p *personstore) Set(ctx context.Context, key int, in *heap.Person) error {
	p.store[key] = *in

	return nil
}
