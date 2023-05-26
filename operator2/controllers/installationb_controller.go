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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	operator1testv1alpha1 "github.com/aliok/two-operators/operator1/api/v1alpha1"
	operator2testv1alpha1 "github.com/aliok/two-operators/operator2/api/v1alpha1"
)

// InstallationBReconciler reconciles a InstallationB object
type InstallationBReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const Backoff = time.Second * 2

//+kubebuilder:rbac:groups=test.aliok,resources=installationbs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=test.aliok,resources=installationbs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=test.aliok,resources=installationbs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the InstallationB object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *InstallationBReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciling InstallationB", "namespace", req.Namespace, "name", req.Name)

	instB := &operator2testv1alpha1.InstallationB{}
	err := r.Get(ctx, req.NamespacedName, instB)
	if err != nil {
		logger.Error(err, "unable to fetch InstallationB")
		return ctrl.Result{}, err
	}
	if instB.Status.Good {
		logger.Info("InstallationB is good")
		return ctrl.Result{}, nil
	}

	// TODO: need to watch operator1's installation CR actually

	instA := &operator1testv1alpha1.InstallationA{}
	// Hardcoded name
	err = r.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: "operator1-installation"}, instA)

	// if there's an error and the error is not NotFound, return the error
	if err != nil && !errors.IsNotFound(err) {
		logger.Error(err, "unable to fetch InstallationA")
		return ctrl.Result{}, err
	}

	if err != nil && errors.IsNotFound(err) {
		// create installation for operator1
		logger.Info("InstallationA not found, creating it")
		err = r.Create(ctx, instA)
		if err != nil {
			logger.Error(err, "unable to create InstallationA")
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: Backoff}, nil
	}

	// if there's no error, check the status of operator1's installation
	if instA.Status.Good != true {
		logger.Info(fmt.Sprintf("InstallationA not ready yet, requeuing in %f seconds", Backoff.Seconds()))
		return ctrl.Result{RequeueAfter: Backoff}, nil
	}

	// if there's no error and the status is good, set the status of operator2's installation to good
	instB.Status.Good = true
	err = r.Status().Update(ctx, instB)
	if err != nil {
		logger.Error(err, "unable to update InstallationB status")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InstallationBReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operator2testv1alpha1.InstallationB{}).
		Complete(r)
}
