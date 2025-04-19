# Copyright (c) HashiCorp, Inc.

resource "brickbybrick_strategy" "my_rapid_progress_strategy" {
  display_name             = "My Rapid Progress Strategy"
  overload_rate            = 5
  target_sets_per_exercise = 10
  target_reps_per_set      = 20
  exercises_per_workout    = 4
}
