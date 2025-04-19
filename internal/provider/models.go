// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

type Exercise struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	DefaultWeight float32 `json:"default_weight"`
}

type Strategy struct {
	ID                    int     `json:"id"`
	DisplayName           string  `json:"display_name"`
	OverloadRate          float32 `json:"overload_rate"`
	ExercisesPerWorkout   int32   `json:"exercises_per_workout"`
	TargetRepsPerSet      int32   `json:"target_reps_per_set"`
	TargetSetsPerExercise int32   `json:"target_sets_per_exercise"`
}

type CreateStrategyPayload struct {
	DisplayName           string  `json:"display_name"`
	OverloadRate          float32 `json:"overload_rate"`
	ExercisesPerWorkout   int32   `json:"exercises_per_workout"`
	TargetRepsPerSet      int32   `json:"target_reps_per_set"`
	TargetSetsPerExercise int32   `json:"target_sets_per_exercise"`
}
