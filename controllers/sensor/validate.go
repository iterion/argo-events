/*
Copyright 2018 BlackRock, Inc.

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

package sensor

import (
	"fmt"
	"time"

	"github.com/argoproj/argo-events/common"
	"github.com/argoproj/argo-events/pkg/apis/sensor/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ValidateSensor accepts a sensor and performs validation against it
// we return an error so that it can be logged as a message on the sensor status
// the error is ignored by the operation context as subsequent re-queues would produce the same error.
// Exporting this function so that external APIs can use this to validate sensor resource.
func ValidateSensor(s *v1alpha1.Sensor) error {
	if err := validateSignals(s.Spec.Signals); err != nil {
		return err
	}
	if err := validateTriggers(s.Spec.Triggers); err != nil {
		return err
	}
	return nil
}

func validateTriggers(triggers []v1alpha1.Trigger) error {
	if len(triggers) < 1 {
		return fmt.Errorf("no triggers found")
	}

	for _, trigger := range triggers {
		if trigger.Name == "" {
			return fmt.Errorf("trigger must define a name")
		}
		// each trigger must have a message or a resource
		if trigger.Resource == nil {
			return fmt.Errorf("trigger '%s' does not contain an absolute action", trigger.Name)
		}
	}
	return nil
}

// perform a check to see that each signal defines one of and at most one of:
// (stream, artifact, calendar, resource, webhook)
func validateSignals(signals []v1alpha1.Signal) error {
	if len(signals) < 1 {
		return fmt.Errorf("no signals found")
	}
	for _, signal := range signals {
		if signal.Name == "" {
			return fmt.Errorf("signal must define a name")
		}
		if err := validateSignalFilter(signal.Filters); err != nil {
			return err
		}
	}
	return nil
}

func validateSignalFilter(filter v1alpha1.SignalFilter) error {
	if filter.Time != nil {
		if err := validateSignalTimeFilter(filter.Time); err != nil {
			return err
		}
	}
	return nil
}

func validateSignalTimeFilter(tFilter *v1alpha1.TimeFilter) error {
	currentT := metav1.Time{Time: time.Now().UTC()}
	currentMonth := fmt.Sprintf("%d", int(currentT.Month()))
	if int(currentT.Month()) < 10 {
		currentMonth = "0" + currentMonth
	}
	currentDay := fmt.Sprintf("%d", int(currentT.Day()))
	if int(currentT.Day()) < 10 {
		currentDay = "0" + currentDay
	}
	currentTStr := fmt.Sprintf("%d-%s-%s", currentT.Year(), currentMonth, currentDay)
	if tFilter.Start != "" && tFilter.Stop != "" {
		startTime, err := time.Parse(common.StandardTimeFormat, currentTStr+" "+tFilter.Start)
		if err != nil {
			return err
		}
		stopTime, err := time.Parse(common.StandardTimeFormat, currentTStr+" "+tFilter.Stop)
		if err != nil {
			return err
		}
		if stopTime.Before(startTime) || startTime.Equal(stopTime) {
			return fmt.Errorf("invalid signal time filter: stop '%s' is before or equal to start '%s", tFilter.Stop, tFilter.Start)
		}
	}
	if tFilter.Stop != "" {
		stopTime, err := time.Parse(common.StandardTimeFormat, currentTStr+" "+tFilter.Stop)
		if err != nil {
			return err
		}
		stopTime = stopTime.UTC()
		if stopTime.Before(currentT.Time) {
			return fmt.Errorf("invalid signal time filter: stop '%s' is before the current time '%s'", tFilter.Stop, currentT)
		}
	}
	return nil
}
