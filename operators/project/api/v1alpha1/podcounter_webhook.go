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

package v1alpha1

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	validationutils "k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var podcounterlog = logf.Log.WithName("podcounter-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *PodCounter) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(&PodCounterCustomValidator{}).
		WithDefaulter(&PodCounterCustomDefaulter{
			DefaultNamespace: "default",
		}).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-counters-cloudprog-polito-it-v1alpha1-podcounter,mutating=true,failurePolicy=fail,sideEffects=None,groups=counters.cloudprog.polito.it,resources=podcounters,verbs=create;update,versions=v1alpha1,name=mpodcounter.kb.io,admissionReviewVersions=v1

// PodCounterCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind PodCounter when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type PodCounterCustomDefaulter struct {
	DefaultNamespace string
}

var _ webhook.CustomDefaulter = &PodCounterCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind podcounter.
func (d *PodCounterCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	podcounter, ok := obj.(*PodCounter)

	if !ok {
		return fmt.Errorf("expected a PodCounter object but got %T", obj)
	}
	podcounterlog.Info("Defaulting for PodCounter", "name", podcounter.GetName())

	// Set default values
	d.applyDefaults(podcounter)
	return nil
}

// applyDefaults applies default values to podcounter fields.
func (d *PodCounterCustomDefaulter) applyDefaults(podcounter *PodCounter) {
	if podcounter.Spec.Namespace == "" {
		podcounter.Spec.Namespace = d.DefaultNamespace
	}
}

/*
We can validate our CRD beyond what's possible with declarative
validation. Generally, declarative validation should be sufficient, but
sometimes more advanced use cases call for complex validation.

For instance, we'll see below that we use this to validate a well-formed cron
schedule without making up a long regular expression.

If `webhook.CustomValidator` interface is implemented, a webhook will automatically be
served that calls the validation.

The `ValidateCreate`, `ValidateUpdate` and `ValidateDelete` methods are expected
to validate its receiver upon creation, update and deletion respectively.
We separate out ValidateCreate from ValidateUpdate to allow behavior like making
certain fields immutable, so that they can only be set on creation.
ValidateDelete is also separated from ValidateUpdate to allow different
validation behavior on deletion.
Here, however, we just use the same shared validation for `ValidateCreate` and
`ValidateUpdate`. And we do nothing in `ValidateDelete`, since we don't need to
validate anything on deletion.
*/

/*
This marker is responsible for generating a validation webhook manifest.
*/
// +kubebuilder:webhook:path=/validate-counters-cloudprog-polito-it-v1alpha1-podcounter,mutating=false,failurePolicy=fail,sideEffects=None,groups=counters.cloudprog.polito.it,resources=podcounters,verbs=create;update,versions=v1alpha1,name=vpodcounter-v1alpha1.kb.io,admissionReviewVersions=v1

// PodCounterCustomValidator struct is responsible for validating the podcounter resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type PodCounterCustomValidator struct {
}

var _ webhook.CustomValidator = &PodCounterCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type PodCounter.
func (v *PodCounterCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	podcounter, ok := obj.(*PodCounter)
	if !ok {
		return nil, fmt.Errorf("expected a PodCounter object but got %T", obj)
	}
	podcounterlog.Info("Validation for PodCounter upon creation", "name", podcounter.GetName())

	return nil, validatePodCounter(podcounter)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type podcounter.
func (v *PodCounterCustomValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	podcounter, ok := newObj.(*PodCounter)
	if !ok {
		return nil, fmt.Errorf("expected a PodCounter object for the newObj but got %T", newObj)
	}
	podcounterlog.Info("Validation for PodCounter upon update", "name", podcounter.GetName())

	return nil, validatePodCounter(podcounter)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type podcounter.
func (v *PodCounterCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	podcounter, ok := obj.(*PodCounter)
	if !ok {
		return nil, fmt.Errorf("expected a PodCounter object but got %T", obj)
	}
	podcounterlog.Info("Validation for PodCounter upon deletion", "name", podcounter.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}

/*
We validate the name and the spec of the podcounter.
*/

// validatePodCounter validates the fields of a podcounter object.
func validatePodCounter(podcounter *PodCounter) error {
	var allErrs field.ErrorList
	if err := validatePodCounterName(podcounter); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := validatePodCounterSpec(podcounter); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "counters.cloudprog.polito.it", Kind: "PodCounter"},
		podcounter.Name, allErrs)
}

/*
Some fields are declaratively validated by OpenAPI schema.
You can find kubebuilder validation markers (prefixed
with `// +kubebuilder:validation`) in the
[Designing an API](api-design.md) section.
You can find all of the kubebuilder supported markers for
declaring validation by running `controller-gen crd -w`,
or [here](/reference/markers/crd-validation.md).
*/

func validatePodCounterSpec(podcounter *PodCounter) *field.Error {
	// The field helpers from the kubernetes API machinery help us return nicely
	// structured validation errors.
	podcounterlog.Info("Validating spec", "namespace", podcounter.Spec.Namespace)
	return validateNamespace(podcounter.Spec.Namespace)
}

func validateNamespace(namespace string) *field.Error {
	// namespaces cannot contain "kube-" prefix
	if len(namespace) > 4 && namespace[:5] == "kube-" {
		return field.Invalid(field.NewPath("spec").Child("namespace"), namespace, "cannot contain 'kube-' prefix")
	}
	return nil
}

/*
Validating the length of a string field can be done declaratively by
the validation schema.

But the `ObjectMeta.Name` field is defined in a shared package under
the apimachinery repo, so we can't declaratively validate it using
the validation schema.
*/

func validatePodCounterName(podcounter *PodCounter) *field.Error {
	if len(podcounter.ObjectMeta.Name) > validationutils.DNS1035LabelMaxLength-11 {
		// The job name length is 63 characters like all Kubernetes objects
		// (which must fit in a DNS subdomain). The podcounter controller appends
		// a 11-character suffix to the podcounter (`-$TIMESTAMP`) when creating
		// a job. The job name length limit is 63 characters. Therefore podcounter
		// names must have length <= 63-11=52. If we don't validate this here,
		// then job creation will fail later.
		return field.Invalid(field.NewPath("metadata").Child("name"), podcounter.ObjectMeta.Name, "must be no more than 52 characters")
	}
	return nil
}

// +kubebuilder:docs-gen:collapse=Validate object name
