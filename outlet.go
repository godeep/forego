package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"io"
	"os"
	"sync"
)

type OutletFactory struct {
	Outlets map[string]*Outlet
	Padding int

	sync.Mutex
}

type Outlet struct {
	Name    string
	Color   ct.Color
	IsError bool
	Factory *OutletFactory
}

var colors = []ct.Color{
	ct.Cyan,
	ct.Yellow,
	ct.Green,
	ct.Magenta,
	ct.Red,
	ct.Blue,
}

func NewOutletFactory() (of *OutletFactory) {
	of = new(OutletFactory)
	of.Outlets = make(map[string]*Outlet)
	return
}

func (o *Outlet) Write(b []byte) (num int, err error) {
	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		o.Factory.WriteLine(o.Name, scanner.Text(), ct.White, ct.None, o.IsError)
	}
	num = len(b)
	return
}

func ProcessOutput(w io.Writer, str string) {
	w.Write([]byte(str))
}

func (of *OutletFactory) CreateOutlet(name string, index int, isError bool) *Outlet {
	of.Outlets[name] = &Outlet{name, colors[index%len(colors)], isError, of}
	return of.Outlets[name]
}

func (of *OutletFactory) SystemOutput(str string) {
	of.WriteLine("forego", str, ct.White, ct.None, false)
}

func (of *OutletFactory) ErrorOutput(str string) {
	fmt.Printf("ERROR: %s\n", str)
	os.Exit(1)
}

// Write out a single coloured line
func (of *OutletFactory) WriteLine(left, right string, leftC, rightC ct.Color, isError bool) {
	of.Lock()
	defer of.Unlock()

	ct.ChangeColor(leftC, true, ct.None, false)
	formatter := fmt.Sprintf("%%-%ds | ", of.Padding)
	fmt.Printf(formatter, left)

	if isError {
		ct.ChangeColor(ct.Red, true, ct.None, true)
	} else {
		ct.ResetColor()
	}
	fmt.Println(right)
}
