/*
Copyright 2022.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"sort"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var aloyslog = logf.Log.WithName("aloys-resource")

func (r *Aloys) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-zy-tech-v1-aloys,mutating=true,failurePolicy=fail,sideEffects=None,groups=zy.tech,resources=aloys,verbs=create;update,versions=v1,name=maloys.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Aloys{}

// 经过测试先验证Validate在进行Default

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Aloys) Default() {
	aloyslog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
	// 用来做修改的控制器，修改参数，设置默认参数等等

	if r.Spec.Ingress.Path == "" && r.Spec.Ingress.Enable {
		r.Spec.Ingress.Path = "/"
	}
	if strconv.Itoa(int(r.Spec.Deployment.Replicas)) == "" {
		r.Spec.Deployment.Replicas = 1
	}

	for _, v := range r.Spec.Deployment.Containers {
		if v.Limits.Cpu.String() == "0" {
			v.Limits.Cpu = v.Request.Cpu
		}

		if v.Limits.Memory.String() == "0" {
			v.Limits.Memory = v.Request.Memory
		}
		// 这里判断传入的是不是Limits就行了,但是有一个误区,我创建 deploy 的模板，就是contraction创建的时候Request没设置就是 0了,而不是没设置
		// 这里针对这种情况显示的设置一下
		// 添加这个设置 QOS 是Guaranteed，否则是Burstable
		// 其实可以添加 namespace或者 tag 来额外判断是否这样设置
		if v.Request.Cpu.String() == "0" {
			v.Request.Cpu = v.Limits.Cpu
		}
		if v.Request.Memory.String() == "0" {
			v.Request.Memory = v.Limits.Memory
		}
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:path=/validate-zy-tech-v1-aloys,mutating=false,failurePolicy=fail,sideEffects=None,groups=zy.tech,resources=aloys,verbs=create;update,versions=v1,name=valoys.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Aloys{}

// 用来做校验的控制器，比如效验设置是否正确或者合理，通常都是写函数后这里面对函数进行调用

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Aloys) ValidateCreate() error {
	aloyslog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	// RC创建的时候触发，其实就是执行创建 RC 的时候触发
	if err := r.ValidateService(); err != nil {
		return err
	}
	if err := r.ValidateIngress(); err != nil {
		return err
	}
	if err := r.ValidateConfigMap(); err != nil {
		return err
	}
	if err := r.ValidateResource(); err != nil {
		return err
	}
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Aloys) ValidateUpdate(old runtime.Object) error {
	aloyslog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	// RC更新的时候触发，其实就是执行更新 RC 的时候触发
	if err := r.ValidateService(); err != nil {
		return err
	}
	if err := r.ValidateIngress(); err != nil {
		return err
	}
	if err := r.ValidateConfigMap(); err != nil {
		return err
	}
	if err := r.ValidateResource(); err != nil {
		return err
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Aloys) ValidateDelete() error {
	aloyslog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	// RC更新的时候触发，其实就是执行更新 RC 的时候触发
	if err := r.ValidateService(); err != nil {
		return err
	}
	if err := r.ValidateIngress(); err != nil {
		return err
	}
	if err := r.ValidateConfigMap(); err != nil {
		return err
	}
	if err := r.ValidateResource(); err != nil {
		return err
	}
	return nil
}

// 这是判断如果 service 没有开启 ingress 开启也报错

func (r *Aloys) ValidateService() error {
	if !r.Spec.Service.Enable && r.Spec.Ingress.Enable {
		return errors.NewInvalid(GroupVersion.WithKind("Aloys").GroupKind(), r.Name,
			field.ErrorList{
				field.Invalid(field.NewPath("enable_service"),
					r.Spec.Service.Enable,
					"enable_service should be true when enable_ingress is true"),
			},
		)
	}
	return nil
}

func (r *Aloys) ValidateIngress() error {
	if r.Spec.Ingress.Enable && r.Spec.Ingress.Host == "" {
		return errors.NewInvalid(GroupVersion.WithKind("Aloys").GroupKind(), r.Name,
			field.ErrorList{
				field.Invalid(field.NewPath("host"),
					r.Spec.Ingress.Host,
					"host should be set when enable_ingress is true"),
			},
		)
	}
	return nil
}

func (r *Aloys) ValidateConfigMap() error {
	// 	这个修改后，应该反着找，格努Containers的 mountpath找 cm，如果 cmkey不存在，就报错
	// 	应该试试取出cm key 的值放在一个列表中，然后和每个 contraction中的 mountpath 的值进行匹配，当匹配不到的时候，就是 cm 没设置
	var cmKey []string
	for _, v := range r.Spec.ConfigMap {
		cmKey = append(cmKey, v.CmKey)
	}
	// 	这样就是判断一个值是不是在列表中了
	for _, k := range r.Spec.Deployment.Containers {
		for _, m := range k.MountPath {
			fileName := strings.Split(m, "/")
			fileSubPath := fileName[len(fileName)-1]
			volumeMountName := strings.Split(fileSubPath, ".")[0]
			// 在 Golang 中，有一个排序模块sort，它里面有一个sort.Strings()函数，可以对字符串数组进行排序。同时，还有一个sort.SearchStrings()[1]函数，会用二分法在一个有序字符串数组中寻找特定字符串的索引
			// 排序，加快查找速度
			sort.Strings(cmKey)
			index := sort.SearchStrings(cmKey, volumeMountName)
			if index < len(cmKey) && cmKey[index] == volumeMountName {
				continue
			} else {
				return errors.NewInvalid(GroupVersion.WithKind("Aloys").GroupKind(), r.Name, field.ErrorList{
					field.Invalid(field.NewPath("MountPath"),
						r.Spec.Ingress.Host,
						"MountPath should be set when configMap key set "),
				},
				)
			}
		}
	}
	return nil
}

func (r *Aloys) ValidateResource() error {
	for _, v := range r.Spec.Deployment.Containers {
		if v.Limits.Cpu.String() == "0" && v.Request.Cpu.String() == "0" {
			return errors.NewInvalid(GroupVersion.WithKind("Aloys").GroupKind(), r.Name, field.ErrorList{field.Invalid(field.NewPath("CPU"), v.Limits.Cpu, "Limits.Cpu or Request.Cpu must be set to one")})
		}
		if v.Limits.Memory.String() == "0" && v.Request.Memory.String() == "0" {
			return errors.NewInvalid(GroupVersion.WithKind("Aloys").GroupKind(), r.Name, field.ErrorList{field.Invalid(field.NewPath("Memory"), v.Limits.Cpu, "Limits.Memory or Request.Memory must be set to one")})
		}
	}
	return nil
}
