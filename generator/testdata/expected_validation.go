// Code generated by protoc-gen-grpc-federation. DO NOT EDIT!
package federation

import (
	"context"
	"io"
	"log/slog"
	"reflect"
	"runtime/debug"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/google/cel-go/cel"
	celtypes "github.com/google/cel-go/common/types"
	grpcfed "github.com/mercari/grpc-federation/grpc/federation"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/singleflight"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

// Org_Federation_GetPostResponseArgument is argument for "org.federation.GetPostResponse" message.
type Org_Federation_GetPostResponseArgument[T any] struct {
	Id     string
	Post   *Post
	Client T
}

// Org_Federation_PostArgument is argument for "org.federation.Post" message.
type Org_Federation_PostArgument[T any] struct {
	Client T
}

// FederationServiceConfig configuration required to initialize the service that use GRPC Federation.
type FederationServiceConfig struct {
	// ErrorHandler Federation Service often needs to convert errors received from downstream services.
	// If an error occurs during method execution in the Federation Service, this error handler is called and the returned error is treated as a final error.
	ErrorHandler grpcfed.ErrorHandler
	// Logger sets the logger used to output Debug/Info/Error information.
	Logger *slog.Logger
}

// FederationServiceClientFactory provides a factory that creates the gRPC Client needed to invoke methods of the gRPC Service on which the Federation Service depends.
type FederationServiceClientFactory interface {
}

// FederationServiceClientConfig information set in `dependencies` of the `grpc.federation.service` option.
// Hints for creating a gRPC Client.
type FederationServiceClientConfig struct {
	// Service returns the name of the service on Protocol Buffers.
	Service string
	// Name is the value set for `name` in `dependencies` of the `grpc.federation.service` option.
	// It must be unique among the services on which the Federation Service depends.
	Name string
}

// FederationServiceDependentClientSet has a gRPC client for all services on which the federation service depends.
// This is provided as an argument when implementing the custom resolver.
type FederationServiceDependentClientSet struct {
}

// FederationServiceResolver provides an interface to directly implement message resolver and field resolver not defined in Protocol Buffers.
type FederationServiceResolver interface {
}

// FederationServiceUnimplementedResolver a structure implemented to satisfy the Resolver interface.
// An Unimplemented error is always returned.
// This is intended for use when there are many Resolver interfaces that do not need to be implemented,
// by embedding them in a resolver structure that you have created.
type FederationServiceUnimplementedResolver struct{}

// FederationService represents Federation Service.
type FederationService struct {
	*UnimplementedFederationServiceServer
	cfg          FederationServiceConfig
	logger       *slog.Logger
	errorHandler grpcfed.ErrorHandler
	env          *cel.Env
	tracer       trace.Tracer
	client       *FederationServiceDependentClientSet
}

// NewFederationService creates FederationService instance by FederationServiceConfig.
func NewFederationService(cfg FederationServiceConfig) (*FederationService, error) {
	logger := cfg.Logger
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	}
	errorHandler := cfg.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(ctx context.Context, methodName string, err error) error { return err }
	}
	celHelper := grpcfed.NewCELTypeHelper(map[string]map[string]*celtypes.FieldType{
		"grpc.federation.private.GetPostResponseArgument": {
			"id": grpcfed.NewCELFieldType(celtypes.StringType, "Id"),
		},
		"grpc.federation.private.PostArgument": {},
	})
	env, err := cel.NewCustomEnv(
		cel.StdLib(),
		cel.CustomTypeAdapter(celHelper.TypeAdapter()),
		cel.CustomTypeProvider(celHelper.TypeProvider()),
	)
	if err != nil {
		return nil, err
	}
	return &FederationService{
		cfg:          cfg,
		logger:       logger,
		errorHandler: errorHandler,
		env:          env,
		tracer:       otel.Tracer("org.federation.FederationService"),
		client:       &FederationServiceDependentClientSet{},
	}, nil
}

