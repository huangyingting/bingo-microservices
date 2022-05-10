package alias

import (
	"errors"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"bingo/app/bs/internal/conf"
	"bingo/app/bs/internal/lock"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jxskiss/base62"
	"go.etcd.io/etcd/client/v3/concurrency"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const MSALIAS_TIME_FORMAT string = "2006-01-01"
const MSALIAS_MAX_BIT_LEN uint32 = 63
const MSALIAS_TIME_UNIT = 1e7 // nsec, i.e. 10 msec

var ErrBitLenOverflow = errors.New("(time + sequence + machine id) bit len > 63")
var ErrTimeInFuture = errors.New("start time cann't be in the future")
var ErrNoAvailableMachineId = errors.New("no available machine id")
var ErrOverTimeLimit = errors.New("over time limit")

// MsAlias is a distributed unique ID generator.
type MsAlias struct {
	mutex       *sync.Mutex
	c           *conf.Alias
	h           *log.Helper
	l           lock.IDistributedLock
	machineId   uint16
	elapsedTime int64
	sequence    uint16
}

func toMsAliasTime(t time.Time) int64 {
	return t.UTC().UnixNano() / MSALIAS_TIME_UNIT
}

func currentElapsedTime(startTime time.Time) int64 {
	return toMsAliasTime(time.Now()) - toMsAliasTime(startTime)
}

func sleepTime(overtime int64) time.Duration {
	return time.Duration(overtime)*10*time.Millisecond -
		time.Duration(time.Now().UTC().UnixNano()%MSALIAS_TIME_UNIT)*time.Nanosecond
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func NewMsAlias(c *conf.Alias, h *log.Helper) (*MsAlias, error) {
	alias := new(MsAlias)
	alias.mutex = new(sync.Mutex)
	alias.c = c
	alias.h = h
	alias.l = nil
	if l, _ := lock.NewEtcdDistributedLock(c.EtcdAddr, h); l != nil {
		alias.l = l
	}

	if c.BitLenTime+c.BitLenSequence+c.BitLenMachineId > MSALIAS_MAX_BIT_LEN {
		return nil, ErrBitLenOverflow
	}

	startTime := alias.c.StartTime.AsTime()
	if startTime.After(time.Now()) {
		return nil, ErrTimeInFuture
	}

	if startTime.IsZero() {
		alias.c.StartTime = timestamppb.New(
			time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC),
		)
	}

	// Generate a randomized machine id if etcd is not enabled
	if alias.l == nil {
		rand.Seed(time.Now().UTC().UnixNano())
		alias.machineId = uint16(rand.Intn(1 << c.BitLenMachineId))
		return alias, nil
	}

	h.Debug("use etcd distributed lock to obtain a machine id")
	var machineId uint16 = 0
	for ; machineId < (1 << c.BitLenMachineId); machineId++ {
		err := alias.l.Lock("/machine/" + strconv.Itoa(int(machineId)))
		if err == nil {
			h.Debugf("machine id %d is obtained", machineId)
			alias.machineId = machineId
			break
		}
		if err == concurrency.ErrLocked {
			h.Debugf("machine id %d is not available", machineId)
		}
	}
	if machineId == (1 << c.BitLenMachineId) {
		h.Debug("no available machine id")
		return nil, ErrNoAvailableMachineId
	}
	return alias, nil
}

// nextId generates a unique id.
func (alias *MsAlias) nextId() (uint64, error) {
	var maskSequence = uint16(1<<alias.c.BitLenSequence - 1)

	alias.mutex.Lock()
	defer alias.mutex.Unlock()

	current := currentElapsedTime(alias.c.StartTime.AsTime())
	if alias.elapsedTime < current {
		alias.elapsedTime = current
		alias.sequence = 0
	} else { // sf.elapsedTime >= current
		alias.sequence = (alias.sequence + 1) & maskSequence
		if alias.sequence == 0 {
			alias.elapsedTime++
			overtime := alias.elapsedTime - current
			time.Sleep(sleepTime((overtime)))
		}
	}

	if alias.elapsedTime >= 1<<alias.c.BitLenTime {
		alias.h.Error("over time limit")
		return 0, ErrOverTimeLimit
	}

	return uint64(alias.elapsedTime)<<(alias.c.BitLenSequence+alias.c.BitLenMachineId) |
		uint64(alias.sequence)<<alias.c.BitLenMachineId |
		uint64(alias.machineId), nil
}

// Next generates next alias
func (alias *MsAlias) Next() (string, error) {
	id, err := alias.nextId()
	if err == nil {
		return reverse(string(base62.FormatUint(id))), nil
	}
	return "", err
}

// Validate validates if alias is base62 encoded
func (a *MsAlias) Validate(alias string) bool {
	_, err := base62.ParseUint([]byte(reverse(alias)))
	if err != nil {
		return false
	}
	return true
}

// decompose returns a set of MsAlias Id parts.
func (alias *MsAlias) decompose(id uint64) map[string]uint64 {
	var maskSequence = uint64((1<<alias.c.BitLenSequence - 1) << alias.c.BitLenMachineId)
	var maskMachineId = uint64(1<<alias.c.BitLenMachineId - 1)

	msb := id >> (alias.c.BitLenTime + alias.c.BitLenSequence + alias.c.BitLenMachineId)
	time := id >> (alias.c.BitLenSequence + alias.c.BitLenMachineId)
	sequence := id & maskSequence >> alias.c.BitLenMachineId
	machineId := id & maskMachineId

	alias.h.Debugf("decompose - id: %d msb: %d time: %d sequence: %d machine-id: %d",
		id, msb, time, sequence, machineId)

	return map[string]uint64{
		"id":         id,
		"msb":        msb,
		"time":       time,
		"sequence":   sequence,
		"machine-id": machineId,
	}
}

func (alias *MsAlias) validate(id uint64) bool {
	id_map := alias.decompose(id)
	return (id_map["machine-id"] < (1 << alias.c.BitLenMachineId)) &&
		(id_map["msb"] == 0) &&
		(id_map["sequence"] < (1 << alias.c.BitLenSequence)) &&
		(id_map["time"] > 0) &&
		(id_map["time"] < uint64(currentElapsedTime(alias.c.StartTime.AsTime())))
}

// Close is called to release resources
func (alias *MsAlias) Close() {
	if alias.l != nil {
		alias.l.Close()
	}
}
