// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/float32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &strategyResource{}
	_ resource.ResourceWithConfigure   = &strategyResource{}
	_ resource.ResourceWithImportState = &strategyResource{}
)

// NewStrategyResource is a helper function to simplify the provider implementation.
func NewStrategyResource() resource.Resource {
	return &strategyResource{}
}

// strategyResource is the resource implementation.
type strategyResource struct {
	client *BrickByBrickClient
}

type strategyResourceModel struct {
	ID                    types.String  `tfsdk:"id"`
	DisplayName           types.String  `tfsdk:"display_name"`
	OverloadRate          types.Float32 `tfsdk:"overload_rate"`
	ExercisesPerWorkout   types.Int32   `tfsdk:"exercises_per_workout"`
	TargetSetsPerExercise types.Int32   `tfsdk:"target_sets_per_exercise"`
	TargetRepsPerSet      types.Int32   `tfsdk:"target_reps_per_set"`
}

// Metadata returns the resource type name.
func (r *strategyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_strategy"
}

// Create a new resource.
func (r *strategyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan strategyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan

	newStrategy := CreateStrategyPayload{
		DisplayName:           plan.DisplayName.ValueString(),
		OverloadRate:          *plan.OverloadRate.ValueFloat32Pointer(),
		TargetRepsPerSet:      *plan.TargetRepsPerSet.ValueInt32Pointer(),
		TargetSetsPerExercise: *plan.TargetSetsPerExercise.ValueInt32Pointer(),
		ExercisesPerWorkout:   *plan.ExercisesPerWorkout.ValueInt32Pointer(),
	}

	// Create new order
	createdStrategy, err := r.client.CreateStrategy(newStrategy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating strategy",
			"Could not create strategy, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(strconv.Itoa(createdStrategy.ID))
	plan.OverloadRate = types.Float32Value(createdStrategy.OverloadRate)
	plan.DisplayName = types.StringValue(createdStrategy.DisplayName)
	plan.ExercisesPerWorkout = types.Int32Value(createdStrategy.ExercisesPerWorkout)
	plan.TargetRepsPerSet = types.Int32Value(createdStrategy.TargetRepsPerSet)
	plan.TargetSetsPerExercise = types.Int32Value(createdStrategy.TargetSetsPerExercise)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
// Read resource information.
func (r *strategyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state strategyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from BrickByBrick
	refreshedStrategy, err := r.client.GetStrategy(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading BrickByBrick Strategy",
			"Could not read BrickByBrick strategy ID "+state.ID.String()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.DisplayName = types.StringValue(refreshedStrategy.DisplayName)
	state.OverloadRate = types.Float32Value(refreshedStrategy.OverloadRate)
	state.ExercisesPerWorkout = types.Int32Value(refreshedStrategy.ExercisesPerWorkout)
	state.TargetRepsPerSet = types.Int32Value(refreshedStrategy.TargetRepsPerSet)
	state.TargetRepsPerSet = types.Int32Value(refreshedStrategy.TargetRepsPerSet)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *strategyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan strategyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updatedStrategy := CreateStrategyPayload{
		DisplayName:           plan.DisplayName.ValueString(),
		OverloadRate:          plan.OverloadRate.ValueFloat32(),
		ExercisesPerWorkout:   plan.ExercisesPerWorkout.ValueInt32(),
		TargetRepsPerSet:      plan.TargetRepsPerSet.ValueInt32(),
		TargetSetsPerExercise: plan.TargetSetsPerExercise.ValueInt32(),
	}

	// Update existing order
	_, err := r.client.UpdateStrategy(plan.ID.ValueString(), updatedStrategy)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating strategy",
			"Could not update strategy, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated items from GetOrder as UpdateOrder items are not
	// populated.
	strategy, err := r.client.GetStrategy(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading strategy",
			"Could not read strategy ID "+plan.ID.String()+": "+err.Error(),
		)
		return
	}

	plan.DisplayName = types.StringValue(strategy.DisplayName)
	plan.OverloadRate = types.Float32Value(strategy.OverloadRate)
	plan.ExercisesPerWorkout = types.Int32Value(strategy.ExercisesPerWorkout)
	plan.TargetRepsPerSet = types.Int32Value(strategy.TargetRepsPerSet)
	plan.TargetSetsPerExercise = types.Int32Value(strategy.TargetSetsPerExercise)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *strategyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state strategyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing strategy
	err := r.client.DeleteStrategy(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting BrickByBrick Strategy",
			"Could not delete strategy, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *strategyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*BrickByBrickClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *BrickByBrickClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Schema defines the schema for the resource.
func (r *strategyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for the strategy",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the strategy",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(64),
				},
			},
			"overload_rate": schema.Float32Attribute{
				Description: "The amount of resistance or weight to add (in lbs) to each rep per session.",
				Required:    true,
				Validators: []validator.Float32{
					float32validator.AtLeast(0),
					float32validator.AtMost(100),
				},
			},
			"exercises_per_workout": schema.Int32Attribute{
				Description: "The number of exercises that each workout should have.",
				Required:    true,
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
					int32validator.AtMost(100),
				},
			},
			"target_sets_per_exercise": schema.Int32Attribute{
				Description: "The goal for the number of sets you eventually want to do for each exercise in a workout.",
				Required:    true,
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
					int32validator.AtMost(10000),
				},
			},
			"target_reps_per_set": schema.Int32Attribute{
				Description: "The goal for the number of reps that you eventually want to do in a set.",
				Required:    true,
				Validators: []validator.Int32{
					int32validator.AtLeast(1),
					int32validator.AtMost(100000),
				},
			},
		},
	}
}

func (r *strategyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
