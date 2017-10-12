package auctioncellrep

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"code.cloudfoundry.org/bbs/models"
	"code.cloudfoundry.org/executor"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/rep"
	"code.cloudfoundry.org/rep/evacuation/evacuation_context"
)

//go:generate counterfeiter . AuctionCellClient

type AuctionCellClient interface {
	State(logger lager.Logger) (rep.CellState, bool, error)
	Perform(logger lager.Logger, work rep.Work) (rep.Work, error)
	Reset() error
}

var ErrPreloadedRootFSNotFound = errors.New("preloaded rootfs path not found")
var ErrCellUnhealthy = errors.New("internal cell healthcheck failed")

type AuctionCellRep struct {
	cellID                string
	stackPathMap          rep.StackPathMap
	rootFSProviders       rep.RootFSProviders
	stack                 string
	zone                  string
	generateInstanceGuid  func() (string, error)
	client                executor.Client
	evacuationReporter    evacuation_context.EvacuationReporter
	placementTags         []string
	optionalPlacementTags []string
}

func New(
	cellID string,
	preloadedStackPathMap rep.StackPathMap,
	arbitraryRootFSes []string,
	zone string,
	generateInstanceGuid func() (string, error),
	client executor.Client,
	evacuationReporter evacuation_context.EvacuationReporter,
	placementTags []string,
	optionalPlacementTags []string,
) *AuctionCellRep {
	return &AuctionCellRep{
		cellID:                cellID,
		stackPathMap:          preloadedStackPathMap,
		rootFSProviders:       rootFSProviders(preloadedStackPathMap, arbitraryRootFSes),
		zone:                  zone,
		generateInstanceGuid:  generateInstanceGuid,
		client:                client,
		evacuationReporter:    evacuationReporter,
		placementTags:         placementTags,
		optionalPlacementTags: optionalPlacementTags,
	}
}

func rootFSProviders(preloaded rep.StackPathMap, arbitrary []string) rep.RootFSProviders {
	rootFSProviders := rep.RootFSProviders{}
	for _, scheme := range arbitrary {
		rootFSProviders[scheme] = rep.ArbitraryRootFSProvider{}
	}

	stacks := make([]string, 0, len(preloaded))
	for stack, _ := range preloaded {
		stacks = append(stacks, stack)
	}
	rootFSProviders[models.PreloadedRootFSScheme] = rep.NewFixedSetRootFSProvider(stacks...)
	rootFSProviders[models.PreloadedOCIRootFSScheme] = rep.NewFixedSetRootFSProvider(stacks...)

	return rootFSProviders
}

func PathForRootFS(rootFS string, stackPathMap rep.StackPathMap) (string, error) {
	if rootFS == "" {
		return rootFS, nil
	}

	url, err := url.Parse(rootFS)
	if err != nil {
		return "", err
	}

	if url.Scheme == models.PreloadedRootFSScheme {
		path, ok := stackPathMap[url.Opaque]
		if !ok {
			return "", ErrPreloadedRootFSNotFound
		}
		return path, nil
	} else if url.Scheme == models.PreloadedOCIRootFSScheme {
		path, ok := stackPathMap[url.Opaque]
		if !ok {
			return "", ErrPreloadedRootFSNotFound
		}

		return fmt.Sprintf("%s:%s?%s", url.Scheme, path, url.RawQuery), nil
	}

	return rootFS, nil
}

