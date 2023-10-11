// Code generated by protoc-gen-grpc-federation. DO NOT EDIT!
package federation

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/cel-go/cel"
	celtypes "github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/singleflight"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	post "example/post"
)

// FederationServiceConfig configuration required to initialize the service that use GRPC Federation.
type FederationServiceConfig struct {
	// Client provides a factory that creates the gRPC Client needed to invoke methods of the gRPC Service on which the Federation Service depends.
	// If this interface is not provided, an error is returned during initialization.
	Client FederationServiceClientFactory // required
	// ErrorHandler Federation Service often needs to convert errors received from downstream services.
	// If an error occurs during method execution in the Federation Service, this error handler is called and the returned error is treated as a final error.
	ErrorHandler FederationServiceErrorHandler
	// Logger sets the logger used to output Debug/Info/Error information.
	Logger *slog.Logger
}

// FederationServiceClientFactory provides a factory that creates the gRPC Client needed to invoke methods of the gRPC Service on which the Federation Service depends.
type FederationServiceClientFactory interface {
	// Org_Post_PostServiceClient create a gRPC Client to be used to call methods in org.post.PostService.
	Org_Post_PostServiceClient(FederationServiceClientConfig) (post.PostServiceClient, error)
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

// FederationServiceDependencyServiceClient has a gRPC client for all services on which the federation service depends.
// This is provided as an argument when implementing the custom resolver.
type FederationServiceDependencyServiceClient struct {
	Org_Post_PostServiceClient post.PostServiceClient
}

// FederationServiceResolver provides an interface to directly implement message resolver and field resolver not defined in Protocol Buffers.
type FederationServiceResolver interface {
}

// FederationServiceUnimplementedResolver a structure implemented to satisfy the Resolver interface.
// An Unimplemented error is always returned.
// This is intended for use when there are many Resolver interfaces that do not need to be implemented,
// by embedding them in a resolver structure that you have created.
type FederationServiceUnimplementedResolver struct{}

// FederationServiceErrorHandler Federation Service often needs to convert errors received from downstream services.
// If an error occurs during method execution in the Federation Service, this error handler is called and the returned error is treated as a final error.
type FederationServiceErrorHandler func(ctx context.Context, methodName string, err error) error

const (
	FederationService_DependentMethod_Org_Post_PostService_GetPost = "/org.post.PostService/GetPost"
)

// FederationServiceRecoveredError represents recovered error.
type FederationServiceRecoveredError struct {
	Message string
	Stack   []string
}

func (e *FederationServiceRecoveredError) Error() string {
	return fmt.Sprintf("recovered error: %s", e.Message)
}

// FederationService represents Federation Service.
type FederationService struct {
	*UnimplementedFederationServiceServer
	cfg          FederationServiceConfig
	logger       *slog.Logger
	errorHandler FederationServiceErrorHandler
	env          *cel.Env
	client       *FederationServiceDependencyServiceClient
}

// Org_Federation_GetPostResponseArgument is argument for "org.federation.GetPostResponse" message.
type Org_Federation_GetPostResponseArgument struct {
	Id                  string
	XOrgFederation_Post *Post
	Client              *FederationServiceDependencyServiceClient
}

// Org_Federation_PostArgument is argument for "org.federation.Post" message.
type Org_Federation_PostArgument struct {
	Id                           string
	XOrgFederation_User          *User
	XOrgPost_PostService_GetPost *post.Post
	Client                       *FederationServiceDependencyServiceClient
}

// Org_Federation_UserArgument is argument for "org.federation.User" message.
type Org_Federation_UserArgument struct {
	UserId string
	Client *FederationServiceDependencyServiceClient
}

// FederationServiceCELTypeHelper
type FederationServiceCELTypeHelper struct {
	celRegistry    *celtypes.Registry
	structFieldMap map[string]map[string]*celtypes.FieldType
	mapMu          sync.Mutex
}

func (h *FederationServiceCELTypeHelper) TypeProvider() celtypes.Provider {
	return h
}

func (h *FederationServiceCELTypeHelper) TypeAdapter() celtypes.Adapter {
	return h.celRegistry
}

func (h *FederationServiceCELTypeHelper) EnumValue(enumName string) ref.Val {
	return h.celRegistry.EnumValue(enumName)
}

func (h *FederationServiceCELTypeHelper) FindIdent(identName string) (ref.Val, bool) {
	return h.celRegistry.FindIdent(identName)
}

func (h *FederationServiceCELTypeHelper) FindStructType(structType string) (*celtypes.Type, bool) {
	if st, found := h.celRegistry.FindStructType(structType); found {
		return st, found
	}
	h.mapMu.Lock()
	defer h.mapMu.Unlock()
	if _, exists := h.structFieldMap[structType]; exists {
		return celtypes.NewObjectType(structType), true
	}
	return nil, false
}

func (h *FederationServiceCELTypeHelper) FindStructFieldNames(structType string) ([]string, bool) {
	if names, found := h.celRegistry.FindStructFieldNames(structType); found {
		return names, found
	}

	h.mapMu.Lock()
	defer h.mapMu.Unlock()
	fieldMap, exists := h.structFieldMap[structType]
	if !exists {
		return nil, false
	}
	fieldNames := make([]string, 0, len(fieldMap))
	for fieldName := range fieldMap {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)
	return fieldNames, true
}

func (h *FederationServiceCELTypeHelper) FindStructFieldType(structType, fieldName string) (*celtypes.FieldType, bool) {
	if field, found := h.celRegistry.FindStructFieldType(structType, fieldName); found {
		return field, found
	}

	h.mapMu.Lock()
	defer h.mapMu.Unlock()
	fieldMap, exists := h.structFieldMap[structType]
	if !exists {
		return nil, false
	}
	field, found := fieldMap[fieldName]
	return field, found
}

func (h *FederationServiceCELTypeHelper) NewValue(structType string, fields map[string]ref.Val) ref.Val {
	return h.celRegistry.NewValue(structType, fields)
}

func newFederationServiceCELTypeHelper() *FederationServiceCELTypeHelper {
	celRegistry := celtypes.NewEmptyRegistry()
	protoregistry.GlobalFiles.RangeFiles(func(f protoreflect.FileDescriptor) bool {
		if err := celRegistry.RegisterDescriptor(f); err != nil {
			return false
		}
		return true
	})
	newFieldType := func(typ *celtypes.Type, fieldName string) *celtypes.FieldType {
		isSet := func(v any, fieldName string) bool {
			rv := reflect.ValueOf(v)
			if rv.Kind() == reflect.Pointer {
				rv = rv.Elem()
			}
			if rv.Kind() != reflect.Struct {
				return false
			}
			return rv.FieldByName(fieldName).IsValid()
		}
		getFrom := func(v any, fieldName string) (any, error) {
			rv := reflect.ValueOf(v)
			if rv.Kind() == reflect.Pointer {
				rv = rv.Elem()
			}
			if rv.Kind() != reflect.Struct {
				return nil, fmt.Errorf("%T is not struct type", v)
			}
			value := rv.FieldByName(fieldName)
			return value.Interface(), nil
		}
		return &celtypes.FieldType{
			Type: typ,
			IsSet: func(v any) bool {
				return isSet(v, fieldName)
			},
			GetFrom: func(v any) (any, error) {
				return getFrom(v, fieldName)
			},
		}
	}
	return &FederationServiceCELTypeHelper{
		celRegistry: celRegistry,
		structFieldMap: map[string]map[string]*celtypes.FieldType{
			"grpc.federation.private.GetPostResponseArgument": map[string]*celtypes.FieldType{
				"id": newFieldType(celtypes.StringType, "Id"),
			},
			"grpc.federation.private.PostArgument": map[string]*celtypes.FieldType{
				"id": newFieldType(celtypes.StringType, "Id"),
			},
			"grpc.federation.private.UserArgument": map[string]*celtypes.FieldType{
				"user_id": newFieldType(celtypes.StringType, "UserId"),
			},
		},
	}
}

// NewFederationService creates FederationService instance by FederationServiceConfig.
func NewFederationService(cfg FederationServiceConfig) (*FederationService, error) {
	if err := validateFederationServiceConfig(cfg); err != nil {
		return nil, err
	}
	Org_Post_PostServiceClient, err := cfg.Client.Org_Post_PostServiceClient(FederationServiceClientConfig{
		Service: "org.post.PostService",
		Name:    "",
	})
	if err != nil {
		return nil, err
	}
	logger := cfg.Logger
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	}
	errorHandler := cfg.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(ctx context.Context, methodName string, err error) error { return err }
	}
	celHelper := newFederationServiceCELTypeHelper()
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
		client: &FederationServiceDependencyServiceClient{
			Org_Post_PostServiceClient: Org_Post_PostServiceClient,
		},
	}, nil
}

