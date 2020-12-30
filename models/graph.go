package models

import (
	"context"
	"encoding/json"
	"os"

	dgo "github.com/dgraph-io/dgo/v2"
	api "github.com/dgraph-io/dgo/v2/protos/api"
	"google.golang.org/grpc"
)

/* This is for consumption of database*/

// -- this struct represents the user payload
type CreateTaskPlInt struct {
	Title       string `json:"task.title"`
	Description string `json:"task.description"`
	Parent      string `json:"task.parent"`
	Assignee    Person `json:"task.person"`
	Time        string `json:"task.time"`
}

//-- this struct represents additional fields response to user
type CreateTaskDB struct {
	ID      string `json:"task.id"`
	Type    string `json:"dgraph.type"`
	SubTask []Task `json:"task.task"`
}

type Task struct {
	CreateTaskPlInt
	CreateTaskDB
}

// This struct is for user

// -- this struct represents the user payload

type CreateTaskPl struct {
	Title       string `json:"title" example:"app title"`
	Description string `json:"description" example:"new app"`
	Parent      string `json:"-" `
	Assignee    Person `json:"task.person"`
	Time        string `json:"time" example:"20202020"`
}

//-- this struct represents additional fields response to user
type CreateTaskDBPl struct {
	ID      string         `json:"id"`
	Type    string         `json:"-"`
	SubTask []CreateTaskPl `json:"subtask"`
}

type TaskPl struct {
	CreateTaskPl
	CreateTaskDBPl
}

type Person struct {
	ID    string `json:"person.id"`
	Name  string `json:"person.name"`
	Email string `json:"person.email"`
	Type  string `json:"dgraph.type"`
}

// Task represents all the teasks and subtasks
//type Task struct {
//ID          string `json:"task.id"`
//Title       string `json:"task.title"`
//Description string `json:"task.description"`
//Time        string `json:"task.time"`
//Subtask     []Task `json:"task.task"`
//Parent      string `json:"task.parent"`
//Type        string `json:"dgraph.type,omitempty"`
//}

var (
	dg   *dgo.Dgraph
	ctx  context.Context
	conn *grpc.ClientConn
)

func init() {
	c, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		os.Exit(1)
	}
	conn = c
	dg = dgo.NewDgraphClient(api.NewDgraphClient(conn))
	ctx = context.Background()
	return
}

// Link function links the parent and child node
func Link(parent, child string) (res *api.Response, err error) {
	q := `<` + parent + `> <task.task> <` + child + `> .`
	link := &api.Mutation{
		SetNquads: []byte(q),
		CommitNow: true,
	}
	res, err = dg.NewTxn().Mutate(ctx, link)
	return
}

func GetPersonUID(id string) (res *api.Response, err error) {
	q := `query gettask($id:string)
	{ getTasks(func: eq(person.id,$id)){
			uid
		}
	}`
	variables := map[string]string{"$id": id}
	res, err = dg.NewReadOnlyTxn().QueryWithVars(ctx, q, variables)
	return

}

// GetUID gets the uid of given id node
func GetUID(id string) (res *api.Response, err error) {
	q := `query gettask($id:string)
	{ getTasks(func: eq(task.id,$id)){
			uid
		}
	}`

	variables := map[string]string{"$id": id}
	res, err = dg.NewReadOnlyTxn().QueryWithVars(ctx, q, variables)
	return
}

func GetAll() (res *api.Response, err error) {
	// get all task
	q := `query { 
	    getTasks(func: has(task.id)){
		id:	task.id
		title:	task.title
		description:	task.description
		time:	task.time
		subtask: task.task{
			    id:task.id
		    	    title:task.title
		    	    description:task.description
		    	    time:task.time
			}
	    }
	}`
	res, err = dg.NewReadOnlyTxn().Query(ctx, q)
	if err != nil {
		return
	}

	return
}

