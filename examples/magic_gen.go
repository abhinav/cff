// @generated by CFF

package example

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"go.uber.org/cff"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

// Request TODO
type Request struct {
	LDAPGroup string
}

// Response TODO
type Response struct {
	MessageIDs []string
}

type fooHandler struct {
	mgr    *ManagerRepository
	users  *UserRepository
	ses    *SESClient
	scope  tally.Scope
	logger *zap.Logger
}

func (h *fooHandler) HandleFoo(ctx context.Context, req *Request) (*Response, error) {
	var res *Response
	err := func(
		ctx context.Context,
		emitter cff.Emitter,
		v1 *Request,
	) (err error) {
		var (
			flowInfo = &cff.FlowInfo{
				Flow:   "HandleFoo",
				File:   "go.uber.org/cff/examples/magic.go",
				Line:   34,
				Column: 9,
			}
			flowEmitter = emitter.FlowInit(flowInfo)

			// possibly unused
			_ = flowInfo
		)

		startTime := time.Now()
		defer func() { flowEmitter.FlowDone(ctx, time.Since(startTime)) }()

		type task struct {
			emitter cff.TaskEmitter
			ran     bool
			run     func(context.Context) error
		}

		var tasks []*task
		defer func() {
			for _, t := range tasks {
				if !t.ran {
					t.emitter.TaskSkipped(ctx, err)
				}
			}

			if err != nil {
				flowEmitter.FlowSkipped(ctx, err)
			}
		}()

		var (
			v2 *GetManagerRequest   // from task0
			v3 *ListUsersRequest    // from task0
			v4 *GetManagerResponse  // from task1
			v5 []*SendEmailResponse // from task2
			v6 *Response            // from task3
			v7 *ListUsersResponse   // from task4
			v8 []*SendEmailRequest  // from task5

		)

		// go.uber.org/cff/examples/magic.go:42:4
		task0 := new(task)
		tasks = append(tasks, task0)
		task0.emitter = cff.NopTaskEmitter()
		task0.run = func(ctx context.Context) (err error) {
			taskEmitter := task0.emitter
			startTime := time.Now()
			defer func() { taskEmitter.TaskDone(ctx, time.Since(startTime)) }()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			v2, v3 = func(req *Request) (*GetManagerRequest, *ListUsersRequest) {
				return &GetManagerRequest{
						LDAPGroup: req.LDAPGroup,
					}, &ListUsersRequest{
						LDAPGroup: req.LDAPGroup,
					}
			}(v1)
			task0.ran = true

			taskEmitter.TaskSuccess(ctx)

			return
		}

		// go.uber.org/cff/examples/magic.go:50:4
		task1 := new(task)
		tasks = append(tasks, task1)
		task1.emitter = cff.NopTaskEmitter()
		task1.run = func(ctx context.Context) (err error) {
			taskEmitter := task1.emitter
			startTime := time.Now()
			defer func() { taskEmitter.TaskDone(ctx, time.Since(startTime)) }()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			v4, err = h.mgr.Get(v2)
			task1.ran = true

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		// go.uber.org/cff/examples/magic.go:51:12
		task2 := new(task)
		tasks = append(tasks, task2)
		task2.emitter = cff.NopTaskEmitter()
		task2.run = func(ctx context.Context) (err error) {
			taskEmitter := task2.emitter
			startTime := time.Now()
			defer func() { taskEmitter.TaskDone(ctx, time.Since(startTime)) }()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			v5, err = h.ses.BatchSendEmail(v8)
			task2.ran = true

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		// go.uber.org/cff/examples/magic.go:53:4
		task3 := new(task)
		tasks = append(tasks, task3)
		task3.emitter = cff.NopTaskEmitter()
		task3.run = func(ctx context.Context) (err error) {
			taskEmitter := task3.emitter
			startTime := time.Now()
			defer func() { taskEmitter.TaskDone(ctx, time.Since(startTime)) }()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			v6 = func(responses []*SendEmailResponse) *Response {
				var r Response
				for _, res := range responses {
					r.MessageIDs = append(r.MessageIDs, res.MessageID)
				}
				return &r
			}(v5)
			task3.ran = true

			taskEmitter.TaskSuccess(ctx)

			return
		}

		// go.uber.org/cff/examples/magic.go:62:4
		task4 := new(task)
		tasks = append(tasks, task4)
		task4.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Task:   "FormSendEmailRequest",
				File:   "go.uber.org/cff/examples/magic.go",
				Line:   62,
				Column: 4,
			},
			flowInfo,
		)
		task4.run = func(ctx context.Context) (err error) {
			taskEmitter := task4.emitter
			startTime := time.Now()
			defer func() { taskEmitter.TaskDone(ctx, time.Since(startTime)) }()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanicRecovered(ctx, recovered)
					v7, err = &ListUsersResponse{}, nil
				}
			}()

			v7, err = h.users.List(v3)
			task4.ran = true

			if err != nil {
				taskEmitter.TaskErrorRecovered(ctx, err)
				v7, err = &ListUsersResponse{}, nil
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		// go.uber.org/cff/examples/magic.go:67:4
		task5 := new(task)
		tasks = append(tasks, task5)
		task5.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Task:   "FormSendEmailRequest",
				File:   "go.uber.org/cff/examples/magic.go",
				Line:   67,
				Column: 4,
			},
			flowInfo,
		)
		task5.run = func(ctx context.Context) (err error) {
			taskEmitter := task5.emitter
			startTime := time.Now()
			defer func() { taskEmitter.TaskDone(ctx, time.Since(startTime)) }()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			if !(func(req *GetManagerRequest) bool {
				return req.LDAPGroup != "everyone"
			}(v2)) {
				return nil
			}

			v8 = func(mgr *GetManagerResponse, users *ListUsersResponse) []*SendEmailRequest {
				var reqs []*SendEmailRequest
				for _, u := range users.Emails {
					reqs = append(reqs, &SendEmailRequest{Address: u})
				}
				return reqs
			}(v4, v7)
			task5.ran = true

			taskEmitter.TaskSuccess(ctx)

			return
		}

		schedule := [][]*task{
			{task0},
			{task1, task4},
			{task5},
			{task2},
			{task3},
		}

		for _, taskGroup := range schedule {
			if err := ctx.Err(); err != nil {
				return err
			}

			if len(taskGroup) == 1 {
				if err := taskGroup[0].run(ctx); err != nil {
					flowEmitter.FlowError(ctx, err)
					return err
				}
				continue
			}

			var (
				wg   sync.WaitGroup
				once sync.Once
				err  error
			)

			wg.Add(len(taskGroup))
			for _, t := range taskGroup {
				go func(t *task) {
					defer wg.Done()
					if terr := t.run(ctx); terr != nil {
						once.Do(func() {
							err = terr
						})
					}
				}(t)
			}

			wg.Wait()

			if err != nil {
				flowEmitter.FlowError(ctx, err)
				return err
			}
		}

		*(&res) = v6 // *go.uber.org/cff/examples.Response

		flowEmitter.FlowSuccess(ctx)
		return nil
	}(
		ctx,
		cff.EmitterStack(cff.TallyEmitter(h.scope), cff.LogEmitter(h.logger)),
		req,
	)
	return res, err
}