// State currently does not return tasks or lrp rootfs, because the
// auctioneer currently does not need them.
func (a *AuctionCellRep) State(logger lager.Logger) (rep.CellState, bool, error) {
	logger = logger.Session("auction-state")
	logger.Info("providing")

	containers, err := a.client.ListContainers(logger)
	if err != nil {
		logger.Error("failed-to-fetch-containers", err)
		return rep.CellState{}, false, err
	}

	totalResources, err := a.client.TotalResources(logger)
	if err != nil {
		logger.Error("failed-to-get-total-resources", err)
		return rep.CellState{}, false, err
	}

	availableResources, err := a.client.RemainingResources(logger)
	if err != nil {
		logger.Error("failed-to-get-remaining-resource", err)
		return rep.CellState{}, false, err
	}

	volumeDrivers, err := a.client.VolumeDrivers(logger)
	if err != nil {
		logger.Error("failed-to-get-volume-drivers", err)
		return rep.CellState{}, false, err
	}

	lrps := []rep.LRP{}
	tasks := []rep.Task{}
	startingContainerCount := 0

	for i := range containers {
		container := &containers[i]

		if containerIsStarting(container) {
			startingContainerCount++
		}

		if container.Tags == nil {
			logger.Error("failed-to-extract-container-tags", nil)
			continue
		}

		placementTagsJSON := container.Tags[rep.PlacementTagsTag]
		var placementTags []string
		err := json.Unmarshal([]byte(placementTagsJSON), &placementTags)
		if err != nil {
			logger.Error("cannot-unmarshal-placement-tags", err, lager.Data{"placement-tags": placementTagsJSON})
		}

		volumeDriversJSON := container.Tags[rep.VolumeDriversTag]
		var volumeDrivers []string
		err = json.Unmarshal([]byte(volumeDriversJSON), &volumeDrivers)
		if err != nil {
			logger.Error("cannot-unmarshal-volume-drivers", err, lager.Data{"volume-drivers": volumeDriversJSON})
		}

		resource := rep.Resource{MemoryMB: int32(container.MemoryMB), DiskMB: int32(container.DiskMB), MaxPids: int32(container.MaxPids)}
		placementConstraint := rep.PlacementConstraint{RootFs: container.RootFSPath,
			VolumeDrivers: volumeDrivers,
			PlacementTags: placementTags,
		}

		switch container.Tags[rep.LifecycleTag] {
		case rep.LRPLifecycle:
			key, err := rep.ActualLRPKeyFromTags(container.Tags)
			if err != nil {
				logger.Error("failed-to-extract-key", err)
				continue
			}
			instanceKey, err := rep.ActualLRPInstanceKeyFromContainer(*container, a.cellID)
			if err != nil {
				logger.Error("failed-to-extract-key", err)
				continue
			}
			lrps = append(lrps, rep.NewLRP(instanceKey.InstanceGuid, *key, resource, placementConstraint))
		case rep.TaskLifecycle:
			domain := container.Tags[rep.DomainTag]
			tasks = append(tasks, rep.NewTask(container.Guid, domain, resource, placementConstraint))
		}
	}

	state := rep.NewCellState(
		a.rootFSProviders,
		a.convertResources(availableResources),
		a.convertResources(totalResources),
		lrps,
		tasks,
		a.zone,
		startingContainerCount,
		a.evacuationReporter.Evacuating(),
		volumeDrivers,
		a.placementTags,
		a.optionalPlacementTags,
	)

	healthy := a.client.Healthy(logger)
	if !healthy {
		logger.Error("failed-garden-health-check", nil)
	}

	logger.Info("provided", lager.Data{
		"available-resources": state.AvailableResources,
		"total-resources":     state.TotalResources,
		"num-lrps":            len(state.LRPs),
		"zone":                state.Zone,
		"evacuating":          state.Evacuating,
	})

	return state, healthy, nil
}

func containerIsStarting(container *executor.Container) bool {
	return container.State == executor.StateReserved ||
		container.State == executor.StateInitializing ||
		container.State == executor.StateCreated
}

func (a *AuctionCellRep) Perform(logger lager.Logger, work rep.Work) (rep.Work, error) {
	var failedWork = rep.Work{}

	logger = logger.Session("auction-work", lager.Data{
		"lrp-starts": len(work.LRPs),
		"tasks":      len(work.Tasks),
	})

	if a.evacuationReporter.Evacuating() {
		return work, nil
	}

	if len(work.LRPs) > 0 {
		lrpLogger := logger.Session("lrp-allocate-instances")

		requests, lrpMap, untranslatedLRPs := a.lrpsToAllocationRequest(work.LRPs)
		if len(untranslatedLRPs) > 0 {
			lrpLogger.Info("failed-to-translate-lrps-to-containers", lager.Data{"num-failed-to-translate": len(untranslatedLRPs)})
			failedWork.LRPs = untranslatedLRPs
		}

		lrpLogger.Info("requesting-container-allocation", lager.Data{"num-requesting-allocation": len(requests)})
		failures, err := a.client.AllocateContainers(logger, requests)
		if err != nil {
			lrpLogger.Error("failed-requesting-container-allocation", err)
			failedWork.LRPs = work.LRPs
		} else {
			lrpLogger.Info("succeeded-requesting-container-allocation", lager.Data{"num-failed-to-allocate": len(failures)})
			for i := range failures {
				failure := &failures[i]
				lrpLogger.Error("container-allocation-failure", failure, lager.Data{"failed-request": &failure.AllocationRequest})
				if lrp, found := lrpMap[failure.Guid]; found {
					failedWork.LRPs = append(failedWork.LRPs, *lrp)
				}
			}
		}
	}

	if len(work.Tasks) > 0 {
		taskLogger := logger.Session("task-allocate-instances")

		requests, taskMap, failedTasks := a.tasksToAllocationRequests(work.Tasks)
		if len(failedTasks) > 0 {
			taskLogger.Info("failed-to-translate-tasks-to-containers", lager.Data{"num-failed-to-translate": len(failedTasks)})
			failedWork.Tasks = failedTasks
		}

		taskLogger.Info("requesting-container-allocation", lager.Data{"num-requesting-allocation": len(requests)})
		failures, err := a.client.AllocateContainers(logger, requests)
		if err != nil {
			taskLogger.Error("failed-requesting-container-allocation", err)
			failedWork.Tasks = work.Tasks
		} else {
			taskLogger.Info("succeeded-requesting-container-allocation", lager.Data{"num-failed-to-allocate": len(failures)})
			for i := range failures {
				failure := &failures[i]
				taskLogger.Error("container-allocation-failure", failure, lager.Data{"failed-request": &failure.AllocationRequest})
				if task, found := taskMap[failure.Guid]; found {
					failedWork.Tasks = append(failedWork.Tasks, *task)
				}
			}
		}
	}

	return failedWork, nil
}

