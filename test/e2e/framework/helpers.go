// Copyright (c) 2017 Chef Software Inc. and/or applicable contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package framework

import (
	"time"

	crv1 "github.com/kinvolk/habitat-operator/pkg/habitat/apis/cr/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
)

// NewStandaloneHabitat returns a new Standalone Habitat.
func NewStandaloneHabitat(habitatName, group, image string) *crv1.Habitat {
	return &crv1.Habitat{
		ObjectMeta: metav1.ObjectMeta{
			Name: habitatName,
		},
		Spec: crv1.HabitatSpec{
			Image: image,
			Count: 1,
			Service: crv1.Service{
				Group:    group,
				Topology: crv1.TopologyStandalone,
			},
		},
	}
}

// AddConfigToHabitat adds a ConfigSecretName field to the Habitat.
func AddConfigToHabitat(habitat *crv1.Habitat) {
	habitat.Spec.Service.ConfigSecretName = habitat.ObjectMeta.Name
}

// AddBindToHabitat appends bind fields to the Habitat.
func AddBindToHabitat(habitat *crv1.Habitat, bindName, bindService string) {
	habitat.Spec.Service.Bind = append(habitat.Spec.Service.Bind, crv1.Bind{
		Name:    bindName,
		Service: bindService,
		Group:   habitat.Spec.Service.Group,
	})
}

// CreateHabitat creates a Habitat.
func (f *Framework) CreateHabitat(habitat *crv1.Habitat) error {
	return f.Client.Post().
		Namespace(TestNs).
		Resource(crv1.HabitatResourcePlural).
		Body(habitat).
		Do().
		Error()
}

// WaitForResources waits until numPods are in the "Running" state.
// We wait for pods, because those take the longest to create.
// Waiting for anything else would be already testing.
func (f *Framework) WaitForResources(habitatName string, numPods int) error {
	return wait.Poll(2*time.Second, 5*time.Minute, func() (bool, error) {
		fs := fields.SelectorFromSet(fields.Set{
			"status.phase": "Running",
		})

		ls := labels.SelectorFromSet(labels.Set{
			crv1.HabitatNameLabel: habitatName,
		})

		pods, err := f.KubeClient.CoreV1().Pods(TestNs).List(metav1.ListOptions{FieldSelector: fs.String(), LabelSelector: ls.String()})
		if err != nil {
			return false, err
		}

		if len(pods.Items) != numPods {
			return false, nil
		}

		return true, nil
	})
}

func (f *Framework) WaitForEndpoints(habitatName string) error {
	return wait.Poll(time.Second, time.Minute*5, func() (bool, error) {
		ep, err := f.KubeClient.CoreV1().Endpoints(TestNs).Get(habitatName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if len(ep.Subsets) == 0 && len(ep.Subsets[0].Addresses) == 0 {
			return false, nil
		}

		return true, nil
	})
}

// DeleteHabitat deletes a Habitat as a user would.
func (f *Framework) DeleteHabitat(habitatName string) error {
	return f.Client.Delete().
		Namespace(TestNs).
		Resource(crv1.HabitatResourcePlural).
		Name(habitatName).
		Do().
		Error()
}

func (f *Framework) DeleteService(service string) error {
	return f.KubeClient.CoreV1().Services(TestNs).Delete(service, &metav1.DeleteOptions{})
}