func validateFederationServiceConfig(cfg FederationServiceConfig) error {
	if cfg.Client == nil {
		return fmt.Errorf("Client field in FederationServiceConfig is not set. this field must be set")
	}
	return nil
}

func withTimeoutFederationService[T any](ctx context.Context, method string, timeout time.Duration, fn func(context.Context) (*T, error)) (*T, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var (
		ret   *T
		errch = make(chan error)
	)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errch <- recoverErrorFederationService(r, debug.Stack())
			}
		}()

		res, err := fn(ctx)
		ret = res
		errch <- err
	}()
	select {
	case <-ctx.Done():
		status := grpcstatus.New(grpccodes.DeadlineExceeded, ctx.Err().Error())
		withDetails, err := status.WithDetails(&errdetails.ErrorInfo{
			Metadata: map[string]string{
				"method":  method,
				"timeout": timeout.String(),
			},
		})
		if err != nil {
			return nil, status.Err()
		}
		return nil, withDetails.Err()
	case err := <-errch:
		return ret, err
	}
}

func withRetryFederationService[T any](b backoff.BackOff, fn func() (*T, error)) (*T, error) {
	var res *T
	if err := backoff.Retry(func() (err error) {
		res, err = fn()
		return
	}, b); err != nil {
		return nil, err
	}
	return res, nil
}

