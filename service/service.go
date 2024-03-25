package service

import (
	"gosnmp/config"
	"gosnmp/snmp"
	"os"
	"os/signal"
	"syscall"

	g "github.com/gosnmp/gosnmp"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	conns   []*g.GoSNMP
	snmpCfg *snmp.SnmpConfig
}

func New(cfgFile string) (*Service, error) {
	c, err := config.Load(cfgFile)
	if err != nil {
		return nil, err
	}

	conns := make([]*g.GoSNMP, 0)
	for _, h := range c.Hosts {
		conns = append(conns, snmp.New(h.Addr))
	}

	snmpCfg := snmp.NewSnmpConfig(
		c.Hosts,
		c.Threshold,
		c.Interval,
	)

	return &Service{
		conns:   conns,
		snmpCfg: snmpCfg,
	}, nil
}

func (s *Service) Run() error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	for _, x := range s.conns {
		if err := x.Connect(); err == nil {
			go snmp.Run(x, s.snmpCfg)
		} else {
			log.Errorf("connect %s snmp failed: %v\n", x.Target, err)
		}
	}

	sig := <-sigCh
	log.Infof("Caught signal %v, exiting...", sig)

	return nil
}
