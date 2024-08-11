package producer

import (
	"hash/fnv"

	"github.com/IBM/sarama"
)

// customPartitioner implements sarama.Partitioner interface
type customPartitioner struct{}

func (p *customPartitioner) Partition(msg *sarama.ProducerMessage, numPartitions int32) (int32, error) {
	h := fnv.New32a()
	if msg.Key != nil {
		_, _ = h.Write(msg.Key.(sarama.ByteEncoder))
	}
	return int32(h.Sum32()) % numPartitions, nil
}

func (p *customPartitioner) RequiresConsistency() bool {
	return true
}

// NewCustomPartitioner is a PartitionerConstructor function that returns a new custom partitioner
func NewCustomPartitioner(topic string) sarama.Partitioner {
	return &customPartitioner{}
}
