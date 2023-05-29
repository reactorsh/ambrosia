// You don't want to use this. It's good enough, but it's not good.
package limiter

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var (
	nilStruct struct{}
)

type Limiter struct {
	rpm            int
	tpm            int
	rpmC           chan struct{}
	tpmC           chan struct{}
	logger         *zerolog.Logger
	statMutex      sync.Mutex
	statRequestEst int
	statTokenUsage int
}

func New(rpm, tpm int, logger *zerolog.Logger) *Limiter {
	ret := &Limiter{
		rpm:    rpm,
		tpm:    tpm,
		rpmC:   make(chan struct{}, rpm),
		tpmC:   make(chan struct{}, tpm),
		logger: logger,
	}

	if logger != nil {
		go ret.stats()
	}

	ret.start()
	return ret
}

func (l *Limiter) start() {
	go filler(l.rpm, l.rpmC)
	go filler(l.tpm, l.tpmC)
}

func (l *Limiter) stats() {
	for range time.Tick(1 * time.Minute) {
		l.statMutex.Lock()
		l.logger.Debug().
			Int("rpm_chan_len", len(l.rpmC)).
			Int("tpm_chan_len", len(l.tpmC)).
			Int("request_est_1m", l.statRequestEst).
			Int("token_usage_1m", l.statTokenUsage).
			Msg("limiter stats")
		l.statRequestEst = 0
		l.statTokenUsage = 0
		l.statMutex.Unlock()
	}
}

func (l *Limiter) Wait(est int) {
	for i := 0; i < est; i++ {
		<-l.tpmC
	}
	<-l.rpmC

	if l.logger != nil {
		l.statMutex.Lock()
		l.statRequestEst += est
		l.statMutex.Unlock()
	}
}

func (l *Limiter) TPMReconcile(est, tokens int) {
	// Take out the actual token count.
	for i := 0; i < tokens; i++ {
		<-l.tpmC
	}

	if l.logger != nil {
		l.statMutex.Lock()
		l.statTokenUsage += tokens
		l.statMutex.Unlock()
	}

	// Return the byte count.
	for i := 0; i < est; i++ {
		select {
		case l.tpmC <- nilStruct:
		default:
		}
	}
}

func filler(fpm int, c chan<- struct{}) {
	ticker := time.NewTicker(time.Minute / time.Duration(fpm))
	for range ticker.C {
		select {
		case c <- nilStruct:
		default:
		}
	}
}
