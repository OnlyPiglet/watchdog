package watchdog

import (
	"fmt"
	"testing"
	"time"
)

type MySheep struct {
	food chan struct{}
	name string
}

func (ms *MySheep) RunAway() {
	println(fmt.Sprintf("%s is running away", ms.name))
}

func (ms *MySheep) GetFood() chan struct{} {
	return ms.food
}

func (ms *MySheep) Feed() error {
	println("ms feed")
	return nil
}

func (ms *MySheep) Name() string {
	return ms.name
}

func (ms *MySheep) DeadTime() time.Duration {
	return 5 * time.Second
}

func TestNewWatchDog(t *testing.T) {

	dog := NewWatchDog()

	sheep := MySheep{
		food: make(chan struct{}, 1),
		name: "my sheep",
	}

	err := dog.AddSheep(&sheep)
	if err != nil {
		panic(err)
	}

	watcherStopped, _ := dog.StartWatching(sheep.name)

	tt := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-tt.C:
			sheep.food <- struct{}{}
		case <-watcherStopped:
			println("watcher stop")
			goto stopped
		}
	}
stopped:
}
