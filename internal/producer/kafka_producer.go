package producer

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log-tracer/internal/config"
	"log-tracer/internal/pkg/logger"
	"os"

	"github.com/IBM/sarama"
)

// KafkaProducer is a wrapper around Sarama's SyncProducer to manage the Kafka producer connection
type KafkaProducer struct {
	Producer sarama.SyncProducer
	Config   *config.ProducerConfig
}

// NewKafkaProducer initializes a new KafkaProducer with the provided configuration
func NewKafkaProducer(cfg *config.ProducerConfig) (*KafkaProducer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.RequiredAcks(cfg.KafkaProducerConfig.Acks)
	saramaConfig.Producer.Retry.Max = cfg.KafkaProducerConfig.Retries
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Compression = sarama.CompressionGZIP
	saramaConfig.Producer.Partitioner = NewCustomPartitioner

	// Setup SSL/TLS if enabled
	if cfg.KafkaProducerConfig.SSL.Enabled {
		tlsConfig, err := createTLSConfiguration(cfg.KafkaProducerConfig.SSL)
		if err != nil {
			logger.Error("Failed to create TLS configuration", err)
			return nil, err
		}
		saramaConfig.Net.TLS.Enable = true
		saramaConfig.Net.TLS.Config = tlsConfig
	}

	producer, err := sarama.NewSyncProducer(cfg.KafkaProducerConfig.Brokers, saramaConfig)
	if err != nil {
		logger.Error("Failed to create Kafka producer", err)
		return nil, err
	}

	return &KafkaProducer{
		Producer: producer,
		Config:   cfg,
	}, nil
}

// createTLSConfiguration creates a TLS configuration based on the provided SSL configuration
func createTLSConfiguration(sslConfig config.SSLConfig) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: sslConfig.InsecureSkipVerify,
	}

	if sslConfig.CACert != "" {
		caCert, err := os.ReadFile(sslConfig.CACert)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	if sslConfig.ClientCert != "" && sslConfig.ClientKey != "" {
		clientCert, err := tls.LoadX509KeyPair(sslConfig.ClientCert, sslConfig.ClientKey)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{clientCert}
	}

	return tlsConfig, nil
}

// SendMessage sends a message to the specified Kafka topic
func (p *KafkaProducer) SendMessage(ctx context.Context, key, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: p.Config.KafkaProducerConfig.Topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(message),
	}

	// Send the message and log the result
	partition, offset, err := p.Producer.SendMessage(msg)
	if err != nil {
		logger.Error("Failed to send message to Kafka broker:", err)
		return err
	}

	logger.Info("Message sent", "partition", partition, "offset", offset)
	return nil
}

// Close closes the Kafka producer connection
func (p *KafkaProducer) Close() error {
	return p.Producer.Close()
}
