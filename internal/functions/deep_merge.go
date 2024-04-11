package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = DeepMerge{}

func NewDeepMerge() function.Function {
	return DeepMerge{}
}

type DeepMerge struct{}

func (r DeepMerge) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "deep_merge"
}

func (r DeepMerge) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Merges an arbitrary number of maps or objects into a single map or object.",
		Description:         "Deeply merges an arbitrary number of maps or objects into a single map or object.",
		MarkdownDescription: "Deeply merges an arbitrary number of maps or objects into a single map or object.",
		VariadicParameter: function.DynamicParameter{
			Name:                "maps",
			Description:         "An arbitrary number of maps or objects to merge.",
			MarkdownDescription: "An arbitrary number of maps or objects to merge.",
			AllowNullValue:      true,
			AllowUnknownValues:  true,
		},
		Return: function.DynamicReturn{},
	}
}

func (r DeepMerge) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var inputs types.Tuple

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &inputs))

	if resp.Error != nil {
		return
	}

	firstInput := inputs.Elements()[0] // basetypes.DynamicType

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, firstInput))
}
