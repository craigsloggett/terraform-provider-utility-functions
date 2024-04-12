package functions

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
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
	var inputs []types.Dynamic

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &inputs))

	if resp.Error != nil {
		return
	}

	output := make(map[string]interface{})

	for _, input := range inputs {

		if input.IsNull() {
			continue
		}

		terraformValue, err := input.ToTerraformValue(ctx)

		if err != nil {
			resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("ERROR: %s", err)))
			return
		}

		if terraformValue.IsKnown() {
			intermediateMap := make(map[string]tftypes.Value)

			err = terraformValue.As(&intermediateMap)

			if err != nil {
				resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("ERROR: %s", err)))
				return
			}
			resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("DEBUG: %T", intermediateMap)))
			resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("DEBUG: %s", intermediateMap)))
		} else {
			resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("ERROR: terraformValue is not known.")))
			return
		}
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(fmt.Sprintf("ERROR: %T", output)))
	return

	//resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, types.DynamicValue(mapValue)))
}