func recoverErrorFederationService(v interface{}, rawStack []byte) *FederationServiceRecoveredError {
	msg := fmt.Sprint(v)
	lines := strings.Split(msg, "\n")
	if len(lines) <= 1 {
		lines := strings.Split(string(rawStack), "\n")
		stack := make([]string, 0, len(lines))
		for _, line := range lines {
			if line == "" {
				continue
			}
			stack = append(stack, strings.TrimPrefix(line, "\t"))
		}
		return &FederationServiceRecoveredError{
			Message: msg,
			Stack:   stack,
		}
	}
	// If panic occurs under singleflight, singleflight's recover catches the error and gives a stack trace.
	// Therefore, once the stack trace is removed.
	stack := make([]string, 0, len(lines))
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		stack = append(stack, strings.TrimPrefix(line, "\t"))
	}
	return &FederationServiceRecoveredError{
		Message: lines[0],
		Stack:   stack,
	}
}

func (s *FederationService) evalCEL(expr string, vars []cel.EnvOption, args map[string]any, outType reflect.Type) (any, error) {
	env, err := s.env.Extend(vars...)
	if err != nil {
		return nil, err
	}
	expr = strings.Replace(expr, "$", "__ARG__", -1)
	ast, iss := env.Compile(expr)
	if iss.Err() != nil {
		return nil, iss.Err()
	}
	program, err := env.Program(ast)
	if err != nil {
		return nil, err
	}
	out, _, err := program.Eval(args)
	if err != nil {
		return nil, err
	}
	if outType != nil {
		return out.ConvertToNative(outType)
	}
	return out.Value(), nil
}

func (s *FederationService) goWithRecover(eg *errgroup.Group, fn func() (interface{}, error)) {
	eg.Go(func() (e error) {
		defer func() {
			if r := recover(); r != nil {
				e = recoverErrorFederationService(r, debug.Stack())
			}
		}()
		_, err := fn()
		return err
	})
}

func (s *FederationService) outputErrorLog(ctx context.Context, err error) {
	if err == nil {
		return
	}
	if status, ok := grpcstatus.FromError(err); ok {
		s.logger.ErrorContext(ctx, status.Message(),
			slog.Group("grpc_status",
				slog.String("code", status.Code().String()),
				slog.Any("details", status.Details()),
			),
		)
		return
	}
	var recoveredErr *FederationServiceRecoveredError
	if errors.As(err, &recoveredErr) {
		trace := make([]interface{}, 0, len(recoveredErr.Stack))
		for idx, stack := range recoveredErr.Stack {
			trace = append(trace, slog.String(fmt.Sprint(idx+1), stack))
		}
		s.logger.ErrorContext(ctx, recoveredErr.Message, slog.Group("stack_trace", trace...))
		return
	}
	s.logger.ErrorContext(ctx, err.Error())
}

// GetPost implements "org.federation.FederationService/GetPost" method.
func (s *FederationService) GetPost(ctx context.Context, req *GetPostRequest) (res *GetPostResponse, e error) {
	defer func() {
		if r := recover(); r != nil {
			e = recoverErrorFederationService(r, debug.Stack())
			s.outputErrorLog(ctx, e)
		}
	}()
	res, err := s.resolve_Org_Federation_GetPostResponse(ctx, &Org_Federation_GetPostResponseArgument{
		Client: s.client,
		Id:     req.Id,
	})
	if err != nil {
		s.outputErrorLog(ctx, err)
		return nil, err
	}
	return res, nil
}

