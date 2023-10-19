/*

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

package keystone

import (
	keystonev1 "github.com/openstack-k8s-operators/keystone-operator/api/v1beta1"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

// CronJob func
func CronJob(
	instance *keystonev1.KeystoneAPI,
	labels map[string]string,
	annotations map[string]string,
) *batchv1.CronJob {

	entrypoint := []string{
		"dumb-init",
		"--single-child",
		"--",
	}
	command := []string{
		"/usr/bin/keystone-manage",
		"trust_flush",
	}

	cronjob := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceName + "-cron",
			Namespace: instance.Namespace,
		},
		Spec: batchv1.CronJobSpec{
			Schedule:          instance.Spec.TrustFlushSchedule,
			Suspend:           &instance.Spec.TrustFlushSuspend,
			ConcurrencyPolicy: batchv1.ForbidConcurrent,
			JobTemplate: batchv1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: annotations,
					Labels:      labels,
				},
				Spec: batchv1.JobSpec{
					Parallelism: ptr.To(int32(1)),
					Completions: ptr.To(int32(1)),
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:         ServiceName + "-cron",
									Image:        instance.Spec.ContainerImage,
									Command:      entrypoint,
									Args:         command,
									VolumeMounts: getVolumeMounts(),
									SecurityContext: &corev1.SecurityContext{
										AllowPrivilegeEscalation: ptr.To(false),
										Capabilities: &corev1.Capabilities{
											Drop: []corev1.Capability{"ALL"},
										},
									},
								},
							},
							RestartPolicy: corev1.RestartPolicyNever,
							SecurityContext: &corev1.PodSecurityContext{
								RunAsNonRoot: ptr.To(true),
								SeccompProfile: &corev1.SeccompProfile{
									Type: corev1.SeccompProfileTypeRuntimeDefault,
								},
							},
							ServiceAccountName: instance.RbacResourceName(),
							Volumes:            getVolumes(instance.Name),
						},
					},
				},
			},
		},
	}
	if instance.Spec.NodeSelector != nil && len(instance.Spec.NodeSelector) > 0 {
		cronjob.Spec.JobTemplate.Spec.Template.Spec.NodeSelector = instance.Spec.NodeSelector
	}

	return cronjob
}
