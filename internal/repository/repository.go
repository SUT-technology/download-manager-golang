package repository

type QueryFunc = func(r *Repo) error

type Pool interface {
	Query(f QueryFunc) error
	Close() error
}

type Tables struct {
	Downloads DownloadRepository
}

type Repo struct {
	Tables
}
