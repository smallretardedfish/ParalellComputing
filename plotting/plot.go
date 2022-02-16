package plotting

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"sync"
	"time"
)

type Stat struct {
	Size int
	T    time.Duration
}

type Storage struct {
	Mtx  sync.RWMutex
	Dict map[string][]Stat
}

func (s *Storage) Add(count string, stat Stat) {
	s.Mtx.Lock()
	if _, ok := s.Dict[count]; !ok {
		s.Dict[count] = make([]Stat, 0)
	}
	s.Dict[count] = append(s.Dict[count], stat)
	s.Mtx.Unlock()
}

func (s *Storage) String() string {
	return fmt.Sprintf("%v", s.Dict)
}

func NewStorage() *Storage {
	return &Storage{
		Mtx:  sync.RWMutex{},
		Dict: make(map[string][]Stat),
	}
}
func CreatePlot(storage *Storage) {

	p := plot.New()

	p.Title.Text = "cool stat of my garbage program"
	p.X.Label.Text = "size of matrix"
	p.Y.Label.Text = "rime (ns)"
	plotItems := make([]interface{}, 0)
	for key := range storage.Dict {
		plotItems = append(plotItems, key, plotPoints(storage, key))
	}
	err := plotutil.AddLinePoints(p, plotItems...)
	if err != nil {
		panic(err)
	}
	// Save the plot to a PNG file.
	if err := p.Save(10*vg.Inch, 10*vg.Inch, "plot.png"); err != nil {
		panic(err)
	}
}

// plotPoints  size, time for  certain number of goroutines.
func plotPoints(storage *Storage, numOfThreads string) plotter.XYs {
	pts := make(plotter.XYs, len(storage.Dict[numOfThreads]))
	for i, item := range storage.Dict[numOfThreads] {
		pts[i].X = float64(item.Size)
		pts[i].Y = float64(item.T)
	}
	return pts
}
