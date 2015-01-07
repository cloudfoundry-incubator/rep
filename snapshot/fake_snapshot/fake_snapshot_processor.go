// This file was generated by counterfeiter
package fake_snapshot

import (
	"sync"

	"github.com/cloudfoundry-incubator/rep/snapshot"
	"github.com/pivotal-golang/lager"
)

type FakeSnapshotProcessor struct {
	ProcessStub        func(lager.Logger, snapshot.Snapshot)
	processMutex       sync.RWMutex
	processArgsForCall []struct {
		arg1 lager.Logger
		arg2 snapshot.Snapshot
	}
}

func (fake *FakeSnapshotProcessor) Process(arg1 lager.Logger, arg2 snapshot.Snapshot) {
	fake.processMutex.Lock()
	fake.processArgsForCall = append(fake.processArgsForCall, struct {
		arg1 lager.Logger
		arg2 snapshot.Snapshot
	}{arg1, arg2})
	fake.processMutex.Unlock()
	if fake.ProcessStub != nil {
		fake.ProcessStub(arg1, arg2)
	}
}

func (fake *FakeSnapshotProcessor) ProcessCallCount() int {
	fake.processMutex.RLock()
	defer fake.processMutex.RUnlock()
	return len(fake.processArgsForCall)
}

func (fake *FakeSnapshotProcessor) ProcessArgsForCall(i int) (lager.Logger, snapshot.Snapshot) {
	fake.processMutex.RLock()
	defer fake.processMutex.RUnlock()
	return fake.processArgsForCall[i].arg1, fake.processArgsForCall[i].arg2
}

var _ snapshot.SnapshotProcessor = new(FakeSnapshotProcessor)
