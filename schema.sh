curl -X POST localhost:8080/alter -d '{"drop_all": true}'
#curl -X POST localhost:8080/alter -d ‘{“drop_op”: “DATA”}’
curl -X POST localhost:8080/alter -d '   task.id: string @index(hash) . 
    task.title: string @index(fulltext).
    task.description: string .
    task.task: [uid] @reverse .

type Task{
   task.id 
   task.title 
   task.description 
   task.task 
}
'