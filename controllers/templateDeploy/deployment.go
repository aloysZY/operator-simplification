package templateDeployment

import (
	"strings"

	zyv1 "operator-simplification/api/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewDeployment(aloys *zyv1.Aloys) *appsv1.Deployment {
	var label = map[string]string{"aloys": aloys.Name}
	var aloysContainers []corev1.Container

	var aloysContainersPort []corev1.ContainerPort
	p := corev1.ContainerPort{ContainerPort: aloys.Spec.Deployment.Port}
	aloysContainersPort = append(aloysContainersPort, p)

	var aloysVolumeMount []corev1.VolumeMount
	fileName := strings.Split(aloys.Spec.Deployment.MountPath, "/")
	vm := &corev1.VolumeMount{
		Name:      "vm-" + aloys.Name,
		ReadOnly:  true,
		MountPath: aloys.Spec.Deployment.MountPath,
		SubPath:   fileName[len(fileName)-1],
	}
	aloysVolumeMount = append(aloysVolumeMount, *vm)

	var aloysVolume []corev1.Volume
	v := &corev1.Volume{
		Name: "vm-" + aloys.Name,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: "cm-" + aloys.Name}}},
	}
	aloysVolume = append(aloysVolume, *v)
	cName := strings.Split(aloys.Spec.Deployment.Image, ":")[0]
	cName2 := strings.Split(cName, "/")
	cname3 := cName2[len(cName2)-1]

	d := &corev1.Container{

		Name:            cname3,
		Image:           aloys.Spec.Deployment.Image,
		Ports:           aloysContainersPort,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    aloys.Spec.Deployment.Limits.Cpu,
				corev1.ResourceMemory: aloys.Spec.Deployment.Limits.Memory,
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    aloys.Spec.Deployment.Request.Cpu,
				corev1.ResourceMemory: aloys.Spec.Deployment.Request.Memory,
			},
		},
		VolumeMounts: aloysVolumeMount,
	}
	aloysContainers = append(aloysContainers, *d)

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deploy-" + aloys.Name,
			Namespace: aloys.Namespace,
			Labels:    label,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &aloys.Spec.Deployment.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: label},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: label},
				Spec: corev1.PodSpec{
					Containers: aloysContainers,
					Volumes:    aloysVolume,
				},
			},
		},
	}
	return deploy
}
