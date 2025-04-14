// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &exerciseResource{}
	_ resource.ResourceWithConfigure   = &exerciseResource{}
	_ resource.ResourceWithImportState = &exerciseResource{}
)

// NewExerciseResource is a helper function to simplify the provider implementation.
func NewExerciseResource() resource.Resource {
	return &exerciseResource{}
}

// exerciseResource is the resource implementation.
type exerciseResource struct {
	client *BrickByBrickClient
}

type exerciseResourceModel struct {
	ID            types.String  `tfsdk:"id"`
	Name          types.String  `tfsdk:"name"`
	DefaultWeight types.Float32 `tfsdk:"default_weight"`
}

// Metadata returns the resource type name.
func (r *exerciseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_exercise"
}

// Create a new resource.
func (r *exerciseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan exerciseResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan

	newExercise := Exercise{
		Name:          plan.Name.ValueString(),
		DefaultWeight: *plan.DefaultWeight.ValueFloat32Pointer(),
	}

	// Create new order
	createdExercise, err := r.client.CreateExercise(newExercise)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating exercise",
			"Could not create exercise, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(strconv.Itoa(createdExercise.ID))
	plan.DefaultWeight = types.Float32Value(createdExercise.DefaultWeight)
	plan.Name = types.StringValue(createdExercise.Name)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
// Read resource information.
func (r *exerciseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state exerciseResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from BrickByBrick
	refreshedExercise, err := r.client.GetExercise(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading BrickByBrick Exercise",
			"Could not read BrickByBrick exercise ID "+state.ID.String()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.Name = types.StringValue(refreshedExercise.Name)
	state.DefaultWeight = types.Float32Value(refreshedExercise.DefaultWeight)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *exerciseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan exerciseResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	updatedExercise := Exercise{
		Name:          plan.Name.ValueString(),
		DefaultWeight: plan.DefaultWeight.ValueFloat32(),
	}

	// Update existing order
	_, err := r.client.UpdateExercise(plan.ID.ValueString(), updatedExercise)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating exercise",
			"Could not update exercise, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated items from GetOrder as UpdateOrder items are not
	// populated.
	exercise, err := r.client.GetExercise(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading exercise",
			"Could not read exercise ID "+plan.ID.String()+": "+err.Error(),
		)
		return
	}

	plan.Name = types.StringValue(exercise.Name)
	plan.DefaultWeight = types.Float32Value(exercise.DefaultWeight)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *exerciseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state exerciseResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing exercise
	err := r.client.DeleteExercise(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting BrickByBrick Exercise",
			"Could not delete exercise, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *exerciseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *exerciseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for the exercise",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the exercise",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"default_weight": schema.Float32Attribute{
				Description: "The starting weight for the first session of this exercise. Measured in lbs. Defaults to 5.",
				Optional:    true,
				Computed:    true,
				Default:     float32default.StaticFloat32(5),
			},
		},
	}
}

func (r *exerciseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
