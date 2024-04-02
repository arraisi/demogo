package metrics

import (
	"demogo/pkg/logger"
	"github.com/DataDog/datadog-go/v5/statsd"
	"os"
	"time"
)

var m *metrics

type Client interface {
	Close() error
	Histogram(name string, value float64, tags []string, rate float64) error
	Gauge(name string, value float64, tags []string, rate float64) error
	Count(name string, value int64, tags []string, rate float64) error
	Event(e *statsd.Event) error
	Distribution(name string, value float64, tags []string, rate float64) error
	Timing(name string, value time.Duration, tags []string, rate float64) error
}

type metrics struct {
	client Client
}

func initStatsd(addr, serviceName string) error {
	client, err := statsd.New(
		addr,
		statsd.WithNamespace(serviceName),
	)
	if err != nil {
		return err
	}

	m = &metrics{
		client: client,
	}
	return nil
}

// Init by default using datadog statsd
// if we want to integrate another client such as new relic
// we can modify this function
func Init(addr, serviceName string) error {
	err := initStatsd(addr, serviceName)
	return err
}

// InitWithClient used when invoker already define their client
func InitWithClient(client Client) {
	m = &metrics{
		client: client,
	}
}

func Close() {
	if err := m.client.Close(); err != nil {
		logger.Log.Errorf("[dogstatsd.close] cannot close connection: %v", err)
	}
}

func Event(priority, title, text string, tags []string) {
	dPriority := statsd.Success
	switch priority {
	case "error":
		dPriority = statsd.Error
	case "warning":
		dPriority = statsd.Warning
	case "info":
		dPriority = statsd.Info
	}

	host, _ := os.Hostname()
	if host == "" {
		host = "localhost"
	}

	event := statsd.Event{
		Title:     title,
		Text:      text,
		Timestamp: time.Now(),
		Hostname:  host,
		Priority:  statsd.Normal,
		AlertType: dPriority,
		Tags:      tags,
	}

	if err := m.client.Event(&event); err != nil {
		logger.Log.Errorf("[event] cannot fire event (%v): %v", event, err)
	}
}

func Gauge(name string, value float64, tags []string) {
	if err := m.client.Gauge(name, value, tags, 1); err != nil {
		logger.Log.Errorf("[count] cannot send metric %v", err)
	}
}

func Count(name string, value int64, tags []string) {
	if err := m.client.Count(name, value, tags, 1); err != nil {
		logger.Log.Errorf("[count] cannot send metric %v", err)
	}
}

func Histogram(name string, value float64, tags []string) {
	if err := m.client.Histogram(name, value, tags, 1); err != nil {
		logger.Log.Errorf("[histogram] cannot send metric %v", err)
	}
}

func Distribution(name string, value float64, tags []string) {
	if err := m.client.Distribution(name, value, tags, 1); err != nil {
		logger.Log.Errorf("[distribution] cannot send metric %v", err)
	}
}

func Timing(name string, value time.Duration, tags []string) {
	if err := m.client.Timing(name, value, tags, 1); err != nil {
		logger.Log.Errorf("[timing] cannot send metric %v", err)
	}
}
