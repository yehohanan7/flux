package cqrs

import (
	"flag"
	"os"
	"testing"

	"github.com/golang/glog"
)

func TestMain(m *testing.M) {
	flag.Parse()
	glog.Info("before test")
	result := m.Run()
	glog.Info("after test")
	os.Exit(result)
}
