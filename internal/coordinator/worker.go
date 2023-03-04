package coordinator

import "time"

func (c *Coordinator) startWorker() {
	go c.start10SecondWorker()
	go c.startMinuteWorker()
}

func (c *Coordinator) start10SecondWorker() {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		go c.pullServers()
	}
}

func (c *Coordinator) startMinuteWorker() {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		go c.syncMetrics()
		go c.pushServers()
	}
}