// GetPost implements "org.federation.FederationService/GetPost" method.
func (s *FederationService) GetPost(ctx context.Context, req *GetPostRequest) (res *GetPostResponse, e error) {
	ctx, span := s.tracer.Start(ctx, "org.federation.FederationService/GetPost")
	defer span.End()

	ctx = grpcfed.WithLogger(ctx, s.logger)
	defer func() {
		if r := recover(); r != nil {
			e = grpcfed.RecoverError(r, debug.Stack())
			grpcfed.OutputErrorLog(ctx, s.logger, e)
		}
	}()
	res, err := s.resolve_Org_Federation_GetPostResponse(ctx, &Org_Federation_GetPostResponseArgument[*FederationServiceDependentClientSet]{
		Client: s.client,
		Id:     req.Id,
	})
	if err != nil {
		grpcfed.RecordErrorToSpan(ctx, err)
		grpcfed.OutputErrorLog(ctx, s.logger, err)
		return nil, err
	}
	return res, nil
}

// resolve_Org_Federation_GetPostResponse resolve "org.federation.GetPostResponse" message.
func (s *FederationService) resolve_Org_Federation_GetPostResponse(ctx context.Context, req *Org_Federation_GetPostResponseArgument[*FederationServiceDependentClientSet]) (*GetPostResponse, error) {
	ctx, span := s.tracer.Start(ctx, "org.federation.GetPostResponse")
	defer span.End()

	s.logger.DebugContext(ctx, "resolve org.federation.GetPostResponse", slog.Any("message_args", s.logvalue_Org_Federation_GetPostResponseArgument(req)))
	var (
		sg        singleflight.Group
		valueMu   sync.RWMutex
		valuePost *Post
	)
	envOpts := []cel.EnvOption{cel.Variable(grpcfed.MessageArgumentVariableName, cel.ObjectType("grpc.federation.private.GetPostResponseArgument"))}
	evalValues := map[string]any{grpcfed.MessageArgumentVariableName: req}
	// A tree view of message dependencies is shown below.
	/*
	   post ─┐
	         _def1 ─┐
	   post ─┐      │
	         _def2 ─┤
	*/
	eg, ctx1 := errgroup.WithContext(ctx)

	grpcfed.GoWithRecover(eg, func() (any, error) {

		// This section's codes are generated by the following proto definition.
		/*
		   def {
		     name: "post"
		     message {
		       name: "Post"
		     }
		   }
		*/
		{
			valueIface, err, _ := sg.Do("post", func() (any, error) {
				valueMu.RLock()
				args := &Org_Federation_PostArgument[*FederationServiceDependentClientSet]{
					Client: s.client,
				}
				valueMu.RUnlock()
				return s.resolve_Org_Federation_Post(ctx1, args)
			})
			if err != nil {
				return nil, err
			}
			value := valueIface.(*Post)
			valueMu.Lock()
			valuePost = value // { name: "post", message: "Post" ... }
			envOpts = append(envOpts, cel.Variable("post", cel.ObjectType("org.federation.Post")))
			evalValues["post"] = valuePost
			valueMu.Unlock()
		}

		// This section's codes are generated by the following proto definition.
		/*
		   def {
		     name: "_def1"
		     validation {
		       error {
		         code: FAILED_PRECONDITION
		         rule: "post.id == 'some-id'"
		       }
		     }
		   }
		*/
		{
			{
				err := func() error {
					valueMu.RLock()
					value, err := grpcfed.EvalCEL(s.env, "post.id == 'some-id'", envOpts, evalValues, reflect.TypeOf(false))
					valueMu.RUnlock()
					if err != nil {
						return err
					}
					if !value.(bool) {
						return grpcstatus.Error(grpccodes.FailedPrecondition, "validation failure")
					}
					return nil
				}()
				if err != nil {
					if _, ok := grpcstatus.FromError(err); ok {
						return nil, err
					}
					s.logger.ErrorContext(ctx, "failed running validations", slog.String("error", err.Error()))
					return nil, grpcstatus.Errorf(grpccodes.Internal, "failed running validations: %s", err)
				}
			}
		}
		return nil, nil
	})

	grpcfed.GoWithRecover(eg, func() (any, error) {

		// This section's codes are generated by the following proto definition.
		/*
		   def {
		     name: "post"
		     message {
		       name: "Post"
		     }
		   }
		*/
		{
			valueIface, err, _ := sg.Do("post", func() (any, error) {
				valueMu.RLock()
				args := &Org_Federation_PostArgument[*FederationServiceDependentClientSet]{
					Client: s.client,
				}
				valueMu.RUnlock()
				return s.resolve_Org_Federation_Post(ctx1, args)
			})
			if err != nil {
				return nil, err
			}
			value := valueIface.(*Post)
			valueMu.Lock()
			valuePost = value // { name: "post", message: "Post" ... }
			envOpts = append(envOpts, cel.Variable("post", cel.ObjectType("org.federation.Post")))
			evalValues["post"] = valuePost
			valueMu.Unlock()
		}

		// This section's codes are generated by the following proto definition.
		/*
		   def {
		     name: "_def2"
		     validation {
		       error {
		         code: FAILED_PRECONDITION
		         details {
		           rule: "post.title == 'some-title'"
		           precondition_failure {...}
		           bad_request {...}
		           localized_message {...}
		         }
		       }
		     }
		   }
		*/
		{
			{
				err := func() error {
					success := true
					var details []proto.Message
					{
						valueMu.RLock()
						value, err := grpcfed.EvalCEL(s.env, "post.title == 'some-title'", envOpts, evalValues, reflect.TypeOf(false))
						valueMu.RUnlock()
						if err != nil {
							return err
						}
						if !value.(bool) {
							success = false
							{
								var violations []*errdetails.PreconditionFailure_Violation
								{
									func() {
										valueMu.RLock()
										typ, err := grpcfed.EvalCEL(s.env, "'some-type'", envOpts, evalValues, reflect.TypeOf(""))
										valueMu.RUnlock()
										if err != nil {
											s.logger.ErrorContext(ctx, "failed evaluating PreconditionFailure violation type", slog.Int("index", 0), slog.String("error", err.Error()))
											return
										}
										valueMu.RLock()
										subject, err := grpcfed.EvalCEL(s.env, "'some-subject'", envOpts, evalValues, reflect.TypeOf(""))
										valueMu.RUnlock()
										if err != nil {
											s.logger.ErrorContext(ctx, "failed evaluating PreconditionFailure violation subject", slog.Int("index", 0), slog.String("error", err.Error()))
											return
										}
										valueMu.RLock()
										description, err := grpcfed.EvalCEL(s.env, "'some-description'", envOpts, evalValues, reflect.TypeOf(""))
										valueMu.RUnlock()
										if err != nil {
											s.logger.ErrorContext(ctx, "failed evaluating PreconditionFailure violation description", slog.Int("index", 0), slog.String("error", err.Error()))
											return
										}
										violations = append(violations, &errdetails.PreconditionFailure_Violation{
											Type:        typ.(string),
											Subject:     subject.(string),
											Description: description.(string),
										})
									}()
								}
								details = append(details, &errdetails.PreconditionFailure{
									Violations: violations,
								})
							}
							{
								var violations []*errdetails.BadRequest_FieldViolation
								{
									func() {
										valueMu.RLock()
										field, err := grpcfed.EvalCEL(s.env, "'some-field'", envOpts, evalValues, reflect.TypeOf(""))
										valueMu.RUnlock()
										if err != nil {
											s.logger.ErrorContext(ctx, "failed evaluating BadRequest field violation field", slog.Int("index", 0), slog.String("error", err.Error()))
											return
										}
										valueMu.RLock()
										description, err := grpcfed.EvalCEL(s.env, "'some-description'", envOpts, evalValues, reflect.TypeOf(""))
										valueMu.RUnlock()
										if err != nil {
											s.logger.ErrorContext(ctx, "failed evaluating BadRequest field violation description", slog.Int("index", 0), slog.String("error", err.Error()))
											return
										}
										violations = append(violations, &errdetails.BadRequest_FieldViolation{
											Field:       field.(string),
											Description: description.(string),
										})
									}()
								}
								details = append(details, &errdetails.BadRequest{
									FieldViolations: violations,
								})
							}
							{
								func() {
									valueMu.RLock()
									message, err := grpcfed.EvalCEL(s.env, "'some-message'", envOpts, evalValues, reflect.TypeOf(""))
									valueMu.RUnlock()
									if err != nil {
										s.logger.ErrorContext(ctx, "failed evaluating LocalizedMessage message", slog.String("error", err.Error()))
										return
									}
									details = append(details, &errdetails.LocalizedMessage{
										Locale:  "en-US",
										Message: message.(string),
									})
								}()
							}
						}
					}
					if !success {
						status := grpcstatus.New(grpccodes.FailedPrecondition, "validation failure")
						statusWithDetails, err := status.WithDetails(details...)
						if err != nil {
							s.logger.ErrorContext(ctx, "failed setting error details", slog.String("error", err.Error()))
							return status.Err()
						}
						return statusWithDetails.Err()
					}
					return nil
				}()
				if err != nil {
					if _, ok := grpcstatus.FromError(err); ok {
						return nil, err
					}
					s.logger.ErrorContext(ctx, "failed running validations", slog.String("error", err.Error()))
					return nil, grpcstatus.Errorf(grpccodes.Internal, "failed running validations: %s", err)
				}
			}
		}
		return nil, nil
	})

	if err := eg.Wait(); err != nil {
		grpcfed.RecordErrorToSpan(ctx, err)
		return nil, err
	}

	// create a message value to be returned.
	ret := &GetPostResponse{}

	// field binding section.
	// (grpc.federation.field).by = "post"
	{
		value, err := grpcfed.EvalCEL(s.env, "post", envOpts, evalValues, reflect.TypeOf((*Post)(nil)))
		if err != nil {
			grpcfed.RecordErrorToSpan(ctx, err)
			return nil, err
		}
		ret.Post = value.(*Post)
	}

	s.logger.DebugContext(ctx, "resolved org.federation.GetPostResponse", slog.Any("org.federation.GetPostResponse", s.logvalue_Org_Federation_GetPostResponse(ret)))
	return ret, nil
}

