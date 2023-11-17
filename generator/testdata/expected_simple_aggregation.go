// Code generated by protoc-gen-grpc-federation. DO NOT EDIT!
package federation

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"runtime/debug"
	"sync"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/cel-go/cel"
	celtypes "github.com/google/cel-go/common/types"
	grpcfed "github.com/mercari/grpc-federation/grpc/federation"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/singleflight"
	grpccodes "google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"

	post "example/post"
	user "example/user"
)

// Org_Federation_GetPostResponseArgument is argument for "org.federation.GetPostResponse" message.
type Org_Federation_GetPostResponseArgument[T any] struct {
	Id     string
	Post   *Post
	Client T
}

// Org_Federation_MArgument is argument for "org.federation.M" message.
type Org_Federation_MArgument[T any] struct {
	Client T
}

// Org_Federation_PostArgument is argument for "org.federation.Post" message.
type Org_Federation_PostArgument[T any] struct {
	Id     string
	M      *M
	Res    *post.GetPostResponse
	User   *User
	Client T
}

// Org_Federation_UserArgument is argument for "org.federation.User" message.
type Org_Federation_UserArgument[T any] struct {
	Content string
	Id      string
	Res     *user.GetUserResponse
	Title   string
	UserId  string
	Client  T
}

// Org_Federation_User_AgeArgument is custom resolver's argument for "age" field of "org.federation.User" message.
type Org_Federation_User_AgeArgument[T any] struct {
	*Org_Federation_UserArgument[T]
	Client T
}

// Org_Federation_ZArgument is argument for "org.federation.Z" message.
type Org_Federation_ZArgument[T any] struct {
	Client T
}

// FederationServiceConfig configuration required to initialize the service that use GRPC Federation.
type FederationServiceConfig struct {
	// Client provides a factory that creates the gRPC Client needed to invoke methods of the gRPC Service on which the Federation Service depends.
	// If this interface is not provided, an error is returned during initialization.
	Client FederationServiceClientFactory // required
	// Resolver provides an interface to directly implement message resolver and field resolver not defined in Protocol Buffers.
	// If this interface is not provided, an error is returned during initialization.
	Resolver FederationServiceResolver // required
	// ErrorHandler Federation Service often needs to convert errors received from downstream services.
	// If an error occurs during method execution in the Federation Service, this error handler is called and the returned error is treated as a final error.
	ErrorHandler grpcfed.ErrorHandler
	// Logger sets the logger used to output Debug/Info/Error information.
	Logger *slog.Logger
}

// FederationServiceClientFactory provides a factory that creates the gRPC Client needed to invoke methods of the gRPC Service on which the Federation Service depends.
type FederationServiceClientFactory interface {
	// Org_Post_PostServiceClient create a gRPC Client to be used to call methods in org.post.PostService.
	Org_Post_PostServiceClient(FederationServiceClientConfig) (post.PostServiceClient, error)
	// Org_User_UserServiceClient create a gRPC Client to be used to call methods in org.user.UserService.
	Org_User_UserServiceClient(FederationServiceClientConfig) (user.UserServiceClient, error)
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
	Org_Post_PostServiceClient post.PostServiceClient
	Org_User_UserServiceClient user.UserServiceClient
}

// FederationServiceResolver provides an interface to directly implement message resolver and field resolver not defined in Protocol Buffers.
type FederationServiceResolver interface {
	// Resolve_Org_Federation_User_Age implements resolver for "org.federation.User.age".
	Resolve_Org_Federation_User_Age(context.Context, *Org_Federation_User_AgeArgument[*FederationServiceDependentClientSet]) (uint64, error)
	// Resolve_Org_Federation_Z implements resolver for "org.federation.Z".
	Resolve_Org_Federation_Z(context.Context, *Org_Federation_ZArgument[*FederationServiceDependentClientSet]) (*Z, error)
}

// FederationServiceUnimplementedResolver a structure implemented to satisfy the Resolver interface.
// An Unimplemented error is always returned.
// This is intended for use when there are many Resolver interfaces that do not need to be implemented,
// by embedding them in a resolver structure that you have created.
type FederationServiceUnimplementedResolver struct{}

