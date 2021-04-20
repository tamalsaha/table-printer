package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ColumnConverter interface {
	GVK() schema.GroupVersionKind
	Convert(obj runtime.Object) (map[string]interface{}, error)
}

var printers = map[schema.GroupVersionKind]ColumnConverter{}

func Register(c ColumnConverter) {
	printers[c.GVK()] = c
}

func Convert(obj runtime.Object) (map[string]interface{}, error) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	c, ok := printers[gvk]
	if !ok {
		return nil, fmt.Errorf("no column converter registered for %+v", gvk)
	}
	return c.Convert(obj)
}
