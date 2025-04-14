// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &exercisesDataSource{}
	_ datasource.DataSourceWithConfigure = &exercisesDataSource{}
)

func NewExercisesDataSource() datasource.DataSource {
	return &exercisesDataSource{}
}

type exercisesDataSource struct {
	client *BrickByBrickClient
}

type exercisesModel struct {
	ID            types.Int64   `tfsdk:"id"`
	Name          types.String  `tfsdk:"name"`
	DefaultWeight types.Float32 `tfsdk:"default_weight"`
}

type exercisesDataSourceModel struct {
	Exercises []exercisesModel `tfsdk:"exercises"`
}

// Configure adds the provider configured client to the data source.
func (d *exercisesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*BrickByBrickClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *exercisesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_exercises"
}

// Schema defines the schema for the data source.
func (d *exercisesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"exercises": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A flat list of exercises associated with your account.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "The unique identifier of the exercise.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the exercise.",
							Computed:    true,
						},
						"default_weight": schema.Float32Attribute{
							Description: "The starting weight for this exercise. Measured in lbs.",
							Computed:    true,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *exercisesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state exercisesDataSourceModel

	exercises, err := d.client.GetExercises()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read BrickByBrick Exercises",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, exercise := range exercises {
		exerciseState := exercisesModel{
			ID:            types.Int64Value(int64(exercise.ID)),
			Name:          types.StringValue(exercise.Name),
			DefaultWeight: types.Float32Value(exercise.DefaultWeight),
		}

		state.Exercises = append(state.Exercises, exerciseState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
