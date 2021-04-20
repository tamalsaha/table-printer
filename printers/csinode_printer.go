package printers

import (
	"fmt"
	"reflect"

	storage "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	Register(CSINodePrinter{})
}

// ref: https://github.com/kubernetes/kubernetes/blob/v1.21.0/pkg/printers/internalversion/printers.go#L164-L174

type CSINodePrinter struct{}

var _ ColumnConverter = CSINodePrinter{}

func (_ CSINodePrinter) GVK() schema.GroupVersionKind {
	return storage.SchemeGroupVersion.WithKind("CSINode")
}

func (p CSINodePrinter) Convert(o runtime.Object) (map[string]interface{}, error) {
	obj, ok := o.(*storage.CSINode)
	if !ok {
		return nil, fmt.Errorf("expected %v, received %v", p.GVK().Kind, reflect.TypeOf(o))
	}

	row := map[string]interface{}{}

	row["Name"] = obj.Name
	row["Drivers"] = len(obj.Spec.Drivers)
	row["Age"] = translateTimestampSince(obj.CreationTimestamp)

	return row, nil
}
