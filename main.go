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
var  numCashier int
var  numLessCashier int
var globalmutex sync.Mutex
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


func (c Cashier) doCheckout() {
		for{
			time.Sleep(time.Millisecond*50)
			if len(c.WaitGroup)>0{
				consumerID := <- c.WaitGroup
				// fmt.Println(consumers[consumerID].Checkspeed)
				time.Sleep(time.Millisecond * time.Duration(consumers[consumerID].Checkspeed*consumers[consumerID].ProductNumber))
				fmt.Printf("checkout compelete consumer %d cashier %d item %d\n",consumerID,c.Id,consumers[consumerID].ProductNumber)
			}
		}
}
func (c *Consumer) setWaitTime() {
	num :=rand.Intn(60)
	c.WaitTime=num
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
	c.Checkspeed  =rand.Intn(5) + 1
	// c.Checkspeed = rand.Float64()*5.5+0.5
}
func (c *Consumer) findQueue(cashiers []*Cashier,wg *sync.WaitGroup){
	defer wg.Done()
	var num int = -1;
	var cap int = 6 ;
	time.Sleep(time.Millisecond *10* time.Duration(c.WaitTime))
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
	fmt.Println("please int number of cashier between 1 and 8：")
	fmt.Scanln(&numCashier)
	fmt.Println("please int number of restriction cashier between 1 and 2:")
	fmt.Scanln(&numLessCashier)
	cashiers := make([]*Cashier, numCashier)
	var wgcashier sync.WaitGroup
	for i := 0; i < numCashier; i++ {
		cashiers[i] = &Cashier{
			Id: i,WaitGroup:make(chan int,6)}
			if i<numLessCashier{
				cashiers[i].setLimit(5)
			}else{
				cashiers[i].setLimit(200)
			}
			wgcashier.Add(1)
			go cashiers[i].doCheckout()
	}

	// for i := 0; i < numCashier; i++ {
	// 	locks:=append(locks, make(chan bool, 1))
    //     locks[i] <- true
    // } 
	
	// for i := 0; i < numCashier; i++ {
	// 	fmt.Println(<-locks[i])
	// }
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		time.Sleep(time.Millisecond *10)
		consumers = append(consumers,&Consumer{
			Id: i})
			consumers[i].setProductNumber() 
			consumers[i].setCheckspeed()
			consumers[i].setWaitTime()
			wg.Add(1)
			go consumers[i].findQueue(cashiers,&wg)

	}
	wg.Wait()
	wgcashier.Wait()
	// for i := 0; i < 10; i++ {
	// 	fmt.Println(consumers[i].String())
	// }
}