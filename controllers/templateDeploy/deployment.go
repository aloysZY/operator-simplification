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
	var aloysVolume []corev1.Volume

	for _, y := range aloys.Spec.Deployment.Containers {
		var aloysContainersPort []corev1.ContainerPort
		// 这不需要在进行 for 循环了，因为每次都是一个单独的
		// 设置端口
		p := corev1.ContainerPort{ContainerPort: y.Port}
		aloysContainersPort = append(aloysContainersPort, p)

		var aloysVolumeMount []corev1.VolumeMount
		// 如果没设置MountPath就不挂载
		for _, v := range y.MountPath {
			if v != "" {
				// 容器内挂载路径和参数
				fileName := strings.Split(v, "/")
				fileSubPath := fileName[len(fileName)-1]
				volumeMountName := strings.Split(fileSubPath, ".")[0]
				vm := &corev1.VolumeMount{
					// 如果一个 deploy 多个 contraction，挂在的文件名称相同就有问题，所以这里添加一个 contraction的名字
					Name:      "vm-" + y.Name + "-" + volumeMountName,
					ReadOnly:  true,
					MountPath: v,
					SubPath:   fileSubPath,
				}
				aloysVolumeMount = append(aloysVolumeMount, *vm)

				v := &corev1.Volume{
					Name: "vm-" + y.Name + "-" + volumeMountName,
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								// cm的名字
								Name: "cm-" + aloys.Name + "-" + volumeMountName,
							}}},
				}
				aloysVolume = append(aloysVolume, *v)
			}
		}

		d := &corev1.Container{
			Name:            y.Name,
			Image:           y.Image,
			Ports:           aloysContainersPort,
			ImagePullPolicy: corev1.PullIfNotPresent,
			Resources: corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceCPU:    y.Limits.Cpu,
					corev1.ResourceMemory: y.Limits.Memory,
				},
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    y.Request.Cpu,
					corev1.ResourceMemory: y.Request.Memory,
				},
			},
			VolumeMounts: aloysVolumeMount,
		}
		aloysContainers = append(aloysContainers, *d)
	}

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
