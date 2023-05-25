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
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	testv1alpha1 "github.com/aliok/two-operators/operator1/api/v1alpha1"
)

// InstallationAReconciler reconciles a InstallationA object
type InstallationAReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const Backoff = time.Second * 2

//+kubebuilder:rbac:groups=test.aliok,resources=installationas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=test.aliok,resources=installationas/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=test.aliok,resources=installationas/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the InstallationA object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *InstallationAReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	instA := &testv1alpha1.InstallationA{}
	err := r.Get(ctx, req.NamespacedName, instA)
	if err != nil {
		logger.Error(err, "unable to fetch InstallationA")
		return ctrl.Result{}, err
	}
	if instA.Status.Good {
		logger.Info("InstallationA is good")
		return ctrl.Result{}, nil
	}

	// initialize Start if not already set
	if instA.Status.Start == 0 {
		instA.Status.Start = time.Now().UnixNano()
		err = r.Status().Update(ctx, instA)
		if err != nil {
			logger.Error(err, "unable to update InstallationA status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// if it is already set, check if it is time to set Good=true
	if instA.Status.Start+(int64(instA.Spec.Delay)*int64(time.Second)) < time.Now().UnixNano() {
		instA.Status.Good = true
		err = r.Status().Update(ctx, instA)
		if err != nil {
			logger.Error(err, "unable to update InstallationA status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	logger.Info(fmt.Sprintf("InstallationA not ready yet, requeuing in %f seconds", Backoff.Seconds()))
	return ctrl.Result{RequeueAfter: Backoff}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InstallationAReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testv1alpha1.InstallationA{}).
		Complete(r)
}
