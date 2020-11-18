package main

import (
	"fmt"
	"sync"
	"time"
	"bytes"
	 "encoding/json"
	 "math/rand"
	//  "strconv"
)
var lostnum int = 0  //The number of lost customers 
var finishnum int = 0  //the number of finish checkout consumer
var totalWait float64 = 0 //total wait time for all consumer
var totalItem int = 0	//total products for the consumer dont leave
var  numCashier int
var  numLessCashier int
var consumernum int
var globalmutex sync.Mutex	//lock the num of consumer
var globalmutexwaittime sync.Mutex  //lock the total wait time
var globalCashierGroup sync.Mutex //lock the cashier group
var consumers []*Consumer




type Weather struct {
	BeginTime int64
	weatherEffect int
}

type Cashier struct {
	Id int
	WaitGroup chan int
	Limit int
	// lock sync.Mutex
	BeginTime int64
	checkTime int
}
type Consumer struct {
	Id int
	ProductNumber int
	Checkspeed int
	WaitTime int
	BeginTime int64
}
// func mutiply(a float, b int) (int c){
// 	c := int64(decimal.NewFromFloat(num1).Mul(decimal.NewFromFloat(float64(num2))))
// 	return
// }
func (w *Weather) changeWeather() {
	for{
		time.Sleep(time.Millisecond*100)
		// var cstZone = time.FixedZone("CST", 8*3600)
		time:=time.Now().UnixNano()
		gap:=time-w.BeginTime 
		// change the weathereffect
		if gap%(3*1e10)>2*1e10{
			w.setWeather(3)
		}else if gap%(3*1e10)>1e10{
			w.setWeather(2)
		}else{
			w.setWeather(1)
		}
	}
}
func (w *Weather) setWeather(value int){
	w.weatherEffect=value
}
func (c *Cashier) setLimit(inVal int) {
	c.Limit = inVal
}
func (c *Cashier) addConsumer(inVal int) {
	c.WaitGroup <- inVal
	// fmt.Println(c.Id,len(c.WaitGroup))
}
func (c *Cashier) setBeginTime() {
	// var cstZone = time.FixedZone("CST", 8*3600)
	time:=time.Now().UnixNano()
	c.BeginTime=time
}

func (c *Cashier) doCheckout(wg *sync.WaitGroup) {
		for{
			// time.Sleep(time.Millisecond*50)
			if len(c.WaitGroup)>0{
				consumerID := <- c.WaitGroup
				// loc, _ := time.LoadLocation("Local")
				TimeNow:=time.Now().UnixNano()
				left:=float64(TimeNow-consumers[consumerID].BeginTime)
				waitTime:=left/1e9
				globalmutexwaittime.Lock()
				totalWait=totalWait+waitTime
				globalmutexwaittime.Unlock()
				// begintime:=time.Unix(consumers[consumerID].BeginTime, 0).Format("2006-01-02 15:04:05" )
				// TimeBegin, _ := time.ParseInLocation("2006-01-02 15:04:05",begintime,loc)
				// left := TimeNow.Sub(TimeBegin)
				consumerchecktime:=consumers[consumerID].Checkspeed*consumers[consumerID].ProductNumber
				c.checkTime=c.checkTime+consumerchecktime
				time.Sleep(time.Second * time.Duration(consumerchecktime))
				globalmutexwaittime.Lock()
				finishnum++
				fmt.Printf("checkout compelete consumer %d cashier %d item %d wait time %f  finshpeople %d\n",consumerID,c.Id,consumers[consumerID].ProductNumber,waitTime,finishnum)
				globalmutexwaittime.Unlock()
				wg.Done()
			}
		}
}
func (c *Consumer) setWaitTime() {
	num :=rand.Intn(20)
	c.WaitTime=num
}
func (c *Consumer) setBeginTime() {
	// var cstZone = time.FixedZone("CST", 8*3600)
	time:=time.Now().UnixNano()
	c.BeginTime=time
}

//pseudorandom generate product number
func (c *Consumer) setProductNumber() {
	num :=rand.Float64();
	// fmt.Println(num>0.2)
	if  num >= 0.2 {
		c.ProductNumber =rand.Intn(194)+6
		test:=rand.Intn(100)
		if test>=99{
			c.ProductNumber=rand.Intn(194)+6
		}else if num>=98{
			c.ProductNumber=rand.Intn(174)+6
		}else if num >=97{
			c.ProductNumber=rand.Intn(154)+6
		}else if num>= 95{
			c.ProductNumber=rand.Intn(134)+6
		}else if num>=90{
			c.ProductNumber=rand.Intn(114)+6
		}else if num>=50{
			c.ProductNumber=rand.Intn(70)+6
		}else{
			c.ProductNumber=rand.Intn(20)+6
		}
	}else{
		c.ProductNumber =rand.Intn(4)+1
		// fmt.Println(c.ProductNumber)
	}
}

