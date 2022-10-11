package templateCm

import (
	zyv1 "operator-simplification/api/v1"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewConfigMap(aloys *zyv1.Aloys) *corev1.ConfigMap {
	dataKey := strings.Split(aloys.Spec.Deployment.MountPath, "/")
	c := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cm-" + aloys.Name,
			Namespace: aloys.Namespace,
		},
		Data: map[string]string{dataKey[len(dataKey)-1]: aloys.Spec.ConfigMap.CmDate},
	}
	return c
}
