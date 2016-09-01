package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/hashicorp/hcl"
)

type configuration struct {
	PrometheusUrl string  `hcl:"prometheus_url"`
	Metrics       metrics `hcl:"metrics"`
	Region        string  `hcl:"aws_region"`
}

var awsUnits = []string{
	"Seconds",
	"Microseconds",
	"Milliseconds",
	"Bytes",
	"Kilobytes",
	"Megabytes",
	"Gigabytes",
	"Terabytes",
	"Bits",
	"Kilobits",
	"Megabits",
	"Gigabits",
	"Terabits",
	"Percent",
	"Count",
	"Bytes/Second",
	"Kilobytes/Second",
	"Megabytes/Second",
	"Gigabytes/Second",
	"Terabytes/Second",
	"Bits/Second",
	"Kilobits/Second",
	"Megabits/Second",
	"Gigabits/Second",
	"Terabits/Second",
	"Count/Second",
	"None",
}

type metrics map[string]metric

type metric struct {
	Query     string `hcl:"query"`
	Asg       string `hcl:"asg"`
	Namespace string `hcl:"namespace"`
	Unit      string `hcl:"unit"`
	Interval  int

	quitChan <-chan struct{}
	wg       *sync.WaitGroup
	name     string
}

// error definitions should be put here
var (
	queryMissing         = errors.New("query missing")
	asgMissing           = errors.New("asg missing")
	namespaceMissing     = errors.New("namespace missing")
	unitMissing          = errors.New("unit missing")
	invalidUnit          = errors.New("invalid unit")
	intervalMissing      = errors.New("interval missing or has a value of 0")
	awsRegionMissing     = errors.New("aws_region missing")
	prometheusUrlMissing = errors.New("prometheus_url missing")
)

func (c configuration) validate() error {
	if c.Region == "" {
		return awsRegionMissing
	}

	if c.PrometheusUrl == "" {
		return prometheusUrlMissing
	}

	return nil
}

func (m metric) validate() error {
	if m.Query == "" {
		return queryMissing
	}

	if m.Asg == "" {
		return asgMissing
	}

	if m.Namespace == "" {
		return namespaceMissing
	}

	if m.Unit == "" {
		return unitMissing
	}

	if !contains(m.Unit, awsUnits) {
		return invalidUnit
	}

	if m.Interval == 0 {
		return intervalMissing
	}

	return nil
}

func NewConfig(filename string) (*configuration, error) {
	var (
		conf *configuration
		data []byte
		err  error
	)

	data, err = ioutil.ReadFile(filename)
	if err != nil {
		return conf, err
	}

	if err := hcl.Decode(&conf, string(data)); err != nil {
		return conf, err
	}

	if err := conf.validate(); err != nil {
		return conf, err
	}

	// simple check if all are metrics
	for m, _ := range conf.Metrics {
		var x metric = conf.Metrics[m]
		if err := x.validate(); err != nil {
			return conf, errors.New(fmt.Sprintf("%s: %s", m, err))
		}
	}

	return conf, err
}

func contains(item string, items []string) bool {
	for _, i := range items {
		if item == i {
			return true
		}
	}
	return false
}
