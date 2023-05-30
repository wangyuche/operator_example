package controllers

import (
	examplev1 "github.com/wangyuche/operator_example/example1/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *ExampleReconciler) pvcpg_data(m *examplev1.Example) *corev1.PersistentVolumeClaim {
	var SSD string = "standard"
	dep := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql-data",
			Namespace: m.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			StorageClassName: &SSD,
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("60Gi"),
				},
			},
		},
	}
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}

func (r *ExampleReconciler) svcpg_standalone(m *examplev1.Example) *corev1.Service {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql",
			Namespace: m.Namespace,
			Labels:    map[string]string{"app": "postgresql"},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "tcp",
					Protocol: "TCP",
					Port:     5432,
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: 5432,
					},
				},
			},
			Selector: map[string]string{"app": "postgresql"},
			Type:     corev1.ServiceType("ClusterIP"),
		},
	}
	ctrl.SetControllerReference(m, svc, r.Scheme)
	return svc
}

func (r *ExampleReconciler) stspg_standalone(m *examplev1.Example) *appsv1.StatefulSet {
	var Replicas int32 = 1
	dep := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql",
			Namespace: m.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "postgresql"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "postgresql"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: GetImages(m.Spec.Images),
							Name:  "postgresql",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 5432,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "POSTGRES_PASSWORD",
									Value: GetPassword(m.Spec.Password),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "postgresql-data",
									MountPath: "/var/lib/postgresql/data",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "postgresql-data",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "postgresql-data",
								},
							},
						},
					},
				},
			},
		},
	}
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}

func GetImages(images *string) string {
	if images == nil {
		return "postgres:13.11"
	}
	return *images
}

func GetPassword(password *string) string {
	if password == nil {
		return "sdffhgn"
	}
	return *password
}
