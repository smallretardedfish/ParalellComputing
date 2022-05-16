package async

type Promise[T any] struct {
	ch chan T
}

func (p *Promise[T]) Await() T {
	return <-p.ch
}

//DoAsync runs given function asynchronously and returns Promise
func DoAsync[T any](f func() T) *Promise[T] {
	res := make(chan T)
	go func() {
		res <- f()
	}()
	return &Promise[T]{ch: res}
}