// Resolve_Org_Federation_User_Age resolve "org.federation.User.age".
// This method always returns Unimplemented error.
func (FederationServiceUnimplementedResolver) Resolve_Org_Federation_User_Age(context.Context, *Org_Federation_User_AgeArgument[*FederationServiceDependentClientSet]) (ret uint64, e error) {
	e = grpcstatus.Errorf(grpccodes.Unimplemented, "method Resolve_Org_Federation_User_Age not implemented")
	return
}

// Resolve_Org_Federation_Z resolve "org.federation.Z".
// This method always returns Unimplemented error.
func (FederationServiceUnimplementedResolver) Resolve_Org_Federation_Z(context.Context, *Org_Federation_ZArgument[*FederationServiceDependentClientSet]) (ret *Z, e error) {
	e = grpcstatus.Errorf(grpccodes.Unimplemented, "method Resolve_Org_Federation_Z not implemented")
	return
}

const (
	FederationService_DependentMethod_Org_Post_PostService_GetPost = "/org.post.PostService/GetPost"
	FederationService_DependentMethod_Org_User_UserService_GetUser = "/org.user.UserService/GetUser"
)

// FederationService represents Federation Service.
type FederationService struct {
	*UnimplementedFederationServiceServer
	cfg          FederationServiceConfig
	logger       *slog.Logger
	errorHandler grpcfed.ErrorHandler
	env          *cel.Env
	tracer       trace.Tracer
	resolver     FederationServiceResolver
	client       *FederationServiceDependentClientSet
}

