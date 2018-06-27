package snowflake

import (
	"errors"
	"sync"
	"time"
)

/*
	snowflake 64int
	1bit    不可用
	41bit   时间戳    time
	10bit   工作机器ID worker
	12bit   序列号 number
*/
const (
	numberBits  uint8 = 12
	workerBits  uint8 = 10
	workerMax   int64 = -1 ^ (-1 << workerBits)
	numberMax   int64 = -1 ^ (-1 << numberBits)
	timeShift   uint8 = workerBits + numberBits //22
	workerShift uint8 = numberBits
	epoch       int64 = 1530097294
)

// Worker in snowflake ID
type Worker struct {
	mu        sync.Mutex
	timeStamp int64
	workID    int64
	number    int64
}

// NewWorker is create a worker whit ID
func NewWorker(id int64) (*Worker, error) {
	if id < 0 || id > workerMax {
		return nil, errors.New("WorkID greater than workerMax")
	}
	return &Worker{
		timeStamp: 0,
		workID:    id,
		number:    0,
	}, nil
}

// GetID is return the snowflake ID
func (w *Worker) GetID() int64 {
	// 加锁
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now().UnixNano() / 1e6
	if w.timeStamp == now {
		w.number++
		// 如果当前毫秒的number用完了，等下一毫秒
		if w.number > numberMax {
			for now <= w.timeStamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		w.number = 0
		w.timeStamp = now
	}
	ID := int64((now-epoch)<<timeShift | (w.workID << workerShift) | (w.number))
	return ID
}
