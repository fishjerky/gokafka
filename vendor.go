package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/davecgh/go-spew"
	"github.com/eapache/go-resiliency"
	"github.com/eapache/go-xerial-snappy"
	"github.com/eapache/queue"
	"github.com/pierrec/lz4"
	"github.com/rcrowley/go-metrics"
)

func main() {
	fmt.Println("vim-go")
	config := sarama.NewConfig()
}
