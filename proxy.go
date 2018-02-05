package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"net/http"
)

//kafka event struct
type Event struct {
	Topic, Message string
}

type ApiResponse struct {
	Success bool
}

var listenPort int
var kafkaIp string

func main() {
	flag.IntVar(&listenPort, "port", 8000, "Proxy hosting port")
	flag.StringVar(&kafkaIp, "kafka", "172.17.35.199:9092", "Kafka ip")
	flag.Parse()

	fmt.Println("=== Kafka proxy setting ===")
	fmt.Printf("Proxy config [Port]=%v, [Kafka ip]:%v\n", listenPort, kafkaIp)

	listen := fmt.Sprintf(":%v", listenPort)

	http.HandleFunc("/goKafka", handleRequest)
	http.ListenAndServe(listen, nil)
	fmt.Println("Ready to handle http request")
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if r := recover(); r != nil {
			response := fmt.Sprintf(
				`{
																																					"Success":false, 
																																										"Response": { 
																																																	"msg": "Parameter error, should be {"Topic":"test", "Message": "test"}
																																																							"devMsg": "%s"}`, r)

			httpResponse(w, response)
		}
	}()

	//1. input
	//jsonString := `{"Topic":"Search", "Message":"test1"}`
	//reader := strings.NewReader(jsonString)
	reader := r.Body

	decoder := json.NewDecoder(reader)
	var e Event
	err := decoder.Decode(&e)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Receive event Topic:%s Message:%s\n", e.Topic, e.Message)

	defer r.Body.Close()

	//2. handle
	sendKafka(e)

	//3. response
	response := fmt.Sprintf(`{"Success":true, "Response":{ "Topic":"%s", "Message":"%s"}}`, e.Topic, e.Message)
	httpResponse(w, response)
}

func httpResponse(w http.ResponseWriter, value string) {

	/*stringb := new(bytes.Buffer)
	body, err := json.NewEncoder(w).Encode(value)
			if err != nil {
							panic(err)
									}
	*/
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(value))
}

func sendKafka(event Event) {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true //bin/config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{kafkaIp}, config)
	if err != nil {
		panic(err)
	}

	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic:     event.Topic,
		Partition: int32(-1),
		Key:       sarama.StringEncoder("key"),
	}

	value := event.Message
	msg.Value = sarama.ByteEncoder(value)

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Send message Fail")
	}
	fmt.Printf("Partition = %d, offset=%d\n", partition, offset)
}