// GetTask gets the task details provided ID string
func GetTask(id string) (res *api.Response, err error) {
	// get all task
	q := `query gettask($id:string)
	{ getTasks(func: eq(task.id,$id)){
		title:task.title
	    	description:task.description
	    	time:task.time
	    	task:task.task{
		    title:task.title
		    description:task.description
		    time:task.time
			}
	    }
	}`

	variables := map[string]string{"$id": id}
	res, err = dg.NewReadOnlyTxn().QueryWithVars(ctx, q, variables)
	return
}

// UpdateTask function updates fields of task
// TODO: update subtask using uid
func UpdateTask(task Task) (res *api.Response, err error) {
	// update a task, for given task id
	query := `query {
	    task as var(func: eq(task.id,"` + task.ID + `"))

	}`

	m1 := &api.Mutation{
		SetNquads: []byte(`uid(task) <task.description> "` + task.Description + `" .`),
	}

	m2 := &api.Mutation{
		SetNquads: []byte(`uid(task) <task.title> "` + task.Title + `" .`),
	}
	m3 := &api.Mutation{
		SetNquads: []byte(`uid(task) <task.time> "` + task.Time + `" .`),
	}

	req := &api.Request{
		Query:     query,
		Mutations: []*api.Mutation{m1, m2, m3},
		CommitNow: true,
	}

	res, err = dg.NewTxn().Do(ctx, req)
	return
}

func DeleteEdges(uid string, edges ...string) (res *api.Response, err error) {
	mu := &api.Mutation{}
	for _, predicate := range edges {
		mu.Del = append(mu.Del, &api.NQuad{
			Subject:     uid,
			Predicate:   predicate,
			ObjectValue: &api.Value{Val: &api.Value_DefaultVal{DefaultVal: "_STAR_ALL"}},
		})
	}

	mu.CommitNow = true
	res, err = dg.NewTxn().Mutate(ctx, mu)
	return
}

//DeleteTask is for deleting nodes
func DeleteTask(id []string) (res *api.Response, err error) {
	// get all task
	var qr []map[string]string
	for _, val := range id {
		qr = append(qr, map[string]string{"uid": val})
	}
	qry, err := json.Marshal(qr)
	del := &api.Mutation{CommitNow: true, DeleteJson: qry}
	res, err = dg.NewTxn().Mutate(ctx, del)
	return
}

// CreateTask function creates a node
//func CreateTask(ctx context.Context, dg *dgo.Dgraph, t Task) (err error, res *api.Response) {
func CreateTask(t Task) (res *api.Response, err error) {
	// make a mutation
	t.Type = "Task"
	tj, _ := json.Marshal(t) // ignoring error because struct is created by app
	create := &api.Mutation{
		SetJson:   tj,
		CommitNow: true,
	}
	res, err = dg.NewTxn().Mutate(ctx, create)
	return
}

func AssignTask(userId string, taskid string) (res *api.Response, err error) {
	// assign task from email
	query := `query {
	    task as var(func: eq(task.id,"` + taskid + `"))
	    person as var(func: eq(person.id,"` + userId + `"))
	}`

	m1 := &api.Mutation{
		SetNquads: []byte(`<uid(task)> <task.person> <uid(person)>.`),
	}
	req := &api.Request{
		Query:     query,
		Mutations: []*api.Mutation{m1},
		CommitNow: true,
	}

	res, err = dg.NewTxn().Do(ctx, req)
	return
}

func UpdatePerson(ID string, email string) (res *api.Response, err error) {
	query := `query {
	    person as var(func: eq(person.id,"` + ID + `"))

	}`

	m1 := &api.Mutation{
		SetNquads: []byte(`uid(person) <person.email> "` + email + `" .`),
	}

	req := &api.Request{
		Query:     query,
		Mutations: []*api.Mutation{m1},
		CommitNow: true,
	}
	res, err = dg.NewTxn().Do(ctx, req)
	return
}

func CreatePerson(p Person) (res *api.Response, err error) {
	// create person
	p.Type = "Person"
	pj, _ := json.Marshal(p)
	create := &api.Mutation{
		SetJson:   pj,
		CommitNow: true,
	}
	res, err = dg.NewTxn().Mutate(ctx, create)
	return
}
