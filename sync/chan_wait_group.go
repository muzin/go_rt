package sync

type ChanWaitGroup struct {
	channel chan bool

	count int
}

func (this *ChanWaitGroup) init() {
	this.channel = make(chan bool)
}

func NewChanWaitGroup() *ChanWaitGroup {
	c := &ChanWaitGroup{}
	c.init()
	return c
}

func (this *ChanWaitGroup) Add(n int) {
	this.count += n
}

func (this *ChanWaitGroup) Done() {
	this.channel <- true
}

func (this *ChanWaitGroup) Wait() {
	for i := 0; i < this.count; i++ {
		<-this.channel
	}
}

func (this *ChanWaitGroup) Close() {
	close(this.channel)
}
