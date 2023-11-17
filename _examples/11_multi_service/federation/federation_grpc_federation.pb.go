// Code generated by protoc-gen-grpc-federation. DO NOT EDIT!
package federation

import (
	"context"
	"io"
	"log/slog"
	"reflect"
	"runtime/debug"
	"sync"

	"github.com/google/cel-go/cel"
	celtypes "github.com/google/cel-go/common/types"
	grpcfed "github.com/mercari/grpc-federation/grpc/federation"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/singleflight"
)

// Federation_GetPostResponseArgument is argument for "federation.GetPostResponse" message.
type Federation_GetPostResponseArgument[T any] struct {
	Id     string
	P      *Post
	Client T
}

// Federation_GetStatusResponseArgument is argument for "federation.GetStatusResponse" message.
type Federation_GetStatusResponseArgument[T any] struct {
	U      *User
	Client T
}

// Federation_PostArgument is argument for "federation.Post" message.
type Federation_PostArgument[T any] struct {
	U      *User
	Client T
}

// Federation_UserArgument is argument for "federation.User" message.
type Federation_UserArgument[T any] struct {
	Id     string
	Name   string
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
		"grpc.federation.private.UserArgument": {
			"id":   grpcfed.NewCELFieldType(celtypes.StringType, "Id"),
			"name": grpcfed.NewCELFieldType(celtypes.StringType, "Name"),
		},
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
		tracer:       otel.Tracer("federation.FederationService"),
		client:       &FederationServiceDependentClientSet{},
	}, nil
}

