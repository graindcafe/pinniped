// Copyright 2021-2022 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package conditionsutil

import (
	"sort"

	"k8s.io/apimachinery/pkg/api/equality"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configv1alpha1 "go.pinniped.dev/generated/latest/apis/supervisor/config/v1alpha1"
	idpv1alpha1 "go.pinniped.dev/generated/latest/apis/supervisor/idp/v1alpha1"
	"go.pinniped.dev/internal/plog"
)

// MergeIDPConditions merges conditions into conditionsToUpdate. If returns true if it merged any error conditions.
func MergeIDPConditions(conditions []*idpv1alpha1.Condition, observedGeneration int64, conditionsToUpdate *[]idpv1alpha1.Condition, log plog.MinLogger) bool {
	hadErrorCondition := false
	for i := range conditions {
		cond := conditions[i].DeepCopy()
		cond.LastTransitionTime = v1.Now()
		cond.ObservedGeneration = observedGeneration
		if mergeIDPCondition(conditionsToUpdate, cond) {
			log.Info("updated condition", "type", cond.Type, "status", cond.Status, "reason", cond.Reason, "message", cond.Message)
		}
		if cond.Status == idpv1alpha1.ConditionFalse {
			hadErrorCondition = true
		}
	}
	sort.SliceStable(*conditionsToUpdate, func(i, j int) bool {
		return (*conditionsToUpdate)[i].Type < (*conditionsToUpdate)[j].Type
	})
	return hadErrorCondition
}

// mergeIDPCondition merges a new idpv1alpha1.Condition into a slice of existing conditions. It returns true
// if the condition has meaningfully changed.
func mergeIDPCondition(existing *[]idpv1alpha1.Condition, new *idpv1alpha1.Condition) bool {
	// Find any existing condition with a matching type.
	var old *idpv1alpha1.Condition
	for i := range *existing {
		if (*existing)[i].Type == new.Type {
			old = &(*existing)[i]
			continue
		}
	}

	// If there is no existing condition of this type, append this one and we're done.
	if old == nil {
		*existing = append(*existing, *new)
		return true
	}

	// Set the LastTransitionTime depending on whether the status has changed.
	new = new.DeepCopy()
	if old.Status == new.Status {
		new.LastTransitionTime = old.LastTransitionTime
	}

	// If anything has actually changed, update the entry and return true.
	if !equality.Semantic.DeepEqual(old, new) {
		*old = *new
		return true
	}

	// Otherwise the entry is already up to date.
	return false
}

// MergeConfigConditions merges conditions into conditionsToUpdate. If returns true if it merged any error conditions.
func MergeConfigConditions(conditions []*configv1alpha1.Condition, observedGeneration int64, conditionsToUpdate *[]configv1alpha1.Condition, log plog.MinLogger) bool {
	hadErrorCondition := false
	for i := range conditions {
		cond := conditions[i].DeepCopy()
		cond.LastTransitionTime = v1.Now()
		cond.ObservedGeneration = observedGeneration
		if mergeConfigCondition(conditionsToUpdate, cond) {
			log.Info("updated condition", "type", cond.Type, "status", cond.Status, "reason", cond.Reason, "message", cond.Message)
		}
		if cond.Status == configv1alpha1.ConditionFalse {
			hadErrorCondition = true
		}
	}
	sort.SliceStable(*conditionsToUpdate, func(i, j int) bool {
		return (*conditionsToUpdate)[i].Type < (*conditionsToUpdate)[j].Type
	})
	return hadErrorCondition
}

// mergeConfigCondition merges a new idpv1alpha1.Condition into a slice of existing conditions. It returns true
// if the condition has meaningfully changed.
func mergeConfigCondition(existing *[]configv1alpha1.Condition, new *configv1alpha1.Condition) bool {
	// Find any existing condition with a matching type.
	var old *configv1alpha1.Condition
	for i := range *existing {
		if (*existing)[i].Type == new.Type {
			old = &(*existing)[i]
			continue
		}
	}

	// If there is no existing condition of this type, append this one and we're done.
	if old == nil {
		*existing = append(*existing, *new)
		return true
	}

	// Set the LastTransitionTime depending on whether the status has changed.
	new = new.DeepCopy()
	if old.Status == new.Status {
		new.LastTransitionTime = old.LastTransitionTime
	}

	// If anything has actually changed, update the entry and return true.
	if !equality.Semantic.DeepEqual(old, new) {
		*old = *new
		return true
	}

	// Otherwise the entry is already up to date.
	return false
}
