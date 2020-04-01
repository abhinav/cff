package cff

import (
	"context"
	"sync"
	"time"

	"github.com/uber-go/tally"
)

// cacheKey uniquely identifies a task or a flow based on the position information.
type cacheKey struct {
	TaskName             string // name of the task
	TaskFile             string // file where task is defined
	TaskLine, TaskColumn int    // line and column in the file where the task is defined
	FlowName             string // name of the flow
	FlowFile             string // file where flow is defined
	FlowLine, FlowColumn int    // line and column in the file where the flow is defined
}

type tallyEmitter struct {
	scope tally.Scope

	flows *sync.Map // map[cacheKey]FlowEmitter
	tasks *sync.Map // map[cacheKey]TaskEmitter
}

// TallyEmitter is a CFF2 emitter that emits metrics to Tally.
//
// A full list of metrics published by TallyEmitter can be found at
// https://eng.uberinternal.com/docs/cff2/observability/#metrics.
func TallyEmitter(scope tally.Scope) Emitter {
	return &tallyEmitter{
		scope: scope,
		flows: new(sync.Map),
		tasks: new(sync.Map),
	}
}

func (e *tallyEmitter) TaskInit(taskInfo *TaskInfo, flowInfo *FlowInfo) TaskEmitter {
	cacheKey := cacheKey{
		TaskName:   taskInfo.Task,
		TaskFile:   taskInfo.File,
		TaskLine:   taskInfo.Line,
		TaskColumn: taskInfo.Column,
		FlowName:   flowInfo.Flow,
		FlowFile:   flowInfo.File,
		FlowLine:   flowInfo.Line,
		FlowColumn: flowInfo.Column,
	}
	// Note: this lookup is an optimization to avoid the expensive Tagged call.
	if v, ok := e.tasks.Load(cacheKey); ok {
		return v.(TaskEmitter)
	}
	tags := map[string]string{
		"task": taskInfo.Task,
	}
	if flowInfo.Flow != "" {
		tags["flow"] = flowInfo.Flow
	}

	scope := e.scope.Tagged(tags)
	te := &tallyTaskEmitter{
		scope: scope,
	}
	v, _ := e.tasks.LoadOrStore(cacheKey, te)

	return v.(TaskEmitter)
}

func (e *tallyEmitter) FlowInit(info *FlowInfo) FlowEmitter {
	cacheKey := cacheKey{
		FlowName:   info.Flow,
		FlowFile:   info.File,
		FlowLine:   info.Line,
		FlowColumn: info.Column,
	}
	// Note: this lookup is an optimization to avoid the expensive Tagged call.
	if v, ok := e.flows.Load(cacheKey); ok {
		return v.(FlowEmitter)
	}
	scope := e.scope.Tagged(map[string]string{"flow": info.Flow})
	fe := &tallyFlowEmitter{
		scope: scope,
	}
	v, _ := e.flows.LoadOrStore(cacheKey, fe)

	return v.(FlowEmitter)
}

type tallyFlowEmitter struct {
	scope tally.Scope
}

func (e *tallyFlowEmitter) FlowError(context.Context, error) {
	e.scope.Counter("taskflow.error").Inc(1)
}

func (e *tallyFlowEmitter) FlowSkipped(context.Context, error) {
	e.scope.Counter("taskflow.skipped").Inc(1)
}

func (e *tallyFlowEmitter) FlowSuccess(context.Context) {
	e.scope.Counter("taskflow.success").Inc(1)
}

func (e *tallyFlowEmitter) FlowFailedTask(_ context.Context, task string, _ error) FlowEmitter {
	return &tallyFlowEmitter{
		scope: e.scope.Tagged(map[string]string{
			"failedtask": task,
		})}
}

func (e *tallyFlowEmitter) FlowDone(_ context.Context, d time.Duration) {
	e.scope.Timer("taskflow.timing").Record(d)
}

type tallyTaskEmitter struct {
	scope tally.Scope
}

func (e *tallyTaskEmitter) TaskError(context.Context, error) {
	e.scope.Counter("task.error").Inc(1)
}

func (e *tallyTaskEmitter) TaskErrorRecovered(_ context.Context, err error) {
	e.scope.Counter("task.error").Inc(1)
	e.scope.Counter("task.recovered").Inc(1)
}

func (e *tallyTaskEmitter) TaskPanic(_ context.Context, x interface{}) {
	e.scope.Counter("task.panic").Inc(1)
}

func (e *tallyTaskEmitter) TaskPanicRecovered(_ context.Context, x interface{}) {
	e.scope.Counter("task.panic").Inc(1)
	e.scope.Counter("task.recovered").Inc(1)
}

func (e *tallyTaskEmitter) TaskSkipped(context.Context, error) {
	e.scope.Counter("task.skipped").Inc(1)
}

func (e *tallyTaskEmitter) TaskSuccess(context.Context) {
	e.scope.Counter("task.success").Inc(1)
}

func (e *tallyTaskEmitter) TaskDone(_ context.Context, d time.Duration) {
	e.scope.Timer("task.timing").Record(d)
}
