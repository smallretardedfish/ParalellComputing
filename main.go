package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

type proc struct {
	LeftoverTime time.Duration
}

func (p *proc) Run(ctx context.Context) { //TODO check case when process is interrupted and processor is finished simultaneously
	start := time.Now()
	select {
	case <-time.After(p.LeftoverTime):
		log.Println("Process finished")
		return
	case <-ctx.Done():
		p.LeftoverTime = p.LeftoverTime - time.Since(start)
		fmt.Println("Process interrupted, time left:", p.LeftoverTime)
		return
	}
}
func main() {
	//ctx, cancel := context.WithCancel(context.Background())
	//go func() {
	//	signCh := make(chan os.Signal, 1)
	//	signal.Notify(signCh, os.Interrupt)
	//	for {
	//		sig := <-signCh
	//		switch sig {
	//		case os.Interrupt:
	//			cancel()
	//			return
	//		}
	//	}
	//}()
	//proc := &proc{LeftoverTime: 2 * time.Second}
	//proc.Run(ctx)
}
