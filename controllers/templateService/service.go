package templateService

import (
	zyv1 "operator-simplification/api/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewService(aloys *zyv1.Aloys) *corev1.Service {
	var label = map[string]string{"aloys": aloys.Name}
	var aloysPort []corev1.ServicePort
	for _, v := range aloys.Spec.Deployment.Containers {
		if aloys.Spec.Service.Name == v.Name {
			s := corev1.ServicePort{
				// Name是固定的，之后选择那个 contration需要暴露都ing 都不需要修改
				Name:       "http",
				Protocol:   "TCP",
				Port:       v.Port,
				TargetPort: intstr.IntOrString{IntVal: v.Port},
			}
			aloysPort = append(aloysPort, s)
		}
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "svc-" + aloys.Name,
			Namespace: aloys.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: label,
			Ports:    aloysPort,
		},
	}

	return service
}
