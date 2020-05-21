package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

//recomended configs
/*
3 9 102 40
*/

/*
2 7 102 80
*/

//antenne amount
const TransmittingAntennae = 3

//antenne delay (s)
const SleepTime = 9

//node amount 102 max
const EthernetLength = 102

//ethernet delay (ms)
const EthernetSpeed = 40

const EmptySpace = " "
const JamSignal = "_"

//ew clear
var ewClear = true

//must be synchronized
var ew = make([]string, EthernetLength)

var conflictCtr int32 = 0

var readConfCtr = func() int {
	m := sync.Mutex{}
	m.Lock()
	defer m.Unlock()
	return (int)(conflictCtr)
}

var readEwClear = func() bool {
	m := sync.Mutex{}
	m.Lock()
	defer m.Unlock()
	return ewClear
}

var wg sync.WaitGroup

//antenaes
var antennae = make([]antenna, TransmittingAntennae)

func main() {
	wg.Add(1)
	go func() {
		for {
			select {
			case x := <-startTransmitting:
				fstartTransmitting(x.l, x.r, x.id)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-printPropagationStat:
				fprintPropagationStat()
			}
		}
	}()

	go func() {
		for {
			select {
			case <-clearWire:
				m := sync.Mutex{}
				m.Lock()
				for i := range ew {
					ew[i] = EmptySpace
				}
				m.Unlock()
			}
		}
	}()

	clearWire <- true
	runAntennae()

	for {
		time.Sleep(10 * time.Second)
	}

}

func runAntennae() {
	for i := 0; i < TransmittingAntennae; i++ {
		antennae[i].id = i + 1
		go antennae[i].start()
	}
}

/*
#######################################################
ANTENNA
#######################################################
*/
type antenna struct {
	id       int
	conflict bool
}

func (a *antenna) start() {
	transPosition := (EthernetLength / (TransmittingAntennae + 1)) * a.id
	var leftSrc, rightSrc, msgRange int
	a.conflict = false

	for {
		time.Sleep(time.Duration(rand.Intn(SleepTime)) * time.Second)

		if getEthernetWire()[transPosition] == EmptySpace {

			leftSrc = transPosition
			rightSrc = transPosition

			//send packet
			msgRange = 0
			for msgRange < EthernetLength {
				startTransmitting <- transStr{leftSrc, rightSrc, a.id}
				leftSrc--
				rightSrc++
				msgRange++
				time.Sleep(time.Duration(EthernetSpeed) * time.Millisecond)

				printPropagationStat <- true

				if !(getEthernetWire()[transPosition] == fmt.Sprintf("%d", a.id)) {
					a.conflict = true
					atomic.AddInt32(&conflictCtr, 1)
					break
				}
			}
			if msgRange == EthernetLength {
				a.conflict = false
				clearWire <- true
				atomic.StoreInt32(&conflictCtr, 0)
			}
			if a.conflict {
				for i := 0; i < EthernetLength; i++ {
					startTransmitting <- transStr{leftSrc, rightSrc, a.id}
					leftSrc--
					rightSrc++
					time.Sleep(time.Duration(EthernetSpeed) * time.Millisecond)
					printPropagationStat <- true
				}

				//elongation delay time if conflictCtr will
				if readConfCtr() < 10 {
					rcc := readConfCtr()
					time.Sleep(time.Duration(rand.Intn(pow(2, rcc))) * time.Millisecond)
				} else {
					time.Sleep(time.Duration(rand.Intn(2^10)) * time.Millisecond)
				}
			}
			if ewClear {
				clearWire <- true
			}
			ewClear = !ewClear
		}
	}
}

/*
#######################################################
other functions
#######################################################
*/
func getEthernetWire() []string {
	m := sync.Mutex{}
	m.Lock()
	getEW := ew
	m.Unlock()
	return getEW
}

func setEthernetWire(position int, sign string) {
	m := sync.Mutex{}
	m.Lock()
	ew[position] = sign
	m.Unlock()
}

//start transmi
type transStr struct {
	l, r, id int
}

var startTransmitting = make(chan transStr)

func fstartTransmitting(leftProp, rightProp, id int) {
	m := sync.Mutex{}
	m.Lock()

	if leftProp >= 0 {
		if ew[leftProp] == EmptySpace {
			ew[leftProp] = fmt.Sprintf("%d", id)
		} else {
			ew[leftProp] = JamSignal
		}
	}
	if rightProp < EthernetLength {
		if ew[rightProp] == EmptySpace || ew[rightProp] == fmt.Sprintf("%d", id) {
			ew[rightProp] = fmt.Sprintf("%d", id)
		} else {
			ew[rightProp] = JamSignal
		}
	}

	m.Unlock()
}

//print propagation
var printPropagationStat = make(chan bool)

func fprintPropagationStat() {
	ethW := getEthernetWire()
	for i := 0; i < EthernetLength; i++ {
		fmt.Print(ethW[i], "|")
	}
	fmt.Println()
}

//clear wire
var clearWire = make(chan bool)

//conflict wait

func conflictWait(sleep int, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(sleep) * time.Millisecond)
}

func pow(x, y int) int {
	if y == 0 {
		return 1
	} else {
		return x * pow(x, y-1)
	}
}
