package functions

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = GetEnvironmentVariable{}

func NewGetEnvironmentVariable() function.Function {
	return GetEnvironmentVariable{}
}

type GetEnvironmentVariable struct{}

func (r GetEnvironmentVariable) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "get_env"
}

func (r GetEnvironmentVariable) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Return the value of an environment variable given the variable name.",
		Description:         "Get the value of an environment variable.",
		MarkdownDescription: "Get the value of an environment variable.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "name",
				Description:         "The name of an environment variable.",
				MarkdownDescription: "The name of an environment variable.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (r GetEnvironmentVariable) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var data string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &data))

	if resp.Error != nil {
		return
	}

	value := os.Getenv(data)

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, value))
}