//pseudorandom generate Checkspeed
func (c *Consumer) setCheckspeed() {
	num:=rand.Intn(100)
	if num>=99{
		c.Checkspeed=6
	}else if num>=98{
		c.Checkspeed=5
	}else if num >=97{
		c.Checkspeed=4
	}else if num>= 95{
		c.Checkspeed=3
	}else if num>=90{
		c.Checkspeed=2
	}else{
		c.Checkspeed=1
	}
	// c.Checkspeed  =rand.Intn(5) + 1
	// c.Checkspeed = rand.Float64()*5.5+0.5
}
func (c *Consumer) findQueue(cashiers []*Cashier,w *Weather,wg *sync.WaitGroup){
	count :=rand.Intn(100)	
	// time.Sleep(time.Millisecond *time.Duration(count))	
	time.Sleep(time.Second *time.Duration(count))//set some time for the consumer before they enter the supermarket
	var num int = -1;
	var cap int = 6 ;
	time.Sleep(time.Second * time.Duration(c.WaitTime*w.weatherEffect))//time arrive at the checkouts 
	// fmt.Println("WEATHER",w.weatherEffect)
	globalCashierGroup.Lock()
	if(c.ProductNumber<=5){
		for i := 0; i < numLessCashier; i++ {
			if len(cashiers[i].WaitGroup)<cap{
				num = i
				cap = len(cashiers[i].WaitGroup)
				// fmt.Println("cashiers",len(cashiers[i].WaitGroup))
			}
		}
	}
	
	if(c.ProductNumber>5 || num<0){
		for i := numLessCashier; i < numCashier; i++ {
			if len(cashiers[i].WaitGroup)<cap{
				num = i
				cap = len(cashiers[i].WaitGroup)
				// fmt.Println("cashiers",len(cashiers[i].WaitGroup))
			}
		}
	}
	if num<0 {
		globalmutex.Lock()
		lostnum++
		fmt.Printf("consumer %d leave total:%d\n",c.Id,lostnum)
		globalmutex.Unlock()
		wg.Done()
	}else{
		c.setBeginTime()
		globalmutex.Lock()
		totalItem=totalItem+c.ProductNumber
		cashiers[num].addConsumer(c.Id)
		globalmutex.Unlock()
	}
	globalCashierGroup.Unlock()
	// fmt.Println(num)
	// for _, value := range cashier {
		
	// }
	// fmt.Println(c.ProductNumber)
}

func (conf *Consumer) String() string {
    b, err := json.Marshal(*conf)
    if err != nil {
        return fmt.Sprintf("%+v", *conf)
    }
    var out bytes.Buffer
    err = json.Indent(&out, b, "", "    ")
    if err != nil {
        return fmt.Sprintf("%+v", *conf)
    }
    return out.String()
}

func main()  {
	// fmt.Println(time.Now())
	fmt.Println("please input number of cashier between 1 and 8：")
	fmt.Scanln(&numCashier)
	fmt.Println("please input number of restriction cashier between 1 and 2:")
	fmt.Scanln(&numLessCashier)
	fmt.Println("please input number of consumer:")
	fmt.Scanln(&consumernum)
	// var cstZone = time.FixedZone("CST", 8*3600)
	now:=time.Now().UnixNano()
	randomweather := Weather{BeginTime:now,weatherEffect:1}
	go randomweather.changeWeather()//start weather agent
	cashiers := make([]*Cashier, numCashier)
	// var wgcashier sync.WaitGroup
	var wg sync.WaitGroup
	for i := 0; i < numCashier; i++ {
		cashiers[i] = &Cashier{
			Id: i,WaitGroup:make(chan int,6),checkTime:0.0}
			if i<numLessCashier{
				cashiers[i].setLimit(5)
			}else{
				cashiers[i].setLimit(200)
			}
			cashiers[i].setBeginTime()
			// wgcashier.Add(1)
			go cashiers[i].doCheckout(&wg)//start cashier agent
	}

	// for i := 0; i < numCashier; i++ {
	// 	locks:=append(locks, make(chan bool, 1))
    //     locks[i] <- true
    // } 
	
	// for i := 0; i < numCashier; i++ {
	// 	fmt.Println(<-locks[i])
	// }
	for i := 0; i < consumernum; i++ {
		consumers = append(consumers,&Consumer{
			Id: i})
			consumers[i].setProductNumber() 
			consumers[i].setCheckspeed()
			consumers[i].setWaitTime()
			// consumers[i].setBeginTime()
			wg.Add(1)
			go consumers[i].findQueue(cashiers,&randomweather,&wg)

	}
	wg.Wait()
	var averageutilization int = 0
	TimeNow:=time.Now().UnixNano()
	fmt.Println("simulation compelete")
	fmt.Printf("total products processed %d average products per trolley %f\n",totalItem,float64(totalItem)/float64(finishnum))
	fmt.Printf("average customer wait time %f\n",totalWait/float64(finishnum))
	totalCashierTime :=0.0
	for i := 0; i < numCashier; i++ {
		left:=float64(TimeNow-cashiers[i].BeginTime)/1e9
		totalCashierTime=totalCashierTime+left
		averageutilization=averageutilization+cashiers[i].checkTime
		fmt.Printf("Cashier%d utilization %f\n",i,float64(cashiers[i].checkTime)*100.0/left)
	}
	fmt.Printf("average utilization %f\n",float64(averageutilization*100.0)/totalCashierTime)
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(consumers[i].String())
	// }
}