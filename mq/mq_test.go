package mq

import "testing"

var mq *MessageQueue

func forever(t *testing.T) {
	for {
		msg := <-mq.Pop()
		mq.Ack(msg, msg.Data)
	}
}

func TestMQ(t *testing.T) {
	mq = New(1 << 10)
	go forever(t)
	func() {
		for i := 0; i < 100; i++ {
			if i%2 == 0 {
				t.Log("Sync", mq.Push(i))
			} else {
				mq.PushAsync(i, func(r interface{}) {
					t.Log("Async", r.(int))
				})
			}
		}
	}()
	mq.Done()
}