// resolve_Org_Federation_Post resolve "org.federation.Post" message.
func (s *FederationService) resolve_Org_Federation_Post(ctx context.Context, req *Org_Federation_PostArgument[*FederationServiceDependentClientSet]) (*Post, error) {
	ctx, span := s.tracer.Start(ctx, "org.federation.Post")
	defer span.End()

	s.logger.DebugContext(ctx, "resolve org.federation.Post", slog.Any("message_args", s.logvalue_Org_Federation_PostArgument(req)))

	// create a message value to be returned.
	ret := &Post{}

	// field binding section.
	ret.Id = "some-id"           // (grpc.federation.field).string = "some-id"
	ret.Title = "some-title"     // (grpc.federation.field).string = "some-title"
	ret.Content = "some-content" // (grpc.federation.field).string = "some-content"

	s.logger.DebugContext(ctx, "resolved org.federation.Post", slog.Any("org.federation.Post", s.logvalue_Org_Federation_Post(ret)))
	return ret, nil
}

func (s *FederationService) logvalue_Org_Federation_GetPostResponse(v *GetPostResponse) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.Any("post", s.logvalue_Org_Federation_Post(v.GetPost())),
	)
}

func (s *FederationService) logvalue_Org_Federation_GetPostResponseArgument(v *Org_Federation_GetPostResponseArgument[*FederationServiceDependentClientSet]) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.String("id", v.Id),
	)
}

func (s *FederationService) logvalue_Org_Federation_Post(v *Post) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.String("id", v.GetId()),
		slog.String("title", v.GetTitle()),
		slog.String("content", v.GetContent()),
	)
}

func (s *FederationService) logvalue_Org_Federation_PostArgument(v *Org_Federation_PostArgument[*FederationServiceDependentClientSet]) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue()
}
