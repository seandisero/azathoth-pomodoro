package azathoth

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const AzathothUpdateTime = 100 * time.Millisecond

type AzState int

const (
	WORK AzState = iota
	REST
	PAUSE
)

type Azathoth struct {
	state            AzState
	prevState        AzState
	time             time.Time
	ticker           *time.Ticker
	showMilliseconds bool
	nextInterval     chan struct{}
	stop             chan struct{}
	pause            bool
	resume           bool
	pauseTime        time.Time
	workInterval     time.Time
	restInterval     time.Time
	intervalCount    int
}

func NewAzathoth(options ...AzOption) *Azathoth {
	az := &Azathoth{
		state:         WORK,
		prevState:     WORK,
		time:          time.Time{},
		ticker:        time.NewTicker(AzathothUpdateTime),
		nextInterval:  make(chan struct{}),
		stop:          make(chan struct{}),
		intervalCount: 0,
	}
	for _, option := range options {
		option(az)
	}
	return az
}

func (a *Azathoth) countDown(t time.Time) {
	a.ticker.Reset(AzathothUpdateTime)
	a.time = t

	if a.time.IsZero() {
		a.nextInterval <- struct{}{}
		return
	}

outer:
	for {
		select {
		case <-a.stop:
			a.state = PAUSE
			break outer
		case <-a.ticker.C:
			if a.time.IsZero() {
				break outer
			}
			a.time = a.time.Add(-1 * AzathothUpdateTime)
			a.printAzathoth()
		}
	}
	a.nextInterval <- struct{}{}
}

func (a *Azathoth) inputHandler() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if input == "" {
			a.stop <- struct{}{}
		}
	}
}

func (a *Azathoth) workTime() time.Time {
	if a.resume {
		a.resume = false
		return a.time
	} else {
		return a.workInterval
	}
}

func (a *Azathoth) restTime() time.Time {
	if a.resume {
		a.resume = false
		return a.time
	} else {
		return a.restInterval
	}
}

func (a *Azathoth) Start() {
	go handleTerm()
	DisableCursor()
	go a.inputHandler()
	for {
		switch a.state {
		case WORK:
			go a.countDown(a.workTime())
			<-a.nextInterval
			if a.state == WORK {
				a.intervalCount++
				a.prevState = WORK
				a.state = REST
			}
		case REST:
			go a.countDown(a.restTime())
			<-a.nextInterval
			if a.state == REST {
				a.prevState = REST
				a.state = WORK
			}
		case PAUSE:
			<-a.stop
			a.pause = false
			a.resume = true
			a.state = a.prevState
		}
	}
}

func (a *Azathoth) shouldAlert() bool {
	if a.state == WORK {
		return a.time.Before(time.Time{}.Add(5 * time.Second))
	}
	return false
}

func milliseconds(t time.Time) float64 {
	conversion := 0.0000001
	return float64(t.Nanosecond()) * conversion
}

func (a *Azathoth) printTime(text_color, bg_color AzColor) {
	formattedTime := fmt.Sprintf("~ %d:%d", a.time.Minute(), a.time.Second())
	if a.showMilliseconds {
		formattedTime += fmt.Sprintf(":%02g", milliseconds(a.time))
	}
	formattedTime += " ~"
	fmt.Printf("%s%s%s%s", text_color, bg_color, formattedTime, NC)
}

func (a *Azathoth) printTimer() {
	text_color := NC
	background_color := BG_NC
	if a.state == WORK {
		text_color = T_GREEN
	} else {
		text_color = T_RED
	}
	if a.shouldAlert() {
		a.printTime(NC, BG_RED)
	} else {
		a.printTime(text_color, background_color)
	}
}

func (a *Azathoth) printAzathoth() {
	fmt.Printf("\033[2J\033[H")
	intervalNum := ": "
	for i := 0; i < a.intervalCount; i++ {
		intervalNum += "• "
	}
	intervalNum += " :\n"
	PrintWithColor(" - - - azathoth - - - \n", T_GREEN, BG_NC)
	PrintWithColor(intervalNum, T_RED, BG_NC)
	a.printTimer()
}

func handleTerm() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
	EnableCursor()
	fmt.Println()
	os.Exit(0)
}
