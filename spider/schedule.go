package spider

var schedule chan *pageTask

func init() {
	schedule = make(chan *pageTask, 1)
}

func RunSchedule() {
	for {
		select {
		case p := <-schedule:
			p.Run(3)
		}
	}
}