// GetPost implements "federation.FederationService/GetPost" method.
func (s *FederationService) GetPost(ctx context.Context, req *GetPostRequest) (res *GetPostResponse, e error) {
	ctx, span := s.tracer.Start(ctx, "federation.FederationService/GetPost")
	defer span.End()

	ctx = grpcfed.WithLogger(ctx, s.logger)
	defer func() {
		if r := recover(); r != nil {
			e = grpcfed.RecoverError(r, debug.Stack())
			grpcfed.OutputErrorLog(ctx, s.logger, e)
		}
	}()
	res, err := s.resolve_Federation_GetPostResponse(ctx, &Federation_GetPostResponseArgument[*FederationServiceDependentClientSet]{
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

// resolve_Federation_GetPostResponse resolve "federation.GetPostResponse" message.
func (s *FederationService) resolve_Federation_GetPostResponse(ctx context.Context, req *Federation_GetPostResponseArgument[*FederationServiceDependentClientSet]) (*GetPostResponse, error) {
	ctx, span := s.tracer.Start(ctx, "federation.GetPostResponse")
	defer span.End()

	s.logger.DebugContext(ctx, "resolve federation.GetPostResponse", slog.Any("message_args", s.logvalue_Federation_GetPostResponseArgument(req)))
	var (
		sg      singleflight.Group
		valueMu sync.RWMutex
		valueP  *Post
	)
	envOpts := []cel.EnvOption{cel.Variable(grpcfed.MessageArgumentVariableName, cel.ObjectType("grpc.federation.private.GetPostResponseArgument"))}
	evalValues := map[string]any{grpcfed.MessageArgumentVariableName: req}

	// This section's codes are generated by the following proto definition.
	/*
	   {
	     name: "p"
	     message: "Post"
	   }
	*/
	{
		valueIface, err, _ := sg.Do("p_federation.Post", func() (any, error) {
			valueMu.RLock()
			args := &Federation_PostArgument[*FederationServiceDependentClientSet]{
				Client: s.client,
			}
			valueMu.RUnlock()
			return s.resolve_Federation_Post(ctx, args)
		})
		if err != nil {
			return nil, err
		}
		value := valueIface.(*Post)
		valueMu.Lock()
		valueP = value // { name: "p", message: "Post" ... }
		envOpts = append(envOpts, cel.Variable("p", cel.ObjectType("federation.Post")))
		evalValues["p"] = valueP
		valueMu.Unlock()
	}

	// assign named parameters to message arguments to pass to the custom resolver.
	req.P = valueP

	// create a message value to be returned.
	ret := &GetPostResponse{}

	// field binding section.
	// (grpc.federation.field).by = "p"
	{
		value, err := grpcfed.EvalCEL(s.env, "p", envOpts, evalValues, reflect.TypeOf((*Post)(nil)))
		if err != nil {
			grpcfed.RecordErrorToSpan(ctx, err)
			return nil, err
		}
		ret.Post = value.(*Post)
	}

	s.logger.DebugContext(ctx, "resolved federation.GetPostResponse", slog.Any("federation.GetPostResponse", s.logvalue_Federation_GetPostResponse(ret)))
	return ret, nil
}

// resolve_Federation_Post resolve "federation.Post" message.
func (s *FederationService) resolve_Federation_Post(ctx context.Context, req *Federation_PostArgument[*FederationServiceDependentClientSet]) (*Post, error) {
	ctx, span := s.tracer.Start(ctx, "federation.Post")
	defer span.End()

	s.logger.DebugContext(ctx, "resolve federation.Post", slog.Any("message_args", s.logvalue_Federation_PostArgument(req)))
	var (
		sg      singleflight.Group
		valueMu sync.RWMutex
		valueU  *User
	)
	envOpts := []cel.EnvOption{cel.Variable(grpcfed.MessageArgumentVariableName, cel.ObjectType("grpc.federation.private.PostArgument"))}
	evalValues := map[string]any{grpcfed.MessageArgumentVariableName: req}

	// This section's codes are generated by the following proto definition.
	/*
	   {
	     name: "u"
	     message: "User"
	     args: [
	       { name: "id", string: "foo" },
	       { name: "name", string: "bar" }
	     ]
	   }
	*/
	{
		valueIface, err, _ := sg.Do("u_federation.User", func() (any, error) {
			valueMu.RLock()
			args := &Federation_UserArgument[*FederationServiceDependentClientSet]{
				Client: s.client,
				Id:     "foo", // { name: "id", string: "foo" }
				Name:   "bar", // { name: "name", string: "bar" }
			}
			valueMu.RUnlock()
			return s.resolve_Federation_User(ctx, args)
		})
		if err != nil {
			return nil, err
		}
		value := valueIface.(*User)
		valueMu.Lock()
		valueU = value // { name: "u", message: "User" ... }
		envOpts = append(envOpts, cel.Variable("u", cel.ObjectType("federation.User")))
		evalValues["u"] = valueU
		valueMu.Unlock()
	}

	// assign named parameters to message arguments to pass to the custom resolver.
	req.U = valueU

	// create a message value to be returned.
	ret := &Post{}

	// field binding section.
	ret.Id = "post-id"      // (grpc.federation.field).string = "post-id"
	ret.Title = "title"     // (grpc.federation.field).string = "title"
	ret.Content = "content" // (grpc.federation.field).string = "content"
	// (grpc.federation.field).by = "u"
	{
		value, err := grpcfed.EvalCEL(s.env, "u", envOpts, evalValues, reflect.TypeOf((*User)(nil)))
		if err != nil {
			grpcfed.RecordErrorToSpan(ctx, err)
			return nil, err
		}
		ret.User = value.(*User)
	}

	s.logger.DebugContext(ctx, "resolved federation.Post", slog.Any("federation.Post", s.logvalue_Federation_Post(ret)))
	return ret, nil
}

// resolve_Federation_User resolve "federation.User" message.
func (s *FederationService) resolve_Federation_User(ctx context.Context, req *Federation_UserArgument[*FederationServiceDependentClientSet]) (*User, error) {
	ctx, span := s.tracer.Start(ctx, "federation.User")
	defer span.End()

	s.logger.DebugContext(ctx, "resolve federation.User", slog.Any("message_args", s.logvalue_Federation_UserArgument(req)))
	envOpts := []cel.EnvOption{cel.Variable(grpcfed.MessageArgumentVariableName, cel.ObjectType("grpc.federation.private.UserArgument"))}
	evalValues := map[string]any{grpcfed.MessageArgumentVariableName: req}

	// create a message value to be returned.
	ret := &User{}

	// field binding section.
	// (grpc.federation.field).by = "$.id"
	{
		value, err := grpcfed.EvalCEL(s.env, "$.id", envOpts, evalValues, reflect.TypeOf(""))
		if err != nil {
			grpcfed.RecordErrorToSpan(ctx, err)
			return nil, err
		}
		ret.Id = value.(string)
	}
	// (grpc.federation.field).by = "$.name"
	{
		value, err := grpcfed.EvalCEL(s.env, "$.name", envOpts, evalValues, reflect.TypeOf(""))
		if err != nil {
			grpcfed.RecordErrorToSpan(ctx, err)
			return nil, err
		}
		ret.Name = value.(string)
	}

	s.logger.DebugContext(ctx, "resolved federation.User", slog.Any("federation.User", s.logvalue_Federation_User(ret)))
	return ret, nil
}

func (s *FederationService) logvalue_Federation_GetPostResponse(v *GetPostResponse) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.Any("post", s.logvalue_Federation_Post(v.GetPost())),
	)
}

func (s *FederationService) logvalue_Federation_GetPostResponseArgument(v *Federation_GetPostResponseArgument[*FederationServiceDependentClientSet]) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.String("id", v.Id),
	)
}

