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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	metakubev1alpha1 "tilmaneggers.de/k8s-meta-ressource-manager/api/v1alpha1"
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

	//foundSecret := &corev1.Secret{}
	//err := r.Get(ctx, req.NamespacedName, foundSecret)
	//if err != nil {
	//	// If a configMap name is provided, then it must exist
	//	// You will likely want to create an Event for the user to understand why their reconcile is failing.
	//	return ctrl.Result{}, err
	//}

	//// https://itnext.io/generically-working-with-kubernetes-resources-in-go-53bce678f887
	//resourceId := schema.GroupVersionResource{
	//	Group:    group,
	//	Version:  version,
	//	Resource: resource,
	//}
	//
	//dynamic.
	//
	//list, err := dynamic.Resource(resourceId).Namespace(namespace).
	//	List(ctx, metav1.ListOptions{})
	//
	//if err != nil {
	//	return nil, err
	//}

	// https://stackoverflow.com/questions/61200605/generic-client-get-for-custom-kubernetes-go-operator
	u := &unstructured.Unstructured{}

	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "Secret",
	})

	err := r.Get(ctx, client.ObjectKey{
		Namespace: "default",
		Name:      "dummy-secret",
	}, u)

	if err != nil {
		return ctrl.Result{}, err
	} else {
		log.Info("Object identified")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MetaRessourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metakubev1alpha1.MetaRessource{}).
		Complete(r)
}