// NewFederationService creates FederationService instance by FederationServiceConfig.
func NewFederationService(cfg FederationServiceConfig) (*FederationService, error) {
	if cfg.Client == nil {
		return nil, fmt.Errorf("Client field in FederationServiceConfig is not set. this field must be set")
	}
	if cfg.Resolver == nil {
		return nil, fmt.Errorf("Resolver field in FederationServiceConfig is not set. this field must be set")
	}
	Org_Post_PostServiceClient, err := cfg.Client.Org_Post_PostServiceClient(FederationServiceClientConfig{
		Service: "org.post.PostService",
		Name:    "",
	})
	if err != nil {
		return nil, err
	}
	Org_User_UserServiceClient, err := cfg.Client.Org_User_UserServiceClient(FederationServiceClientConfig{
		Service: "org.user.UserService",
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
	celHelper := grpcfed.NewCELTypeHelper(map[string]map[string]*celtypes.FieldType{
		"grpc.federation.private.GetPostResponseArgument": {
			"id": grpcfed.NewCELFieldType(celtypes.StringType, "Id"),
		},
		"grpc.federation.private.MArgument": {},
		"grpc.federation.private.PostArgument": {
			"id": grpcfed.NewCELFieldType(celtypes.StringType, "Id"),
		},
		"grpc.federation.private.UserArgument": {
			"id":      grpcfed.NewCELFieldType(celtypes.StringType, "Id"),
			"title":   grpcfed.NewCELFieldType(celtypes.StringType, "Title"),
			"content": grpcfed.NewCELFieldType(celtypes.StringType, "Content"),
			"user_id": grpcfed.NewCELFieldType(celtypes.StringType, "UserId"),
		},
		"grpc.federation.private.ZArgument": {},
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
		resolver:     cfg.Resolver,
		client: &FederationServiceDependentClientSet{
			Org_Post_PostServiceClient: Org_Post_PostServiceClient,
			Org_User_UserServiceClient: Org_User_UserServiceClient,
		},
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
	res, err := grpcfed.WithTimeout[GetPostResponse](ctx, "org.federation.FederationService/GetPost", 60000000000 /* 1m0s */, func(ctx context.Context) (*GetPostResponse, error) {
		return s.resolve_Org_Federation_GetPostResponse(ctx, &Org_Federation_GetPostResponseArgument[*FederationServiceDependentClientSet]{
			Client: s.client,
			Id:     req.Id,
		})
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

	// This section's codes are generated by the following proto definition.
	/*
	   def {
	     name: "post"
	     message {
	       name: "Post"
	       args { name: "id", by: "$.id" }
	     }
	   }
	*/
	{
		valueIface, err, _ := sg.Do("post", func() (any, error) {
			valueMu.RLock()
			args := &Org_Federation_PostArgument[*FederationServiceDependentClientSet]{
				Client: s.client,
			}
			// { name: "id", by: "$.id" }
			{
				value, err := grpcfed.EvalCEL(s.env, "$.id", envOpts, evalValues, reflect.TypeOf(""))
				if err != nil {
					grpcfed.RecordErrorToSpan(ctx, err)
					return nil, err
				}
				args.Id = value.(string)
			}
			valueMu.RUnlock()
			return s.resolve_Org_Federation_Post(ctx, args)
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
	ret.Const = "foo" // (grpc.federation.field).string = "foo"

	s.logger.DebugContext(ctx, "resolved org.federation.GetPostResponse", slog.Any("org.federation.GetPostResponse", s.logvalue_Org_Federation_GetPostResponse(ret)))
	return ret, nil
}

// resolve_Org_Federation_M resolve "org.federation.M" message.
func (s *FederationService) resolve_Org_Federation_M(ctx context.Context, req *Org_Federation_MArgument[*FederationServiceDependentClientSet]) (*M, error) {
	ctx, span := s.tracer.Start(ctx, "org.federation.M")
	defer span.End()

	s.logger.DebugContext(ctx, "resolve org.federation.M", slog.Any("message_args", s.logvalue_Org_Federation_MArgument(req)))

	// create a message value to be returned.
	ret := &M{}

	// field binding section.
	ret.Foo = "foo" // (grpc.federation.field).string = "foo"
	ret.Bar = 1     // (grpc.federation.field).int64 = 1

	s.logger.DebugContext(ctx, "resolved org.federation.M", slog.Any("org.federation.M", s.logvalue_Org_Federation_M(ret)))
	return ret, nil
}

// resolve_Org_Federation_Post resolve "org.federation.Post" message.
func (s *FederationService) resolve_Org_Federation_Post(ctx context.Context, req *Org_Federation_PostArgument[*FederationServiceDependentClientSet]) (*Post, error) {
	ctx, span := s.tracer.Start(ctx, "org.federation.Post")
	defer span.End()

	s.logger.DebugContext(ctx, "resolve org.federation.Post", slog.Any("message_args", s.logvalue_Org_Federation_PostArgument(req)))
	var (
		sg        singleflight.Group
		valueM    *M
		valueMu   sync.RWMutex
		valuePost *post.Post
		valueRes  *post.GetPostResponse
		valueUser *User
	)
	envOpts := []cel.EnvOption{cel.Variable(grpcfed.MessageArgumentVariableName, cel.ObjectType("grpc.federation.private.PostArgument"))}
	evalValues := map[string]any{grpcfed.MessageArgumentVariableName: req}
	// A tree view of message dependencies is shown below.
	/*
	                 m ─┐
	   res ─┐           │
	        post ─┐     │
	              user ─┤
	                 z ─┤
	*/
	eg, ctx1 := errgroup.WithContext(ctx)

	grpcfed.GoWithRecover(eg, func() (any, error) {

		// This section's codes are generated by the following proto definition.
		/*
		   def {
		     name: "m"
		     autobind: true
		     message {
		       name: "M"
		     }
		   }
		*/
		{
			valueIface, err, _ := sg.Do("m", func() (any, error) {
				valueMu.RLock()
				args := &Org_Federation_MArgument[*FederationServiceDependentClientSet]{
					Client: s.client,
				}
				valueMu.RUnlock()
				return s.resolve_Org_Federation_M(ctx1, args)
			})
			if err != nil {
				return nil, err
			}
			value := valueIface.(*M)
			valueMu.Lock()
			valueM = value // { name: "m", message: "M" ... }
			envOpts = append(envOpts, cel.Variable("m", cel.ObjectType("org.federation.M")))
			evalValues["m"] = valueM
			valueMu.Unlock()
		}
		return nil, nil
	})

	grpcfed.GoWithRecover(eg, func() (any, error) {

		// This section's codes are generated by the following proto definition.
		/*
		   def {
		     name: "res"
		     call {
		       method: "org.post.PostService/GetPost"
		       request { field: "id", by: "$.id" }
		     }
		   }
		*/
		{
			valueIface, err, _ := sg.Do("res", func() (any, error) {
				valueMu.RLock()
				args := &post.GetPostRequest{}
				// { field: "id", by: "$.id" }
				{
					value, err := grpcfed.EvalCEL(s.env, "$.id", envOpts, evalValues, reflect.TypeOf(""))
					if err != nil {
						grpcfed.RecordErrorToSpan(ctx, err)
						return nil, err
					}
					args.Id = value.(string)
				}
				valueMu.RUnlock()
				return grpcfed.WithTimeout[post.GetPostResponse](ctx1, "org.post.PostService/GetPost", 10000000000 /* 10s */, func(ctx context.Context) (*post.GetPostResponse, error) {
					var b backoff.BackOff = backoff.NewConstantBackOff(2000000000 /* 2s */)
					b = backoff.WithMaxRetries(b, 3)
					b = backoff.WithContext(b, ctx1)
					return grpcfed.WithRetry[post.GetPostResponse](b, func() (*post.GetPostResponse, error) {
						return s.client.Org_Post_PostServiceClient.GetPost(ctx1, args)
					})
				})
			})
			if err != nil {
				if err := s.errorHandler(ctx1, FederationService_DependentMethod_Org_Post_PostService_GetPost, err); err != nil {
					grpcfed.RecordErrorToSpan(ctx, err)
					return nil, err
				}
			}
			value := valueIface.(*post.GetPostResponse)
			valueMu.Lock()
			valueRes = value
			envOpts = append(envOpts, cel.Variable("res", cel.ObjectType("org.post.GetPostResponse")))
			evalValues["res"] = valueRes
			valueMu.Unlock()
		}

		// This section's codes are generated by the following proto definition.
		/*
		   def {
		     name: "post"
		     autobind: true
		     by: "res.post"
		   }
		*/
		{
			valueIface, err, _ := sg.Do("post", func() (any, error) {
				valueMu.RLock()
				valueMu.RUnlock()
				return grpcfed.EvalCEL(s.env, "res.post", envOpts, evalValues, reflect.TypeOf((*post.Post)(nil)))
			})
			if err != nil {
				return nil, err
			}
			value := valueIface.(*post.Post)
			valueMu.Lock()
			valuePost = value
			envOpts = append(envOpts, cel.Variable("post", cel.ObjectType("org.post.Post")))
			evalValues["post"] = valuePost
			valueMu.Unlock()
		}

		// This section's codes are generated by the following proto definition.
		/*
		   def {
		     name: "user"
		     message {
		       name: "User"
		       args { inline: "post" }
		     }
		   }
		*/
		{
			valueIface, err, _ := sg.Do("user", func() (any, error) {
				valueMu.RLock()
				args := &Org_Federation_UserArgument[*FederationServiceDependentClientSet]{
					Client: s.client,
				}
				// { inline: "post" }
				{
					value, err := grpcfed.EvalCEL(s.env, "post", envOpts, evalValues, reflect.TypeOf((*post.Post)(nil)))
					if err != nil {
						grpcfed.RecordErrorToSpan(ctx, err)
						return nil, err
					}
					inlineValue := value.(*post.Post)
					args.Id = inlineValue.GetId()
					args.Title = inlineValue.GetTitle()
					args.Content = inlineValue.GetContent()
					args.UserId = inlineValue.GetUserId()
				}
				valueMu.RUnlock()
				return s.resolve_Org_Federation_User(ctx1, args)
			})
			if err != nil {
				return nil, err
			}
			value := valueIface.(*User)
			valueMu.Lock()
			valueUser = value // { name: "user", message: "User" ... }
			envOpts = append(envOpts, cel.Variable("user", cel.ObjectType("org.federation.User")))
			evalValues["user"] = valueUser
			valueMu.Unlock()
		}
		return nil, nil
	})

	grpcfed.GoWithRecover(eg, func() (any, error) {

		// This section's codes are generated by the following proto definition.
		/*
		   def {
		     name: "z"
		     message {
		       name: "Z"
		     }
		   }
		*/
		{
			if _, err, _ := sg.Do("z", func() (any, error) {
				valueMu.RLock()
				args := &Org_Federation_ZArgument[*FederationServiceDependentClientSet]{
					Client: s.client,
				}
				valueMu.RUnlock()
				return s.resolve_Org_Federation_Z(ctx1, args)
			}); err != nil {
				return nil, err
			}
			valueMu.Lock()
			valueMu.Unlock()
		}
		return nil, nil
	})

	if err := eg.Wait(); err != nil {
		grpcfed.RecordErrorToSpan(ctx, err)
		return nil, err
	}

	// create a message value to be returned.
	ret := &Post{}

	// field binding section.
	ret.Id = valuePost.GetId()           // { name: "post", autobind: true }
	ret.Title = valuePost.GetTitle()     // { name: "post", autobind: true }
	ret.Content = valuePost.GetContent() // { name: "post", autobind: true }
	// (grpc.federation.field).by = "user"
	{
		value, err := grpcfed.EvalCEL(s.env, "user", envOpts, evalValues, reflect.TypeOf((*User)(nil)))
		if err != nil {
			grpcfed.RecordErrorToSpan(ctx, err)
			return nil, err
		}
		ret.User = value.(*User)
	}
	ret.Foo = valueM.GetFoo() // { name: "m", autobind: true }
	ret.Bar = valueM.GetBar() // { name: "m", autobind: true }

	s.logger.DebugContext(ctx, "resolved org.federation.Post", slog.Any("org.federation.Post", s.logvalue_Org_Federation_Post(ret)))
	return ret, nil
}

// resolve_Org_Federation_User resolve "org.federation.User" message.
func (s *FederationService) resolve_Org_Federation_User(ctx context.Context, req *Org_Federation_UserArgument[*FederationServiceDependentClientSet]) (*User, error) {
	ctx, span := s.tracer.Start(ctx, "org.federation.User")
	defer span.End()

	s.logger.DebugContext(ctx, "resolve org.federation.User", slog.Any("message_args", s.logvalue_Org_Federation_UserArgument(req)))
	var (
		sg        singleflight.Group
		valueMu   sync.RWMutex
		valueRes  *user.GetUserResponse
		valueUser *user.User
	)
	envOpts := []cel.EnvOption{cel.Variable(grpcfed.MessageArgumentVariableName, cel.ObjectType("grpc.federation.private.UserArgument"))}
	evalValues := map[string]any{grpcfed.MessageArgumentVariableName: req}

	// This section's codes are generated by the following proto definition.
	/*
	   def {
	     name: "res"
	     call {
	       method: "org.user.UserService/GetUser"
	       request { field: "id", by: "$.user_id" }
	     }
	   }
	*/
	{
		valueIface, err, _ := sg.Do("res", func() (any, error) {
			valueMu.RLock()
			args := &user.GetUserRequest{}
			// { field: "id", by: "$.user_id" }
			{
				value, err := grpcfed.EvalCEL(s.env, "$.user_id", envOpts, evalValues, reflect.TypeOf(""))
				if err != nil {
					grpcfed.RecordErrorToSpan(ctx, err)
					return nil, err
				}
				args.Id = value.(string)
			}
			valueMu.RUnlock()
			return grpcfed.WithTimeout[user.GetUserResponse](ctx, "org.user.UserService/GetUser", 20000000000 /* 20s */, func(ctx context.Context) (*user.GetUserResponse, error) {
				eb := backoff.NewExponentialBackOff()
				eb.InitialInterval = 1000000000 /* 1s */
				eb.RandomizationFactor = 0.7
				eb.Multiplier = 1.7
				eb.MaxInterval = 30000000000    /* 30s */
				eb.MaxElapsedTime = 20000000000 /* 20s */

				var b backoff.BackOff = eb
				b = backoff.WithMaxRetries(b, 3)
				b = backoff.WithContext(b, ctx)
				return grpcfed.WithRetry[user.GetUserResponse](b, func() (*user.GetUserResponse, error) {
					return s.client.Org_User_UserServiceClient.GetUser(ctx, args)
				})
			})
		})
		if err != nil {
			if err := s.errorHandler(ctx, FederationService_DependentMethod_Org_User_UserService_GetUser, err); err != nil {
				grpcfed.RecordErrorToSpan(ctx, err)
				return nil, err
			}
		}
		value := valueIface.(*user.GetUserResponse)
		valueMu.Lock()
		valueRes = value
		envOpts = append(envOpts, cel.Variable("res", cel.ObjectType("org.user.GetUserResponse")))
		evalValues["res"] = valueRes
		valueMu.Unlock()
	}

	// This section's codes are generated by the following proto definition.
	/*
	   def {
	     name: "user"
	     autobind: true
	     by: "res.user"
	   }
	*/
	{
		valueIface, err, _ := sg.Do("user", func() (any, error) {
			valueMu.RLock()
			valueMu.RUnlock()
			return grpcfed.EvalCEL(s.env, "res.user", envOpts, evalValues, reflect.TypeOf((*user.User)(nil)))
		})
		if err != nil {
			return nil, err
		}
		value := valueIface.(*user.User)
		valueMu.Lock()
		valueUser = value
		envOpts = append(envOpts, cel.Variable("user", cel.ObjectType("org.user.User")))
		evalValues["user"] = valueUser
		valueMu.Unlock()
	}

	// create a message value to be returned.
	ret := &User{}

	// field binding section.
	ret.Id = valueUser.GetId()                                                            // { name: "user", autobind: true }
	ret.Type = s.cast_Org_User_UserType__to__Org_Federation_UserType(valueUser.GetType()) // { name: "user", autobind: true }
	ret.Name = valueUser.GetName()                                                        // { name: "user", autobind: true }
	{
		// (grpc.federation.field).custom_resolver = true
		var err error
		ret.Age, err = s.resolver.Resolve_Org_Federation_User_Age(ctx, &Org_Federation_User_AgeArgument[*FederationServiceDependentClientSet]{
			Client:                      s.client,
			Org_Federation_UserArgument: req,
		})
		if err != nil {
			grpcfed.RecordErrorToSpan(ctx, err)
			return nil, err
		}
	}
	ret.Desc = valueUser.GetDesc()                                                                    // { name: "user", autobind: true }
	ret.MainItem = s.cast_Org_User_Item__to__Org_Federation_Item(valueUser.GetMainItem())             // { name: "user", autobind: true }
	ret.Items = s.cast_repeated_Org_User_Item__to__repeated_Org_Federation_Item(valueUser.GetItems()) // { name: "user", autobind: true }
	ret.Profile = valueUser.GetProfile()                                                              // { name: "user", autobind: true }

	switch {
	case s.cast_Org_User_User_AttrA___to__Org_Federation_User_AttrA_(valueUser.GetAttrA()) != nil:

		ret.Attr = s.cast_Org_User_User_AttrA___to__Org_Federation_User_AttrA_(valueUser.GetAttrA())
	case s.cast_Org_User_User_B__to__Org_Federation_User_B(valueUser.GetB()) != nil:

		ret.Attr = s.cast_Org_User_User_B__to__Org_Federation_User_B(valueUser.GetB())
	}

	s.logger.DebugContext(ctx, "resolved org.federation.User", slog.Any("org.federation.User", s.logvalue_Org_Federation_User(ret)))
	return ret, nil
}

// cast_Org_User_Item_ItemType__to__Org_Federation_Item_ItemType cast from "org.user.Item.ItemType" to "org.federation.Item.ItemType".
func (s *FederationService) cast_Org_User_Item_ItemType__to__Org_Federation_Item_ItemType(from user.Item_ItemType) Item_ItemType {
	switch from {
	case user.Item_ITEM_TYPE_1:
		return Item_ITEM_TYPE_1
	case user.Item_ITEM_TYPE_2:
		return Item_ITEM_TYPE_2
	case user.Item_ITEM_TYPE_3:
		return Item_ITEM_TYPE_3
	default:
		return 0
	}
}

// cast_Org_User_Item__to__Org_Federation_Item cast from "org.user.Item" to "org.federation.Item".
func (s *FederationService) cast_Org_User_Item__to__Org_Federation_Item(from *user.Item) *Item {
	if from == nil {
		return nil
	}

	return &Item{
		Name:  from.GetName(),
		Type:  s.cast_Org_User_Item_ItemType__to__Org_Federation_Item_ItemType(from.GetType()),
		Value: from.GetValue(),
	}
}

// cast_Org_User_User_AttrA__to__Org_Federation_User_AttrA cast from "org.user.User.AttrA" to "org.federation.User.AttrA".
func (s *FederationService) cast_Org_User_User_AttrA__to__Org_Federation_User_AttrA(from *user.User_AttrA) *User_AttrA {
	if from == nil {
		return nil
	}

	return &User_AttrA{
		Foo: from.GetFoo(),
	}
}

// cast_Org_User_User_AttrB__to__Org_Federation_User_AttrB cast from "org.user.User.AttrB" to "org.federation.User.AttrB".
func (s *FederationService) cast_Org_User_User_AttrB__to__Org_Federation_User_AttrB(from *user.User_AttrB) *User_AttrB {
	if from == nil {
		return nil
	}

	return &User_AttrB{
		Bar: from.GetBar(),
	}
}

// cast_Org_User_User_AttrA___to__Org_Federation_User_AttrA_ cast from "org.user.User.attr_a" to "org.federation.User.attr_a".
func (s *FederationService) cast_Org_User_User_AttrA___to__Org_Federation_User_AttrA_(from *user.User_AttrA) *User_AttrA_ {
	if from == nil {
		return nil
	}
	return &User_AttrA_{
		AttrA: s.cast_Org_User_User_AttrA__to__Org_Federation_User_AttrA(from),
	}
}

// cast_Org_User_User_B__to__Org_Federation_User_B cast from "org.user.User.b" to "org.federation.User.b".
func (s *FederationService) cast_Org_User_User_B__to__Org_Federation_User_B(from *user.User_AttrB) *User_B {
	if from == nil {
		return nil
	}
	return &User_B{
		B: s.cast_Org_User_User_AttrB__to__Org_Federation_User_AttrB(from),
	}
}

// cast_Org_User_UserType__to__Org_Federation_UserType cast from "org.user.UserType" to "org.federation.UserType".
func (s *FederationService) cast_Org_User_UserType__to__Org_Federation_UserType(from user.UserType) UserType {
	switch from {
	case user.UserType_USER_TYPE_1:
		return UserType_USER_TYPE_1
	case user.UserType_USER_TYPE_2:
		return UserType_USER_TYPE_2
	default:
		return 0
	}
}

// cast_repeated_Org_User_Item__to__repeated_Org_Federation_Item cast from "repeated org.user.Item" to "repeated org.federation.Item".
func (s *FederationService) cast_repeated_Org_User_Item__to__repeated_Org_Federation_Item(from []*user.Item) []*Item {
	ret := make([]*Item, 0, len(from))
	for _, v := range from {
		ret = append(ret, s.cast_Org_User_Item__to__Org_Federation_Item(v))
	}
	return ret
}

// resolve_Org_Federation_Z resolve "org.federation.Z" message.
func (s *FederationService) resolve_Org_Federation_Z(ctx context.Context, req *Org_Federation_ZArgument[*FederationServiceDependentClientSet]) (*Z, error) {
	ctx, span := s.tracer.Start(ctx, "org.federation.Z")
	defer span.End()

	s.logger.DebugContext(ctx, "resolve org.federation.Z", slog.Any("message_args", s.logvalue_Org_Federation_ZArgument(req)))

	// create a message value to be returned.
	// `custom_resolver = true` in "grpc.federation.message" option.
	ret, err := s.resolver.Resolve_Org_Federation_Z(ctx, req)
	if err != nil {
		grpcfed.RecordErrorToSpan(ctx, err)
		return nil, err
	}

	s.logger.DebugContext(ctx, "resolved org.federation.Z", slog.Any("org.federation.Z", s.logvalue_Org_Federation_Z(ret)))
	return ret, nil
}

func (s *FederationService) logvalue_Google_Protobuf_Any(v *anypb.Any) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.String("type_url", v.GetTypeUrl()),
		slog.String("value", string(v.GetValue())),
	)
}

func (s *FederationService) logvalue_Org_Federation_GetPostResponse(v *GetPostResponse) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.Any("post", s.logvalue_Org_Federation_Post(v.GetPost())),
		slog.String("const", v.GetConst()),
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

func (s *FederationService) logvalue_Org_Federation_Item(v *Item) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.String("name", v.GetName()),
		slog.String("type", s.logvalue_Org_Federation_Item_ItemType(v.GetType()).String()),
		slog.Int64("value", v.GetValue()),
	)
}

func (s *FederationService) logvalue_Org_Federation_Item_ItemType(v Item_ItemType) slog.Value {
	switch v {
	case Item_ITEM_TYPE_1:
		return slog.StringValue("ITEM_TYPE_1")
	case Item_ITEM_TYPE_2:
		return slog.StringValue("ITEM_TYPE_2")
	case Item_ITEM_TYPE_3:
		return slog.StringValue("ITEM_TYPE_3")
	}
	return slog.StringValue("")
}

func (s *FederationService) logvalue_Org_Federation_M(v *M) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.String("foo", v.GetFoo()),
		slog.Int64("bar", v.GetBar()),
	)
}

func (s *FederationService) logvalue_Org_Federation_MArgument(v *Org_Federation_MArgument[*FederationServiceDependentClientSet]) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue()
}

func (s *FederationService) logvalue_Org_Federation_Post(v *Post) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.String("id", v.GetId()),
		slog.String("title", v.GetTitle()),
		slog.String("content", v.GetContent()),
		slog.Any("user", s.logvalue_Org_Federation_User(v.GetUser())),
		slog.String("foo", v.GetFoo()),
		slog.Int64("bar", v.GetBar()),
	)
}

