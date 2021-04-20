package main

import (
	"fmt"
	"reflect"

	"gomodules.xyz/pointer"
	apps "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	Register(ReplicaSetPrinter{})
}

// ref: https://github.com/kubernetes/kubernetes/blob/v1.21.0/pkg/printers/internalversion/printers.go#L135-L146

type ReplicaSetPrinter struct{}

var _ ColumnConverter = ReplicaSetPrinter{}

func (_ ReplicaSetPrinter) GVK() schema.GroupVersionKind {
	return apps.SchemeGroupVersion.WithKind("ReplicaSet")
}

func (p ReplicaSetPrinter) Convert(o runtime.Object) (map[string]interface{}, error) {
	obj, ok := o.(*apps.ReplicaSet)
	if !ok {
		return nil, fmt.Errorf("expected %v, received %v", p.GVK().Kind, reflect.TypeOf(o))
	}

	row := map[string]interface{}{}

	desiredReplicas := obj.Spec.Replicas
	currentReplicas := obj.Status.Replicas
	readyReplicas := obj.Status.ReadyReplicas

	row["Name"] = obj.Name
	row["Desired"] = int64(pointer.Int32(desiredReplicas))
	row["Current"] = int64(currentReplicas)
	row["Ready"] = int64(readyReplicas)
	row["Age"] = translateTimestampSince(obj.CreationTimestamp)

	names, images := layoutContainerCells(obj.Spec.Template.Spec.Containers)
	row["Containers"] = names
	row["Images"] = images
	row["Selector"] = metav1.FormatLabelSelector(obj.Spec.Selector)

	return row, nil
}
