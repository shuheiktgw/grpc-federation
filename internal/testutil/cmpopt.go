package testutil

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/mercari/grpc-federation/resolver"
)

func ResolverCmpOpts() []cmp.Option {
	return []cmp.Option{
		cmpopts.IgnoreFields(resolver.File{}, "Messages", "Services", "Enums", "Desc"),
		cmpopts.IgnoreFields(resolver.Package{}, "Files"),
		cmpopts.IgnoreFields(resolver.Method{}, "Service"),
		cmpopts.IgnoreFields(resolver.Message{}, "File", "ParentMessage"),
		cmpopts.IgnoreFields(resolver.Enum{}, "File", "Message.Rule"),
		cmpopts.IgnoreFields(resolver.EnumValue{}, "Enum"),
		cmpopts.IgnoreFields(resolver.EnumRule{}, "Alias.Rule"),
		cmpopts.IgnoreFields(resolver.MessageResolver{}, "VariableDefinition"),
		cmpopts.IgnoreFields(resolver.MessageRule{}, "Alias.Rule"),
		cmpopts.IgnoreFields(resolver.MessageDependencyGraph{}, "MessageRule", "FieldOneofRule", "Roots"),
		cmpopts.IgnoreFields(resolver.MessageDependencyGraphNode{}, "BaseMessage", "VariableDefinition", "Parent", "ParentMap", "Children", "ChildrenMap", "Message.Rule"),
		cmpopts.IgnoreFields(resolver.AllMessageDependencyGraph{}),
		cmpopts.IgnoreFields(resolver.AllMessageDependencyGraphNode{}, "Parent", "Children", "Message.Rule"),
		cmpopts.IgnoreFields(resolver.ValidationErrorDetail{}, "DependencyGraph", "Resolvers"),
		cmpopts.IgnoreFields(resolver.AutoBindField{}, "VariableDefinition"),
		cmpopts.IgnoreFields(resolver.Type{}, "Ref.Rule", "Enum.Rule", "OneofField"),
		cmpopts.IgnoreFields(resolver.Oneof{}, "Message"),
		cmpopts.IgnoreFields(resolver.Field{}, "Oneof.Message", "Oneof.Fields"),
		cmpopts.IgnoreFields(resolver.Value{}, "CEL", "Const"),
		cmpopts.IgnoreFields(resolver.CELValue{}, "CheckedExpr"),
		cmpopts.IgnoreFields(resolver.MessageExpr{}, "Message.Rule", "Message.Fields"),
	}
}
