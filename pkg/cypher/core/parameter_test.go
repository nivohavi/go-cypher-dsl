package core

import (
	"testing"
)

func TestNewParameters(t *testing.T) {
	params := NewParameters()
	if params == nil {
		t.Error("NewParameters() returned nil")
	}
}

func TestParametersAdd(t *testing.T) {
	params := NewParameters()
	paramExpr := params.Add("value")

	// Add() returns a parameter expression and adds it to params
	if paramExpr == nil {
		t.Error("Add() returned nil")
	}
	
	// Check that the parameter was added
	result := params.Get()
	if len(result) == 0 {
		t.Error("Add() did not add parameter")
	}
}

func TestParametersGet(t *testing.T) {
	params := NewParameters()
	params.Add("value1")
	params.Add("value2")

	result := params.Get()
	if len(result) != 2 {
		t.Errorf("Get() length = %d, want 2", len(result))
	}
}

func TestParametersMerge(t *testing.T) {
	params1 := NewParameters()
	params1.AddNamed("key1", "value1")

	params2 := NewParameters()
	params2.AddNamed("key2", "value2")

	params1.Merge(params2)
	result := params1.Get()

	// After merge, should have both parameters
	if len(result) != 2 {
		t.Errorf("Merge() length = %d, want 2", len(result))
	}
	
	// Verify both values are present
	if result["key1"] != "value1" {
		t.Errorf("Merge() key1 = %v, want 'value1'", result["key1"])
	}
	if result["key2"] != "value2" {
		t.Errorf("Merge() key2 = %v, want 'value2'", result["key2"])
	}
}

func TestParametersAddNamed(t *testing.T) {
	params := NewParameters()
	paramExpr := params.AddNamed("key", "value")

	if paramExpr == nil {
		t.Error("AddNamed() returned nil")
	}
	if paramExpr.Name() != "key" {
		t.Errorf("AddNamed() name = %q, want 'key'", paramExpr.Name())
	}
	
	result := params.Get()
	if result["key"] != "value" {
		t.Errorf("AddNamed() value = %v, want 'value'", result["key"])
	}
}

