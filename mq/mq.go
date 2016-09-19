package mq

import "sync"

type Message struct {
	response chan interface{}
	callback func(interface{})
	Async    bool
	Data     interface{}
}

type MessageQueue struct {
	queue chan Message
	Size  uint32
	wg    sync.WaitGroup
}

func New(size uint32) *MessageQueue {
	return &MessageQueue{
		queue: make(chan Message, size),
		Size:  size,
	}
}

func (self *MessageQueue) Push(data interface{}) interface{} {
	self.wg.Add(1)
	msg := Message{Data: data, response: make(chan interface{})}
	self.queue <- msg
	return <-msg.response
}

func (self *MessageQueue) PushAsync(data interface{}, callback func(interface{})) {
	self.wg.Add(1)
	msg := Message{Data: data, Async: true, callback: callback}
	self.queue <- msg
}

func (self *MessageQueue) Pop() <-chan Message {
	return self.queue
}

func (self *MessageQueue) Ack(msg Message, r interface{}) {
	defer self.wg.Done()
	if msg.Async {
		if msg.callback != nil {
			msg.callback(r)
		}
	} else {
		msg.response <- r
	}
}

func (self *MessageQueue) Done() {
	self.wg.Wait()
}
