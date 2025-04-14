// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &brickbybrickProvider{}
)

type brickbybrickProviderModel struct {
	ApiKey types.String `tfsdk:"api_key"`
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &brickbybrickProvider{
			version: version,
		}
	}
}

// brickbybrickProvider is the provider implementation.
type brickbybrickProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *brickbybrickProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "brickbybrick"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
// Schema defines the provider-level schema for configuration data.
func (p *brickbybrickProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with the BrickByBrick Fitness API via Terraform.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Required:    true,
				Description: "Your BrickByBrick Fitness API Key",
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a brickbybrick API client for data sources and resources.
func (p *brickbybrickProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring BrickByBrick client")
	// Retrieve provider data from configuration
	var config brickbybrickProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown BrickByBrick API Key",
			"The provider cannot create the BrickByBrick API client as there is an unknown configuration value for the BrickByBrick API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BRICKBYBRICK_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	apiKey := os.Getenv("BRICKBYBRICK_API_KEY")

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing BrickByBrick API Key",
			"The provider cannot create the BrickByBrick API client as there is a missing or empty value for the BrickByBrick API host. "+
				"Set the host value in the configuration or use the BRICKBYBRICK_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "brickbybrick_api_key", apiKey)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "brickbybrick_api_key")

	tflog.Debug(ctx, "Creating BrickByBrick client")

	// Create a new HashiCups client using the configuration values
	client, err := NewClient(&apiKey)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create BrickByBrick API Client",
			"An unexpected error occurred when creating the BrickByBrick API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"BrickByBrick Client Error: "+err.Error(),
		)
		return
	}

	// Make the HashiCups client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured BrickByBrick client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *brickbybrickProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExercisesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *brickbybrickProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExerciseResource,
	}
}
