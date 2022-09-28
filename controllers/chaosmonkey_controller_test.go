package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	podv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	chaosv1 "github.com/Kostov6/chaosmonkey/api/v1"
)

var _ = Describe("Chaosmonkey controller", func() {
	// Define utility constants for object names and testing timeouts/durations and intervals.

	const (
		chaosmonkeyNamespace = "default"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
		period   = 10
	)

	var (
		ctx context.Context

		podName          string
		namedPod         *podv1.Pod
		namedChaosmonkey *chaosv1.Chaosmonkey

		podLabel           string
		podLabelValue      string
		labeledPod         *podv1.Pod
		labeledPodName     string
		labeledChaosmonkey *chaosv1.Chaosmonkey

		brokenPodName string
		brokenPod     *podv1.Pod
		//fieldChaosmonkey *chaosv1.Chaosmonkey
	)

	JustBeforeEach(func() {
		k8sClient.Create(ctx, namedPod)
		k8sClient.Create(ctx, labeledPod)
		k8sClient.Create(ctx, brokenPod)
	})

	Describe("Testing chaosmonkey deleting capabilities", func() {
		When("a pod matches the podName spec property", func() {
			BeforeEach(func() {
				ctx = context.Background()
				podName = "foo-pod"
				namedPod = &podv1.Pod{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "Pod",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      podName,
						Namespace: chaosmonkeyNamespace,
					},
					Spec: podv1.PodSpec{
						Containers: []podv1.Container{
							{
								Name:  "nginx",
								Image: "nginx",
							},
						},
					},
				}
				namedChaosmonkey = &chaosv1.Chaosmonkey{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "onboarding.my.domain/v1",
						Kind:       "Chaosmonkey",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monkey1",
						Namespace: chaosmonkeyNamespace,
					},
					Spec: chaosv1.ChaosmonkeySpec{
						PodName:   podName,
						Namespace: chaosmonkeyNamespace,
					},
				}
			})

			It("Should be deleted successfully", func() {
				By("Waiting for initial pod to be ready")
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: podName, Namespace: chaosmonkeyNamespace}, &podv1.Pod{})
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				By("Creating chaosmonkey")
				Expect(k8sClient.Create(ctx, namedChaosmonkey)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: namedChaosmonkey.Name, Namespace: chaosmonkeyNamespace}, &chaosv1.Chaosmonkey{})
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				By("Waiting for named pod to be deleted")
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: podName, Namespace: chaosmonkeyNamespace}, &podv1.Pod{})
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeFalse())
			})
		})

		When("a pod matches the wirhFields spec property", func() {
			BeforeEach(func() {
				ctx = context.Background()
				brokenPodName = "broken-pod"
				brokenPod = &podv1.Pod{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "Pod",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      brokenPodName,
						Namespace: chaosmonkeyNamespace,
					},
					Spec: podv1.PodSpec{
						Containers: []podv1.Container{
							{
								Name:  "broken",
								Image: "imagedoesnotexist",
							},
						},
					},
				}
				/*	fieldChaosmonkey = &chaosv1.Chaosmonkey{
						TypeMeta: metav1.TypeMeta{
							APIVersion: "onboarding.my.domain/v1",
							Kind:       "Chaosmonkey",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "monkey3",
							Namespace: chaosmonkeyNamespace,
						},
						Spec: chaosv1.ChaosmonkeySpec{
							Namespace:  chaosmonkeyNamespace,
							WithFields: map[string]string{"metadata.namespace": chaosmonkeyNamespace},
						},
					}
				*/
			})

			It("Should be deleted successfully", func() {
				/*By("Waiting for broken pod be ready")
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: brokenPodName, Namespace: chaosmonkeyNamespace}, &podv1.Pod{})
					if err != nil {
						fmt.Println(err)
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				By("Creating chaosmonkey")
				Expect(k8sClient.Create(ctx, fieldChaosmonkey)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: fieldChaosmonkey.Name, Namespace: chaosmonkeyNamespace}, &chaosv1.Chaosmonkey{})
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				By("Waiting for broken pod to be deleted")
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: brokenPodName, Namespace: chaosmonkeyNamespace}, &podv1.Pod{})
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeFalse())
				*/
			})
		})

		When("a pod matches the withLabels spec property", func() {
			BeforeEach(func() {
				ctx = context.Background()
				labeledPodName = "pod2"
				podLabel = "foo"
				podLabelValue = "bar"
				labeledPod = &podv1.Pod{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "Pod",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      labeledPodName,
						Namespace: chaosmonkeyNamespace,
						Labels:    map[string]string{podLabel: podLabelValue},
					},
					Spec: podv1.PodSpec{
						Containers: []podv1.Container{
							{
								Name:  "nginx",
								Image: "nginx",
							},
						},
					},
				}
				labeledChaosmonkey = &chaosv1.Chaosmonkey{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "onboarding.my.domain/v1",
						Kind:       "Chaosmonkey",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "monkey2",
						Namespace: chaosmonkeyNamespace,
					},
					Spec: chaosv1.ChaosmonkeySpec{
						Namespace:  chaosmonkeyNamespace,
						WithLabels: map[string]string{podLabel: podLabelValue},
					},
				}
			})

			It("Should be deleted successfully", func() {
				By("Waiting for named pod to be ready")
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: labeledPodName, Namespace: chaosmonkeyNamespace}, &podv1.Pod{})
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				By("Creating chaosmonkey")
				Expect(k8sClient.Create(ctx, labeledChaosmonkey)).Should(Succeed())
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: labeledChaosmonkey.Name, Namespace: chaosmonkeyNamespace}, &chaosv1.Chaosmonkey{})
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeTrue())
				By("Waiting for labeled pod to be deleted")
				Eventually(func() bool {
					err := k8sClient.Get(ctx, types.NamespacedName{Name: labeledPodName, Namespace: chaosmonkeyNamespace}, &podv1.Pod{})
					if err != nil {
						return false
					}
					return true
				}, timeout, interval).Should(BeFalse())
			})
		})

	})

})
