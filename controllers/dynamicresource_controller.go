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
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/util/jsonpath"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"k8s.io/kubectl/pkg/cmd/get"

	dynamickubev1alpha1 "github.com/tiegs/k8s-dynamic-resources/api/v1alpha1"
)

// DynamicResourceReconciler reconciles a DynamicResource object
type DynamicResourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=dynamic.kube,resources=dynamicresources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dynamic.kube,resources=dynamicresources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dynamic.kube,resources=dynamicresources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *DynamicResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("Reconciling...")

	// Retrieve DynamicResource
	var dynamicResource dynamickubev1alpha1.DynamicResource
	if err := r.Get(ctx, req.NamespacedName, &dynamicResource); err != nil {
		//logger.Error(err, "unable to fetch MetaRessource")

		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	u := &unstructured.Unstructured{}

	// https://stackoverflow.com/questions/61200605/generic-client-get-for-custom-kubernetes-go-operator

	// Prepare Target object
	u.SetUnstructuredContent(dynamicResource.Spec.Target.Object)

	// Define owner reference
	gvk, err := apiutil.GVKForObject(&dynamicResource, r.Scheme)
	if err != nil {
		return ctrl.Result{}, err
	}

	ref := *metav1.NewControllerRef(&dynamicResource, gvk)
	refs := u.GetOwnerReferences()
	refs = append(refs, ref)

	u.SetOwnerReferences(append(u.GetOwnerReferences(), ref))

	// Resolve Transformations
	for _, trans := range dynamicResource.Spec.Transformations {

		// Handle fieldFrom transfomation
		src := &unstructured.Unstructured{}

		src.SetAPIVersion(trans.FieldFrom.APIVersion)
		src.SetKind(trans.FieldFrom.Kind)

		// Todo: More elaborate matchers
		key := client.ObjectKey{Namespace: dynamicResource.Namespace, Name: trans.FieldFrom.Name}

		err := r.Get(ctx, key, src)
		if err != nil {
			//logger.Error(err, "Failed to retrieve fieldFrom source object")
			return ctrl.Result{}, err
		}

		// https://iximiuz.com/en/posts/kubernetes-api-go-types-and-common-machinery/

		// Parse jsonpath
		// https://kubernetes.io/docs/reference/kubectl/jsonpath/
		fields, err := get.RelaxedJSONPathExpression(trans.FieldFrom.FieldSpec)
		if err != nil {
			return ctrl.Result{}, errors.WithMessage(err, "Invalid FieldSpec (needs to be a valid jsonpath)")
		}

		j := jsonpath.New("")
		err = j.Parse(fields)
		if err != nil {
			return ctrl.Result{}, errors.WithMessage(err, "Failed to parse FieldSpec (needs to be a valid jsonpath)")
		}

		values, err := j.FindResults(src.Object)
		if err != nil {
			return ctrl.Result{}, errors.WithMessage(err, "Failed to execute FieldSpec")

		}

		// Allow only single-result jsonpaths
		var data string

		if len(values) == 0 {
			return ctrl.Result{}, errors.New(fmt.Sprintf("JSONPath '%s' did not yield any result", trans.FieldFrom.FieldSpec))
		} else if len(values) > 1 {
			return ctrl.Result{}, errors.New(fmt.Sprintf("JSONPath '%s' yield '%d' result", trans.FieldFrom.FieldSpec, len(values)))
		} else {
			buf := &bytes.Buffer{}
			err = j.PrintResults(buf, values[0])
			if err != nil {
				return ctrl.Result{}, err
			}

			data = buf.String()
		}

		// Inject into target field
		spec := strings.Split(trans.TargetField, ".")
		err = unstructured.SetNestedField(u.Object, data, spec...)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	err = r.Update(ctx, u)
	if err != nil {
		return ctrl.Result{RequeueAfter: 10 * time.Second}, err
	}

	logger.Info("Dynamic resource reconciled!", "resource", client.ObjectKeyFromObject(u))

	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DynamicResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dynamickubev1alpha1.DynamicResource{}).
		Complete(r)
}
