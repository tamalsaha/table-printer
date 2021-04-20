package printers

import (
	"fmt"
	"reflect"

	"gomodules.xyz/pointer"

	apps "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	Register(StatefulSetPrinter{})
}

// ref: https://github.com/kubernetes/kubernetes/blob/v1.21.0/pkg/printers/internalversion/printers.go#L223-L231

type StatefulSetPrinter struct{}

var _ ColumnConverter = StatefulSetPrinter{}

func (_ StatefulSetPrinter) GVK() schema.GroupVersionKind {
	return apps.SchemeGroupVersion.WithKind("StatefulSet")
}

func (p StatefulSetPrinter) Convert(o runtime.Object) (map[string]interface{}, error) {
	obj, ok := o.(*apps.StatefulSet)
	if !ok {
		return nil, fmt.Errorf("expected %v, received %v", p.GVK().Kind, reflect.TypeOf(o))
	}

	row := map[string]interface{}{}

	desiredReplicas := obj.Spec.Replicas
	readyReplicas := obj.Status.ReadyReplicas
	createTime := translateTimestampSince(obj.CreationTimestamp)

	row["Name"] = obj.Name
	row["Ready"] = fmt.Sprintf("%d/%d", int64(readyReplicas), int64(pointer.Int32(desiredReplicas)))
	row["Age"] = createTime

	names, images := layoutContainerCells(obj.Spec.Template.Spec.Containers)
	row["Containers"] = names
	row["Images"] = images

	return row, nil
}