func (s *FederationService) logvalue_Federation_Post(v *Post) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.String("id", v.GetId()),
		slog.String("title", v.GetTitle()),
		slog.String("content", v.GetContent()),
		slog.Any("user", s.logvalue_Federation_User(v.GetUser())),
	)
}

func (s *FederationService) logvalue_Federation_PostArgument(v *Federation_PostArgument[*FederationServiceDependentClientSet]) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue()
}

func (s *FederationService) logvalue_Federation_User(v *User) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.String("id", v.GetId()),
		slog.String("name", v.GetName()),
	)
}

func (s *FederationService) logvalue_Federation_UserArgument(v *Federation_UserArgument[*FederationServiceDependentClientSet]) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.String("id", v.Id),
		slog.String("name", v.Name),
	)
}

// DebugServiceConfig configuration required to initialize the service that use GRPC Federation.
type DebugServiceConfig struct {
	// ErrorHandler Federation Service often needs to convert errors received from downstream services.
	// If an error occurs during method execution in the Federation Service, this error handler is called and the returned error is treated as a final error.
	ErrorHandler grpcfed.ErrorHandler
	// Logger sets the logger used to output Debug/Info/Error information.
	Logger *slog.Logger
}

// DebugServiceClientFactory provides a factory that creates the gRPC Client needed to invoke methods of the gRPC Service on which the Federation Service depends.
type DebugServiceClientFactory interface {
}

// DebugServiceClientConfig information set in `dependencies` of the `grpc.federation.service` option.
// Hints for creating a gRPC Client.
type DebugServiceClientConfig struct {
	// Service returns the name of the service on Protocol Buffers.
	Service string
	// Name is the value set for `name` in `dependencies` of the `grpc.federation.service` option.
	// It must be unique among the services on which the Federation Service depends.
	Name string
}

// DebugServiceDependentClientSet has a gRPC client for all services on which the federation service depends.
// This is provided as an argument when implementing the custom resolver.
type DebugServiceDependentClientSet struct {
}

// DebugServiceResolver provides an interface to directly implement message resolver and field resolver not defined in Protocol Buffers.
type DebugServiceResolver interface {
}

