/*
Copyright 2022.

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

package controllers

import (
	"context"
	v1 "k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	cachev1alpha1 "github.com/quan930/sample-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateDeploymentReconciler reconciles a CreateDeployment object
type CreateDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cache.lilqcn,resources=createdeployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.lilqcn,resources=createdeployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.lilqcn,resources=createdeployments/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CreateDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *CreateDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	createDeployment := &cachev1alpha1.CreateDeployment{}
	err := r.Get(ctx, req.NamespacedName, createDeployment)
	if err != nil {
		if errors.IsNotFound(err) {
			klog.Info(err, ", CreateDeployment resource not found. Ignoring since object must be deleted: ", createDeployment)
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		klog.Error(err, "Failed to get CreateDeployment")
		return ctrl.Result{}, err
	}
	klog.Info("CreateDeployment: ", createDeployment)

	found := &v1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: createDeployment.Spec.Name, Namespace: createDeployment.Spec.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.deploymentForDemo(createDeployment)

		klog.Info("Creating a new Deployment, ", "Deployment.Namespace:", dep.Namespace, "\tDeployment.Name:", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			klog.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		klog.Info("create deployment successful")
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		klog.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CreateDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.CreateDeployment{}).
		Owns(&v1.Deployment{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		Complete(r)
}

//deploymentForDemo
func (r *CreateDeploymentReconciler) deploymentForDemo(m *cachev1alpha1.CreateDeployment) *v1.Deployment {
	ls := labelsForMyBook(m.Name)
	replicas := int32(1)
	boolTure := true

	dep := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Spec.Name,
			Namespace: m.Spec.Namespace,
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion:         "cache.lilqcn/v1alpha1",
				Kind:               "CreateDeployment",
				Name:               m.Name,
				UID:                m.UID,
				BlockOwnerDeletion: &boolTure,
			}},
		},
		Spec: v1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: v12.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: v12.PodSpec{
					Containers: []v12.Container{{
						Image: "nginx:stable-alpine",
						Name:  "book-server",
						//Command: []string{"memcached", "-m=64", "-o", "modern", "-v"},
						Ports: []v12.ContainerPort{{
							ContainerPort: 8080,
							Name:          "nginx",
						}},
					}},
				},
			},
		},
	}
	// Set Memcached instance as the owner and controller
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}

func labelsForMyBook(name string) map[string]string {
	return map[string]string{"app": "sample-operator", "sample-operator": name}
}
