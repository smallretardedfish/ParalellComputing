package main

import (
	"context"
	"github.com/EbicHecker/queue"
	"log"

	"math/rand"
	"sync"
	"time"
)

//Програма моделює обслуговування двох потоків процесів з різними параметрами одним центральним процесором і трьома чергами.
//Якщо згенеровано процес другого потоку і процесор занятий, процес надходить в другу чергу.
//Якщо згенеровано процес першого потоку, то, якщо процесор обробляє процес другого потоку,
//процес другого потоку призупиняється і поміщається в третю чергу (з часом, що залишився на обробку).
//Якщо процесор обробляє процес першого потоку, то процес надходить в першу чергу.
//При звільнені процесора черги проглядаються в наступному порядку: перша черга, третя черга, друга черга.
//Визначити максимальну довжину черг і кількість перерваних процесів другого потоку.

const (
	Gen1Processes = 23
	Gen2Processes = 10
)

var (
	MaxProcessTime = (500 * time.Millisecond).Nanoseconds()
)

type Stats struct {
	MaxQueue1Length                  int64
	MaxQueue2Length                  int64
	MaxQueue3Length                  int64
	SecondStreamProcessesInterrupted int64
	Mu                               sync.Mutex
}

type CPU struct {
	*sync.Mutex
	IsBusy         bool
	CurrentProcess *Process

	Q1 *Queue
	Q2 *Queue
	Q3 *Queue
}

func (c *CPU) Run(ctx context.Context) {
	for {
		if p, err := c.Q1.q.TryDequeue(); err != nil {
			c.RunProcess(ctx, p)
		}
		if p, err := c.Q3.q.TryDequeue(); err != nil {
			ctx, cancel := context.WithCancel(ctx)
			p.Cancel = cancel
			c.RunProcess(ctx, p)
		}
		if p, err := c.Q2.q.TryDequeue(); err != nil {
			ctx, cancel := context.WithCancel(ctx)
			p.Cancel = cancel
			c.RunProcess(ctx, p)
		}
	}
}

func (c *CPU) RunProcess(ctx context.Context, process *Process) bool {
	c.Lock()
	c.IsBusy = true
	c.CurrentProcess = process
	c.Unlock()
	process.Run(ctx) // running a process

	c.Lock()
	c.IsBusy = false
	c.Unlock()
	return true
}

func (c *CPU) GetCurrentProcess() (*Process, bool) {
	return c.CurrentProcess, c.IsBusy
}

type Scheduler struct {
	CPU   *CPU
	Stats *Stats

	Gen1 *Generator
	Gen2 *Generator
}

func (s *Scheduler) Run() {
	gen1Ch := make(chan Process)
	gen2Ch := make(chan Process)

	go func() {
		for i := 0; i < Gen1Processes; i++ {
			gen1Ch <- s.Gen1.Generate()
		}
		close(gen1Ch)
	}()
	go func() {
		for i := 0; i < Gen2Processes; i++ {
			gen2Ch <- s.Gen2.Generate()
		}
		close(gen2Ch)
	}()

}

func (s *Scheduler) schedule(g1ch, g2ch chan Process) {

	for g1ch != nil || g2ch != nil {
		proc, busy := s.CPU.GetCurrentProcess()
		if proc == nil {
			continue
		}
		select {
		case p, ok := <-g1ch:
			if !ok { // no more processes (chan was closed)
				g1ch = nil
			}
			// Якщо процесор обробляє процес першого потоку, то процес надходить в першу чергу.
			if busy && proc.GenID == 1 {
				s.CPU.Q1.q.Enqueue(p)
				continue
			}
			//Якщо згенеровано процес першого потоку, то, якщо процесор обробляє процес другого потоку,
			//процес другого потоку призупиняється і поміщається в третю чергу (з часом, що залишився на обробку).
			if busy && proc.GenID == 2 {
				proc.Cancel()
				s.CPU.Q1.q.Enqueue(p)
				s.CPU.Q3.q.Enqueue(*proc)
				continue
			}
		case p, ok := <-g2ch:

			if !ok { // no more processes (chan was closed)
				g1ch = nil
			}
			//Якщо згенеровано процес другого потоку і процесор занятий, процес надходить в другу чергу.
			if busy && proc.GenID == 2 {
				s.CPU.Q2.q.Enqueue(p)
				continue
			}
		}
	}
}

type Process struct {
	ID           int64
	GenID        int64
	LeftoverTime time.Duration
	Cancel       context.CancelFunc
}

//Run method making kind of job for a process during some time.
//In case of interrupting, time which is leftover for a process is written to LeftOverTime field.
func (p *Process) Run(ctx context.Context) { //TODO check case when process is interrupted and processor is finished simultaneously
	start := time.Now()
	select {
	case <-time.After(p.LeftoverTime):
		log.Println("Process", p.ID, "of Gen", p.GenID, " is finished")
		return
	case <-ctx.Done():
		p.LeftoverTime = p.LeftoverTime - time.Since(start)
		return
	}
}

type Generator struct {
	ID         int64
	Mu         sync.Mutex
	counter    int64
	MaxGenTime int64
	MinGenTime int64
}

func (g *Generator) Generate() Process {
	g.Mu.Lock()
	defer g.Mu.Unlock()
	g.counter++

	randSource := rand.NewSource(time.Now().UnixNano())
	r := *rand.New(randSource)
	timeToWait := (r.Int() % (int(g.MaxGenTime) - int(g.MinGenTime))) + int(g.MinGenTime)
	time.Sleep(time.Millisecond * time.Duration(timeToWait)) // kinda of delay for generator
	return Process{
		ID:    g.counter,
		GenID: g.ID,
	} // generating a process
}

type Queue struct {
	q *queue.ConcurrentQueue[Process]
}

func main() {

	Gen1 := Generator{
		ID:      1,
		Mu:      sync.Mutex{},
		counter: 0,
	}
	Gen2 := Generator{
		ID:      2,
		Mu:      sync.Mutex{},
		counter: 0,
	}

	cpu := CPU{
		Mutex:  &sync.Mutex{},
		IsBusy: false,
		Q1: &Queue{
			q: queue.NewConcurrentQueue[Process](),
		},
		Q2: &Queue{
			q: queue.NewConcurrentQueue[Process]()},
		Q3: &Queue{
			q: queue.NewConcurrentQueue[Process](),
		},
		CurrentProcess: nil,
	}

	stats := Stats{
		MaxQueue1Length: 0, // max length of Q1

	}

	scheduler := Scheduler{
		CPU:   &cpu,
		Stats: &stats,
		Gen1:  &Gen1,
		Gen2:  &Gen2,
	}
	scheduler.Run()

	cpu.Run(context.Background())
}
