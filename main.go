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
var lostnum int = 0
var finishnum int = 0
var totalWait float64 = 0
var  numCashier int
var  numLessCashier int
var globalmutex sync.Mutex
var globalmutexwaittime sync.Mutex
var consumers []*Consumer

type Cashier struct {
	Id int
	WaitGroup chan int
	Limit int
	lock sync.Mutex
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

func (c *Cashier) setLimit(inVal int) {
	c.Limit = inVal
}
func (c *Cashier) addConsumer(inVal int) {
	c.WaitGroup <- inVal
	// fmt.Println(c.Id,len(c.WaitGroup))
}


func (c Cashier) doCheckout(wg *sync.WaitGroup) {
		for{
			// time.Sleep(time.Millisecond*50)
			if len(c.WaitGroup)>0{
				consumerID := <- c.WaitGroup
				 var cstZone = time.FixedZone("CST", 8*3600)
				// loc, _ := time.LoadLocation("Local")
				TimeNow:=time.Now().In(cstZone).UnixNano()
				left:=float64(TimeNow-consumers[consumerID].BeginTime)
				waitTime:=left/1e9
				globalmutexwaittime.Lock()
				totalWait=totalWait+waitTime
				globalmutexwaittime.Unlock()
				// begintime:=time.Unix(consumers[consumerID].BeginTime, 0).Format("2006-01-02 15:04:05" )
				// TimeBegin, _ := time.ParseInLocation("2006-01-02 15:04:05",begintime,loc)
				// left := TimeNow.Sub(TimeBegin)
				time.Sleep(time.Second * time.Duration(consumers[consumerID].Checkspeed*consumers[consumerID].ProductNumber))
				globalmutexwaittime.Lock()
				finishnum++
				fmt.Printf("checkout compelete consumer %d cashier %d item %d wait time %f totalwait %f finshpeople %d\n",consumerID,c.Id,consumers[consumerID].ProductNumber,waitTime,totalWait,finishnum)
				globalmutexwaittime.Unlock()
				wg.Done()
			}
		}
}
func (c *Consumer) setWaitTime() {
	num :=rand.Intn(60)
	c.WaitTime=num
}
func (c *Consumer) setBeginTime() {
	var cstZone = time.FixedZone("CST", 8*3600)
	time:=time.Now().In(cstZone).UnixNano()
	c.BeginTime=time
}

func (c *Consumer) setProductNumber() {
	num :=rand.Float64();
	// fmt.Println(num)
	// fmt.Println(num>0.2)
	if  num >= 0.2 {
		c.ProductNumber =rand.Intn(194)+6
		// fmt.Println(c.ProductNumber)
	}else{
		c.ProductNumber =rand.Intn(4)+1
		// fmt.Println(c.ProductNumber)
	}
}
func (c *Consumer) setCheckspeed() {
	num:=rand.Intn(100)
	if num>=99{
		c.Checkspeed=6
	}else if num>=95{
		c.Checkspeed=5
	}else if num >=90{
		c.Checkspeed=4
	}else if num>= 85{
		c.Checkspeed=3
	}else if num>=80{
		c.Checkspeed=2
	}else if num>=70{
		c.Checkspeed=1
	}else{
		c.Checkspeed=1
	}
	// c.Checkspeed  =rand.Intn(5) + 1
	// c.Checkspeed = rand.Float64()*5.5+0.5
}
func (c *Consumer) findQueue(cashiers []*Cashier,wg *sync.WaitGroup){
	var num int = -1;
	var cap int = 6 ;
	time.Sleep(time.Second * time.Duration(c.WaitTime))
	if(c.ProductNumber<=5){
		for i := 0; i < numLessCashier; i++ {
			cashiers[i].lock.Lock()
			if len(cashiers[i].WaitGroup)<cap{
				num = i
				cap = len(cashiers[i].WaitGroup)
				// fmt.Println("cashiers",len(cashiers[i].WaitGroup))
			}
			cashiers[i].lock.Unlock()
		}
	}
	
	if(c.ProductNumber>5 || num<0){
		for i := numLessCashier; i < numCashier; i++ {
			cashiers[i].lock.Lock()
			if len(cashiers[i].WaitGroup)<cap{
				num = i
				cap = len(cashiers[i].WaitGroup)
				// fmt.Println("cashiers",len(cashiers[i].WaitGroup))
			}
			cashiers[i].lock.Unlock()
		}
	}
	if num<0 {
		globalmutex.Lock()
		lostnum++
		fmt.Printf("consumer %d leave total:%d\n",c.Id,lostnum)
		globalmutex.Unlock()
		defer wg.Done()
	}else{
		cashiers[num].lock.Lock()
		cashiers[num].addConsumer(c.Id)
		cashiers[num].lock.Unlock()
	}
	// fmt.Println(num)
	// for _, value := range cashier {
		
	// }
	// fmt.Println(c.ProductNumber)
}
// func calculateQueue(cashiers []*Cashier,)
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
	fmt.Println(time.Now())
	fmt.Println("please int number of cashier between 1 and 8：")
	fmt.Scanln(&numCashier)
	fmt.Println("please int number of restriction cashier between 1 and 2:")
	fmt.Scanln(&numLessCashier)
	cashiers := make([]*Cashier, numCashier)
	// var wgcashier sync.WaitGroup
	var wg sync.WaitGroup
	for i := 0; i < numCashier; i++ {
		cashiers[i] = &Cashier{
			Id: i,WaitGroup:make(chan int,6)}
			if i<numLessCashier{
				cashiers[i].setLimit(5)
			}else{
				cashiers[i].setLimit(200)
			}
			// wgcashier.Add(1)
			go cashiers[i].doCheckout(&wg)
	}

	// for i := 0; i < numCashier; i++ {
	// 	locks:=append(locks, make(chan bool, 1))
    //     locks[i] <- true
    // } 
	
	// for i := 0; i < numCashier; i++ {
	// 	fmt.Println(<-locks[i])
	// }
	for i := 0; i < 1000; i++ {
		num :=rand.Intn(5)+1
		
		time.Sleep(time.Millisecond *time.Duration(num))
		consumers = append(consumers,&Consumer{
			Id: i})
			consumers[i].setProductNumber() 
			consumers[i].setCheckspeed()
			consumers[i].setWaitTime()
			consumers[i].setBeginTime()
			wg.Add(1)
			go consumers[i].findQueue(cashiers,&wg)

	}
	wg.Wait()
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(consumers[i].String())
	// }
}