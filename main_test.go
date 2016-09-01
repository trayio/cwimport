package main

import (
	"fmt"
	"testing"
)

func wrongErr(expectedErr, receivedErr error) string {
	return fmt.Sprintf("Wrong error, expecting %s, got %s", expectedErr, receivedErr)
}

func TestMetricValidation(t *testing.T) {
	// metric.validate() returns on first error and checks keys in order
	// - Query (query)
	// - Asg (asg)
	// - Namespace (namespace)
	// - Unit (unit)
	// - Unit (valid unit name)
	// - Interval (interval)
	var m metric = metric{}
	var noErrMsg = "No error returned while validating incomplete metric"

	if err := m.validate(); err == nil {
		t.Fatal(noErrMsg)
	} else {
		if err != queryMissing {
			t.Fatal(wrongErr(queryMissing, err))
		}
	}

	m.Query = "query"

	if err := m.validate(); err == nil {
		t.Fatal(noErrMsg)
	} else {
		if err != asgMissing {
			t.Fatal(wrongErr(asgMissing, err))
		}
	}

	m.Asg = "asg"

	if err := m.validate(); err == nil {
		t.Fatal(noErrMsg)
	} else {
		if err != namespaceMissing {
			t.Fatal(wrongErr(namespaceMissing, err))
		}
	}

	m.Namespace = "namespace"

	if err := m.validate(); err == nil {
		t.Fatal(noErrMsg)
	} else {
		if err != unitMissing {
			t.Fatal(wrongErr(unitMissing, err))
		}
	}

	m.Unit = "unit"

	if err := m.validate(); err == nil {
		t.Fatal(noErrMsg)
	} else {
		if err != invalidUnit {
			t.Fatal(wrongErr(invalidUnit, err))
		}
	}

	m.Unit = "None"

	m.Interval = 1

	if err := m.validate(); err != nil {
		t.Fatal("Error while validating valid metric")
	}
}

func TestConfigValidation(t *testing.T) {
	// configuration.validate() returns on first error and checks keys in order:
	// - Region (aws_region)
	// - PrometheusUrl (prometheus_url)
	var c configuration = configuration{}
	var noErr = "No error returned while validating incomplete configuration"

	if err := c.validate(); err == nil {
		t.Fatal(noErr)
	} else {
		if err != awsRegionMissing {
			t.Fatal(wrongErr(awsRegionMissing, err))
		}
	}

	c.Region = "region"

	if err := c.validate(); err == nil {
		t.Fatal(noErr)
	} else {
		if err != prometheusUrlMissing {
			t.Fatal(wrongErr(prometheusUrlMissing, err))
		}
	}

	c.PrometheusUrl = "url"

	if err := c.validate(); err != nil {
		t.Fatal("Error while validating valid configuration:", err)
	}
}
