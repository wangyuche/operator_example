/*
Copyright 2023.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	examplev1 "github.com/wangyuche/operator_example/example1/api/v1"
)

// ExampleReconciler reconciles a Example object
type ExampleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=example.example,resources=examples,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=example.example,resources=examples/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=example.example,resources=examples/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Example object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ExampleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	config := &examplev1.Example{}
	err := r.Get(ctx, req.NamespacedName, config)
	if err != nil {
		l.Error(err, "get config err")
		return ctrl.Result{}, err
	}
	pg_pvc := r.pvcpg_data(config)
	err = r.Create(ctx, pg_pvc)
	if err != nil {
		l.Error(err, "pg_pvc")
		return ctrl.Result{}, err
	}
	pg_svc := r.svcpg_standalone(config)
	err = r.Create(ctx, pg_svc)
	if err != nil {
		l.Error(err, "pg_svc")
		return ctrl.Result{}, err
	}
	pg_sts := r.stspg_standalone(config)
	err = r.Create(ctx, pg_sts)
	if err != nil {
		l.Error(err, "pg_sts")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExampleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1.Example{}).
		Complete(r)
}
