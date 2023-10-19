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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
)

// getVolumes - service volumes
func getVolumes(name string) []corev1.Volume {

	return []corev1.Volume{
		{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: ptr.To(int32(0440)),
					SecretName:  name + "-config",
				},
			},
		},
		{
			Name: "apache-config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: ptr.To(int32(0440)),
					SecretName:  name + "-apache-config",
				},
			},
		},
		{
			Name: "fernet-keys",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  ServiceName,
					DefaultMode: ptr.To(int32(0440)),
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
					SecretName:  ServiceName,
					DefaultMode: ptr.To(int32(0440)),
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

}

// getVolumeMounts - general VolumeMounts
func getVolumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			MountPath: "/etc/keystone/fernet-keys",
			ReadOnly:  true,
			Name:      "fernet-keys",
		},
		{
			MountPath: "/etc/keystone/credential-keys",
			ReadOnly:  true,
			Name:      "credential-keys",
		},
		{
			MountPath: "/etc/keystone.conf.d",
			ReadOnly:  true,
			Name:      "config",
		},
		{
			MountPath: "/etc/httpd/conf/httpd.conf",
			ReadOnly:  true,
			Name:      "apache-config",
			SubPath:   "httpd.conf",
		},
		{
			MountPath: "/etc/httpd/conf.d/keystone-api.conf",
			ReadOnly:  true,
			Name:      "apache-config",
			SubPath:   "keystone-api.conf",
		},
	}
}
