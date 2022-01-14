/*
Copyright 2021.

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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	webappv1 "mydev.org/guestbook/api/v1"
)

// RestAPIReconciler reconciles a RestAPI object
type RestAPIReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webapp.mydev.org,resources=restapis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webapp.mydev.org,resources=restapis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webapp.mydev.org,resources=restapis/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;get;patch
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;get;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RestAPI object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *RestAPIReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.
		FromContext(ctx).
		WithValues("restapi", req.NamespacedName)

	log.Info("reconciling RestAPI")

	var restAPI webappv1.RestAPI
	if err := r.Get(ctx, req.NamespacedName, &restAPI); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("fetching Redis object")
	var redis webappv1.Redis
	redisName := client.ObjectKey{Name: restAPI.Spec.RedisName, Namespace: req.Namespace}
	if err := r.Get(ctx, redisName, &redis); err != nil {
		return ctrl.Result{}, err
	}

	log.Info("creating Deployment object")
	deployment, err := r.desiredDeployment(restAPI, redis)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("creating Service object")
	svc, err := r.desiredService(restAPI)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("applying changes for RestAPI")

	applyOpts := []client.PatchOption{client.ForceOwnership, client.FieldOwner("guestbook-controller")}

	err = r.Patch(ctx, &deployment, client.Apply, applyOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.Patch(ctx, &svc, client.Apply, applyOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}

	restAPI.Status.URL = urlForService(svc, restAPI.Spec.Frontend.ServingPort)

	err = r.Status().Update(ctx, &restAPI)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("reconciled RestAPI")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RestAPIReconciler) SetupWithManager(mgr ctrl.Manager) error {
	mgr.GetFieldIndexer().IndexField(
		context.TODO(),
		&webappv1.RestAPI{},
		".spec.redisName",
		func(obj client.Object) []string {
			redisName := obj.(*webappv1.RestAPI).Spec.RedisName
			if redisName == "" {
				return nil
			}
			return []string{redisName}
		})

	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.RestAPI{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.Deployment{}).
		Watches(
			&source.Kind{Type: &webappv1.Redis{}},
			handler.EnqueueRequestsFromMapFunc(r.restAPIsUsingRedis)).
		Complete(r)
}
