package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"sync"
	"time"
)

var Address = []string{"192.168.190.161:9092","192.168.190.161:9093","192.168.190.161:9094"}

//同步产生数据
func syncProducer(address []string,wg *sync.WaitGroup)  {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewSyncProducer(address, config)
	if err != nil {
		log.Printf("sarama.NewSyncProducer err, message=%s \n", err)
		return
	}
	defer p.Close()
	topic := "topic"
	srcValue := "sync: this is a message. index=%d"
	for i:=0; i<10; i++ {
		wg.Add(1)
		value := fmt.Sprintf(srcValue, i)
		msg := &sarama.ProducerMessage{
			Topic:topic,
			Value:sarama.ByteEncoder(value),
		}
		part, offset, err := p.SendMessage(msg)															//随机的发送数据到不同的分区
		if err != nil {
			log.Printf("send message(%s) err=%s \n", value, err)
		}else {
			fmt.Fprintf(os.Stdout, value + "发送成功，partition=%d, offset=%d \n", part, offset)
		}
		wg.Done()
	}
	wg.Wait()
}

func main() {
	//syncProducer(Address,&wg)
	config := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//随机向partition发送消息
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	//注意，版本设置不对的话，kafka会返回很奇怪的错误，并且无法成功发送消息
	config.Version = sarama.V0_10_0_1

	fmt.Println("start make producer")
	//使用配置,新建一个异步生产者
	producer, e := sarama.NewAsyncProducer(Address, config)
	if e != nil {
		fmt.Println(e)
		return
	}
	defer producer.AsyncClose()
	exit := make(chan bool)
	asyncProducer(producer,exit)
	fmt.Println("发送完了1")
	xiaofei(producer)
	fmt.Println("发送完了2")
	<-exit
	fmt.Println("发送完了3")
}

func xiaofei(producer sarama.AsyncProducer)  {
	for{
		select {
		//case suc,ok := (<-producer.Successes():			//这种方式 ok 一直为true
		case suc:= (<-producer.Successes()):
			fmt.Println("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.String(), "partitions: ", suc.Partition)
		case <-time.After(time.Second*3):					//这里采用过期时间来，普通的通道在发送端关闭，这边是可以接收到的，这里主要是不同的通道 	***
			goto out
		case fail := <-producer.Errors():
			fmt.Println("err: ",fail.Err)
			goto out
		}
	}
	out:
		fmt.Println("end")
}

//异步产生消息
func asyncProducer(producer sarama.AsyncProducer,exit chan bool)  {
	fmt.Println("start goroutine")
	var value string
	go func() {
		for i :=0;i<100;i++ {
			//time.Sleep(500 * time.Millisecond)
			time1 := time.Now()
			value = "成功 message"+time1.Format("15:04:05")
			// 发送的消息,主题。
			// 注意：这里的msg必须得是新构建的变量，不然你会发现发送过去的消息内容都是一样的，因为批次发送消息的关系。
			msg := &sarama.ProducerMessage{
				Topic:     "topic",
				Key:       sarama.StringEncoder("成功"),
			}
			//将字符串转化为字节数组
			msg.Value = sarama.ByteEncoder(value)
			producer.Input() <- msg												//将消息传入通道
		}
		exit<-true
		//close(producer.Input())															//发送完了,但是关闭这个没用
	}()
}
