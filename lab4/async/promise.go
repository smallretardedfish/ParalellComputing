package async

type Promise[T any] struct {

}


type DoAsync[T any](f func() T)