// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

type Exercise struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	DefaultWeight float32 `json:"default_weight"`
}
