package functions_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/craigsloggett/terraform-provider-utility-functions/internal/functions"
)

func TestDeepMergeRun(t *testing.T) {
	t.Parallel()

	emptyObject := types.ObjectValueMust(map[string]attr.Type{}, map[string]attr.Value{})

	tests := map[string]struct {
		arguments []attr.Value
		want      types.Dynamic
		wantError bool
	}{
		"flat-merge-with-override": {
			arguments: []attr.Value{
				types.DynamicValue(types.ObjectValueMust(
					map[string]attr.Type{"a": types.StringType, "b": types.StringType},
					map[string]attr.Value{"a": types.StringValue("a1"), "b": types.StringValue("b1")},
				)),
				types.DynamicValue(types.ObjectValueMust(
					map[string]attr.Type{"b": types.StringType, "c": types.StringType},
					map[string]attr.Value{"b": types.StringValue("b2"), "c": types.StringValue("c2")},
				)),
			},
			want: types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{"a": types.StringType, "b": types.StringType, "c": types.StringType},
				map[string]attr.Value{"a": types.StringValue("a1"), "b": types.StringValue("b2"), "c": types.StringValue("c2")},
			)),
		},
		"nested-recursive-merge": {
			arguments: []attr.Value{
				types.DynamicValue(types.ObjectValueMust(
					map[string]attr.Type{
						"tags": types.ObjectType{AttrTypes: map[string]attr.Type{
							"environment": types.StringType,
							"team":        types.StringType,
						}},
					},
					map[string]attr.Value{
						"tags": types.ObjectValueMust(
							map[string]attr.Type{
								"environment": types.StringType,
								"team":        types.StringType,
							},
							map[string]attr.Value{
								"environment": types.StringValue("development"),
								"team":        types.StringValue("platform"),
							},
						),
					},
				)),
				types.DynamicValue(types.ObjectValueMust(
					map[string]attr.Type{
						"tags": types.ObjectType{AttrTypes: map[string]attr.Type{
							"environment": types.StringType,
						}},
					},
					map[string]attr.Value{
						"tags": types.ObjectValueMust(
							map[string]attr.Type{"environment": types.StringType},
							map[string]attr.Value{"environment": types.StringValue("production")},
						),
					},
				)),
			},
			want: types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{
					"tags": types.ObjectType{AttrTypes: map[string]attr.Type{
						"environment": types.StringType,
						"team":        types.StringType,
					}},
				},
				map[string]attr.Value{
					"tags": types.ObjectValueMust(
						map[string]attr.Type{
							"environment": types.StringType,
							"team":        types.StringType,
						},
						map[string]attr.Value{
							"environment": types.StringValue("production"),
							"team":        types.StringValue("platform"),
						},
					),
				},
			)),
		},
		"type-change-string-to-object": {
			arguments: []attr.Value{
				types.DynamicValue(types.ObjectValueMust(
					map[string]attr.Type{"a": types.StringType},
					map[string]attr.Value{"a": types.StringValue("a1")},
				)),
				types.DynamicValue(types.ObjectValueMust(
					map[string]attr.Type{"a": types.ObjectType{AttrTypes: map[string]attr.Type{"b": types.StringType}}},
					map[string]attr.Value{
						"a": types.ObjectValueMust(
							map[string]attr.Type{"b": types.StringType},
							map[string]attr.Value{"b": types.StringValue("b2")},
						),
					},
				)),
			},
			want: types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{"a": types.ObjectType{AttrTypes: map[string]attr.Type{"b": types.StringType}}},
				map[string]attr.Value{
					"a": types.ObjectValueMust(
						map[string]attr.Type{"b": types.StringType},
						map[string]attr.Value{"b": types.StringValue("b2")},
					),
				},
			)),
		},
		"type-change-object-to-string": {
			arguments: []attr.Value{
				types.DynamicValue(types.ObjectValueMust(
					map[string]attr.Type{"a": types.ObjectType{AttrTypes: map[string]attr.Type{"b": types.StringType}}},
					map[string]attr.Value{
						"a": types.ObjectValueMust(
							map[string]attr.Type{"b": types.StringType},
							map[string]attr.Value{"b": types.StringValue("b1")},
						),
					},
				)),
				types.DynamicValue(types.ObjectValueMust(
					map[string]attr.Type{"a": types.StringType},
					map[string]attr.Value{"a": types.StringValue("a2")},
				)),
			},
			want: types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{"a": types.StringType},
				map[string]attr.Value{"a": types.StringValue("a2")},
			)),
		},
		"map-and-object-merge": {
			arguments: []attr.Value{
				types.DynamicValue(types.MapValueMust(
					types.StringType,
					map[string]attr.Value{"a": types.StringValue("a1"), "b": types.StringValue("b1")},
				)),
				types.DynamicValue(types.ObjectValueMust(
					map[string]attr.Type{"b": types.Int64Type},
					map[string]attr.Value{"b": types.Int64Value(2)},
				)),
			},
			want: types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{"a": types.StringType, "b": types.NumberType},
				map[string]attr.Value{"a": types.StringValue("a1"), "b": types.NumberValue(big.NewFloat(2))},
			)),
		},
		"null-argument-is-skipped": {
			arguments: []attr.Value{
				types.DynamicValue(types.ObjectValueMust(
					map[string]attr.Type{"a": types.StringType},
					map[string]attr.Value{"a": types.StringValue("a1")},
				)),
				types.DynamicNull(),
			},
			want: types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{"a": types.StringType},
				map[string]attr.Value{"a": types.StringValue("a1")},
			)),
		},
		"unknown-argument-makes-result-unknown": {
			arguments: []attr.Value{
				types.DynamicValue(types.ObjectValueMust(
					map[string]attr.Type{"a": types.StringType},
					map[string]attr.Value{"a": types.StringValue("a1")},
				)),
				types.DynamicUnknown(),
			},
			want: types.DynamicUnknown(),
		},
		"non-map-argument-is-an-error": {
			arguments: []attr.Value{
				types.DynamicValue(types.StringValue("not a map")),
			},
			wantError: true,
		},
		"no-arguments-yield-an-empty-object": {
			arguments: []attr.Value{},
			want:      types.DynamicValue(emptyObject),
		},
		"all-null-arguments-yield-an-empty-object": {
			arguments: []attr.Value{
				types.DynamicNull(),
				types.DynamicNull(),
			},
			want: types.DynamicValue(emptyObject),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			// The framework packs variadic arguments into a single tuple.
			elementTypes := make([]attr.Type, len(test.arguments))
			for i := range test.arguments {
				elementTypes[i] = types.DynamicType
			}

			req := function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.TupleValueMust(elementTypes, test.arguments),
				}),
			}
			resp := function.RunResponse{
				Result: function.NewResultData(types.DynamicNull()),
			}

			functions.NewDeepMerge().Run(ctx, req, &resp)

			if test.wantError {
				if resp.Error == nil {
					t.Fatal("expected an error, got none")
				}
				return
			}

			if resp.Error != nil {
				t.Fatalf("unexpected error: %s", resp.Error)
			}

			got := resp.Result.Value()

			if !got.Equal(test.want) {
				t.Errorf("got %s, want %s", got, test.want)
			}
		})
	}
}
