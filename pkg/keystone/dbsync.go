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

	"github.com/openstack-k8s-operators/lib-common/modules/common/env"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

// DbSyncJob func
func DbSyncJob(
	instance *keystonev1.KeystoneAPI,
	labels map[string]string,
	annotations map[string]string,
) *batchv1.Job {

	entrypoint := []string{
		"dumb-init",
		"--single-child",
		"--",
	}
	command := []string{
		"/usr/bin/keystone-manage",
		"db_sync",
	}
	if instance.Spec.Debug.DBSync {
		command = []string{
			"/bin/sleep",
			"infinity",
		}
	}

	envVars := map[string]env.Setter{}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceName + "-db-sync",
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: annotations,
				},
				Spec: corev1.PodSpec{
					RestartPolicy:      corev1.RestartPolicyOnFailure,
					ServiceAccountName: instance.RbacResourceName(),
					SecurityContext: &corev1.PodSecurityContext{
						RunAsNonRoot: ptr.To(true),
						SeccompProfile: &corev1.SeccompProfile{
							Type: corev1.SeccompProfileTypeRuntimeDefault,
						},
					},
					Volumes: getVolumes(instance.Name),
					Containers: []corev1.Container{
						{
							Name:         ServiceName + "-db-sync",
							Command:      entrypoint,
							Args:         command,
							Image:        instance.Spec.ContainerImage,
							Env:          env.MergeEnvs([]corev1.EnvVar{}, envVars),
							VolumeMounts: getVolumeMounts(),
							SecurityContext: &corev1.SecurityContext{
								AllowPrivilegeEscalation: ptr.To(false),
								Capabilities: &corev1.Capabilities{
									Drop: []corev1.Capability{"ALL"},
								},
							},
						},
					},
				},
			},
		},
	}

	return job
}