func (s *FederationService) logvalue_Org_Federation_PostArgument(v *Org_Federation_PostArgument[*FederationServiceDependentClientSet]) slog.Value {
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
		slog.String("id", v.GetId()),
		slog.String("type", s.logvalue_Org_Federation_UserType(v.GetType()).String()),
		slog.String("name", v.GetName()),
		slog.Uint64("age", v.GetAge()),
		slog.Any("desc", v.GetDesc()),
		slog.Any("main_item", s.logvalue_Org_Federation_Item(v.GetMainItem())),
		slog.Any("items", s.logvalue_repeated_Org_Federation_Item(v.GetItems())),
		slog.Any("profile", s.logvalue_Org_Federation_User_ProfileEntry(v.GetProfile())),
		slog.Any("attr_a", s.logvalue_Org_Federation_User_AttrA(v.GetAttrA())),
		slog.Any("b", s.logvalue_Org_Federation_User_AttrB(v.GetB())),
	)
}

func (s *FederationService) logvalue_Org_Federation_UserArgument(v *Org_Federation_UserArgument[*FederationServiceDependentClientSet]) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.String("id", v.Id),
		slog.String("title", v.Title),
		slog.String("content", v.Content),
		slog.String("user_id", v.UserId),
	)
}

func (s *FederationService) logvalue_Org_Federation_UserType(v UserType) slog.Value {
	switch v {
	case UserType_USER_TYPE_1:
		return slog.StringValue("USER_TYPE_1")
	case UserType_USER_TYPE_2:
		return slog.StringValue("USER_TYPE_2")
	}
	return slog.StringValue("")
}

