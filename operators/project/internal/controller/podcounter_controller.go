/*
Copyright 2025.

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

package controller

import (
	"context"

	countersv1alpha1 "cloudprog.polito.it/project/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const podCounterFinalizer = "finalizer.counters.cloudprog.polito.it"

// PodCounterReconciler reconciles a PodCounter object
type PodCounterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=counters.cloudprog.polito.it,resources=podcounters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=counters.cloudprog.polito.it,resources=podcounters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=counters.cloudprog.polito.it,resources=podcounters/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;
// +kubebuilder:rbac:groups=core,resources=pods/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodCounter object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *PodCounterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch PodCounter instance
	var podCounter countersv1alpha1.PodCounter
	if err := r.Get(ctx, req.NamespacedName, &podCounter); err != nil {
		log.Error(err, "unable to fetch PodCounter")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("Fetched PodCounter", "namespace", podCounter.Spec.Namespace)

	// List all Pods in the specified namespace
	// wait namespace is set by D
	var podList corev1.PodList
	if err := r.List(ctx, &podList, client.InNamespace(podCounter.Spec.Namespace)); err != nil {
		log.Error(err, "unable to list Pods")
		return ctrl.Result{}, err
	}

	// Update status with pod count
	podCounter.Status.Count = len(podList.Items)

	// Exclude Pods with the label "app=busybox"
	for _, pod := range podList.Items {
		if hasBusyboxLabel(&pod) {
			podCounter.Status.Count--
		}
	}

	if err := r.Status().Update(ctx, &podCounter); err != nil {
		log.Error(err, "unable to update PodCounter status")
		return ctrl.Result{}, err
	}

	// Handle deletion (Finalizer Logic)
	if !podCounter.DeletionTimestamp.IsZero() {
		// Resource is being deleted
		if controllerutil.ContainsFinalizer(&podCounter, podCounterFinalizer) {
			// Perform cleanup tasks
			log.Info("Performing finalizer cleanup", "PodCounter", podCounter.Name)
			r.cleanupPodCounter(ctx, &podCounter)

			// Remove finalizer after cleanup
			controllerutil.RemoveFinalizer(&podCounter, podCounterFinalizer)
			if err := r.Update(ctx, &podCounter); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Ensure the finalizer is added
	if !controllerutil.ContainsFinalizer(&podCounter, podCounterFinalizer) {
		controllerutil.AddFinalizer(&podCounter, podCounterFinalizer)
		if err := r.Update(ctx, &podCounter); err != nil {
			return ctrl.Result{}, err
		}
	}

	log.Info("Updated PodCounter", "namespace", podCounter.Spec.Namespace, "count", podCounter.Status.Count)
	return ctrl.Result{}, nil
}

// Cleanup logic when deleting a PodCounter
func (r *PodCounterReconciler) cleanupPodCounter(ctx context.Context, podCounter *countersv1alpha1.PodCounter) {
	// Example cleanup: Log that we're cleaning up
	log := log.FromContext(ctx)
	log.Info("Cleaning up resources for PodCounter", "name", podCounter.Name)

	// If there were external resources (e.g., deleting a ConfigMap), we would remove them here
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodCounterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&countersv1alpha1.PodCounter{}).
		Watches(
			&corev1.Pod{}, // This logic works in controller-runtime v0.15+
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
				_, ok := obj.(*corev1.Pod)
				if !ok {
					return nil
				}

				// Find associated PodCounter (adjust logic if needed)
				podCounterList := &countersv1alpha1.PodCounterList{}
				err := r.Client.List(ctx, podCounterList)
				if err != nil {
					return nil
				}

				// Enqueue reconcile requests
				var requests []reconcile.Request
				for _, pc := range podCounterList.Items {
					requests = append(requests, reconcile.Request{
						NamespacedName: types.NamespacedName{Name: pc.Name, Namespace: pc.Namespace},
					})
				}
				return requests
			}),
			builder.WithPredicates(
				predicate.Funcs{
					CreateFunc: func(e event.CreateEvent) bool {
						// Only process Pods that DO NOT have the label "app=busybox"
						return !hasBusyboxLabel(e.Object)
					},
					UpdateFunc: func(e event.UpdateEvent) bool {
						// Process updates only if the Pod DOES NOT have "app=busybox"
						return !hasBusyboxLabel(e.ObjectNew)
					},
					DeleteFunc: func(e event.DeleteEvent) bool {
						// Process deletions only if the Pod DID NOT have "app=busybox"
						return !hasBusyboxLabel(e.Object)
					},
					GenericFunc: func(e event.GenericEvent) bool {
						// Ignore generic events
						return false
					},
				}),
		).
		Complete(r)
}

// Helper function to check if a Pod has the "app=busybox" label
func hasBusyboxLabel(obj client.Object) bool {
	if obj == nil {
		return false
	}
	labels := obj.GetLabels()
	val, exists := labels["app"]
	return exists && val == "busybox"
}