// resolve_Org_Federation_GetPostResponse resolve "org.federation.GetPostResponse" message.
func (s *FederationService) resolve_Org_Federation_GetPostResponse(ctx context.Context, req *Org_Federation_GetPostResponseArgument) (*GetPostResponse, error) {
	s.logger.DebugContext(ctx, "resolve  org.federation.GetPostResponse", slog.Any("message_args", s.logvalue_Org_Federation_GetPostResponseArgument(req)))
	var (
		sg                       singleflight.Group
		valueMu                  sync.RWMutex
		value_OrgFederation_Post *Post
	)
	envOpts := []cel.EnvOption{cel.Variable("__ARG__", cel.ObjectType("grpc.federation.private.GetPostResponseArgument"))}
	evalValues := map[string]any{"__ARG__": req}

	// This section's codes are generated by the following proto definition.
	/*
	   {
	     name: "_org_federation_Post"
	     message: "Post"
	     args { name: "id", by: "$.id" }
	     autobind: true
	   }
	*/
	resPostIface, err, _ := sg.Do("_org_federation_Post_org.federation.Post", func() (interface{}, error) {
		valueMu.RLock()
		args := &Org_Federation_PostArgument{
			Client: s.client,
		}
		// { name: "id", by: "$.id" }
		{
			_value, err := s.evalCEL("$.id", envOpts, evalValues, reflect.TypeOf(args.Id))
			if err != nil {
				return nil, err
			}
			args.Id = _value.(string)
		}
		valueMu.RUnlock()
		return s.resolve_Org_Federation_Post(ctx, args)
	})
	if err != nil {
		return nil, err
	}
	resPost := resPostIface.(*Post)
	valueMu.Lock()
	value_OrgFederation_Post = resPost // { name: "_org_federation_Post", message: "Post" ... }
	envOpts = append(envOpts, cel.Variable("_org_federation_Post", cel.ObjectType("org.federation.Post")))
	evalValues["_org_federation_Post"] = value_OrgFederation_Post
	valueMu.Unlock()

	// assign named parameters to message arguments to pass to the custom resolver.
	req.XOrgFederation_Post = value_OrgFederation_Post

	// create a message value to be returned.
	ret := &GetPostResponse{}

	// field binding section.
	ret.Id = value_OrgFederation_Post.GetId()           // { name: "_org_federation_Post", autobind: true }
	ret.Title = value_OrgFederation_Post.GetTitle()     // { name: "_org_federation_Post", autobind: true }
	ret.Content = value_OrgFederation_Post.GetContent() // { name: "_org_federation_Post", autobind: true }

	s.logger.DebugContext(ctx, "resolved org.federation.GetPostResponse", slog.Any("org.federation.GetPostResponse", s.logvalue_Org_Federation_GetPostResponse(ret)))
	return ret, nil
}

