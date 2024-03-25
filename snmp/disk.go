package snmp

import (
	"fmt"
	"gosnmp/alert"
	"gosnmp/config"
	"sync"
	"time"

	g "github.com/gosnmp/gosnmp"
	log "github.com/sirupsen/logrus"
)

type Disk struct {
	MountPoint   string
	AvailPercent int
	AvailSpace   string
}

const (
	dskPath        string = ".1.3.6.1.4.1.2021.9.1.2" // walk
	dskAvail       string = ".1.3.6.1.4.1.2021.9.1.7" // walk
	dskUsedPercent string = ".1.3.6.1.4.1.2021.9.1.9" // walk
)

func (d *Disk) SetMountPoint(mountpoint string) {
	d.MountPoint = mountpoint
}

func (d *Disk) SetAvailSpace(avaSpc string) {
	d.AvailSpace = avaSpc
}

func (d *Disk) SetAvailPercent(avaPerc int) {
	d.AvailPercent = avaPerc
}

func GetDiskMetrics(x *g.GoSNMP) []*Disk {
	// get disk mountpoint
	pathes, err := x.WalkAll(dskPath)
	if err != nil {
		log.Errorf("walk dskPath %v: %v", x.Target, err)
	}

	disks := make([]*Disk, len(pathes))

	for i := 0; i < len(pathes); i++ {
		disks[i] = &Disk{}
	}

	// set disk mountpoint
	for i, path := range pathes {
		name := path.Value.([]uint8)
		disks[i].SetMountPoint(string(name))
	}

	// get available disk space
	das, err := x.WalkAll(dskAvail)
	if err != nil {
		log.Errorf("walk dskAvail %v: %v", x.Target, err)
	}

	// set available disk space
	for i, da := range das {
		disksize := da.Value.(int)
		disksize = disksize / (1024 * 1024)
		disks[i].SetAvailSpace(fmt.Sprintf("%vG", disksize))
	}

	// get available disk percentage
	ups, err := x.WalkAll(dskUsedPercent)
	if err != nil {
		log.Errorf("walk dskUsedPercent %v: %v", x.Target, err)
	}

	// set available disk percentage
	for i, up := range ups {
		percent := up.Value.(int)
		disks[i].SetAvailPercent(100 - percent)
	}

	return disks
}

func (d *Disk) Alert(host *config.Host, thr *config.Threshold) {
	if d.AvailPercent < thr.Disk {
		log.Infof("%s 磁盘空间还剩 %s %d%%, 发送告警", host.Addr, d.AvailSpace, d.AvailPercent)
		body := fmt.Sprintf("磁盘空间 %s 还剩%d%% %s", d.MountPoint, d.AvailPercent, d.AvailSpace)
		qwMsg := alert.NewMessage()
		qwMsg.Alert(host, body)
	}
}

func diskHandler(x *g.GoSNMP, c *SnmpConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	h := config.GetHostByAddr(c.hosts, x.Target)
	t := time.NewTicker(c.interval.Disk)
	for range t.C {
		disks := GetDiskMetrics(x)
		for _, d := range disks {
			d.Alert(h, c.threshold)
		}
	}
}
