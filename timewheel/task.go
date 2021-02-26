package timewheel

import "time"

type Task struct {
	delay time.Duration
	circle int // 当延迟时间超过时间轮最大时间，代表再第几圈
	pos int
	id int64
	run func()
}

func newTask(delay time.Duration, f func(), pos,circle int) *Task {
	task := Task{
		delay:  delay,
		circle: circle,
		run: f,
		pos: pos,
	}
	return &task
}