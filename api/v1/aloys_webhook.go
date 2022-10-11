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

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Aloys) Default() {
	aloyslog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
	// 用来做修改的控制器，修改参数，设置默认参数等等

	if r.Spec.Ingress.Path == "" && r.Spec.Ingress.Enable {
		r.Spec.Ingress.Path = "/"
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
