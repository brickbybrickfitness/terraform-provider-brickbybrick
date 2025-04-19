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
	_ datasource.DataSource              = &strategiesDataSource{}
	_ datasource.DataSourceWithConfigure = &strategiesDataSource{}
)

func NewStrategiesDataSource() datasource.DataSource {
	return &strategiesDataSource{}
}

type strategiesDataSource struct {
	client *BrickByBrickClient
}

type strategiesModel struct {
	ID                    types.Int64   `tfsdk:"id"`
	OverloadRate          types.Float32 `tfsdk:"overload_rate"`
	ExercisesPerWorkout   types.Int32   `tfsdk:"exercises_per_workout"`
	TargetSetsPerExercise types.Int32   `tfsdk:"target_sets_per_exercise"`
	TargetRepsPerSet      types.Int32   `tfsdk:"target_reps_per_set"`
	DisplayName           types.String  `tfsdk:"display_name"`
}

type strategiesDataSourceModel struct {
	Strategies []strategiesModel `tfsdk:"strategies"`
}

// Configure adds the provider configured client to the data source.
func (d *strategiesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *strategiesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_strategies"
}

// Schema defines the schema for the data source.
func (d *strategiesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"strategies": schema.ListNestedAttribute{
				Computed:    true,
				Description: "A list of your progressive overload strategies.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "The unique identifier of the strategy.",
							Computed:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "The name of the strategy.",
							Required:    true,
						},
						"overload_rate": schema.Float32Attribute{
							Description: "The amount of resistance or weight to add (in lbs) to each rep per session.",
							Required:    true,
						},
						"exercises_per_workout": schema.Int32Attribute{
							Description: "The number of exercises that each workout should have.",
							Required:    true,
						},
						"target_sets_per_exercise": schema.Int32Attribute{
							Description: "The goal for the number of sets you eventually want to do for each exercise in a workout.",
							Required:    true,
						},
						"target_reps_per_set": schema.Int32Attribute{
							Description: "The goal for the number of reps that you eventually want to do in a set.",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *strategiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state strategiesDataSourceModel

	strategies, err := d.client.GetStrategies()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read BrickByBrick Strategies",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, strategy := range strategies {
		strategyState := strategiesModel{
			ID:                    types.Int64Value(int64(strategy.ID)),
			DisplayName:           types.StringValue(strategy.DisplayName),
			OverloadRate:          types.Float32Value(strategy.OverloadRate),
			TargetSetsPerExercise: types.Int32Value(strategy.TargetSetsPerExercise),
			TargetRepsPerSet:      types.Int32Value(strategy.TargetRepsPerSet),
			ExercisesPerWorkout:   types.Int32Value(strategy.ExercisesPerWorkout),
		}

		state.Strategies = append(state.Strategies, strategyState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
