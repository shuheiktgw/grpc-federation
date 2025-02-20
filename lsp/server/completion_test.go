package server_test

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/k0kubun/pp/v3"

	"github.com/mercari/grpc-federation/compiler"
	"github.com/mercari/grpc-federation/lsp/server"
	"github.com/mercari/grpc-federation/source"
)

func TestCompletion(t *testing.T) {
	path := filepath.Join("testdata", "service.proto")
	file, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	completer := server.NewCompleter(compiler.New(), log.New(os.Stdout, "", 0), pp.New())
	t.Run("method", func(t *testing.T) {
		// resolver.method value position of Post in service.proto file
		_, candidates, err := completer.Completion(ctx, nil, path, file, source.Position{
			Line: 37,
			Col:  16,
		})
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(candidates, []string{
			"post.PostService/GetPost",
			"post.PostService/GetPosts",
			"user.UserService/GetUser",
			"user.UserService/GetUsers",
		}); diff != "" {
			t.Errorf("(-got, +want)\n%s", diff)
		}
	})

	t.Run("request.field", func(t *testing.T) {
		// resolver.request.field value position of Post in service.proto file
		_, candidates, err := completer.Completion(ctx, nil, path, file, source.Position{
			Line: 39,
			Col:  20,
		})
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(candidates, []string{
			"id",
		}); diff != "" {
			t.Errorf("(-got, +want)\n%s", diff)
		}
	})

	t.Run("request.by", func(t *testing.T) {
		// resolver.request.by value position os Post in service.proto file
		_, candidates, err := completer.Completion(ctx, nil, path, file, source.Position{
			Line: 39,
			Col:  29,
		})
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(candidates, []string{
			"$.id",
		}); diff != "" {
			t.Errorf("(-got, +want)\n%s", diff)
		}
	})

	t.Run("response.field", func(t *testing.T) {
		// resolver.response.field value position of Post in service.proto file
		_, candidates, err := completer.Completion(ctx, nil, path, file, source.Position{
			Line: 41,
			Col:  43,
		})
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(candidates, []string{
			"post",
		}); diff != "" {
			t.Errorf("(-got, +want)\n%s", diff)
		}
	})

	t.Run("message", func(t *testing.T) {
		// messages[0].message value position of Post in service.proto file
		_, candidates, err := completer.Completion(ctx, nil, path, file, source.Position{
			Line: 44,
			Col:  34,
		})
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(candidates, []string{
			// federation package messages.
			"GetPostRequest",
			"GetPostResponse",

			// google messages
			"google.protobuf.DescriptorProto",
			"google.protobuf.DescriptorProto.ExtensionRange",
			"google.protobuf.DescriptorProto.ReservedRange",
			"google.protobuf.EnumDescriptorProto",
			"google.protobuf.EnumDescriptorProto.EnumReservedRange",
			"google.protobuf.EnumOptions",
			"google.protobuf.EnumValueDescriptorProto",
			"google.protobuf.EnumValueOptions",
			"google.protobuf.ExtensionRangeOptions",
			"google.protobuf.ExtensionRangeOptions.Declaration",
			"google.protobuf.FieldDescriptorProto",
			"google.protobuf.FieldOptions",
			"google.protobuf.FileDescriptorProto",
			"google.protobuf.FileDescriptorSet",
			"google.protobuf.FileOptions",
			"google.protobuf.GeneratedCodeInfo",
			"google.protobuf.GeneratedCodeInfo.Annotation",
			"google.protobuf.MessageOptions",
			"google.protobuf.MethodDescriptorProto",
			"google.protobuf.MethodOptions",
			"google.protobuf.OneofDescriptorProto",
			"google.protobuf.OneofOptions",
			"google.protobuf.ServiceDescriptorProto",
			"google.protobuf.ServiceOptions",
			"google.protobuf.SourceCodeInfo",
			"google.protobuf.SourceCodeInfo.Location",
			"google.protobuf.UninterpretedOption",
			"google.protobuf.UninterpretedOption.NamePart",

			// post package messages.
			"post.GetPostRequest",
			"post.GetPostResponse",
			"post.GetPostsRequest",
			"post.GetPostsResponse",
			"post.Post",

			// user package messages.
			"user.GetUserRequest",
			"user.GetUserResponse",
			"user.GetUsersRequest",
			"user.GetUsersResponse",
			"user.User",
		}); diff != "" {
			t.Errorf("(-got, +want)\n%s", diff)
		}
	})

}
