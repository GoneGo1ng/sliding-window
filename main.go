package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	sw := InitWindow(5, 2000)
	for i := 1; i <= 20; i++ {
		time.Sleep(1 * time.Second)
		if i == 10 {
			time.Sleep(5 * time.Second)
		}
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				sw.Write(1)
			}()
		}
		wg.Wait()
		log.Println("----", i, "----")
		log.Println(sw.Read())
	}
}

type SlidingWindow struct {
	mu         sync.RWMutex
	ClockAt    time.Time // 滑动时间
	WindowSize uint      // 窗口大小
	WindowTime uint      // 窗口时间（单位：毫秒）
	WindowList []uint    // 窗口数据列表
}

// 初始化窗口
func InitWindow(windowSize, windowTime uint) *SlidingWindow {
	sw := &SlidingWindow{
		ClockAt:    time.Now(),
		WindowSize: windowSize,
		WindowTime: windowTime,
		WindowList: make([]uint, windowSize),
	}
	return sw
}

// 读取窗口数据列表
func (sw *SlidingWindow) Read() []uint {
	sw.mu.RLock()
	defer sw.mu.RUnlock()

	return sw.WindowList
}

// 窗口中写入一个数值
func (sw *SlidingWindow) Write(num uint) {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	sw.Shift()
	sw.WindowList[len(sw.WindowList)-1] = num + sw.WindowList[len(sw.WindowList)-1]
}

// 滑动窗口
func (sw *SlidingWindow) Shift() {
	now := time.Now()
	// 两次操作的时间差
	passTime := now.Sub(sw.ClockAt).Milliseconds()
	// 需要滑动的窗口格数
	length := uint(passTime) / sw.WindowTime
	if length > 0 {
		// 重置滑动时间
		sw.ClockAt = time.Now()
		// 如果需要滑动的窗口格数大于窗口大小，则初始化窗口数据列表
		if length > sw.WindowSize {
			sw.WindowList = make([]uint, sw.WindowSize)
			return
		}
		// 滑动相应的窗口格数
		sw.Move(length)
	}
}

// 滑动相应的窗口格数
func (sw *SlidingWindow) Move(length uint) {
	zeroList := make([]uint, length)
	sw.WindowList = append(sw.WindowList[length:], zeroList...)

	/*newList := make([]uint, sw.WindowSize)
	index := 0
	for i, num := range sw.WindowList {
		if uint(i) > length-1 {
			newList[index] = num
			index++
		}
	}
	sw.WindowList = newList*/
}
