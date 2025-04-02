// Event-Driven Scheduling

//A single thread that is going to wake when the earliest task is ready to be executed

// We create a task and push it into a pq, our scheduler reads off the top of our pq and calculates time to sleep until the task is ready to be executed

package main

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

// ------------------
//Task
type Task struct {
	task func() interface{};
	time int;
    index int;
}
// ------------------

// ----------------
//pq
type PriorityQueue []*Task

func (pq PriorityQueue) Len() int {
    return len(pq)
}

func (pq PriorityQueue) Less(i,j int) bool{
    if i >= len(pq) || j >= len(pq) {
        fmt.Printf("Index out of bounds: i=%d, j=%d, len=%d\n", i, j, len(pq))
        return false
    }
    return pq[i].time < pq[j].time
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Task)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(t *Task, new_task func()interface{}, new_time int) {
	t.task = new_task
	t.time = new_time
	heap.Fix(pq, t.index)
}

// ------------------


// -------------------
//Scheduler
type Scheduler_Methods interface{
    AddTask(t *Task)
    Run()
}

type Scheduler struct {
    PQ *PriorityQueue
    mutex sync.Mutex
    newTaskCh chan struct{}
}

func (s *Scheduler) AddTask(t *Task){
    //Push the task into the PQ,
    s.mutex.Lock()
    s.PQ.Push(t)
    s.mutex.Unlock()

    select{
    case s.newTaskCh <- struct{}{}:

    default:
    }
    return
}

func (s *Scheduler) Run(){
    //Should keep alive while PQ is active
    startTime := time.Now()

    for {
        s.mutex.Lock()
        if s.PQ.Len() > 0 {
            nextTask := (*s.PQ)[0]

            elapsedTime := time.Since(startTime).Seconds()
            remainingTime := nextTask.time - int(elapsedTime)

            if remainingTime < 0 {
                remainingTime = 0
            }

            waitDuration := time.Duration(remainingTime) * time.Second
            s.mutex.Unlock()

            timer := time.NewTimer(waitDuration)

            select{
            case <- timer.C:

                s.mutex.Lock()
                if s.PQ.Len() > 0 {
                    task := heap.Pop(s.PQ).(*Task)
                    s.mutex.Unlock()

                    go func(t *Task){
                        result := t.task()
                        fmt.Println(result)
                    }(task)
                } else {
                    s.mutex.Unlock()
                }
            case <- s.newTaskCh:
                timer.Stop()
                continue
            }
        } else {
            s.mutex.Unlock()
            select {
            case <- s.newTaskCh:
                continue
            }
        }

    }
}
//--------------------

func main(){
    a := func(c int, b int) int{
        return c + b
    }

    task1 := &Task{
        task: func() interface{} {
            return a(1,2)
        },
        time: 10,
        index: 0,
    }
    task2 := &Task{
        task: func() interface{} {
            return a(3,4)
        },
        time: 7,
        index: 0,
    }
    task3 := &Task{
        task: func() interface{} {
            return a(5,6)
        },
        time: 1,
        index: 0,
    }

    PQ := make(PriorityQueue, 0)

    heap.Init(&PQ)

    event_scheduler := &Scheduler{
        PQ: &PQ,
        newTaskCh: make(chan struct{}, 1),
    }

    go event_scheduler.Run()

    event_scheduler.AddTask(task1)
    event_scheduler.AddTask(task2)
    event_scheduler.AddTask(task3)

     select{}
}
