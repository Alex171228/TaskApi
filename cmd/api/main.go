package main

import "database/sql"

type repo struct {
	s      []int
	nextId int
}

func (r *repo) get() {

}
func (r *repo) create() {

}
func Newrepo() *repo {
	return &repo{}
}

type service struct {
	r db
}
type db interface {
	get()
	create()
}

func (s *service) get() {
	s.r.get()
}
func (s *service) create() {
	s.r.create()
}

func NewService(r db) *service {
	return &service{r: r}
}

type repoDb struct {
	*sql.DB
}

func (repo *repoDb) get() {

}
func (repo *repoDb) create() {

}
func main() {
	r1 := Newrepo()
	r2 := repoDb{}
	s3 := NewService(r1)
}
