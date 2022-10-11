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

package main

import (
	"flag"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	zyv1 "operator-simplification/api/v1"
	"operator-simplification/controllers"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	// 初始化scheme并注册Scheme，这时候就可以通过结构体直接拿到对应的字段了，操作scheme就和操作结构体一样
	// scheme GVK 的映射
	// 获取默认类型	scheme对象
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	// 获取添加类型的scheme对象
	utilruntime.Must(zyv1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	// var x string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	// 参数传入到controller
	// flag.StringVar(&x, "x", "", "test")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	// 本地测试webhook
	// options := ctrl.Options{
	// 	Scheme:                 scheme,
	// 	MetricsBindAddress:     metricsAddr,
	// 	Port:                   9443,
	// 	HealthProbeBindAddress: probeAddr,
	// 	LeaderElection:         enableLeaderElection,
	// 	LeaderElectionID:       "48e422bc.tech",
	// }
	// if os.Getenv("ENVIRONMENT") == "DEV" {
	// 	path, err := os.Getwd()
	// 	if err != nil {
	// 		setupLog.Error(err, "unable to get work dir")
	// 		os.Exit(1)
	// 	}
	// 	options.CertDir = path + "/certs"
	// }
	// mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), options)
	// if err != nil {
	// 	setupLog.Error(err, "unable to start manager")
	// 	os.Exit(1)
	// }

	// 创建Manager，传入scheme和其他参数
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		// 初始化scheme
		Scheme: scheme,
		// 资源探针端口
		MetricsBindAddress: metricsAddr,
		// webhook服务器服务的端口号
		Port: 9443,
		// 健康检查端口
		HealthProbeBindAddress: probeAddr,
		// 高可用设置参数，还不知道什么意思
		LeaderElection:   enableLeaderElection,
		LeaderElectionID: "48e422bc.tech",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Reconciler添加Manager
	// AloysReconciler是实现的具体的 controller
	// controller注册到Manager
	if err = (&controllers.AloysReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		// 在这里写入就可以传入到 controller
		// X: x,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Aloys")
		os.Exit(1)
	}

	// webhook 添加代码
	if err = (&zyv1.Aloys{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "Aloys")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	// 健康检查相关设置
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	// 启动manager
	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
