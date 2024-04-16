package storage

type Storage interface {
	Search(term string) ([]string, error)
}