// DebugServiceUnimplementedResolver a structure implemented to satisfy the Resolver interface.
// An Unimplemented error is always returned.
// This is intended for use when there are many Resolver interfaces that do not need to be implemented,
// by embedding them in a resolver structure that you have created.
type DebugServiceUnimplementedResolver struct{}

// DebugService represents Federation Service.
type DebugService struct {
	*UnimplementedDebugServiceServer
	cfg          DebugServiceConfig
	logger       *slog.Logger
	errorHandler grpcfed.ErrorHandler
	env          *cel.Env
	tracer       trace.Tracer
	client       *DebugServiceDependentClientSet
}

// NewDebugService creates DebugService instance by DebugServiceConfig.
func NewDebugService(cfg DebugServiceConfig) (*DebugService, error) {
	logger := cfg.Logger
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	}
	errorHandler := cfg.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(ctx context.Context, methodName string, err error) error { return err }
	}
	celHelper := grpcfed.NewCELTypeHelper(map[string]map[string]*celtypes.FieldType{
		"grpc.federation.private.GetStatusResponseArgument": {},
		"grpc.federation.private.UserArgument": {
			"id":   grpcfed.NewCELFieldType(celtypes.StringType, "Id"),
			"name": grpcfed.NewCELFieldType(celtypes.StringType, "Name"),
		},
	})
	env, err := cel.NewCustomEnv(
		cel.StdLib(),
		cel.CustomTypeAdapter(celHelper.TypeAdapter()),
		cel.CustomTypeProvider(celHelper.TypeProvider()),
	)
	if err != nil {
		return nil, err
	}
	return &DebugService{
		cfg:          cfg,
		logger:       logger,
		errorHandler: errorHandler,
		env:          env,
		tracer:       otel.Tracer("federation.DebugService"),
		client:       &DebugServiceDependentClientSet{},
	}, nil
}

// GetStatus implements "federation.DebugService/GetStatus" method.
func (s *DebugService) GetStatus(ctx context.Context, req *GetStatusRequest) (res *GetStatusResponse, e error) {
	ctx, span := s.tracer.Start(ctx, "federation.DebugService/GetStatus")
	defer span.End()

	ctx = grpcfed.WithLogger(ctx, s.logger)
	defer func() {
		if r := recover(); r != nil {
			e = grpcfed.RecoverError(r, debug.Stack())
			grpcfed.OutputErrorLog(ctx, s.logger, e)
		}
	}()
	res, err := s.resolve_Federation_GetStatusResponse(ctx, &Federation_GetStatusResponseArgument[*DebugServiceDependentClientSet]{
		Client: s.client,
	})
	if err != nil {
		grpcfed.RecordErrorToSpan(ctx, err)
		grpcfed.OutputErrorLog(ctx, s.logger, err)
		return nil, err
	}
	return res, nil
}

