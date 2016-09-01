package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/net/context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/prometheus/client_golang/api/prometheus"
	"github.com/prometheus/common/model"
)

type Collector interface {
	Collect(string) []float64
}

type PrometheusCollector struct {
	q prometheus.QueryAPI
}

func (p *PrometheusCollector) Collect(query string) []float64 {
	var collection []float64

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	value, err := p.q.Query(ctx, query, time.Now().UTC())
	if err != nil {
		fmt.Printf("Error executing query '%s': %s", query, err)
		return collection
	}

	// https://godoc.org/github.com/prometheus/common/model
	switch t := value.(type) {
	case model.Vector:
		for _, v := range t {
			collection = append(collection, float64(v.Value))
		}

	default:
		fmt.Println("Unknown type:", t)
	}

	return collection
}

func NewPrometheusCollector(url string) (*PrometheusCollector, error) {
	var (
		err    error
		client prometheus.Client
	)

	client, err = prometheus.New(prometheus.Config{Address: url})
	if err != nil {
		return nil, err
	}

	return &PrometheusCollector{
		q: prometheus.NewQueryAPI(client),
	}, nil
}

// https://docs.aws.amazon.com/sdk-for-go/api/service/cloudwatch/CloudWatch.html#PutMetricData-instance_method
func (m metric) Run(c Collector, ch chan<- *cloudwatch.PutMetricDataInput) {
	tick := time.NewTicker(time.Duration(m.Interval) * time.Minute)
	defer m.wg.Done()

	for {
		select {
		case <-tick.C:
			var metricData []*cloudwatch.MetricDatum

			for _, value := range c.Collect(m.Query) {
				metric := &cloudwatch.MetricDatum{
					Value:      aws.Float64(value),
					MetricName: aws.String(m.name),
					Dimensions: []*cloudwatch.Dimension{
						{
							Name:  aws.String("AutoScalingGroupName"),
							Value: aws.String(m.Asg),
						},
					},
					Timestamp: aws.Time(time.Now().UTC()),
					Unit:      aws.String(m.Unit),
				}

				metricData = append(metricData, metric)
			}

			ch <- &cloudwatch.PutMetricDataInput{MetricData: metricData, Namespace: aws.String(m.Namespace)}

		case <-m.quitChan:
			fmt.Printf("%s ", m.name)
			return
		}
	}
}

func NewCloudWatchClient(region string) (*cloudwatch.CloudWatch, error) {
	s := session.New(&aws.Config{
		Region: aws.String(region),
	})

	if _, err := s.Config.Credentials.Get(); err != nil {
		return nil, err
	}

	return cloudwatch.New(s), nil
}

func main() {
	var (
		wg         sync.WaitGroup
		conf       *configuration
		err        error
		pc         *PrometheusCollector
		sc         chan os.Signal                      // signal channel
		qr         chan struct{}                       // quit runner channel
		qm         chan struct{}                       // quit main channel
		mc         chan *cloudwatch.PutMetricDataInput // metric channel
		configFile string
		testOnly   bool
	)

	flag.StringVar(&configFile, "config", "config.hcl", "Configuration file")
	flag.BoolVar(&testOnly, "t", false, "Test configuration and exit")
	flag.Parse()

	conf, err = NewConfig(configFile)
	if err != nil {
		fmt.Println("Configuration error:", err)
		os.Exit(1)
	}

	if testOnly {
		fmt.Printf("Configuration %s OK\n", configFile)
		os.Exit(0)
	}

	fmt.Println("Creating prometheus collector with url:", conf.PrometheusUrl)
	pc, err = NewPrometheusCollector(conf.PrometheusUrl)
	if err != nil {
		fmt.Println("Failed creating prometheus collector:", err)
		os.Exit(1)
	}

	fmt.Println("Creating CloudWatch client")
	cw, err := NewCloudWatchClient(conf.Region)
	if err != nil {
		fmt.Println("Failed to get AWS credentials:", err)
		os.Exit(1)
	}

	// signal catcher
	sc = make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT)

	// quit runner channel
	qr = make(chan struct{})

	// quit main channel
	qm = make(chan struct{})

	// metrics channel
	mc = make(chan *cloudwatch.PutMetricDataInput)

	for name, metric := range conf.Metrics {
		fmt.Println("Launching collector for:", name)
		wg.Add(1)

		metric.quitChan = qr
		metric.name = name
		metric.wg = &wg
		go metric.Run(pc, mc)
	}

	for {
		select {
		case metricData := <-mc:
			_, err := cw.PutMetricData(metricData)
			if err != nil {
				fmt.Println(err.Error())
			}

		case <-sc:
			fmt.Printf("Shutting down: ")
			close(qr)

			// need to wait for workers in a goroutine in case some still have data
			// to send over mc channel which we'd block while waiting
			go func() {
				wg.Wait()
				close(qm)
			}()

		case <-qm:
			fmt.Printf("\nDone\n")
			os.Exit(0)
		}
	}
}
