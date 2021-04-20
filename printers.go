package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// GenerateOptions encapsulates attributes for table generation.
type GenerateOptions struct {
	NoHeaders bool
	Wide      bool
}

// TableGenerator - an interface for generating metav1.Table provided a runtime.Object
type TableGenerator interface {
	GenerateTable(obj runtime.Object, options GenerateOptions) (*metav1.Table, error)
}
