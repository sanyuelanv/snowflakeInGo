package snowflake

import (
	"fmt"
	"testing"
)

func TestSnowflake(t *testing.T) {
	count := 10
	var workerID int64 = 1
	worker, err := NewWorker(workerID)
	if err != nil {
		fmt.Println(err)
		return
	}
	ch := make(chan int64)
	for index := 0; index < count; index++ {
		go func() {
			id := worker.GetID()
			ch <- id
		}()
	}
	defer close(ch)
	for index := 0; index < count; index++ {
		id := <-ch
		fmt.Println(id)
	}
	fmt.Println("OK!")
}
