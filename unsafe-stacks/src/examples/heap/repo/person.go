package repo

import (
	"context"

	"github.com/tigh-latte/heap"

	_ "unsafe"
)

type saferepo struct {
	store heap.Store[heap.Person]
}

func NewSafeRepo(store heap.Store[heap.Person]) heap.PersonRepo {
	return &saferepo{store: store}
}

func (s *saferepo) Person(ctx context.Context, args *heap.PersonArgs) (heap.Person, error) {
	var person heap.Person
	if err := s.store.Get(ctx, args.ID, &person); err != nil {
		return person, err
	}

	return person, nil
}

func (s *saferepo) StorePerson(ctx context.Context, person *heap.Person) (heap.Person, error) {
	if err := s.store.Set(ctx, person.ID, person); err != nil {
		return heap.Person{}, err
	}
	return *person, nil
}

type unsaferepo struct {
	store heap.Store[heap.Person]
}

func NewUnsafeRepo(store heap.Store[heap.Person]) heap.PersonRepo {
	return &unsaferepo{store: store}
}

func (u *unsaferepo) Person(ctx context.Context, args *heap.PersonArgs) (heap.Person, error) {
	var person heap.Person
	if err := __get(ctx, u.store, args.ID, &person); err != nil {
		return heap.Person{}, err
	}

	return person, nil
}

func (u *unsaferepo) StorePerson(ctx context.Context, person *heap.Person) (heap.Person, error) {
	if err := __set(ctx, u.store, person.ID, person); err != nil {
		return heap.Person{}, err
	}
	return *person, nil
}

//go:noescape
//go:linkname __get github.com/tigh-latte/heap/repo.__impl__get
func __get(ctx context.Context, store heap.Store[heap.Person], ID int, out *heap.Person) error
func __impl__get(ctx context.Context, store heap.Store[heap.Person], ID int, out *heap.Person) error {
	return store.Get(ctx, ID, out)
}

//go:noescape
//go:linkname __set github.com/tigh-latte/heap/repo.__impl__set
func __set(ctx context.Context, store heap.Store[heap.Person], key int, in *heap.Person) error
func __impl__set(ctx context.Context, store heap.Store[heap.Person], key int, in *heap.Person) error {
	return store.Set(ctx, key, in)
}
