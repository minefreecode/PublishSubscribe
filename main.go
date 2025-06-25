package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== Паттерн публикация и подписка ===")

	b := newBroadcaster()
	numSubscribers := 3
	var wg sync.WaitGroup

	for i := 1; i <= numSubscribers; i++ {
		ch := b.subscribe()
		wg.Add(1)
		go func(id int, ch <-chan string) {
			defer wg.Done()
			for msg := range ch {
				fmt.Printf("Подписчик %d получен: %s\n", id, msg)
			}
			fmt.Printf("Подписчик %d сделан.\n", id)
		}(i, ch)
	}

	go func() {
		for i := 1; i <= 5; i++ {
			msg := fmt.Sprintf("Сообщение %d", i)
			fmt.Printf("Публикация отправлена: %s\n", msg)
			b.publish(msg)
			time.Sleep(400 * time.Millisecond)
		}
		b.close()
	}()

	wg.Wait()
	fmt.Println("Публикация/Подписка завершена!")

}

type broadcaster struct {
	subscribers []chan string
	closed      bool
	mu          sync.Mutex
}

func newBroadcaster() *broadcaster {
	return &broadcaster{
		subscribers: make([]chan string, 0),
	}
}

func (b *broadcaster) subscribe() <-chan string {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan string, 2)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *broadcaster) publish(msg string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	for _, ch := range b.subscribers {
		ch <- msg
	}
}

func (b *broadcaster) close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	for _, ch := range b.subscribers {
		close(ch)
	}
	b.closed = true
}
