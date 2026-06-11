package functions

import (
	"context"
	"fmt"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ function.Function = DeepMergeFunction{}

type DeepMergeFunction struct{}

func NewDeepMergeFunction() function.Function {
	return DeepMergeFunction{}
}

// mergeable reports whether a value can participate in a recursive merge,
// requiring it to be a known, non-null map or object.
func mergeable(v tftypes.Value) bool {
	if !v.IsKnown() || v.IsNull() {
		return false
	}

	switch v.Type().(type) {
	case tftypes.Object, tftypes.Map:
		return true
	default:
		return false
	}
}

// objectValue rebuilds a merged node as an object so that sibling values of
// different types, which a map cannot represent, remain valid.
func objectValue(m map[string]tftypes.Value) tftypes.Value {
	attributeTypes := make(map[string]tftypes.Type, len(m))

	for key, value := range m {
		attributeTypes[key] = value.Type()
	}

	return tftypes.NewValue(tftypes.Object{AttributeTypes: attributeTypes}, m)
}

// mergeValues merges later over earlier, recursing only when both sides are
// known, non-null maps or objects; otherwise later wins.
func mergeValues(earlier, later tftypes.Value) (tftypes.Value, error) {
	if !mergeable(earlier) || !mergeable(later) {
		return later, nil
	}

	earlierMap := make(map[string]tftypes.Value)

	if err := earlier.As(&earlierMap); err != nil {
		return tftypes.Value{}, err
	}

	laterMap := make(map[string]tftypes.Value)

	if err := later.As(&laterMap); err != nil {
		return tftypes.Value{}, err
	}

	merged := make(map[string]tftypes.Value, len(earlierMap)+len(laterMap))

	maps.Copy(merged, earlierMap)

	for key, value := range laterMap {
		existing, ok := merged[key]

		if !ok {
			merged[key] = value
			continue
		}

		mergedChild, err := mergeValues(existing, value)
		if err != nil {
			return tftypes.Value{}, err
		}

		merged[key] = mergedChild
	}

	return objectValue(merged), nil
}

func (r DeepMergeFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "deep_merge"
}

func (r DeepMergeFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Deeply merges an arbitrary number of maps or objects into a single object.",
		Description:         "Deeply merges an arbitrary number of maps or objects into a single object. When a key collides and both values are maps or objects, they are merged recursively; otherwise the later value replaces the earlier one, including values of a different type.",
		MarkdownDescription: "Deeply merges an arbitrary number of maps or objects into a single object. When a key collides and both values are maps or objects, they are merged recursively; otherwise the later value replaces the earlier one, including values of a different type.",
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

func (r DeepMergeFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var inputs []types.Dynamic

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &inputs))

	if resp.Error != nil {
		return
	}

	output := make(map[string]tftypes.Value)

	for i, input := range inputs {
		if input.IsNull() || input.IsUnderlyingValueNull() {
			continue
		}

		// An unknown argument makes the entire result unknown, matching the
		// behavior of Terraform's builtin merge function.
		if input.IsUnknown() || input.IsUnderlyingValueUnknown() {
			resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, types.DynamicUnknown()))
			return
		}

		terraformValue, err := input.ToTerraformValue(ctx)

		if err != nil {
			resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("Unable to Convert Argument: "+err.Error()))
			return
		}

		switch terraformValue.Type().(type) {
		case tftypes.Object, tftypes.Map:
		default:
			resp.Error = function.ConcatFuncErrors(resp.Error, function.NewArgumentFuncError(int64(i), fmt.Sprintf("Invalid Argument Type: arguments must be maps or objects, got %s", terraformValue.Type())))
			return
		}

		argumentMap := make(map[string]tftypes.Value)

		if err := terraformValue.As(&argumentMap); err != nil {
			resp.Error = function.ConcatFuncErrors(resp.Error, function.NewArgumentFuncError(int64(i), "Invalid Argument: "+err.Error()))
			return
		}

		for key, value := range argumentMap {
			existing, ok := output[key]

			if !ok {
				output[key] = value
				continue
			}

			merged, err := mergeValues(existing, value)
			if err != nil {
				resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("Unable to Merge Values: "+err.Error()))
				return
			}

			output[key] = merged
		}
	}

	result, err := basetypes.DynamicType{}.ValueFromTerraform(ctx, objectValue(output))

	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("Unable to Convert Result: "+err.Error()))
		return
	}

	dynamicResult, ok := result.(types.Dynamic)

	if !ok {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("Unable to Convert Result: unexpected result type."))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, dynamicResult))
}
