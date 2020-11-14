package bus

type Publisher interface {
	Publish(event interface{})
	Subscribe(c Consumer)
}

type Consumer interface {
	Consume(event interface{})
}