// ManagerRepository TODO
type ManagerRepository struct{}

// GetManagerRequest TODO
type GetManagerRequest struct {
	LDAPGroup string
}

// GetManagerResponse TODO
type GetManagerResponse struct {
	Email string
}

// Get TODO
func (*ManagerRepository) Get(req *GetManagerRequest) (*GetManagerResponse, error) {
	return &GetManagerResponse{Email: "boss@example.com"}, nil
}

// UserRepository TODO
type UserRepository struct{}

// ListUsersRequest TODO
type ListUsersRequest struct {
	LDAPGroup string
}

// ListUsersResponse TODO
type ListUsersResponse struct {
	Emails []string
}

// List TODO
func (*UserRepository) List(req *ListUsersRequest) (*ListUsersResponse, error) {
	return &ListUsersResponse{
		Emails: []string{"a@example.com", "b@example.com"},
	}, nil
}

// SESClient TODO
type SESClient struct{}

// SendEmailRequest TODO
type SendEmailRequest struct {
	Address string
}

// SendEmailResponse TODO
type SendEmailResponse struct {
	MessageID string
}

// BatchSendEmail TODO
func (*SESClient) BatchSendEmail(req []*SendEmailRequest) ([]*SendEmailResponse, error) {
	res := make([]*SendEmailResponse, len(req))
	for i := range req {
		res[i] = &SendEmailResponse{MessageID: strconv.Itoa(i)}
	}
	return res, nil
}
