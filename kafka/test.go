package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	//"github.com/bsm/sarama-cluster"
	"sync"
)

var value string

func create()  {
	fmt.Printf("producer_test\n")
	config := sarama.NewConfig()
	// 等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 随机的分区类型：返回一个分区器，该分区器每次选择一个随机分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 是否等待成功和失败后的响应
	config.Producer.Return.Successes = true

	// 使用给定代理地址和配置创建一个同步生产者
	producer, err := sarama.NewSyncProducer([]string{"192.168.190.161:9092"}, config)
	if err != nil {
		panic(err)
	}
	defer producer.Close()
	//构建发送的消息，
	msg := &sarama.ProducerMessage {
		//Topic: "topic",//包含了消息的主题     分区就为topic
		Partition: int32(1),//
		Key:        sarama.StringEncoder("key"),//
	}
	for {
		_, err := fmt.Scanf("%s", &value)													//用户输入的数据
		if err != nil {
			break
		}
		//fmt.Scanf("%s",&msgType)
		//fmt.Println("msgType = ",msgType,",value = ",value)
		msg.Topic = "topic"																   //向这个主题下的所有分区队列加数据
		//将字符串转换为字节数组
		msg.Value = sarama.ByteEncoder(value)
		//fmt.Println(value)
		//SendMessage：该方法是生产者生产给定的消息
		//生产成功的时候返回该消息的分区和所在的偏移量
		//生产失败的时候返回error
		partition, offset, err := producer.SendMessage(msg)								 //可能会产生每个分区数据有几条不同步的情况

		if err != nil {
			fmt.Println("Send message Fail")
		}
		fmt.Printf("Partition = %d, offset=%d\n", partition, offset)
	}
}

var wg sync.WaitGroup
func consume()  {
	// 根据给定的代理地址和配置创建一个消费者
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		panic(err)
	}

	//Partitions(topic):该方法返回了该topic的所有分区id
	partitionList, err := consumer.Partitions("topic")
	if err != nil {
		panic(err)
	}
	//fmt.Println(partitionList)

	for partition := range partitionList {
		//ConsumePartition方法根据主题，分区和给定的偏移量创建创建了相应的分区消费者
		//如果该分区消费者已经消费了该信息将会返回error
		//sarama.OffsetNewest:表明了为最新消息
		pc, err := consumer.ConsumePartition("topic", int32(partition), sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}
		defer pc.AsyncClose()
		wg.Add(1)
		go func(sarama.PartitionConsumer) {
			defer wg.Done()
			//Messages()该方法返回一个消费消息类型的只读通道，由代理产生
			for msg := range pc.Messages() {
				//fmt.Println(35)
				fmt.Printf("%s---Partition:%d, Offset:%d, Key:%s, Value:%s\n", msg.Topic,msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
		}(pc)
	}
	wg.Wait()
	consumer.Close()
}

