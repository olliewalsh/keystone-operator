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

// BootstrapJob func
func BootstrapJob(
	instance *keystonev1.KeystoneAPI,
	labels map[string]string,
	annotations map[string]string,
	endpoints map[string]string,
) *batchv1.Job {

	entrypoint := []string{
		"dumb-init",
		"--single-child",
		"--",
	}
	command := []string{
		"/usr/bin/keystone-manage",
		"bootstrap",
	}

	envVars := map[string]env.Setter{}
	envVars["OS_BOOTSTRAP_USERNAME"] = env.SetValue(instance.Spec.AdminUser)
	envVars["OS_BOOTSTRAP_PROJECT_NAME"] = env.SetValue(instance.Spec.AdminProject)
	envVars["OS_BOOTSTRAP_SERVICE_NAME"] = env.SetValue(ServiceName)
	envVars["OS_BOOTSTRAP_REGION_ID"] = env.SetValue(instance.Spec.Region)

	if _, ok := endpoints["admin"]; ok {
		envVars["OS_BOOTSTRAP_ADMIN_URL"] = env.SetValue(endpoints["admin"])
	}
	if _, ok := endpoints["internal"]; ok {
		envVars["OS_BOOTSTRAP_INTERNAL_URL"] = env.SetValue(endpoints["internal"])
	}
	if _, ok := endpoints["public"]; ok {
		envVars["OS_BOOTSTRAP_PUBLIC_URL"] = env.SetValue(endpoints["public"])
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceName + "-bootstrap",
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
							Name:    ServiceName + "-bootstrap",
							Image:   instance.Spec.ContainerImage,
							Command: entrypoint,
							Args:    command,
							Env: env.MergeEnvs(
								[]corev1.EnvVar{
									{
										Name: "OS_BOOTSTRAP_PASSWORD",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: instance.Spec.Secret,
												},
												Key: "AdminPassword",
											},
										},
									},
								},
								envVars,
							),
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
