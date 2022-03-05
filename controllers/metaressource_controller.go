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
	"errors"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
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

	// Todo: Retrieve resources managed by this metaressource or create new
	u := &unstructured.Unstructured{}

	// https://stackoverflow.com/questions/61200605/generic-client-get-for-custom-kubernetes-go-operator

	// Prepare Target object
	u.SetUnstructuredContent(metaRessource.Spec.Target.Object)

	// Defined owner reference
	gvk, err := apiutil.GVKForObject(&metaRessource, r.Scheme)
	if err != nil {
		return ctrl.Result{}, err
	}

	ref := *metav1.NewControllerRef(&metaRessource, gvk)
	refs := u.GetOwnerReferences()
	refs = append(refs, ref)

	u.SetOwnerReferences(append(u.GetOwnerReferences(), ref))

	// Resolve Transformations
	for _, trans := range metaRessource.Spec.Transformations {

		// Handle fieldFrom transfomation
		src := &unstructured.Unstructured{}

		src.SetAPIVersion(trans.FieldFrom.APIVersion)
		src.SetKind(trans.FieldFrom.Kind)

		// Todo: More elaborate matchers
		key := client.ObjectKey{Namespace: metaRessource.Namespace, Name: trans.FieldFrom.Name}

		err := r.Get(ctx, key, src)
		if err != nil {
			log.Error(err, "Failed to retrieve fieldFrom source object")
			return ctrl.Result{}, err
		}

		// https://iximiuz.com/en/posts/kubernetes-api-go-types-and-common-machinery/

		// Split FieldSpec and retrieve its string value
		spec := strings.Split(trans.FieldFrom.FieldSpec, ".")
		data, found, err := unstructured.NestedString(src.UnstructuredContent(), spec...)
		if err != nil {
			log.Error(err, fmt.Sprintf("Failed to retrieve field '%s' source object", trans.TargetField))
			return ctrl.Result{}, err
		} else if !found {
			return ctrl.Result{}, errors.New(fmt.Sprintf("Field '%s' not found on source object", trans.TargetField))
		}

		// Inject into target field
		spec = strings.Split(trans.TargetField, ".")
		err = unstructured.SetNestedField(u.Object, data, spec...)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

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
