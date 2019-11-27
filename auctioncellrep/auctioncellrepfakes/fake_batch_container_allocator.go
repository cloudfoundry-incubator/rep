// Code generated by counterfeiter. DO NOT EDIT.
package auctioncellrepfakes

import (
	"sync"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/rep"
	"code.cloudfoundry.org/rep/auctioncellrep"
)

type FakeBatchContainerAllocator struct {
	BatchLRPAllocationRequestStub        func(lager.Logger, bool, int, []rep.LRP) []rep.LRP
	batchLRPAllocationRequestMutex       sync.RWMutex
	batchLRPAllocationRequestArgsForCall []struct {
		arg1 lager.Logger
		arg2 bool
		arg3 int
		arg4 []rep.LRP
	}
	batchLRPAllocationRequestReturns struct {
		result1 []rep.LRP
	}
	batchLRPAllocationRequestReturnsOnCall map[int]struct {
		result1 []rep.LRP
	}
	BatchTaskAllocationRequestStub        func(lager.Logger, []rep.Task) []rep.Task
	batchTaskAllocationRequestMutex       sync.RWMutex
	batchTaskAllocationRequestArgsForCall []struct {
		arg1 lager.Logger
		arg2 []rep.Task
	}
	batchTaskAllocationRequestReturns struct {
		result1 []rep.Task
	}
	batchTaskAllocationRequestReturnsOnCall map[int]struct {
		result1 []rep.Task
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBatchContainerAllocator) BatchLRPAllocationRequest(arg1 lager.Logger, arg2 bool, arg3 int, arg4 []rep.LRP) []rep.LRP {
	var arg4Copy []rep.LRP
	if arg4 != nil {
		arg4Copy = make([]rep.LRP, len(arg4))
		copy(arg4Copy, arg4)
	}
	fake.batchLRPAllocationRequestMutex.Lock()
	ret, specificReturn := fake.batchLRPAllocationRequestReturnsOnCall[len(fake.batchLRPAllocationRequestArgsForCall)]
	fake.batchLRPAllocationRequestArgsForCall = append(fake.batchLRPAllocationRequestArgsForCall, struct {
		arg1 lager.Logger
		arg2 bool
		arg3 int
		arg4 []rep.LRP
	}{arg1, arg2, arg3, arg4Copy})
	fake.recordInvocation("BatchLRPAllocationRequest", []interface{}{arg1, arg2, arg3, arg4Copy})
	batchLRPAllocationRequestStubCopy := fake.BatchLRPAllocationRequestStub
	fake.batchLRPAllocationRequestMutex.Unlock()
	if batchLRPAllocationRequestStubCopy != nil {
		return batchLRPAllocationRequestStubCopy(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.batchLRPAllocationRequestReturns
	return fakeReturns.result1
}

func (fake *FakeBatchContainerAllocator) BatchLRPAllocationRequestCallCount() int {
	fake.batchLRPAllocationRequestMutex.RLock()
	defer fake.batchLRPAllocationRequestMutex.RUnlock()
	return len(fake.batchLRPAllocationRequestArgsForCall)
}

func (fake *FakeBatchContainerAllocator) BatchLRPAllocationRequestCalls(stub func(lager.Logger, bool, int, []rep.LRP) []rep.LRP) {
	fake.batchLRPAllocationRequestMutex.Lock()
	defer fake.batchLRPAllocationRequestMutex.Unlock()
	fake.BatchLRPAllocationRequestStub = stub
}

func (fake *FakeBatchContainerAllocator) BatchLRPAllocationRequestArgsForCall(i int) (lager.Logger, bool, int, []rep.LRP) {
	fake.batchLRPAllocationRequestMutex.RLock()
	defer fake.batchLRPAllocationRequestMutex.RUnlock()
	argsForCall := fake.batchLRPAllocationRequestArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeBatchContainerAllocator) BatchLRPAllocationRequestReturns(result1 []rep.LRP) {
	fake.batchLRPAllocationRequestMutex.Lock()
	defer fake.batchLRPAllocationRequestMutex.Unlock()
	fake.BatchLRPAllocationRequestStub = nil
	fake.batchLRPAllocationRequestReturns = struct {
		result1 []rep.LRP
	}{result1}
}

func (fake *FakeBatchContainerAllocator) BatchLRPAllocationRequestReturnsOnCall(i int, result1 []rep.LRP) {
	fake.batchLRPAllocationRequestMutex.Lock()
	defer fake.batchLRPAllocationRequestMutex.Unlock()
	fake.BatchLRPAllocationRequestStub = nil
	if fake.batchLRPAllocationRequestReturnsOnCall == nil {
		fake.batchLRPAllocationRequestReturnsOnCall = make(map[int]struct {
			result1 []rep.LRP
		})
	}
	fake.batchLRPAllocationRequestReturnsOnCall[i] = struct {
		result1 []rep.LRP
	}{result1}
}

func (fake *FakeBatchContainerAllocator) BatchTaskAllocationRequest(arg1 lager.Logger, arg2 []rep.Task) []rep.Task {
	var arg2Copy []rep.Task
	if arg2 != nil {
		arg2Copy = make([]rep.Task, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.batchTaskAllocationRequestMutex.Lock()
	ret, specificReturn := fake.batchTaskAllocationRequestReturnsOnCall[len(fake.batchTaskAllocationRequestArgsForCall)]
	fake.batchTaskAllocationRequestArgsForCall = append(fake.batchTaskAllocationRequestArgsForCall, struct {
		arg1 lager.Logger
		arg2 []rep.Task
	}{arg1, arg2Copy})
	fake.recordInvocation("BatchTaskAllocationRequest", []interface{}{arg1, arg2Copy})
	batchTaskAllocationRequestStubCopy := fake.BatchTaskAllocationRequestStub
	fake.batchTaskAllocationRequestMutex.Unlock()
	if batchTaskAllocationRequestStubCopy != nil {
		return batchTaskAllocationRequestStubCopy(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.batchTaskAllocationRequestReturns
	return fakeReturns.result1
}

func (fake *FakeBatchContainerAllocator) BatchTaskAllocationRequestCallCount() int {
	fake.batchTaskAllocationRequestMutex.RLock()
	defer fake.batchTaskAllocationRequestMutex.RUnlock()
	return len(fake.batchTaskAllocationRequestArgsForCall)
}

func (fake *FakeBatchContainerAllocator) BatchTaskAllocationRequestCalls(stub func(lager.Logger, []rep.Task) []rep.Task) {
	fake.batchTaskAllocationRequestMutex.Lock()
	defer fake.batchTaskAllocationRequestMutex.Unlock()
	fake.BatchTaskAllocationRequestStub = stub
}

func (fake *FakeBatchContainerAllocator) BatchTaskAllocationRequestArgsForCall(i int) (lager.Logger, []rep.Task) {
	fake.batchTaskAllocationRequestMutex.RLock()
	defer fake.batchTaskAllocationRequestMutex.RUnlock()
	argsForCall := fake.batchTaskAllocationRequestArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeBatchContainerAllocator) BatchTaskAllocationRequestReturns(result1 []rep.Task) {
	fake.batchTaskAllocationRequestMutex.Lock()
	defer fake.batchTaskAllocationRequestMutex.Unlock()
	fake.BatchTaskAllocationRequestStub = nil
	fake.batchTaskAllocationRequestReturns = struct {
		result1 []rep.Task
	}{result1}
}

func (fake *FakeBatchContainerAllocator) BatchTaskAllocationRequestReturnsOnCall(i int, result1 []rep.Task) {
	fake.batchTaskAllocationRequestMutex.Lock()
	defer fake.batchTaskAllocationRequestMutex.Unlock()
	fake.BatchTaskAllocationRequestStub = nil
	if fake.batchTaskAllocationRequestReturnsOnCall == nil {
		fake.batchTaskAllocationRequestReturnsOnCall = make(map[int]struct {
			result1 []rep.Task
		})
	}
	fake.batchTaskAllocationRequestReturnsOnCall[i] = struct {
		result1 []rep.Task
	}{result1}
}

func (fake *FakeBatchContainerAllocator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.batchLRPAllocationRequestMutex.RLock()
	defer fake.batchLRPAllocationRequestMutex.RUnlock()
	fake.batchTaskAllocationRequestMutex.RLock()
	defer fake.batchTaskAllocationRequestMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBatchContainerAllocator) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ auctioncellrep.BatchContainerAllocator = new(FakeBatchContainerAllocator)