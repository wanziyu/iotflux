package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"

)
var (
	client sarama.SyncProducer
)
func InitProducer(address []string)(err error){
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	// 连接kafka
	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return err
	}
	fmt.Println("Kafka init success")

	return
}

func NewKafkaConsumer(address []string, config *sarama.Config)( consumer sarama.Consumer, err error ){
	consumer, err = sarama.NewConsumer(address, config)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	fmt.Println("Kafka Consumer Build Success")
	return consumer, err
}


func SendMessage(topic, data string){
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
		return
	}
	fmt.Println("Send Message to Kafka Success")
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
}


func ReadMessage(consumer sarama.Consumer, topic string){

	partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println(partitionList)
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetOldest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}

		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Println("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, string(msg.Value))

			}
		}(pc)
	}
	select {

	}
}