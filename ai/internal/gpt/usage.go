package gpt

import (
	"github.com/finishy1995/go-library/routine"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

type usage struct {
	sync.RWMutex
	promptTokens      int
	completionTokens  int
	totalTokens       int
	startTime         time.Time
	recent15minTokens []tokenRecord
	limit             int
}

type tokenRecord struct {
	tokens int
	time   time.Time
}

const (
	DefaultLimit = 800000
)

var (
	usageInstance *usage
)

func init() {
	usageInstance = &usage{
		promptTokens:      0,
		completionTokens:  0,
		totalTokens:       0,
		startTime:         time.Now(),
		recent15minTokens: make([]tokenRecord, 15, 15),
		limit:             DefaultLimit,
	}
	err := routine.Run(false, usageInstance.process)
	if err != nil {
		panic(err)
	}
}

func (u *usage) process() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		u.Lock()
		usageInstance.updateRecent15minTokens(0)
		sum := u.sum()
		u.Unlock()
		logx.Infof("totalTokens: %d, recent15minTokens: %d, limit: %d, average rate(minute): %d", u.totalTokens, sum, u.limit, u.totalTokens/int(time.Now().Sub(u.startTime).Minutes()))
		logx.Debugf("recent15minTokens: %v", u.recent15minTokens)
		logx.Debugf("promptTokens: %d, completionTokens: %d, totalTokens: %d", u.promptTokens, u.completionTokens, u.totalTokens)
	}
}

func (u *usage) update(promptTokens, completionTokens, totalTokens int) {
	u.Lock()
	defer u.Unlock()
	u.promptTokens += promptTokens
	u.completionTokens += completionTokens
	u.totalTokens += totalTokens

	u.updateRecent15minTokens(totalTokens)
}

func (u *usage) updateRecent15minTokens(tokens int) {
	now := time.Now()
	minute := now.Minute()
	index := minute % 15
	if now.Sub(u.recent15minTokens[index].time) < time.Minute {
		u.recent15minTokens[minute%15].tokens += tokens
		return
	} else {
		u.recent15minTokens[minute%15].tokens = tokens
	}
	u.recent15minTokens[minute%15].time = now
}

// 检查近15分钟是否超过限制
func (u *usage) checkRecent15min() bool {
	u.RLock()
	defer u.RUnlock()

	return u.sum() > u.limit
}

func (u *usage) sum() int {
	// 对 recent15minTokens token 求和
	sum := 0
	for _, v := range u.recent15minTokens {
		sum += v.tokens
	}
	return sum
}
