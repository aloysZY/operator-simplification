package templateIngress

import (
	zyv1 "operator-simplification/api/v1"

	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewIngress(aloys *zyv1.Aloys) *netv1.Ingress {
	var aloysIngress []netv1.IngressRule
	var aloysPath []netv1.HTTPIngressPath
	Prefix := netv1.PathTypePrefix
	p := &netv1.HTTPIngressPath{
		Path:     aloys.Spec.Ingress.Path,
		PathType: &Prefix,
		Backend: netv1.IngressBackend{
			Service: &netv1.IngressServiceBackend{
				Name: "svc-" + aloys.Name,
				Port: netv1.ServiceBackendPort{Name: "http"},
			},
		},
	}
	aloysPath = append(aloysPath, *p)

	i := netv1.IngressRule{
		Host:             aloys.Spec.Ingress.Host,
		IngressRuleValue: netv1.IngressRuleValue{HTTP: &netv1.HTTPIngressRuleValue{Paths: aloysPath}},
	}
	aloysIngress = append(aloysIngress, i)

	var aloysHost []string
	aloysHost = append(aloysHost, aloys.Spec.Ingress.Host)

	var aloysTls []netv1.IngressTLS
	t := netv1.IngressTLS{
		Hosts: aloysHost,
		// 这个名字我们可以随便起，现在的配置其实是根据这个名字生成证书到这个名字下
		SecretName: "secret-" + aloys.Name,
	}
	aloysTls = append(aloysTls, t)

	ing := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ingress-" + aloys.Name,
			Namespace: aloys.Namespace,
			// certmanager 注入的时候需要添加改行注解
			Annotations: map[string]string{"cert-manager.io/cluster-issuer": "letsencrypt-staging"}},
		Spec: netv1.IngressSpec{
			Rules: aloysIngress,
			TLS:   aloysTls,
		},
	}

	return ing
}
