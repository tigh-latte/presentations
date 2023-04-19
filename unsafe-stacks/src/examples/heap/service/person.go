package service

import (
	"context"

	"github.com/tigh-latte/heap"

	_ "unsafe"
)

type safesvc struct {
	repo heap.PersonRepo
}

func NewSafeService(repo heap.PersonRepo) heap.PersonService {
	return &safesvc{repo: repo}
}

func (s *safesvc) Person(ctx context.Context, args *heap.PersonArgs) (heap.Person, error) {
	return s.repo.Person(ctx, args)
}

func (s *safesvc) StorePerson(ctx context.Context, person *heap.Person) (heap.Person, error) {
	return s.repo.StorePerson(ctx, person)
}

type unsafesvc struct {
	repo heap.PersonRepo
}

func NewUnsafeService(repo heap.PersonRepo) heap.PersonService {
	return &unsafesvc{repo: repo}
}

func (s *unsafesvc) Person(ctx context.Context, args *heap.PersonArgs) (heap.Person, error) {
	return __person(ctx, s.repo, args)
}

func (s *unsafesvc) StorePerson(ctx context.Context, person *heap.Person) (heap.Person, error) {
	return __store_person(ctx, s.repo, person)
}

//go:noescape
//go:linkname __person github.com/tigh-latte/heap/service.__impl__person
func __person(ctx context.Context, repo heap.PersonRepo, args *heap.PersonArgs) (heap.Person, error)
func __impl__person(ctx context.Context, repo heap.PersonRepo, args *heap.PersonArgs) (heap.Person, error) {
	return repo.Person(ctx, args)
}

//go:noescape
//go:linkname __store_person github.com/tigh-latte/heap/service.__impl__store_person
func __store_person(ctx context.Context, repo heap.PersonRepo, person *heap.Person) (heap.Person, error)
func __impl__store_person(ctx context.Context, repo heap.PersonRepo, person *heap.Person) (heap.Person, error) {
	return repo.StorePerson(ctx, person)
}
