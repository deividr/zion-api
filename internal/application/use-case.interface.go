package usecases

type UseCase[T any, R any] interface {
	Execute(input T) (R, error)
}
