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

func TestWatchDog_RemoveSheep(t *testing.T) {

	dog := NewWatchDog()

	sheep := MySheep{
		food: make(chan struct{}, 1),
		name: "my sheep",
	}

	getSheep, b := dog.GetSheep((&sheep).name)

	println(fmt.Sprintf("%v", getSheep))
	println(b)

	err := dog.AddSheep(&sheep)
	if err != nil {
		panic(err)
	}

	getSheep, b = dog.GetSheep((&sheep).name)

	println(fmt.Sprintf("%v", getSheep))
	println(b)

	dog.StartWatching((&sheep).name)

	time.NewTicker(3 * time.Second)

	dog.RemoveSheep(sheep.name)

	getSheep, b = dog.GetSheep((&sheep).name)

	println(fmt.Sprintf("%v", getSheep))
	println(b)

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
