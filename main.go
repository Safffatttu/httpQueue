package main

import (
	"container/list"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var queueMutex sync.Mutex

type safeQueue struct {
	list  list.List
	mutex sync.Mutex
}

var queues map[string]*safeQueue

func handler(w http.ResponseWriter, r *http.Request) {
	queueName := r.RequestURI[1:]
	queue, found := queues[queueName]

	if !found {
		queues[queueName] = new(safeQueue)
		log.Default().Print("Created New queue", queueName)
		queue, _ = queues[queueName]
	}

	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	switch r.Method {
	case "GET":
		{
			if queue.list.Len() == 0 {
				http.Error(w, "EMPTY QUEUE", 500)
				log.Default().Println("Empty queue")
				return
			}

			v := queue.list.Front().Value.([]uint8)

			queue.list.Remove(queue.list.Front())
			w.Write(v)
			return
		}
	case "PUT":
		{
			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Fatalln(err)
			}

			queue.list.PushBack(body)
			return
		}
	}

	http.Error(w, "Method not used", 400)
}

func main() {
	queues = make(map[string]*safeQueue)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
