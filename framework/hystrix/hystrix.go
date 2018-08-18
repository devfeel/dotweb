package hystrix

import (
	"time"
	"sync"
)

const(
	status_Hystrix = 1
	status_Alive = 2
	DefaultCheckHystrixInterval = 10 //unit is Second
	DefaultCheckAliveInterval = 60 //unit is Second
	DefaultCleanHistoryInterval = 60 * 5 //unit is Second
	DefaultMaxFailedNumber = 100
	DefaultReserveMinutes = 30

)

type Hystrix interface{
	// Do begin do check
	Do()
	// RegisterAliveCheck register check Alive func
	RegisterAliveCheck(CheckFunc)
	// RegisterHystrixCheck register check Hystrix func
	RegisterHystrixCheck(CheckFunc)
	// IsHystrix return is Hystrix status
	IsHystrix() bool
	// TriggerHystrix trigger Hystrix status
	TriggerHystrix()
	// TriggerAlive trigger Alive status
	TriggerAlive()
	// SetCheckInterval set interval for doCheckHystric and doCheckAlive, unit is Second
	SetCheckInterval(int, int)

	// GetCounter get lasted Counter with time key
	GetCounter() Counter

	// SetMaxFailed set max failed count for hystrix default counter
	SetMaxFailedNumber(int64)
}

type CheckFunc func()bool

type StandHystrix struct{
	status int
	checkHystrixFunc CheckFunc
	checkHystrixInterval int
	checkAliveFunc CheckFunc
	checkAliveInterval int

	maxFailedNumber int64
	counters *sync.Map
}


// NewHystrix create new Hystrix, config with CheckAliveFunc and checkAliveInterval, unit is Minute
func NewHystrix(checkAlive CheckFunc, checkHysrix CheckFunc) Hystrix{
	h := &StandHystrix{
		counters : new(sync.Map),
		status:status_Alive,
		checkAliveFunc: checkAlive,
		checkHystrixFunc:checkHysrix,
		checkAliveInterval:DefaultCheckAliveInterval,
		checkHystrixInterval:DefaultCheckHystrixInterval,
		maxFailedNumber:DefaultMaxFailedNumber,
	}
	if h.checkHystrixFunc == nil{
		h.checkHystrixFunc = h.defaultCheckHystrix
	}
	return h
}

func (h *StandHystrix) Do(){
	go h.doCheck()
	go h.doCleanHistoryCounter()
}

func (h *StandHystrix) SetCheckInterval(hystrixInterval, aliveInterval int){
	h.checkAliveInterval = aliveInterval
	h.checkHystrixInterval = hystrixInterval
}

// SetMaxFailed set max failed count for hystrix default counter
func (h *StandHystrix) SetMaxFailedNumber(number int64){
	h.maxFailedNumber = number
}

// GetCounter get lasted Counter with time key
func (h *StandHystrix) GetCounter() Counter{
	key := getLastedTimeKey()
	var counter Counter
	loadCounter, exists := h.counters.Load(key)
	if !exists{
		counter = NewCounter()
		h.counters.Store(key, counter)
	}else{
		counter = loadCounter.(Counter)
	}
	return counter
}


func (h *StandHystrix) IsHystrix() bool{
	return h.status == status_Hystrix
}

func (h *StandHystrix) RegisterAliveCheck(check CheckFunc){
	h.checkAliveFunc = check
}

func (h *StandHystrix) RegisterHystrixCheck(check CheckFunc){
	h.checkHystrixFunc = check
}

func (h *StandHystrix) TriggerHystrix(){
	h.status = status_Hystrix
}

func (h *StandHystrix) TriggerAlive(){
	h.status = status_Alive
}


// doCheck do checkAlive when status is Hystrix or checkHytrix when status is Alive
func (h *StandHystrix) doCheck(){
	if h.checkAliveFunc == nil || h.checkHystrixFunc == nil {
		return
	}
	if h.IsHystrix() {
		isAlive := h.checkAliveFunc()
		if isAlive {
			h.TriggerAlive()
			h.GetCounter().Clear()
			time.AfterFunc(time.Duration(h.checkHystrixInterval)*time.Second, h.doCheck)
		} else {
			time.AfterFunc(time.Duration(h.checkAliveInterval)*time.Second, h.doCheck)
		}
	}else{
		isHystrix := h.checkHystrixFunc()
		if isHystrix{
			h.TriggerHystrix()
			time.AfterFunc(time.Duration(h.checkAliveInterval)*time.Second, h.doCheck)
		}else{
			time.AfterFunc(time.Duration(h.checkHystrixInterval)*time.Second, h.doCheck)
		}
	}
}

func (h *StandHystrix) doCleanHistoryCounter(){
	var needRemoveKey []string
	now, _ := time.Parse(minuteTimeLayout, time.Now().Format(minuteTimeLayout))
	h.counters.Range(func(k, v interface{}) bool{
		key := k.(string)
		if t, err := time.Parse(minuteTimeLayout, key); err != nil {
			needRemoveKey = append(needRemoveKey, key)
		} else {
			if now.Sub(t) > (DefaultReserveMinutes * time.Minute) {
				needRemoveKey = append(needRemoveKey, key)
			}
		}
		return true
	})
	for _, k := range needRemoveKey {
		//fmt.Println(time.Now(), "hystrix doCleanHistoryCounter remove key",k)
		h.counters.Delete(k)
	}
	time.AfterFunc(time.Duration(DefaultCleanHistoryInterval)*time.Second, h.doCleanHistoryCounter)
}


func (h *StandHystrix) defaultCheckHystrix() bool{
	count := h.GetCounter().Count()
	if count > h.maxFailedNumber{
		return true
	}else{
		return false
	}
}

func getLastedTimeKey() string{
	key :=  time.Now().Format(minuteTimeLayout)
	if time.Now().Minute() / 2 != 0{
		key = time.Now().Add(time.Duration(-1*time.Minute)).Format(minuteTimeLayout)
	}
	return key
}