// resolve_Org_Federation_Post resolve "org.federation.Post" message.
func (s *FederationService) resolve_Org_Federation_Post(ctx context.Context, req *Org_Federation_PostArgument) (*Post, error) {
	s.logger.DebugContext(ctx, "resolve  org.federation.Post", slog.Any("message_args", s.logvalue_Org_Federation_PostArgument(req)))
	var (
		sg                                singleflight.Group
		valueMu                           sync.RWMutex
		value_OrgFederation_User          *User
		value_OrgPost_PostService_GetPost *post.Post
	)
	envOpts := []cel.EnvOption{cel.Variable("__ARG__", cel.ObjectType("grpc.federation.private.PostArgument"))}
	evalValues := map[string]any{"__ARG__": req}
	// A tree view of message dependencies is shown below.
	/*
	   _org_federation_User ─┐
	                GetPost ─┤
	*/
	eg, ctx := errgroup.WithContext(ctx)

	s.goWithRecover(eg, func() (interface{}, error) {

		// This section's codes are generated by the following proto definition.
		/*
		   {
		     name: "_org_federation_User"
		     message: "User"
		     args { name: "user_id", string: "foo" }
		     autobind: true
		   }
		*/
		resUserIface, err, _ := sg.Do("_org_federation_User_org.federation.User", func() (interface{}, error) {
			valueMu.RLock()
			args := &Org_Federation_UserArgument{
				Client: s.client,
				UserId: "foo", // { name: "user_id", string: "foo" }
			}
			valueMu.RUnlock()
			return s.resolve_Org_Federation_User(ctx, args)
		})
		if err != nil {
			return nil, err
		}
		resUser := resUserIface.(*User)
		valueMu.Lock()
		value_OrgFederation_User = resUser // { name: "_org_federation_User", message: "User" ... }
		envOpts = append(envOpts, cel.Variable("_org_federation_User", cel.ObjectType("org.federation.User")))
		evalValues["_org_federation_User"] = value_OrgFederation_User
		valueMu.Unlock()
		return nil, nil
	})

	s.goWithRecover(eg, func() (interface{}, error) {

		// This section's codes are generated by the following proto definition.
		/*
		   resolver: {
		     method: "org.post.PostService/GetPost"
		     request { field: "id", by: "$.id" }
		     response { name: "_org_post_PostService_GetPost", field: "post", autobind: true }
		   }
		*/
		resGetPostResponseIface, err, _ := sg.Do("org.post.PostService/GetPost", func() (interface{}, error) {
			valueMu.RLock()
			args := &post.GetPostRequest{}
			// { field: "id", by: "$.id" }
			{
				_value, err := s.evalCEL("$.id", envOpts, evalValues, reflect.TypeOf(args.Id))
				if err != nil {
					return nil, err
				}
				args.Id = _value.(string)
			}
			valueMu.RUnlock()
			return s.client.Org_Post_PostServiceClient.GetPost(ctx, args)
		})
		if err != nil {
			if err := s.errorHandler(ctx, FederationService_DependentMethod_Org_Post_PostService_GetPost, err); err != nil {
				return nil, err
			}
		}
		resGetPostResponse := resGetPostResponseIface.(*post.GetPostResponse)
		valueMu.Lock()
		value_OrgPost_PostService_GetPost = resGetPostResponse.GetPost() // { name: "_org_post_PostService_GetPost", field: "post", autobind: true }
		envOpts = append(envOpts, cel.Variable("_org_post_PostService_GetPost", cel.ObjectType("org.post.Post")))
		evalValues["_org_post_PostService_GetPost"] = value_OrgPost_PostService_GetPost
		valueMu.Unlock()
		return nil, nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// assign named parameters to message arguments to pass to the custom resolver.
	req.XOrgFederation_User = value_OrgFederation_User
	req.XOrgPost_PostService_GetPost = value_OrgPost_PostService_GetPost

	// create a message value to be returned.
	ret := &Post{}

	// field binding section.
	ret.Id = value_OrgPost_PostService_GetPost.GetId()           // { name: "_org_post_PostService_GetPost", autobind: true }
	ret.Title = value_OrgPost_PostService_GetPost.GetTitle()     // { name: "_org_post_PostService_GetPost", autobind: true }
	ret.Content = value_OrgPost_PostService_GetPost.GetContent() // { name: "_org_post_PostService_GetPost", autobind: true }
	ret.Uid = value_OrgFederation_User.GetUid()                  // { name: "_org_federation_User", autobind: true }

	s.logger.DebugContext(ctx, "resolved org.federation.Post", slog.Any("org.federation.Post", s.logvalue_Org_Federation_Post(ret)))
	return ret, nil
}

// resolve_Org_Federation_User resolve "org.federation.User" message.
func (s *FederationService) resolve_Org_Federation_User(ctx context.Context, req *Org_Federation_UserArgument) (*User, error) {
	s.logger.DebugContext(ctx, "resolve  org.federation.User", slog.Any("message_args", s.logvalue_Org_Federation_UserArgument(req)))
	envOpts := []cel.EnvOption{cel.Variable("__ARG__", cel.ObjectType("grpc.federation.private.UserArgument"))}
	evalValues := map[string]any{"__ARG__": req}

	// create a message value to be returned.
	ret := &User{}

	// field binding section.
	// (grpc.federation.field).by = "$.user_id"
	{
		_value, err := s.evalCEL("$.user_id", envOpts, evalValues, reflect.TypeOf(ret.Uid))
		if err != nil {
			return nil, err
		}
		ret.Uid = _value.(string)
	}

	s.logger.DebugContext(ctx, "resolved org.federation.User", slog.Any("org.federation.User", s.logvalue_Org_Federation_User(ret)))
	return ret, nil
}

func (s *FederationService) logvalue_Org_Federation_GetPostResponse(v *GetPostResponse) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.String("id", v.GetId()),
		slog.String("title", v.GetTitle()),
		slog.String("content", v.GetContent()),
	)
}

func (s *FederationService) logvalue_Org_Federation_GetPostResponseArgument(v *Org_Federation_GetPostResponseArgument) slog.Value {
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
		slog.String("uid", v.GetUid()),
	)
}

func (s *FederationService) logvalue_Org_Federation_PostArgument(v *Org_Federation_PostArgument) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.String("id", v.Id),
	)
}

func (s *FederationService) logvalue_Org_Federation_User(v *User) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.String("uid", v.GetUid()),
	)
}

func (s *FederationService) logvalue_Org_Federation_UserArgument(v *Org_Federation_UserArgument) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.String("user_id", v.UserId),
	)
}
