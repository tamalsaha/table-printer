package main

import (
	"fmt"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func init() {
	Register(ServicePrinter{})
}

// ref: https://github.com/kubernetes/kubernetes/blob/v1.21.0/pkg/printers/internalversion/printers.go#L190-L198

type ServicePrinter struct{}

var _ ColumnConverter = ServicePrinter{}

func (_ ServicePrinter) GVK() schema.GroupVersionKind {
	return core.SchemeGroupVersion.WithKind("Service")
}

func (p ServicePrinter) Convert(o runtime.Object) (map[string]interface{}, error) {
	obj, ok := o.(*core.Service)
	if !ok {
		return nil, fmt.Errorf("expected %v, received %v", p.GVK().Kind, reflect.TypeOf(o))
	}

	row := map[string]interface{}{}

	svcType := obj.Spec.Type
	internalIP := "<none>"
	if len(obj.Spec.ClusterIPs) > 0 {
		internalIP = obj.Spec.ClusterIPs[0]
	}

	externalIP := getServiceExternalIP(obj)
	svcPorts := makePortString(obj.Spec.Ports)
	if len(svcPorts) == 0 {
		svcPorts = "<none>"
	}

	row["Name"] = obj.Name
	row["Type"] = string(svcType)
	row["Cluster-IP"] = internalIP
	row["External-IP"] = externalIP
	row["Port(s)"] = svcPorts
	row["Age"] = translateTimestampSince(obj.CreationTimestamp)
	row["Selector"] = labels.FormatLabels(obj.Spec.Selector)

	return row, nil
}

func getServiceExternalIP(svc *core.Service) string {
	switch svc.Spec.Type {
	case core.ServiceTypeClusterIP:
		if len(svc.Spec.ExternalIPs) > 0 {
			return strings.Join(svc.Spec.ExternalIPs, ",")
		}
		return "<none>"
	case core.ServiceTypeNodePort:
		if len(svc.Spec.ExternalIPs) > 0 {
			return strings.Join(svc.Spec.ExternalIPs, ",")
		}
		return "<none>"
	case core.ServiceTypeLoadBalancer:
		lbIps := loadBalancerStatusStringer(svc.Status.LoadBalancer)
		if len(svc.Spec.ExternalIPs) > 0 {
			results := []string{}
			if len(lbIps) > 0 {
				results = append(results, strings.Split(lbIps, ",")...)
			}
			results = append(results, svc.Spec.ExternalIPs...)
			return strings.Join(results, ",")
		}
		if len(lbIps) > 0 {
			return lbIps
		}
		return "<pending>"
	case core.ServiceTypeExternalName:
		return svc.Spec.ExternalName
	}
	return "<unknown>"
}

// loadBalancerStatusStringer behaves mostly like a string interface and converts the given status to a string.
// `wide` indicates whether the returned value is meant for --o=wide output. If not, it's clipped to 16 bytes.
func loadBalancerStatusStringer(s core.LoadBalancerStatus) string {
	ingress := s.Ingress
	result := sets.NewString()
	for i := range ingress {
		if ingress[i].IP != "" {
			result.Insert(ingress[i].IP)
		} else if ingress[i].Hostname != "" {
			result.Insert(ingress[i].Hostname)
		}
	}

	return strings.Join(result.List(), ",")
}

func makePortString(ports []core.ServicePort) string {
	pieces := make([]string, len(ports))
	for ix := range ports {
		port := &ports[ix]
		pieces[ix] = fmt.Sprintf("%d/%s", port.Port, port.Protocol)
		if port.NodePort > 0 {
			pieces[ix] = fmt.Sprintf("%d:%d/%s", port.Port, port.NodePort, port.Protocol)
		}
	}
	return strings.Join(pieces, ",")
}
