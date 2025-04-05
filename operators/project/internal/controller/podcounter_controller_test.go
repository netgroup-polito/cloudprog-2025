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
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	countersv1alpha1 "cloudprog.polito.it/project/api/v1alpha1"
)

const (
	PodReplicas = 5
)

var _ = Describe("PodCounter Controller", func() {
	Context("When reconciling a resource (testenv)", func() {
		const resourceName = "test-resource"

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // TODO(user):Modify as needed
		}
		podcounter := &countersv1alpha1.PodCounter{}

		BeforeEach(func() {
			err := k8sClient.Get(testContext, typeNamespacedName, podcounter)
			if err != nil && errors.IsNotFound(err) {
				resource := &countersv1alpha1.PodCounter{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
				}
				Expect(k8sClient.Create(testContext, resource)).To(Succeed())
			}

			for i := 0; i < PodReplicas; i++ {
				pod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      fmt.Sprintf("test-pod-%d", i),
						Namespace: "default",
						Labels: map[string]string{
							"app": "test-pod",
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "test-container",
								Image: "nginx",
							},
						},
					},
				}
				err := k8sClient.Create(testContext, pod)
				Expect(err).NotTo(HaveOccurred())
			}

			Eventually(func() bool {
				pods := &corev1.PodList{}
				err := k8sClient.List(testContext, pods, client.InNamespace("default"))
				Expect(err).NotTo(HaveOccurred())
				return len(pods.Items) == PodReplicas
			}, time.Minute, time.Second).Should(BeTrue())
		})

		AfterEach(func() {
			resource := &countersv1alpha1.PodCounter{}
			err := k8sClient.Get(testContext, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			Expect(k8sClient.Delete(testContext, resource)).To(Succeed())

			Eventually(func() bool {
				err := k8sClient.Get(testContext, typeNamespacedName, resource)
				return errors.IsNotFound(err)
			}).Should(BeTrue())
		})

		It("should successfully update the resource", func() {
			By("counting the pod instances")
			Eventually(func() bool {
				podcounter := &countersv1alpha1.PodCounter{}
				err := k8sClient.Get(testContext, typeNamespacedName, podcounter)
				Expect(err).NotTo(HaveOccurred())
				return podcounter.Status.Count == PodReplicas
			}, time.Minute, time.Second).Should(BeTrue())
		})
	})

	Context("When reconciling a resource (fake client)", func() {
		BeforeEach(func() {
			podList := forgePodList("default", PodReplicas)
			k8sClientFake = k8sClientFakeBuilder.
				WithObjects(
					&countersv1alpha1.PodCounter{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-resource",
							Namespace: "default",
						},
					},
				).
				WithObjects(podList...).
				WithStatusSubresource(&countersv1alpha1.PodCounter{}).
				Build()

			podcounter := &countersv1alpha1.PodCounter{}
			err := k8sClientFake.Get(testContext, types.NamespacedName{
				Name:      "test-resource",
				Namespace: "default",
			}, podcounter)
			Expect(err).NotTo(HaveOccurred())

			controllerReconciler := &PodCounterReconciler{
				Client: k8sClientFake,
				Scheme: k8sClientFake.Scheme(),
			}

			_, err = controllerReconciler.Reconcile(testContext, controllerruntime.Request{
				NamespacedName: types.NamespacedName{
					Name:      "test-resource",
					Namespace: "default",
				},
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should successfully update the resource", func() {

			podcounter := &countersv1alpha1.PodCounter{}
			err := k8sClientFake.Get(testContext, types.NamespacedName{
				Name:      "test-resource",
				Namespace: "default",
			}, podcounter)
			Expect(err).NotTo(HaveOccurred())

			Expect(podcounter.Status.Count).To(Equal(PodReplicas))
		})

	})
})

func forgePodList(ns string, replicas int) []client.Object {
	podList := make([]client.Object, replicas)
	for i := range replicas {
		podList[i] = &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("test-pod-%d", i),
				Namespace: ns,
				Labels: map[string]string{
					"app": "test-pod",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "test-container",
						Image: "nginx",
					},
				},
			},
		}
	}
	return podList
}
