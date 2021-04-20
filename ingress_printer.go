package main

import (
	"fmt"
	"reflect"
	"strings"

	networking "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	Register(IngressPrinter{})
}

// ref: https://github.com/kubernetes/kubernetes/blob/v1.21.0/pkg/printers/internalversion/printers.go#L203-L212

type IngressPrinter struct{}

var _ ColumnConverter = IngressPrinter{}

func (_ IngressPrinter) GVK() schema.GroupVersionKind {
	return networking.SchemeGroupVersion.WithKind("Ingress")
}

func (p IngressPrinter) Convert(o runtime.Object) (map[string]interface{}, error) {
	obj, ok := o.(*networking.Ingress)
	if !ok {
		return nil, fmt.Errorf("expected %v, received %v", p.GVK().Kind, reflect.TypeOf(o))
	}

	row := map[string]interface{}{}

	className := "<none>"
	if obj.Spec.IngressClassName != nil {
		className = *obj.Spec.IngressClassName
	}
	hosts := formatHosts(obj.Spec.Rules)
	address := loadBalancerStatusStringer(obj.Status.LoadBalancer)
	ports := formatPorts(obj.Spec.TLS)
	createTime := translateTimestampSince(obj.CreationTimestamp)

	row["Name"] = obj.Name
	row["Class"] = className
	row["Hosts"] = hosts
	row["Address"] = address
	row["Ports"] = ports
	row["Age"] = createTime

	return row, nil
}

func formatHosts(rules []networking.IngressRule) string {
	list := []string{}
	max := 3
	more := false
	for _, rule := range rules {
		if len(list) == max {
			more = true
		}
		if !more && len(rule.Host) != 0 {
			list = append(list, rule.Host)
		}
	}
	if len(list) == 0 {
		return "*"
	}
	ret := strings.Join(list, ",")
	if more {
		return fmt.Sprintf("%s + %d more...", ret, len(rules)-max)
	}
	return ret
}

func formatPorts(tls []networking.IngressTLS) string {
	if len(tls) != 0 {
		return "80, 443"
	}
	return "80"
}
