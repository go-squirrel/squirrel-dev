package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PersistentVolume struct definition
type PersistentVolume struct {
	APIVersion string               `yaml:"apiVersion"`
	Kind       string               `yaml:"kind"`
	Metadata   Metadata             `yaml:"metadata"`
	Spec       PersistentVolumeSpec `yaml:"spec"`
}

type Metadata struct {
	Name   string            `yaml:"name"`
	Labels map[string]string `yaml:"labels,omitempty"`
}

// PersistentVolumeSpec defines PV specification
type PersistentVolumeSpec struct {
	Capacity                      Capacity              `yaml:"capacity"`
	AccessModes                   []string              `yaml:"accessModes"`
	StorageClassName              string                `yaml:"storageClassName,omitempty"`
	HostPath                      *HostPathVolumeSource `yaml:"hostPath,omitempty"`
	PersistentVolumeReclaimPolicy string                `yaml:"persistentVolumeReclaimPolicy,omitempty"`
}

// HostPathVolumeSource defines hostPath volume source
type HostPathVolumeSource struct {
	Path string `yaml:"path"`
}

type Capacity struct {
	Storage string `yaml:"storage"`
}

// size: 10Gi
func (c *Client) GeneratePV(name string,
	size string,
	storageClassName string,
	path string,
	labelKey, labelValue string,
	mode string) (pv *PersistentVolume, err error) {

	pv = &PersistentVolume{
		APIVersion: "v1",
		Kind:       "PersistentVolume",
		Metadata: Metadata{
			Name: name,
			Labels: map[string]string{
				labelKey: labelValue,
			},
		},
	}

	// Set storage capacity
	pv.Spec.Capacity.Storage = size

	// Set access mode
	pv.Spec.AccessModes = []string{mode}

	// Set storage class name
	pv.Spec.StorageClassName = storageClassName

	// Set reclaim policy
	pv.Spec.PersistentVolumeReclaimPolicy = "Retain"

	// Set storage type (using hostPath as an example)
	pv.Spec.HostPath = &HostPathVolumeSource{
		Path: path,
	}
	return pv, err
}

func (c *Client) CreatePersistentVolume(name string, size int64, storageClassName string, path string, labelKey, labelValue string, mode string) (err error) {
	createMode := corev1.ReadWriteOnce
	if mode == "ReadWriteMany" {
		createMode = corev1.ReadWriteMany
	}
	pv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.PersistentVolumeSpec{
			Capacity: corev1.ResourceList{
				corev1.ResourceStorage: *resource.NewQuantity(size*1024*1024*1024, resource.BinarySI), // 5Gi
			},
			AccessModes: []corev1.PersistentVolumeAccessMode{
				createMode,
			},
			PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimDelete,
			StorageClassName:              storageClassName,
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: path,
					Type: func() *corev1.HostPathType {
						pathType := corev1.HostPathDirectoryOrCreate
						return &pathType
					}(),
				},
			},
			NodeAffinity: &corev1.VolumeNodeAffinity{
				Required: &corev1.NodeSelector{
					NodeSelectorTerms: []corev1.NodeSelectorTerm{
						{
							MatchExpressions: []corev1.NodeSelectorRequirement{
								{
									Key:      labelKey,
									Operator: corev1.NodeSelectorOpIn,
									Values:   []string{labelValue},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err = c.ClientSet.CoreV1().PersistentVolumes().Create(context.TODO(), pv, metav1.CreateOptions{})

	return err
}
