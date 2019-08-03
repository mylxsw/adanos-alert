package repository

type Sequence struct {
	Name  string `bson:"name" json:"name"`
	Value int64  `bson:"value" json:"value"`
}

type SequenceRepo interface {
	Next(name string) (seq Sequence, err error)
	Truncate(name string) error
}
