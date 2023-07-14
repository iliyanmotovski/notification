package orm

import "github.com/latolukasz/beeorm"

type Flusher interface {
	Track(entity interface{})
	Flush()
}

type beeORMFlusher struct {
	beeorm.Flusher
}

func (b *beeORMFlusher) Track(entity interface{}) {
	b.Flusher.Track(entity.(beeorm.Entity))
}

func (b *beeORMFlusher) Flush() {
	b.Flusher.Flush()
}
