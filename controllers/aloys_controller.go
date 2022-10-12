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

package controllers

import (
	"context"
	"operator-simplification/controllers/templateCm"
	templateDeployment "operator-simplification/controllers/templateDeploy"
	"operator-simplification/controllers/templateIngress"
	"operator-simplification/controllers/templateService"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	zyv1 "operator-simplification/api/v1"
)

// AloysReconciler reconciles a Aloys object
type AloysReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	// 	添加字段
	// X string
}

// 在这里添加 rbac 权限注解，额外添加了对k8s中 deoloy,svc,ing资源的注解，这样 rbac 对生成相关资源配置，这个 controller就可以操作相关资源了
// +kubebuilder:rbac:groups=zy.tech,resources=aloys,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=zy.tech,resources=aloys/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=zy.tech,resources=aloys/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Aloys object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *AloysReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	aloysLog := log.FromContext(ctx)
	// 具体的逻辑实现
	// TODO(user): your logic here

	// 测试使用
	// fmt.Println(r.X)

	var aloys zyv1.Aloys
	// 先去查找这个kind 资源是否存在，如果不存在返回不存在的错误就结束
	// apiVersion: zy.tech/v1
	// kind: Aloys
	if err := r.Get(ctx, req.NamespacedName, &aloys); err != nil {
		aloysLog.Error(err, "unable to fetch aloys")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 这里开始用模板 Unmarshal进行反向解析，现在直接创建相关资源后补全相关设置信息

	// cm
	// for 循环创建 cm
	cm := templateCm.NewConfigMap(&aloys)
	for _, v := range cm {
		if err := ctrl.SetControllerReference(&aloys, v, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		c := &corev1.ConfigMap{}
		if err := r.Get(ctx, types.NamespacedName{Name: v.Name, Namespace: v.Namespace}, c); err != nil {
			if errors.IsNotFound(err) {
				if err := r.Create(ctx, v); err != nil {
					aloysLog.Error(err, "create cm failed")
					return ctrl.Result{}, err
				}
			}
		} else {
			if err := r.Update(ctx, v); err != nil {
				aloysLog.Error(err, "update cm failed")
				return ctrl.Result{}, err
			}
		}
	}

	// deploy
	deploy := templateDeployment.NewDeployment(&aloys)
	// SetControllerReference是做了Owner设置，k8s GC在删除一个对象时，任何ownerReference是该对象的对象都会被清除
	// 调用 SetControllerReference，设置对象的 owner ，同时设置 contrller 和 blockdelete 来帮助垃圾回收，以及后续要 watch 的能力
	if err := ctrl.SetControllerReference(&aloys, deploy, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	d := &appsv1.Deployment{}
	if err := r.Get(ctx, types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, d); err != nil {
		if errors.IsNotFound(err) {
			if err := r.Create(ctx, deploy); err != nil {
				aloysLog.Error(err, "create deploy failed")
				return ctrl.Result{}, err
			}
		}
	} else {
		if err := r.Update(ctx, deploy); err != nil {
			aloysLog.Error(err, "update deploy failed")
			return ctrl.Result{}, err
		}
	}

	// service
	svc := templateService.NewService(&aloys)
	if err := ctrl.SetControllerReference(&aloys, svc, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	s := &corev1.Service{}
	if err := r.Get(ctx, types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, s); err != nil {
		if errors.IsNotFound(err) && aloys.Spec.Service.Enable {
			if err := r.Create(ctx, svc); err != nil {
				aloysLog.Error(err, "create service failed")
				return ctrl.Result{}, err
			}
		}
	} else if aloys.Spec.Service.Enable {
		if err := r.Update(ctx, svc); err != nil {
			aloysLog.Error(err, "update service failed")
			return ctrl.Result{}, err
		}
	} else {
		if err := r.Delete(ctx, svc); err != nil {
			aloysLog.Error(err, "delete service failed")
			return ctrl.Result{}, err
		}
	}

	// ingress
	// 这里暂时有一些小问题，ingress开启，service 没开启，虽然不会创建，也是也不会报错，想要直接提交失败，需要使用 webhook
	ing := templateIngress.NewIngress(&aloys)
	if err := ctrl.SetControllerReference(&aloys, ing, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	i := &netv1.Ingress{}
	if err := r.Get(ctx, types.NamespacedName{Name: ing.Name, Namespace: ing.Namespace}, i); err != nil {
		if errors.IsNotFound(err) && aloys.Spec.Service.Enable && aloys.Spec.Ingress.Enable {
			if err := r.Create(ctx, ing); err != nil {
				aloysLog.Error(err, "create ingress failed")
				return ctrl.Result{}, err
			}
		}
	} else if aloys.Spec.Service.Enable && aloys.Spec.Ingress.Enable {
		if err := r.Update(ctx, ing); err != nil {
			aloysLog.Error(err, "update ingress failed")
			return ctrl.Result{}, err
		}
	} else {
		if err := r.Delete(ctx, ing); err != nil {
			aloysLog.Error(err, "delete ingress failed")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
// NewControllerManagedBy 设置AloysReconciler 管理mgr
// For(&zyv1.Aloys{} controller 监控什么资源触发Reconcile
func (r *AloysReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&zyv1.Aloys{}).
		// 因为需要操作Deployment，svc,ing，所以要设置
		Owns(&corev1.ConfigMap{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&netv1.Ingress{}).
		Complete(r)
}
