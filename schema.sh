curl -X POST localhost:8080/alter -d '{"drop_all": true}'
#curl -X POST localhost:8080/alter -d ‘{“drop_op”: “DATA”}’
curl -X POST localhost:8080/alter -d '   task.id: string @index(hash) . 
    task.title: string @index(fulltext).
    task.description: string .
    task.task: [uid] @reverse .
    task.person: uid  .
    person.id: string @index(hash).
    person.name: string .
    person.email: string @index(fulltext).
    person.tasks: [uid] .

type Task{
   task.id 
   task.title 
   task.description
   task.person 
   task.task 
}

type Person{
	person.id
	person.name
	person.email
	person.tasks
}
'
