package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &UtilityFunctionsProvider{}

type UtilityFunctionsProvider struct{}

type UtilityFunctionsProviderModel struct {
	Owner types.String `tfsdk:"owner"`
	Token types.String `tfsdk:"token"`
}

func NewUtilityFunctionsProvider() func() provider.Provider {
	return func() provider.Provider {
		return &UtilityFunctionsProvider{}
	}
}

func (p *UtilityFunctionsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "github"
}

func (p *UtilityFunctionsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
}

func (p *UtilityFunctionsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var model UtilityFunctionsProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (p *UtilityFunctionsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		nil,
	}
}

func (p *UtilityFunctionsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		nil,
	}
}

func (p *UtilityFunctionsProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		nil,
	}
}
