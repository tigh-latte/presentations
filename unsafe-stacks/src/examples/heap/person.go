package heap

import "context"

type PersonArgs struct {
	ID int
}

type Person struct {
	Name string
	Age  int16
	ID   int
}

type PersonService interface {
	PersonRepo
}

type PersonRepo interface {
	Person(ctx context.Context, args *PersonArgs) (Person, error)
	StorePerson(ctx context.Context, person *Person) (Person, error)
}

type Store[T any] interface {
	Get(ctx context.Context, key int, out *T) error
	Set(ctx context.Context, key int, in *T) error
}
