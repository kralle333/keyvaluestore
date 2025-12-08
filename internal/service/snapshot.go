package service

import (
	"time"

	"github.com/kralle333/keyvaluestore/internal/model"
	"go.uber.org/zap"
)

type SnapshotService struct {
	communication    *model.KeyValueActorCommunication
	logger           *zap.Logger
	snapshotInterval time.Duration
	shutdownChannel  chan struct{}
}

func NewSnapshotService(communication *model.KeyValueActorCommunication, snapshotIntervalSeconds int64, parentLogger *zap.Logger) *SnapshotService {
	return &SnapshotService{
		communication:    communication,
		logger:           parentLogger.With(zap.String("source", "snapshot service")),
		snapshotInterval: time.Second * time.Duration(snapshotIntervalSeconds),
		shutdownChannel:  make(chan struct{}),
	}
}

func (s *SnapshotService) Spawn() {
	s.logger.Info("Spawning snapshot taker goroutine", zap.Int("snapshot interval seconds", int(s.snapshotInterval)))
	go func() bool {
		ticker := time.NewTicker(s.snapshotInterval)
		for {
			select {
			case <-ticker.C:
				s.logger.Debug("Sending snapshot request")
				s.communication.TakeSnapshot()
			case <-s.shutdownChannel:
				s.logger.Info("Shutting down")
			}
			<-ticker.C
		}
	}()
}
func (s *SnapshotService) Shutdown() {
	s.shutdownChannel <- struct{}{}
}
