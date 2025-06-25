package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== Паттерн издатель и подписка ===")

	b := newBroadcaster() // Менеджер для управленя подписками
	numSubscribers := 3   //Количество подписчиков
	var wg sync.WaitGroup //Группа ожидания

	for i := 1; i <= numSubscribers; i++ { //Цикл по подписчикам
		ch := b.subscribe()                 //Получает канал с подписчиком
		wg.Add(1)                           //Увеличивает счётчик группы для горутины
		go func(id int, ch <-chan string) { //Эта горутина будет ждать сообщений со стороны издателя
			defer wg.Done() //Уменьшение счётчика
			for msg := range ch {
				fmt.Printf("Подписчик %d получен: %s\n", id, msg)
			}
			fmt.Printf("Подписчик %d сделан.\n", id)
		}(i, ch) //Запуск анонимной функции в виде горутины, где 1-й параметр - идентфикатор подписчика, 2-й параметр - канал с подписчиками
	}

	go func() { //Горутина для отправки сообщений
		for i := 1; i <= 5; i++ {
			msg := fmt.Sprintf("Сообщение %d", i) //Отправляемое сообщение
			fmt.Printf("Публикация отправлена: %s\n", msg)
			b.publish(msg)                     //Опубликовать сообщение
			time.Sleep(400 * time.Millisecond) //Задержка по времени
		}
		b.close()
	}()

	wg.Wait()
	fmt.Println("Издатель/Подписка завершена!")
}

// Управляет публикациями и подпиской
type broadcaster struct {
	subscribers []chan string // Массив каналов подписчиков
	closed      bool          // Закрыто
	mu          sync.Mutex    // Для блокировки доступа
}

// Создать новый объект Broadcaster
func newBroadcaster() *broadcaster {
	return &broadcaster{
		subscribers: make([]chan string, 0), //Канал подписчиков
	}
}

// Подписаться
func (b *broadcaster) subscribe() <-chan string {
	b.mu.Lock()                               // Блокировка
	defer b.mu.Unlock()                       //Разблокировка
	ch := make(chan string, 2)                // Создаёт буферизированный канал из 2-х строк
	b.subscribers = append(b.subscribers, ch) //Добавляет канал к массиву каналов подписчиков
	return ch
}

// Публикация
func (b *broadcaster) publish(msg string) {
	b.mu.Lock()         //Блокировка мьютекса
	defer b.mu.Unlock() //Отложенная разблокировка
	if b.closed {       //Если закрытый статус
		return
	}
	for _, ch := range b.subscribers {
		ch <- msg
	}
}

// Закрытие менеджера
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
