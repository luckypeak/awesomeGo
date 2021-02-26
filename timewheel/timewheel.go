package timewheel

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

type Job func(interface{})


type TimeWheel struct {
    interval time.Duration
    ticker *time.Ticker
    time.Timer
    stop chan struct{}
    slots []*list.List
    slotNum int
    currentPos int

    taskId int64
    taskMap map[int64]*Task

    mu sync.Mutex
    once sync.Once
}

func (tw *TimeWheel)getTaskID() int64 {
	tw.mu.Lock()
	tw.taskId++
	tw.mu.Unlock()
	return tw.taskId
}

func New(interval time.Duration, slots int) (*TimeWheel, error)  {
	wheel := new(TimeWheel)
	wheel.taskMap = map[int64]*Task{}
	wheel.interval = interval
	wheel.slots = make([]*list.List, slots)
	wheel.slotNum = slots
	wheel.ticker = time.NewTicker(interval)
	for i := 0; i < slots; i ++{
		wheel.slots[i] = list.New()
	}
	return wheel, nil
}

func (tw *TimeWheel)AfterFunc(delay time.Duration,  f func()) (int64, error) {
	if delay.Seconds() == 0 {
		return -1, errors.New("must big 0")
	}
	circle, pos := tw.findPos(delay)
	task := newTask(delay, f, pos,circle)
	task.id = tw.getTaskID()
	tw.slots[pos].PushBack(task)
	tw.mu.Lock()
	tw.taskMap[task.id] = task
	tw.mu.Unlock()
	return task.id, nil
}


func (tw *TimeWheel) RemoverTimer(id int64)  {
	 tw.mu.Lock()
	 task, ok := tw.taskMap[id]
	 tw.mu.Unlock()
	 if !ok {
		 return
	 }
	 l := tw.slots[task.pos]
	 for e := l.Front(); e!= nil; e = e.Next(){
	 	 t := e.Value.(*Task)
	 	 if t.id == id{
	 	 	l.Remove(e)
			 return
		 }
	 }
}

func (tw *TimeWheel)findPos(delay time.Duration) (circle, pos int) {
	delaySeconds := int(delay.Seconds())
	tickSeconds := int(tw.interval.Seconds())
	circle = delaySeconds/tickSeconds/tw.slotNum
	pos = delaySeconds / tickSeconds %tw.slotNum
	return
}

func (tw *TimeWheel)scanAndRunTask(l *list.List)  {
	for e := l.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.circle > 0{
			task.circle--
			e = e.Next()
			continue
		}
		go task.run()
		tw.mu.Lock()
		delete(tw.taskMap, task.id)
		tw.mu.Unlock()
		next := e.Next()
		l.Remove(e)
		e = next

	}
}
func (tw *TimeWheel)onTicker()  {
	l := tw.slots[tw.currentPos]
	tw.scanAndRunTask(l)
	if tw.currentPos == tw.slotNum-1{
		tw.currentPos = 0
	}else {
		tw.currentPos ++
	}
}

func (tw *TimeWheel) Start() {
	tw.once.Do(func() {
		go tw.run()
	})
}

func (tw *TimeWheel) Stop()  {
	tw.stop <- struct{}{}
}

func (tw *TimeWheel) run() {
	for {
		select {
		 	case <-tw.ticker.C:
		 		tw.onTicker()
			case <- tw.stop:
				tw.ticker.Stop()
		}
	}
}