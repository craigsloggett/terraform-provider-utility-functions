package provider

import (
	"context"

	"github.com/craigsloggett/terraform-provider-utility-functions/internal/functions"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &UtilityFunctionsProvider{}

type UtilityFunctionsProvider struct{}

func NewUtilityFunctionsProvider() func() provider.Provider {
	return func() provider.Provider {
		return &UtilityFunctionsProvider{}
	}
}

func (p *UtilityFunctionsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "utility_functions"
}

func (p *UtilityFunctionsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *UtilityFunctionsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *UtilityFunctionsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *UtilityFunctionsProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func (p *UtilityFunctionsProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		functions.NewGetEnvironmentVariable,
	}
}
