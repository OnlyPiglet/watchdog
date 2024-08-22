package watchdog

import (
	"fmt"
	"sync"
	"time"
)

type Sheep interface {
	//投喂的其他处理
	Feed() error
	//跑飞的处理
	GetFood() chan struct{}
	RunAway()
	Name() string
	DeadTime() time.Duration
}

type SingleSheepMonitor struct {
	ticker         *time.Ticker
	sheep          *Sheep
	food           chan struct{}
	deadtime       time.Duration
	monitorStopped chan struct{}
	once           *sync.Once
}

func (ssm *SingleSheepMonitor) closeWather() {
	println((*ssm.sheep).Name() + " watcher stopped")
	ssm.once.Do(func() {
		ssm.monitorStopped <- struct{}{}
		close(ssm.monitorStopped)
		close(ssm.food)
		ssm.ticker.Stop()
	})

}

type WatchDog struct {
	sheepfold sync.Map
}

func NewWatchDog() *WatchDog {
	return &WatchDog{
		sheepfold: sync.Map{},
	}
}

func (s *SingleSheepMonitor) run() {
	defer func() {
		s.closeWather()
	}()

	for {
		select {
		case <-s.ticker.C:
			(*s.sheep).RunAway()
			return
		case _,ok := <-s.food:
			if (*s.sheep).Feed() != nil || !ok {
				return
			}
			s.ticker.Reset(s.deadtime)
		}
	}
}

func (w *WatchDog) AddSheep(sheep Sheep) error {
	if _, ok := w.sheepfold.Load(sheep.Name()); !ok {
		singleSheepMonitor := &SingleSheepMonitor{
			sheep:          &sheep,
			deadtime:       sheep.DeadTime(),
			ticker:         time.NewTicker(sheep.DeadTime()),
			food:           sheep.GetFood(),
			monitorStopped: make(chan struct{}, 1),
			once:           &sync.Once{},
		}
		w.sheepfold.Store(sheep.Name(), singleSheepMonitor)
		return nil
	} else {
		return fmt.Errorf("shape %s has existed in this sheepfold", sheep.Name())
	}
}

func (w *WatchDog) StartWatching(sheepName string) (chan struct{}, error) {
	if v, ok := w.sheepfold.Load(sheepName); ok {
		ssm := v.(*SingleSheepMonitor)
		go ssm.run()
		return ssm.monitorStopped, nil
	} else {
		return nil, fmt.Errorf("shape %s has existed in this sheepfold ", sheepName)
	}
}

func (w *WatchDog) GetSheep(sheepName string) (*SingleSheepMonitor, bool) {
	value, ok := w.sheepfold.Load(sheepName)
	if ok {
		return value.(*SingleSheepMonitor), ok
	}
	return nil, ok
}

func (w *WatchDog) RemoveSheep(sheepName string) {
	value, ok := w.sheepfold.Load(sheepName)
	if ok {
		ssm := value.(*SingleSheepMonitor)
		ssm.closeWather()
	}
	w.sheepfold.Delete(sheepName)
}

func rangeCloseSingleSheepMonitor(_, singleSheepMonitor interface{}) bool {
	ssm := singleSheepMonitor.(*SingleSheepMonitor)
	ssm.closeWather()
	return true
}

func (w *WatchDog) Close() {
	w.sheepfold.Range(rangeCloseSingleSheepMonitor)
}
