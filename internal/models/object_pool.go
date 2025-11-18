package models

type Projectile interface {
	IsActive() bool
	SetActive(bool)
}

type Pool[T Projectile] struct {
	objects []T
	next    int
}

func (pool *Pool[T]) Get() T {
	if pool.objects[pool.next].IsActive() {
		return *new(T)
	}

	object := pool.objects[pool.next]
	pool.next = (pool.next + 1) % len(pool.objects)
	return object
}

func (pool *Pool[T]) Elements() []T {
	return pool.objects
}

func (pool *Pool[T]) Reset() {
	for _, object := range pool.objects {
		object.SetActive(false)
	}
}

func CreatePool[T Projectile](bullets []T) Pool[T] {
	return Pool[T]{
		objects: bullets,
		next:    0,
	}
}
