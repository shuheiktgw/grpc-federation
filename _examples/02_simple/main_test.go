package main_test

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/anypb"

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

type PostServer struct {
	*post.UnimplementedPostServiceServer
}

func (s *PostServer) GetPost(ctx context.Context, req *post.GetPostRequest) (*post.GetPostResponse, error) {
	return &post.GetPostResponse{
		Post: &post.Post{
			Id:      req.Id,
			Title:   "foo",
			Content: "bar",
			UserId:  fmt.Sprintf("user:%s", req.Id),
		},
	}, nil
}

func (s *PostServer) GetPosts(ctx context.Context, req *post.GetPostsRequest) (*post.GetPostsResponse, error) {
	var posts []*post.Post
	for _, id := range req.Ids {
		posts = append(posts, &post.Post{
			Id:      id,
			Title:   "foo",
			Content: "bar",
			UserId:  fmt.Sprintf("user:%s", id),
		})
	}
	return &post.GetPostsResponse{Posts: posts}, nil
}

type UserServer struct {
	*user.UnimplementedUserServiceServer
}

func (s *UserServer) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	profile, err := anypb.New(&user.User{
		Name: "foo",
	})
	if err != nil {
		return nil, err
	}
	return &user.GetUserResponse{
		User: &user.User{
			Id:   req.Id,
			Name: fmt.Sprintf("name_%s", req.Id),
			Items: []*user.Item{
				{
					Name:  "item1",
					Type:  user.Item_ITEM_TYPE_1,
					Value: 1,
					Location: &user.Item_Location{
						Addr1: "foo",
						Addr2: "bar",
						Addr3: &user.Item_Location_B{
							B: &user.Item_Location_AddrB{Bar: 1},
						},
					},
				},
				{Name: "item2", Type: user.Item_ITEM_TYPE_2, Value: 2},
			},
			Profile: map[string]*anypb.Any{"user": profile},
			Attr: &user.User_B{
				B: &user.User_AttrB{
					Bar: true,
				},
			},
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
		Client: new(clientConfig),
		Logger: logger,
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
	profile, err := anypb.New(&user.User{
		Name: "foo",
	})
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(res, &federation.GetPostResponse{
		Post: &federation.Post{
			Id:      "foo",
			Title:   "foo",
			Content: "bar",
			User: &federation.User{
				Id:   "user:foo",
				Name: "name_user:foo",
				Items: []*federation.Item{
					{
						Name:  "item1",
						Type:  federation.Item_ITEM_TYPE_1,
						Value: 1,
						Location: &federation.Item_Location{
							Addr1: "foo",
							Addr2: "bar",
							Addr3: &federation.Item_Location_B{
								B: &federation.Item_Location_AddrB{
									Bar: 1,
								},
							},
						},
					},
					{
						Name:  "item2",
						Type:  federation.Item_ITEM_TYPE_2,
						Value: 2,
					},
				},
				Profile: map[string]*anypb.Any{
					"user": profile,
				},
				Attr: &federation.User_B{
					B: &federation.User_AttrB{
						Bar: true,
					},
				},
			},
		},
		Str: "hello",
	}, cmpopts.IgnoreUnexported(
		federation.GetPostResponse{},
		federation.Post{},
		federation.User{},
		federation.Item{},
		federation.Item_Location{},
		federation.User_B{},
		federation.User_AttrB{},
		federation.Item_Location_B{},
		federation.Item_Location_AddrB{},
		anypb.Any{},
	)); diff != "" {
		t.Errorf("(-got, +want)\n%s", diff)
	}
}
