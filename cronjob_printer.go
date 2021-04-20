package main

import (
	"fmt"
	"reflect"

	batch "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	Register(CronJobPrinter{})
}

// ref: https://github.com/kubernetes/kubernetes/blob/v1.21.0/pkg/printers/internalversion/printers.go#L176-L188

type CronJobPrinter struct{}

var _ ColumnConverter = CronJobPrinter{}

func (_ CronJobPrinter) GVK() schema.GroupVersionKind {
	return batch.SchemeGroupVersion.WithKind("CronJob")
}

func (p CronJobPrinter) Convert(o runtime.Object) (map[string]interface{}, error) {
	obj, ok := o.(*batch.CronJob)
	if !ok {
		return nil, fmt.Errorf("expected %v, received %v", p.GVK().Kind, reflect.TypeOf(o))
	}

	row := map[string]interface{}{}

	lastScheduleTime := "<none>"
	if obj.Status.LastScheduleTime != nil {
		lastScheduleTime = translateTimestampSince(*obj.Status.LastScheduleTime)
	}

	row["Name"] = obj.Name
	row["Schedule"] = obj.Spec.Schedule
	row["Suspend"] = printBoolPtr(obj.Spec.Suspend)
	row["Active"] = int64(len(obj.Status.Active))
	row["Last Schedule"] = lastScheduleTime
	row["Age"] = translateTimestampSince(obj.CreationTimestamp)

	names, images := layoutContainerCells(obj.Spec.JobTemplate.Spec.Template.Spec.Containers)
	row["Containers"] = names
	row["Images"] = images
	row["Selector"] = metav1.FormatLabelSelector(obj.Spec.JobTemplate.Spec.Selector)

	return row, nil
}
