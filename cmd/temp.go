package main

import (
	"log"
	m "todo/models"

	"github.com/google/uuid"
)

func main() {
	u := uuid.New()
	res, err := m.CreatePerson(m.Person{ID: u.String(), Name: "Kurian", Email: "kurian.c@shopalyst.com"})
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
	log.Println(u)

	res, err = m.AssignTask(u.String(), "d7f5f89d-d5de-41d1-a7b8-b968715b069c")
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
	res, err = m.UpdatePerson(u.String(), "kurian.c@intel.com")
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
	return
}
