package printers

import (
	"fmt"
	"reflect"

	apps "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	Register(DaemonSetPrinter{})
}

// ref: https://github.com/kubernetes/kubernetes/blob/v1.21.0/pkg/printers/internalversion/printers.go#L135-L146

type DaemonSetPrinter struct{}

var _ ColumnConverter = DaemonSetPrinter{}

func (_ DaemonSetPrinter) GVK() schema.GroupVersionKind {
	return apps.SchemeGroupVersion.WithKind("DaemonSet")
}

func (p DaemonSetPrinter) Convert(o runtime.Object) (map[string]interface{}, error) {
	obj, ok := o.(*apps.DaemonSet)
	if !ok {
		return nil, fmt.Errorf("expected %v, received %v", p.GVK().Kind, reflect.TypeOf(o))
	}

	row := map[string]interface{}{}

	desiredScheduled := obj.Status.DesiredNumberScheduled
	currentScheduled := obj.Status.CurrentNumberScheduled
	numberReady := obj.Status.NumberReady
	numberUpdated := obj.Status.UpdatedNumberScheduled
	numberAvailable := obj.Status.NumberAvailable

	row["Name"] = obj.Name
	row["Desired"] = int64(desiredScheduled)
	row["Current"] = int64(currentScheduled)
	row["Ready"] = int64(numberReady)
	row["Up-to-date"] = int64(numberUpdated)
	row["Available"] = int64(numberAvailable)
	row["Node Selector"] = labels.FormatLabels(obj.Spec.Template.Spec.NodeSelector)
	row["Age"] = translateTimestampSince(obj.CreationTimestamp)

	names, images := layoutContainerCells(obj.Spec.Template.Spec.Containers)
	row["Containers"] = names
	row["Images"] = images
	row["Selector"] = metav1.FormatLabelSelector(obj.Spec.Selector)

	return row, nil
}
