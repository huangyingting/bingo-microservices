package ws

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type RealTimeStats struct {
	MemUsed    uint64  `json:"mem_used"`
	MemPercent float64 `json:"mem_percent"`
	CPUPercent float64 `json:"cpu_percent"`
	Uptime     uint64  `json:"uptime"`
}

var (
	// poll stats with this period.
	STAT_PERIOD = 10 * time.Second
)

// Monitor maintains the set of active clients and broadcasts messages to the
// clients.
type Monitor struct {
	cancel chan struct{}
	wg     sync.WaitGroup
	// registered clients
	clients map[*Client]bool
	// register requests from the clients.
	register chan *Client
	// unregister requests from clients.
	unregister chan *Client
	h          *log.Helper
}

func NewMonitor(h *log.Helper) *Monitor {
	return &Monitor{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		h:          h,
	}
}

// Start starts up the load goroutine.
func (m *Monitor) Start() {
	m.cancel = make(chan struct{})
	m.wg.Add(1)
	go m.run()
}

func (m *Monitor) Stop() {
	close(m.cancel)
	m.wg.Wait()
}

func (m *Monitor) run() {
	defer m.wg.Done()
	statTicker := time.NewTicker(STAT_PERIOD)
	defer func() {
		statTicker.Stop()
	}()
	for {
		select {
		case <-m.cancel:
			return
		case client := <-m.register:
			m.h.Debug("client registered")
			m.clients[client] = true
		case client := <-m.unregister:
			m.h.Debug("client unregistered")
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				close(client.send)
			}
		case <-statTicker.C:
			host, _ := host.Info()
			mem, _ := mem.VirtualMemory()
			cpu_percent, _ := cpu.Percent(1000*time.Millisecond, false)
			realTimeStats := RealTimeStats{
				MemUsed:    mem.Used,
				MemPercent: mem.UsedPercent,
				CPUPercent: cpu_percent[0],
				Uptime:     host.Uptime,
			}
			message, _ := json.Marshal(realTimeStats)
			if message != nil {
				for client := range m.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(m.clients, client)
					}
				}
			}
		}
	}
}

// Serve handles websocket requests from the peer.
func (m *Monitor) Serve(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := &Client{m: m, conn: conn, send: make(chan []byte, MAX_MESSAGE_SIZE)}
	client.m.register <- client

	// allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