// resolve_Federation_GetStatusResponse resolve "federation.GetStatusResponse" message.
func (s *DebugService) resolve_Federation_GetStatusResponse(ctx context.Context, req *Federation_GetStatusResponseArgument[*DebugServiceDependentClientSet]) (*GetStatusResponse, error) {
	ctx, span := s.tracer.Start(ctx, "federation.GetStatusResponse")
	defer span.End()

	s.logger.DebugContext(ctx, "resolve federation.GetStatusResponse", slog.Any("message_args", s.logvalue_Federation_GetStatusResponseArgument(req)))
	var (
		sg      singleflight.Group
		valueMu sync.RWMutex
		valueU  *User
	)
	envOpts := []cel.EnvOption{cel.Variable(grpcfed.MessageArgumentVariableName, cel.ObjectType("grpc.federation.private.GetStatusResponseArgument"))}
	evalValues := map[string]any{grpcfed.MessageArgumentVariableName: req}

	// This section's codes are generated by the following proto definition.
	/*
	   {
	     name: "u"
	     message: "User"
	     args: [
	       { name: "id", string: "xxxx" },
	       { name: "name", string: "yyyy" }
	     ]
	   }
	*/
	{
		valueIface, err, _ := sg.Do("u_federation.User", func() (any, error) {
			valueMu.RLock()
			args := &Federation_UserArgument[*DebugServiceDependentClientSet]{
				Client: s.client,
				Id:     "xxxx", // { name: "id", string: "xxxx" }
				Name:   "yyyy", // { name: "name", string: "yyyy" }
			}
			valueMu.RUnlock()
			return s.resolve_Federation_User(ctx, args)
		})
		if err != nil {
			return nil, err
		}
		value := valueIface.(*User)
		valueMu.Lock()
		valueU = value // { name: "u", message: "User" ... }
		envOpts = append(envOpts, cel.Variable("u", cel.ObjectType("federation.User")))
		evalValues["u"] = valueU
		valueMu.Unlock()
	}

	// assign named parameters to message arguments to pass to the custom resolver.
	req.U = valueU

	// create a message value to be returned.
	ret := &GetStatusResponse{}

	// field binding section.
	// (grpc.federation.field).by = "u"
	{
		value, err := grpcfed.EvalCEL(s.env, "u", envOpts, evalValues, reflect.TypeOf((*User)(nil)))
		if err != nil {
			grpcfed.RecordErrorToSpan(ctx, err)
			return nil, err
		}
		ret.User = value.(*User)
	}

	s.logger.DebugContext(ctx, "resolved federation.GetStatusResponse", slog.Any("federation.GetStatusResponse", s.logvalue_Federation_GetStatusResponse(ret)))
	return ret, nil
}

// resolve_Federation_User resolve "federation.User" message.
func (s *DebugService) resolve_Federation_User(ctx context.Context, req *Federation_UserArgument[*DebugServiceDependentClientSet]) (*User, error) {
	ctx, span := s.tracer.Start(ctx, "federation.User")
	defer span.End()

	s.logger.DebugContext(ctx, "resolve federation.User", slog.Any("message_args", s.logvalue_Federation_UserArgument(req)))
	envOpts := []cel.EnvOption{cel.Variable(grpcfed.MessageArgumentVariableName, cel.ObjectType("grpc.federation.private.UserArgument"))}
	evalValues := map[string]any{grpcfed.MessageArgumentVariableName: req}

	// create a message value to be returned.
	ret := &User{}

	// field binding section.
	// (grpc.federation.field).by = "$.id"
	{
		value, err := grpcfed.EvalCEL(s.env, "$.id", envOpts, evalValues, reflect.TypeOf(""))
		if err != nil {
			grpcfed.RecordErrorToSpan(ctx, err)
			return nil, err
		}
		ret.Id = value.(string)
	}
	// (grpc.federation.field).by = "$.name"
	{
		value, err := grpcfed.EvalCEL(s.env, "$.name", envOpts, evalValues, reflect.TypeOf(""))
		if err != nil {
			grpcfed.RecordErrorToSpan(ctx, err)
			return nil, err
		}
		ret.Name = value.(string)
	}

	s.logger.DebugContext(ctx, "resolved federation.User", slog.Any("federation.User", s.logvalue_Federation_User(ret)))
	return ret, nil
}

func (s *DebugService) logvalue_Federation_GetStatusResponse(v *GetStatusResponse) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.Any("user", s.logvalue_Federation_User(v.GetUser())),
	)
}

func (s *DebugService) logvalue_Federation_GetStatusResponseArgument(v *Federation_GetStatusResponseArgument[*DebugServiceDependentClientSet]) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue()
}

func (s *DebugService) logvalue_Federation_User(v *User) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.String("id", v.GetId()),
		slog.String("name", v.GetName()),
	)
}

func (s *DebugService) logvalue_Federation_UserArgument(v *Federation_UserArgument[*DebugServiceDependentClientSet]) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.String("id", v.Id),
		slog.String("name", v.Name),
	)
}