func (a *AuctionCellRep) lrpsToAllocationRequest(lrps []rep.LRP) ([]executor.AllocationRequest, map[string]*rep.LRP, []rep.LRP) {
	requests := make([]executor.AllocationRequest, 0, len(lrps))
	untranslatedLRPs := make([]rep.LRP, 0)
	lrpMap := make(map[string]*rep.LRP, len(lrps))
	for i := range lrps {
		lrp := &lrps[i]
		tags := executor.Tags{}

		instanceGuid, err := a.generateInstanceGuid()
		if err != nil {
			untranslatedLRPs = append(untranslatedLRPs, *lrp)
			continue
		}

		tags[rep.DomainTag] = lrp.Domain
		tags[rep.ProcessGuidTag] = lrp.ProcessGuid
		tags[rep.ProcessIndexTag] = strconv.Itoa(int(lrp.Index))
		tags[rep.LifecycleTag] = rep.LRPLifecycle
		tags[rep.InstanceGuidTag] = instanceGuid

		rootFSPath, err := PathForRootFS(lrp.RootFs, a.stackPathMap)
		if err != nil {
			untranslatedLRPs = append(untranslatedLRPs, *lrp)
			continue
		}

		containerGuid := rep.LRPContainerGuid(lrp.ProcessGuid, instanceGuid)
		lrpMap[containerGuid] = lrp

		resource := executor.NewResource(int(lrp.MemoryMB), int(lrp.DiskMB), int(lrp.MaxPids), rootFSPath)
		requests = append(requests, executor.NewAllocationRequest(containerGuid, &resource, tags))
	}

	return requests, lrpMap, untranslatedLRPs
}

func (a *AuctionCellRep) tasksToAllocationRequests(tasks []rep.Task) ([]executor.AllocationRequest, map[string]*rep.Task, []rep.Task) {
	failedTasks := make([]rep.Task, 0)
	taskMap := make(map[string]*rep.Task, len(tasks))
	requests := make([]executor.AllocationRequest, 0, len(tasks))

	for i := range tasks {
		task := &tasks[i]
		taskMap[task.TaskGuid] = task
		rootFSPath, err := PathForRootFS(task.RootFs, a.stackPathMap)
		if err != nil {
			failedTasks = append(failedTasks, *task)
			continue
		}
		tags := executor.Tags{}
		tags[rep.LifecycleTag] = rep.TaskLifecycle
		tags[rep.DomainTag] = task.Domain

		resource := executor.NewResource(int(task.MemoryMB), int(task.DiskMB), int(task.MaxPids), rootFSPath)
		requests = append(requests, executor.NewAllocationRequest(task.TaskGuid, &resource, tags))
	}

	return requests, taskMap, failedTasks
}

func (a *AuctionCellRep) convertResources(resources executor.ExecutorResources) rep.Resources {
	return rep.Resources{
		MemoryMB:   int32(resources.MemoryMB),
		DiskMB:     int32(resources.DiskMB),
		Containers: resources.Containers,
	}
}

func (a *AuctionCellRep) Reset() error {
	return errors.New("not-a-simulation-rep")
}
