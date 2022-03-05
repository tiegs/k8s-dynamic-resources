/*
Copyright 2022 Tilman Eggers.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	metakubev1alpha1 "tilmaneggers.de/k8s-meta-ressource-manager/api/v1alpha1"
	"time"
)

// MetaRessourceReconciler reconciles a MetaRessource object
type MetaRessourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=meta.kube,resources=metaressources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=meta.kube,resources=metaressources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=meta.kube,resources=metaressources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MetaRessource object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *MetaRessourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// TODO(user): your logic here
	log.Info("Reconciling...")

	// Retrieve MetaRessource
	var metaRessource metakubev1alpha1.MetaRessource
	if err := r.Get(ctx, req.NamespacedName, &metaRessource); err != nil {
		log.Error(err, "unable to fetch MetaRessource")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Retrieve ressources managed by

	// https://stackoverflow.com/questions/61200605/generic-client-get-for-custom-kubernetes-go-operator

	// Get arbitrary object
	//u := &unstructured.Unstructured{}
	//
	//u.SetGroupVersionKind(schema.GroupVersionKind{
	//	Group:   "",
	//	Version: "v1",
	//	Kind:    "Secret",
	//})
	//
	//err := r.Get(ctx, client.ObjectKey{
	//	Namespace: "default",
	//	Name:      "dummy-secret",
	//}, u)
	//
	//if err != nil {
	//	return ctrl.Result{}, err
	//} else {
	//	log.Info("Object identified")
	//}
	//
	//return ctrl.Result{}, nil

	// Create arbitrary object
	u := &unstructured.Unstructured{}

	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Secret",
	})

	//u.SetNamespace("default")
	u.SetName("asdf-test-secret")

	gvk, err := apiutil.GVKForObject(&metaRessource, r.Scheme)
	if err != nil {
		return ctrl.Result{}, err
	}

	ref := *metav1.NewControllerRef(&metaRessource, gvk)

	refs := u.GetOwnerReferences()
	refs = append(refs, ref)

	u.SetOwnerReferences(append(u.GetOwnerReferences(), ref))

	//u.SetUnstructuredContent(map[string] string {"France":"Paris","Italy":"Rome","Japan":"Tokyo","India":"New Delhi"})

	err = r.Create(ctx, u)

	//err := r.Get(ctx, client.ObjectKey{
	//	Namespace: "default",
	//	Name:      "dummy-secret",
	//}, u)

	if err != nil {
		return ctrl.Result{RequeueAfter: 10 * time.Second}, err
	} else {
		log.Info("Object created!")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MetaRessourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metakubev1alpha1.MetaRessource{}).
		Complete(r)
}
