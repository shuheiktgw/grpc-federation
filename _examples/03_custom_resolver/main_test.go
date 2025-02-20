package main_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"example/federation"
	"example/post"
	"example/user"
)

const bufSize = 1024

var (
	listener   *bufconn.Listener
	postClient post.PostServiceClient
	userClient user.UserServiceClient
)

type clientConfig struct{}

func (c *clientConfig) Post_PostServiceClient(cfg federation.FederationServiceClientConfig) (post.PostServiceClient, error) {
	return postClient, nil
}

func (c *clientConfig) User_UserServiceClient(cfg federation.FederationServiceClientConfig) (user.UserServiceClient, error) {
	return userClient, nil
}

type Resolver struct {
	federation.FederationServiceUnimplementedResolver
}

func (r *Resolver) Resolve_Federation_User(ctx context.Context, arg *federation.Federation_UserArgument) (*federation.User, error) {
	return &federation.User{
		Id:   arg.U.Id,
		Name: arg.U.Name,
	}, nil
}

func (r *Resolver) Resolve_Federation_Post_User(ctx context.Context, arg *federation.Federation_Post_UserArgument) (*federation.User, error) {
	return arg.Federation_PostArgument.User, nil
}

func (r *Resolver) Resolve_Federation_Unused(_ context.Context, _ *federation.Federation_UnusedArgument) (*federation.Unused, error) {
	return &federation.Unused{}, nil
}

func (r *Resolver) Resolve_Federation_ForNameless(_ context.Context, _ *federation.Federation_ForNamelessArgument) (*federation.ForNameless, error) {
	return &federation.ForNameless{}, nil
}

func (r *Resolver) Resolve_Federation_User_Name(ctx context.Context, arg *federation.Federation_User_NameArgument) (string, error) {
	return arg.Federation_User.Name, nil
}

type PostServer struct {
	*post.UnimplementedPostServiceServer
}

func (s *PostServer) GetPost(ctx context.Context, req *post.GetPostRequest) (*post.GetPostResponse, error) {
	return nil, errors.New("error!!")
}

func (s *PostServer) GetPosts(ctx context.Context, req *post.GetPostsRequest) (*post.GetPostsResponse, error) {
	return nil, errors.New("error!!")
}

type UserServer struct {
	*user.UnimplementedUserServiceServer
}

func (s *UserServer) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	if req.Id == "" {
		return &user.GetUserResponse{User: &user.User{Id: "anonymous_id", Name: "anonymous"}}, nil
	}
	return &user.GetUserResponse{
		User: &user.User{
			Id:   req.Id,
			Name: fmt.Sprintf("name_%s", req.Id),
		},
	}, nil
}

func (s *UserServer) GetUsers(ctx context.Context, req *user.GetUsersRequest) (*user.GetUsersResponse, error) {
	var users []*user.User
	for _, id := range req.Ids {
		users = append(users, &user.User{
			Id:   id,
			Name: fmt.Sprintf("name_%s", id),
		})
	}
	return &user.GetUsersResponse{Users: users}, nil
}

func dialer(ctx context.Context, address string) (net.Conn, error) {
	return listener.Dial()
}

func TestFederation(t *testing.T) {
	ctx := context.Background()
	listener = bufconn.Listen(bufSize)

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	postClient = post.NewPostServiceClient(conn)
	userClient = user.NewUserServiceClient(conn)

	grpcServer := grpc.NewServer()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	federationServer, err := federation.NewFederationService(federation.FederationServiceConfig{
		Client:   new(clientConfig),
		Resolver: new(Resolver),
		Logger:   logger,
		ErrorHandler: func(ctx context.Context, methodName string, err error) error {
			federationServiceMethodName, _ := grpc.Method(ctx)
			log.Printf("FederationServiceMethod: %s DependentMethod: %s", federationServiceMethodName, methodName)
			switch methodName {
			case federation.FederationService_DependentMethod_Post_PostService_GetPost:
				return nil
			case federation.FederationService_DependentMethod_User_UserService_GetUser:
				return err
			}
			return err
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	post.RegisterPostServiceServer(grpcServer, &PostServer{})
	user.RegisterUserServiceServer(grpcServer, &UserServer{})
	federation.RegisterFederationServiceServer(grpcServer, federationServer)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			t.Fatal(err)
		}
	}()

	client := federation.NewFederationServiceClient(conn)
	res, err := client.GetPost(ctx, &federation.GetPostRequest{
		Id: "foo",
	})
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(res, &federation.GetPostResponse{
		Post: &federation.Post{
			User: &federation.User{
				Id:   "anonymous_id",
				Name: "anonymous",
			},
		},
	}, cmpopts.IgnoreUnexported(federation.GetPostResponse{}, federation.Post{}, federation.User{})); diff != "" {
		t.Errorf("(-got, +want)\n%s", diff)
	}
}
