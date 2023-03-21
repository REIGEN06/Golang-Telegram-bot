package main

import (
	"fmt"

	"github.com/REIGEN06/RPBD/tree/people_service/service/store"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	s := store.NewStore("postgres://postgres:admin@localhost:5432/postgres")
	people, _ := s.ListPeople()
	fmt.Println(people)

	person, _ := s.GetPeopleByID(1)
	fmt.Println(person)
}
