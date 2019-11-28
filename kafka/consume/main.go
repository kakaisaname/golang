package consume

import (
	"fmt"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"log"
	"os"
	"os/signal"
	"sync"
)

func main() {
	topic := []string{"topic"}
	Address :=  []string{"192.168.190.161:9092","192.168.190.161:9093","192.168.190.161:9094"}
	var wg = &sync.WaitGroup{}
	wg.Add(2)
	//广播式消费：消费者1																					可以开多个协程去消费  *****
	go clusterConsumer(wg, Address, topic, "test-consumer-group")
	//广播式消费：消费者2
	go clusterConsumer(wg, Address, topic, "test-consumer-group")
	wg.Wait()
}


//用这个方法  **
func clusterConsumer(wg *sync.WaitGroup,brokers,topics []string,groupId string)  {
	defer wg.Done()
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// init consumer
	consumer, err := cluster.NewConsumer(brokers, groupId, topics, config)
	if err != nil {
		log.Printf("%s: sarama.NewSyncProducer err, message=%s \n", groupId, err)
		return
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("%s:Error: %s\n", groupId, err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("%s:Rebalanced: %+v \n", groupId, ntf)
		}
	}()

	// consume messages, watch signals
	var successes int
Loop:
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				fmt.Fprintf(os.Stdout, "%s:%s/%d/%d\t%s\t%s\n", groupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "")  // mark message as processed
				successes++
			}
		case <-signals:
			break Loop
		}
	}
	fmt.Fprintf(os.Stdout, "%s consume %d messages \n", groupId, successes)
}


//这个只能消费（获取）单分区的数据
func singleConsumer()  {
	topic := []string{"topic"}
	Address :=  []string{"192.168.190.161:9092","192.168.190.161:9093","192.168.190.161:9094"}
	var wg = &sync.WaitGroup{}
	//广播式消费：消费者1
	clusterConsumer(wg, Address, topic, "test-consumer-group")
	//广播式消费：消费者2
	//go clusterConsumer(wg, Address, topic, "group-2")
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V0_11_0_2

	// consumer
	consumer, err := sarama.NewConsumer(Address, config)
	if err != nil {
		fmt.Printf("consumer_test create consumer error %s\n", err.Error())
		return
	}
	defer consumer.Close()
	partitionList,err := consumer.Partitions("topic")
	if err != nil {
		fmt.Printf("create partitionList error %s\n", err.Error())
		return
	}
	for partition := range partitionList {
		fmt.Println("for进入")
		pc, errRet := consumer.ConsumePartition("topic",int32(partition),sarama.OffsetNewest)
		if errRet != nil {
			fmt.Printf("Failed to start consumer for partition %d: %s\n",partition,err)
			return
		}
		defer pc.AsyncClose()
		wg.Add(1)
		go func(pc sarama.PartitionConsumer){
			fmt.Println("进来了")
			for msg := range pc.Messages() {
				fmt.Println("func执行")
				if err != nil {
					fmt.Printf("msg %d: %s\n",partition,err)
				}
				fmt.Printf("获取到msg:%s ",msg.Value)
			}
			wg.Done()
		}(pc)
	}
	wg.Wait()
}
