package templateCm

import (
	zyv1 "operator-simplification/api/v1"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 多个 CM 就要创建多个了
func NewConfigMap(aloys *zyv1.Aloys) []*corev1.ConfigMap {
	var c []*corev1.ConfigMap

	for _, v := range aloys.Spec.ConfigMap {
		for _, y := range aloys.Spec.Deployment.Containers {
			// 去找对应的 contrations
			if v.Name == y.Name {
				dataKey := strings.Split(y.MountPath, "/")
				x := &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cm-" + aloys.Name + "-" + v.Name,
						Namespace: aloys.Namespace,
					},
					Data: map[string]string{dataKey[len(dataKey)-1]: v.CmDate},
				}
				c = append(c, x)
			}
		}
	}
	return c
}
