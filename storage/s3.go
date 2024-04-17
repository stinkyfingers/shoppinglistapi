package storage

type S3 struct{}

func NewS3(profile string) (Storage, error) {
	return &S3{}, nil
}

func (s *S3) Search(term string) ([]string, error) {
	return nil, nil
}
