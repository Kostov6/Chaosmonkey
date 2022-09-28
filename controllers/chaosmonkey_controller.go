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
	"time"

	podv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	onboardingv1 "github.com/Kostov6/chaosmonkey/api/v1"
)

// ChaosmonkeyReconciler reconciles a Chaosmonkey object
type ChaosmonkeyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=onboarding.my.domain,resources=chaosmonkeys,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=onboarding.my.domain,resources=chaosmonkeys/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=onboarding.my.domain,resources=chaosmonkeys/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Chaosmonkey object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ChaosmonkeyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var chaosmonkey onboardingv1.Chaosmonkey
	if err := r.Get(ctx, req.NamespacedName, &chaosmonkey); err != nil {
		log.Log.Error(err, "unable to fetch Chaosmonkey")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	scheduledReconcile := false
	//defaultTime := v1.NewTime(time.Time{})
	// check if pods should be deleted
	if time.Now().After(chaosmonkey.Status.LastDelete.Add(chaosmonkey.Spec.Period.Duration)) {
		if chaosmonkey.Spec.PodName != "" {
			pod := &podv1.Pod{}
			if err := r.Client.Get(ctx, client.ObjectKey{Namespace: chaosmonkey.Spec.Namespace, Name: chaosmonkey.Spec.PodName}, pod); err != nil {
				log.Log.Info("pod with name " + chaosmonkey.Spec.PodName + " does not exist")
			} else if err := r.Client.Delete(ctx, pod); err != nil {
				log.Log.Error(err, "pod with name "+chaosmonkey.Spec.PodName+" cannot be deleted")
				chaosmonkey.Status.State = "Error"
				r.Status().Update(ctx, &chaosmonkey)
				return ctrl.Result{}, client.IgnoreNotFound(err)
			}
		}

		var podList podv1.PodList
		if len(chaosmonkey.Spec.WithFields) > 0 {
			//".metadata.name": "broken-pod"
			if err := r.List(ctx, &podList, client.InNamespace(chaosmonkey.Spec.Namespace), client.MatchingFields(chaosmonkey.Spec.WithFields)); err != nil {
				log.Log.Error(err, "unable to list pods matching by fields")
				chaosmonkey.Status.State = "Error"
				r.Status().Update(ctx, &chaosmonkey)
				return ctrl.Result{}, err
			}
			for _, pod := range podList.Items {
				log.Log.Info("Deleting labeled pod " + pod.Name)
				if err := r.Client.Delete(ctx, &pod); err != nil {
					log.Log.Error(err, "pod with name "+chaosmonkey.Spec.PodName+" cannot be deleted")
					chaosmonkey.Status.State = "Error"
					r.Status().Update(ctx, &chaosmonkey)
					return ctrl.Result{}, client.IgnoreNotFound(err)
				}
			}
		}
		if len(chaosmonkey.Spec.WithLabels) > 0 {
			if err := r.List(ctx, &podList, client.InNamespace(chaosmonkey.Spec.Namespace), client.MatchingLabels(chaosmonkey.Spec.WithLabels)); err != nil {
				log.Log.Error(err, "unable to list pods matching by labels")
				chaosmonkey.Status.State = "Error"
				r.Status().Update(ctx, &chaosmonkey)
				return ctrl.Result{}, err
			}
			for _, pod := range podList.Items {
				log.Log.Info("Deleting labeled pod " + pod.Name)
				if err := r.Client.Delete(ctx, &pod); err != nil {
					log.Log.Error(err, "pod with name "+chaosmonkey.Spec.PodName+" cannot be deleted")
					chaosmonkey.Status.State = "Error"
					r.Status().Update(ctx, &chaosmonkey)
					return ctrl.Result{}, client.IgnoreNotFound(err)
				}
			}
		}

		log.Log.Info("Pod deletion cycle")
		chaosmonkey.Status.LastDelete.Time = time.Now()
		scheduledReconcile = true

	}

	log.Log.Info("Successfull reconciliating")
	chaosmonkey.Status.State = "Running"
	r.Status().Update(ctx, &chaosmonkey)
	if scheduledReconcile {
		return ctrl.Result{RequeueAfter: time.Duration(chaosmonkey.Spec.Period.Nanoseconds())}, nil
	}
	return ctrl.Result{Requeue: false}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChaosmonkeyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&onboardingv1.Chaosmonkey{}).
		Complete(r)
}
