// Code generated by protoc-gen-grain. DO NOT EDIT.
// versions:
//  protoc-gen-grain v0.4.1
//  protoc           v4.25.0
// source: protos.proto

package shared

import (
	errors "errors"
	fmt "fmt"
	actor "github.com/asynkron/protoactor-go/actor"
	cluster "github.com/asynkron/protoactor-go/cluster"
	proto "google.golang.org/protobuf/proto"
	slog "log/slog"
	time "time"
)

var xCalculatorFactory func() Calculator

// CalculatorFactory produces a Calculator
func CalculatorFactory(factory func() Calculator) {
	xCalculatorFactory = factory
}

// GetCalculatorGrainClient instantiates a new CalculatorGrainClient with given Identity
func GetCalculatorGrainClient(c *cluster.Cluster, id string) *CalculatorGrainClient {
	if c == nil {
		panic(fmt.Errorf("nil cluster instance"))
	}
	if id == "" {
		panic(fmt.Errorf("empty id"))
	}
	return &CalculatorGrainClient{Identity: id, cluster: c}
}

// GetCalculatorKind instantiates a new cluster.Kind for Calculator
func GetCalculatorKind(opts ...actor.PropsOption) *cluster.Kind {
	props := actor.PropsFromProducer(func() actor.Actor {
		return &CalculatorActor{
			Timeout: 60 * time.Second,
		}
	}, opts...)
	kind := cluster.NewKind("Calculator", props)
	return kind
}

// GetCalculatorKind instantiates a new cluster.Kind for Calculator
func NewCalculatorKind(factory func() Calculator, timeout time.Duration, opts ...actor.PropsOption) *cluster.Kind {
	xCalculatorFactory = factory
	props := actor.PropsFromProducer(func() actor.Actor {
		return &CalculatorActor{
			Timeout: timeout,
		}
	}, opts...)
	kind := cluster.NewKind("Calculator", props)
	return kind
}

// Calculator interfaces the services available to the Calculator
type Calculator interface {
	Init(ctx cluster.GrainContext)
	Terminate(ctx cluster.GrainContext)
	ReceiveDefault(ctx cluster.GrainContext)
	Add(req *NumberRequest, ctx cluster.GrainContext) (*CountResponse, error)
	Subtract(req *NumberRequest, ctx cluster.GrainContext) (*CountResponse, error)
	GetCurrent(req *Void, ctx cluster.GrainContext) (*CountResponse, error)
}

// CalculatorGrainClient holds the base data for the CalculatorGrain
type CalculatorGrainClient struct {
	Identity string
	cluster  *cluster.Cluster
}

// Add requests the execution on to the cluster with CallOptions
func (g *CalculatorGrainClient) Add(r *NumberRequest, opts ...cluster.GrainCallOption) (*CountResponse, error) {
	bytes, err := proto.Marshal(r)
	if err != nil {
		return nil, err
	}
	reqMsg := &cluster.GrainRequest{MethodIndex: 0, MessageData: bytes}
	resp, err := g.cluster.Request(g.Identity, "Calculator", reqMsg, opts...)
	if err != nil {
		return nil, fmt.Errorf("error request: %w", err)
	}
	switch msg := resp.(type) {
	case *CountResponse:
		return msg, nil
	case *cluster.GrainErrorResponse:
		return nil, errors.New(msg.Err)
	default:
		return nil, fmt.Errorf("unknown response type %T", resp)
	}
}

// Subtract requests the execution on to the cluster with CallOptions
func (g *CalculatorGrainClient) Subtract(r *NumberRequest, opts ...cluster.GrainCallOption) (*CountResponse, error) {
	bytes, err := proto.Marshal(r)
	if err != nil {
		return nil, err
	}
	reqMsg := &cluster.GrainRequest{MethodIndex: 1, MessageData: bytes}
	resp, err := g.cluster.Request(g.Identity, "Calculator", reqMsg, opts...)
	if err != nil {
		return nil, fmt.Errorf("error request: %w", err)
	}
	switch msg := resp.(type) {
	case *CountResponse:
		return msg, nil
	case *cluster.GrainErrorResponse:
		return nil, errors.New(msg.Err)
	default:
		return nil, fmt.Errorf("unknown response type %T", resp)
	}
}

// GetCurrent requests the execution on to the cluster with CallOptions
func (g *CalculatorGrainClient) GetCurrent(r *Void, opts ...cluster.GrainCallOption) (*CountResponse, error) {
	bytes, err := proto.Marshal(r)
	if err != nil {
		return nil, err
	}
	reqMsg := &cluster.GrainRequest{MethodIndex: 2, MessageData: bytes}
	resp, err := g.cluster.Request(g.Identity, "Calculator", reqMsg, opts...)
	if err != nil {
		return nil, fmt.Errorf("error request: %w", err)
	}
	switch msg := resp.(type) {
	case *CountResponse:
		return msg, nil
	case *cluster.GrainErrorResponse:
		return nil, errors.New(msg.Err)
	default:
		return nil, fmt.Errorf("unknown response type %T", resp)
	}
}

// CalculatorActor represents the actor structure
type CalculatorActor struct {
	ctx     cluster.GrainContext
	inner   Calculator
	Timeout time.Duration
}

// Receive ensures the lifecycle of the actor for the received message
func (a *CalculatorActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started: //pass
	case *cluster.ClusterInit:
		a.ctx = cluster.NewGrainContext(ctx, msg.Identity, msg.Cluster)
		a.inner = xCalculatorFactory()
		a.inner.Init(a.ctx)

		if a.Timeout > 0 {
			ctx.SetReceiveTimeout(a.Timeout)
		}
	case *actor.ReceiveTimeout:
		ctx.Poison(ctx.Self())
	case *actor.Stopped:
		a.inner.Terminate(a.ctx)
	case actor.AutoReceiveMessage: // pass
	case actor.SystemMessage: // pass

	case *cluster.GrainRequest:
		switch msg.MethodIndex {
		case 0:
			req := &NumberRequest{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				ctx.Logger().Error("[Grain] Add(NumberRequest) proto.Unmarshal failed.", slog.Any("error", err))
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
				return
			}

			r0, err := a.inner.Add(req, a.ctx)
			if err != nil {
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
				return
			}
			ctx.Respond(r0)
		case 1:
			req := &NumberRequest{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				ctx.Logger().Error("[Grain] Subtract(NumberRequest) proto.Unmarshal failed.", slog.Any("error", err))
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
				return
			}

			r0, err := a.inner.Subtract(req, a.ctx)
			if err != nil {
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
				return
			}
			ctx.Respond(r0)
		case 2:
			req := &Void{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				ctx.Logger().Error("[Grain] GetCurrent(Void) proto.Unmarshal failed.", slog.Any("error", err))
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
				return
			}

			r0, err := a.inner.GetCurrent(req, a.ctx)
			if err != nil {
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
				return
			}
			ctx.Respond(r0)
		}
	default:
		a.inner.ReceiveDefault(a.ctx)
	}
}

// onError should be used in ctx.ReenterAfter
// you can just return error in reenterable method for other errors
func (a *CalculatorActor) onError(err error) {
	resp := &cluster.GrainErrorResponse{Err: err.Error()}
	a.ctx.Respond(resp)
}

func respond[T proto.Message](ctx cluster.GrainContext) func(T) {
	return func(resp T) {
		ctx.Respond(resp)
	}
}
