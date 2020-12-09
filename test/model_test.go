package main

import (
	"github.com/google/uuid"
	"testing"
	m "todo/models"
)

func BenchmarkCreate(b *testing.B) {
	//createtask
	for i := 0; i < b.N; i++ {
		id := uuid.New()
		ts := m.Task{ID: id.String(), Description: "test", Title: "test task", Time: "000000000000"}
		m.CreateTask(ts)
	}

}

func BenchmarkGet(b *testing.B) {
	//get task
	// id is one that is already created
	ts := "1acf4e37-e190-43f6-b142-6373ceba704d"
	for i := 0; i < b.N; i++ {
		m.GetTask(ts)
	}
}

func BenchmarkGetAll(b *testing.B) {
	//get task
	for i := 0; i < b.N; i++ {
		m.GetAll()
	}
}

func BenchmarkUpdate(b *testing.B) {
	//update task
	var ts m.Task
	// id should exist
	ts.ID = "9420e46f-7352-458d-abc9-ccf37f989d35"
	ts.Description = "test benchmark"
	ts.Title = "test"
	for i := 0; i < b.N; i++ {
		m.UpdateTask(ts)
	}

}

//func TestDelete(t *testing.T) {
//// delete task
//m.DeleteTask(ts)

//}
