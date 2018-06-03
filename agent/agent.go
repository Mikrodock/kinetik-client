package agent

import (
	"bytes"
	"encoding/json"
	"kinetik-client/tools"
	"kinetik-server/logger"
	"net/http"
	"os"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/subosito/gotenv"
)

type CPUStats []float64

type ServerReport struct {
	*load.AvgStat
	MemUsedPercent float64         `json:"mem_used_percent,omitempty"`
	MemUsedBytes   uint64          `json:"mem_used_bytes,omitempty"`
	CPUUsedPercent float64         `json:"cpu_used_percent,omitempty"`
	CPUCount       int             `json:"cpu_count,omitempty"`
	DiskUsage      *disk.UsageStat `json:"disk_usage,omitempty"`
}

var urlServer string
var myIP string

var netClient = &http.Client{
	Timeout: time.Second * 10,
}

func init() {
	err := gotenv.Load("/root/.env")
	if err != nil {
		panic(err)
	}
	urlServer = os.Getenv("KINETIK_MASTER")
	myIP, _ = tools.GetIPAddress()
}

func StartSampling() {
	logger.StdLog.Println("Started sampling")
	for {
		time.Sleep(5 * time.Second)
		logger.StdLog.Println("Tick sampling")
		go SampleAndSend()
	}
}

func SampleAndSend() {
	smpl := Sample()
	jsonBytes, _ := json.Marshal(smpl)
	buffer := bytes.NewBuffer(jsonBytes)
	_, err := netClient.Post("http://"+urlServer+"/nodes/"+myIP, "application/json", buffer)
	if err != nil {
		logger.ErrLog.Println(err.Error())
	}
}

func Sample() *ServerReport {
	srvRep := &ServerReport{}
	srvRep.AvgStat, _ = load.Avg()
	srvRep.DiskUsage, _ = disk.Usage("/")
	virtMem, _ := mem.VirtualMemory()
	srvRep.MemUsedPercent = virtMem.UsedPercent
	srvRep.MemUsedBytes = virtMem.Used
	cpuPercents, _ := cpu.Percent(time.Second*3, false)
	cpuCount, _ := cpu.Counts(true)
	cpuPercent := cpuPercents[0] * float64(cpuCount)
	srvRep.CPUUsedPercent = cpuPercent
	srvRep.CPUCount = cpuCount

	return srvRep
}