func (s *FederationService) logvalue_Org_Federation_User_AttrA(v *User_AttrA) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.String("foo", v.GetFoo()),
	)
}

func (s *FederationService) logvalue_Org_Federation_User_AttrB(v *User_AttrB) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.Bool("bar", v.GetBar()),
	)
}

func (s *FederationService) logvalue_Org_Federation_User_ProfileEntry(v map[string]*anypb.Any) slog.Value {
	attrs := make([]slog.Attr, 0, len(v))
	for key, value := range v {
		attrs = append(attrs, slog.Attr{
			Key:   fmt.Sprint(key),
			Value: s.logvalue_Google_Protobuf_Any(value),
		})
	}
	return slog.GroupValue(attrs...)
}

func (s *FederationService) logvalue_Org_Federation_Z(v *Z) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue(
		slog.String("foo", v.GetFoo()),
	)
}

func (s *FederationService) logvalue_Org_Federation_ZArgument(v *Org_Federation_ZArgument[*FederationServiceDependentClientSet]) slog.Value {
	if v == nil {
		return slog.GroupValue()
	}
	return slog.GroupValue()
}

func (s *FederationService) logvalue_repeated_Org_Federation_Item(v []*Item) slog.Value {
	attrs := make([]slog.Attr, 0, len(v))
	for idx, vv := range v {
		attrs = append(attrs, slog.Attr{
			Key:   fmt.Sprint(idx),
			Value: s.logvalue_Org_Federation_Item(vv),
		})
	}
	return slog.GroupValue(attrs...)
}
