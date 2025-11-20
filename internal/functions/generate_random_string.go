package functions

import (
	"context"
	"crypto/rand"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(n int64) (string, error) {
	b := make([]byte, n)
	// For each position, pick a random index in letters using crypto/rand
	for i := range b {
		num := make([]byte, 1)
		if _, err := rand.Read(num); err != nil {
			return "", err
		}
		b[i] = letters[int(num[0])%len(letters)]
	}
	return string(b), nil
}

var _ function.Function = GenerateRandomString{}

func NewGenerateRandomString() function.Function {
	return GenerateRandomString{}
}

type GenerateRandomString struct{}

func (r GenerateRandomString) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "generate_random_string"
}

func (r GenerateRandomString) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Return a random string, storing the output in state to use in subsequent runs.",
		Description:         "Return a random string.",
		MarkdownDescription: "Return a random string.",
		Parameters: []function.Parameter{
			function.Int64Parameter{
				Name:                "length",
				Description:         "The length of the random string to be generated.",
				MarkdownDescription: "The length of the random string to be generated.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (r GenerateRandomString) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var length int64

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &length))
	if resp.Error != nil {
		return
	}

	s, err := RandomString(length)
	if err != nil {
		panic(err)
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, s))
}
