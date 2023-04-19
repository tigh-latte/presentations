package service_test

import (
	"context"
	"testing"

	"github.com/tigh-latte/heap"
	"github.com/tigh-latte/heap/repo"
	"github.com/tigh-latte/heap/service"
	"github.com/tigh-latte/heap/store"
)

type mockRepo struct{}

func (m *mockRepo) Person(ctx context.Context, args *heap.PersonArgs) (heap.Person, error) {
	return heap.Person{Name: "Oonagh", Age: 75, ID: args.ID}, nil
}

func (m *mockRepo) StorePerson(ctx context.Context, person *heap.Person) (heap.Person, error) {
	person.ID = 1
	return *person, nil
}

func Benchmark_SafePerson(b *testing.B) {
	store := store.NewPersonStore()
	repo := repo.NewSafeRepo(store)
	svc := service.NewSafeService(repo)
	for i := 0; i < b.N; i++ {
		_, _ = svc.Person(context.Background(), &heap.PersonArgs{ID: i})
	}
}

func Benchmark_SafeStorePerson(b *testing.B) {
	store := store.NewPersonStore()
	repo := repo.NewSafeRepo(store)
	svc := service.NewSafeService(repo)
	for i := 0; i < b.N; i++ {
		_, _ = svc.StorePerson(context.Background(), &heap.Person{Name: "Bob", Age: 25, ID: i})
	}
}

func Benchmark_UnsafePerson(b *testing.B) {
	store := store.NewPersonStore()
	repo := repo.NewUnsafeRepo(store)
	svc := service.NewUnsafeService(repo)
	for i := 0; i < b.N; i++ {
		_, _ = svc.Person(context.Background(), &heap.PersonArgs{ID: i})
	}
}

func Benchmark_UnsafeStorePerson(b *testing.B) {
	store := store.NewPersonStore()
	repo := repo.NewUnsafeRepo(store)
	svc := service.NewUnsafeService(repo)
	for i := 0; i < b.N; i++ {
		_, _ = svc.StorePerson(context.Background(), &heap.Person{Name: "Bob", Age: 25, ID: i})
	}
}
