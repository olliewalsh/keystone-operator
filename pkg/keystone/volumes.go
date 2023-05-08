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
	corev1 "k8s.io/api/core/v1"
)

// getVolumes - service volumes
func getVolumes(instance *keystonev1.KeystoneAPI) []corev1.Volume {
	var scriptsVolumeDefaultMode int32 = 0755
	var config0644AccessMode int32 = 0644
	var config0640AccessMode int32 = 0640
	var config0600AccessMode int32 = 0600

	var volumes = []corev1.Volume{
		{
			Name: "scripts",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					DefaultMode: &scriptsVolumeDefaultMode,
					LocalObjectReference: corev1.LocalObjectReference{
						Name: instance.Name + "-scripts",
					},
				},
			},
		},
		{
			Name: "config-data",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					DefaultMode: &config0640AccessMode,
					LocalObjectReference: corev1.LocalObjectReference{
						Name: instance.Name + "-config-data",
					},
				},
			},
		},
		{
			Name: "config-data-merged",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{Medium: ""},
			},
		},
		{
			Name: "fernet-keys",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: ServiceName,
					Items: []corev1.KeyToPath{
						{
							Key:  "FernetKeys0",
							Path: "0",
						},
						{
							Key:  "FernetKeys1",
							Path: "1",
						},
					},
				},
			},
		},
		{
			Name: "credential-keys",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: ServiceName,
					Items: []corev1.KeyToPath{
						{
							Key:  "CredentialKeys0",
							Path: "0",
						},
						{
							Key:  "CredentialKeys1",
							Path: "1",
						},
					},
				},
			},
		},
	}

	if instance.Spec.TLS.SecretName != "" {
		volumes = append(
			volumes,
			corev1.Volume{
				Name: "tls-secret",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						DefaultMode: &config0600AccessMode,
						SecretName:  instance.Spec.TLS.SecretName,
					},
				},
			},
		)
	}

	if instance.Spec.TLS.CaSecretName != "" {
		volumes = append(
			volumes,
			corev1.Volume{
				Name: "tlsca-secret",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						DefaultMode: &config0644AccessMode,
						SecretName:  instance.Spec.TLS.CaSecretName,
					},
				},
			},
		)
	}

	return volumes
}

// getInitVolumeMounts - general init task VolumeMounts
func getInitVolumeMounts() []corev1.VolumeMount {
	var volumeMounts = []corev1.VolumeMount{
		{
			Name:      "scripts",
			MountPath: "/usr/local/bin/container-scripts",
			ReadOnly:  true,
		},
		{
			Name:      "config-data",
			MountPath: "/var/lib/config-data/default",
			ReadOnly:  true,
		},
		{
			Name:      "config-data-merged",
			MountPath: "/var/lib/config-data/merged",
			ReadOnly:  false,
		},
	}

	return volumeMounts
}

// getVolumeMounts - general VolumeMounts
func getVolumeMounts(instance *keystonev1.KeystoneAPI) []corev1.VolumeMount {
	var volumeMounts = []corev1.VolumeMount{
		{
			Name:      "scripts",
			MountPath: "/usr/local/bin/container-scripts",
			ReadOnly:  true,
		},
		{
			Name:      "config-data-merged",
			MountPath: "/var/lib/config-data/merged",
			ReadOnly:  false,
		},
		{
			MountPath: "/var/lib/fernet-keys",
			ReadOnly:  true,
			Name:      "fernet-keys",
		},
		{
			MountPath: "/var/lib/credential-keys",
			ReadOnly:  true,
			Name:      "credential-keys",
		},
	}

	if instance.Spec.TLS.SecretName != "" {
		volumeMounts = append(
			volumeMounts,
			corev1.VolumeMount{
				Name:      "tls-secret",
				MountPath: "/var/lib/tls-data",
				ReadOnly:  true,
			},
		)
	}

	if instance.Spec.TLS.CaSecretName != "" {
		volumeMounts = append(
			volumeMounts,
			corev1.VolumeMount{
				Name:      "tlsca-secret",
				MountPath: "/var/lib/tlsca-data",
				ReadOnly:  true,
			},
		)
	}

	return volumeMounts
}
