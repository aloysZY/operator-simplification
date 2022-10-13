package templateCm

import (
	zyv1 "operator-simplification/api/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 多个 CM 就要创建多个了

func NewConfigMap(aloys *zyv1.Aloys) []*corev1.ConfigMap {
	var label = map[string]string{"aloys": aloys.Name}
	var cm []*corev1.ConfigMap

	for _, v := range aloys.Spec.ConfigMap {
		c := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cm-" + aloys.Name + "-" + v.CmKey,
				Namespace: aloys.Namespace,
				Labels:    label,
			},
			Data: map[string]string{v.CmKey: v.CmDate},
		}
		cm = append(cm, c)
	}
	return cm
}
