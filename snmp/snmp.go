package snmp

import (
	"gosnmp/config"
	"sync"
	"time"

	g "github.com/gosnmp/gosnmp"
	log "github.com/sirupsen/logrus"
)

type SnmpConfig struct {
	hosts     []*config.Host
	threshold *config.Threshold
	interval  *config.Interval
}

func NewSnmpConfig(hosts []*config.Host, thr *config.Threshold, intv *config.Interval) *SnmpConfig {
	return &SnmpConfig{
		hosts:     hosts,
		threshold: thr,
		interval:  intv,
	}
}

func New(ipaddr string) *g.GoSNMP {
	return &g.GoSNMP{
		Port:               161,
		Transport:          "udp",
		Community:          "public",
		Version:            g.Version2c,
		Timeout:            time.Duration(2) * time.Second,
		Retries:            3,
		ExponentialTimeout: true,
		MaxOids:            g.MaxOids,
		Target:             ipaddr,
	}
}

func Run(x *g.GoSNMP, c *SnmpConfig) {
	log.Infof("start scrape metrics from %s\n", x.Target)
	var snmpWg = new(sync.WaitGroup)
	snmpWg.Add(3)
	// 1. handle disk
	go diskHandler(x, c, snmpWg)
	// TODO 2. handle memory
	// TODO 3. handle cpu
	snmpWg.Wait()
}
