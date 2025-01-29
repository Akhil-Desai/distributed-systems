// Event-Driven Scheduling

//A single thread that is going to wake when the earliest task is ready to be executed

// We create a task and push it into a pq, our scheduler reads off the top of our pq and calculates time to sleep until the task is ready to be executed

package main

import (
	"fmt"
)

//Task
type Task struct {
	task func(...interface{}) interface{};
	time int;
}


func main(){
    task1 := Task{
        task: func(i ...interface{}) interface{} {
            a := i[0].(int)
            b := i[1].(int)
            return a * b
        },
        time: 10,
    }

    fmt.Println(task1.task(1,2))


}
