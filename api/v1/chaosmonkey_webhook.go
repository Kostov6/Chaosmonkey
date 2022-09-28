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

package v1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var chaosmonkeylog = logf.Log.WithName("chaosmonkey-resource")

func (r *Chaosmonkey) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-onboarding-my-domain-v1-chaosmonkey,mutating=true,failurePolicy=fail,sideEffects=None,groups=onboarding.my.domain,resources=chaosmonkeys,verbs=create;update,versions=v1,name=mchaosmonkey.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Chaosmonkey{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Chaosmonkey) Default() {
	chaosmonkeylog.Info("default", "name", r.Name)

	if r.Spec.Namespace == "" {
		r.Spec.Namespace = "default"
	}
	if r.Spec.Period.Duration == 0 {
		r.Spec.Period.Duration = 10
	}
	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-onboarding-my-domain-v1-chaosmonkey,mutating=false,failurePolicy=fail,sideEffects=None,groups=onboarding.my.domain,resources=chaosmonkeys,verbs=create;update,versions=v1,name=vchaosmonkey.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Chaosmonkey{}

var _ webhook.Admission

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Chaosmonkey) ValidateCreate() error {
	chaosmonkeylog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return r.validateChaosmonkey()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Chaosmonkey) ValidateUpdate(old runtime.Object) error {
	chaosmonkeylog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return r.validateChaosmonkey()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Chaosmonkey) ValidateDelete() error {
	chaosmonkeylog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return r.validateChaosmonkey()
}

func (r *Chaosmonkey) validateChaosmonkey() error {
	if r.Spec.Period.Duration < 0 {
		return fmt.Errorf("Period is negative number")
	}

	return nil
}
