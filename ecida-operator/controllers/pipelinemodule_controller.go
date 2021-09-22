/*
Copyright 2021.

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
	"fmt"
	"reflect"

	ecidav1alpha "ecida-operator/api/v1alpha1"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PipelineModuleReconciler reconciles a PipelineModule object
type PipelineModuleReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ecida.researchable.nl,resources=pipelinemodules,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ecida.researchable.nl,resources=pipelinemodules/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ecida.researchable.nl,resources=pipelinemodules/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PipelineModule object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *PipelineModuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("pipelinemodule", req.NamespacedName)

	// your logic here
	pipelineModule := &ecidav1alpha.PipelineModule{}

	// Fetch an instance of the PipelineModule CRD that the operator is watching
	err := r.Get(ctx, req.NamespacedName, pipelineModule)

	if err != nil {
		if errors.IsNotFound(err) {
			// There is no instance of the CRD, don't reconcile further
			log.Info("PipelineModule custom resource not found, ignoring reconcile")
			return ctrl.Result{}, nil
		}

		// Some other error
		log.Error(err, "Failed to get instance of PipelineModule CRD")
		return ctrl.Result{}, err
	}

	// There is an instance of the CRD
	// This operator manages a deployment instance that is owned by the custom resource.
	// Check if it exist and fetch it
	dep := &appsv1.Deployment{}
	err = r.Get(ctx, r.pipelineModuleDepName(req.NamespacedName), dep)

	if err != nil {
		if errors.IsNotFound(err) {
			// The deployment does not yet exist, create it
			depSpec := r.deploymentSpecPM(pipelineModule, req.NamespacedName)
			log.Info("Creating new deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", depSpec.Name)
			err = r.Create(ctx, depSpec)
			if err != nil {
				log.Error(err, "Failed to create new deployment", "Deployment.Namespace", depSpec.Namespace, "Deployment.Name", depSpec.Name)
				return ctrl.Result{}, err
			}

			log.Info("Created new deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{Requeue: true}, nil
		}
	}

	// On the deployment, check if the Image, Port and Command are still correct.
	deployAttrsChanged := false
	moduleContainerSpec := &dep.Spec.Template.Spec.Containers[0]

	if moduleContainerSpec.Image != pipelineModule.Spec.Image {
		log.Info("Deployment Image has changed", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		moduleContainerSpec.Image = pipelineModule.Spec.Image
		deployAttrsChanged = true
	}

	if moduleContainerSpec.Ports[0].ContainerPort != pipelineModule.Spec.Port {
		log.Info("Deployment Port has changed", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		moduleContainerSpec.Ports[0].ContainerPort = pipelineModule.Spec.Port
		deployAttrsChanged = true
	}

	if !reflect.DeepEqual(moduleContainerSpec.Command, pipelineModule.Spec.Command) {
		log.Info("Deployment Command has changed", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		moduleContainerSpec.Command = pipelineModule.Spec.Command
		deployAttrsChanged = true
	}

	if deployAttrsChanged {
		log.Info("Updating deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Update(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to update deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}

		return ctrl.Result{Requeue: true}, nil
	}

	// Get the pods of the deployment to update the status
	podList := &corev1.PodList{}
	podMatcher := []client.ListOption{
		client.InNamespace(pipelineModule.Namespace),
		client.MatchingLabels(r.pipelineModuleLabels(pipelineModule.Name)),
	}

	if err = r.List(ctx, podList, podMatcher...); err != nil {
		log.Error(err, "Failed to get PodList", "PipelineModule.Namespace", pipelineModule.Namespace, "PipelineModule.Name", pipelineModule.Name)
		return ctrl.Result{}, err
	}

	podNames := r.podNames(podList.Items)

	if !reflect.DeepEqual(podNames, pipelineModule.Status.Nodes) {
		pipelineModule.Status.Nodes = podNames

		if err = r.Update(ctx, pipelineModule); err != nil {
			log.Error(err, "Failed to update pipelineModule status", "PipelineModule.Namespace", pipelineModule.Namespace, "PipelineModule.Name", pipelineModule.Name)
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// Defines the name of the deployment name that is owned by the custom resource
func (r *PipelineModuleReconciler) pipelineModuleDepName(crname types.NamespacedName) types.NamespacedName {
	return types.NamespacedName{
		Name:      fmt.Sprintf("%s-deployment", crname.Name),
		Namespace: crname.Namespace,
	}
}

func (r *PipelineModuleReconciler) podNames(pods []corev1.Pod) []string {
	var names []string

	for _, pod := range pods {
		names = append(names, pod.Name)
	}

	return names
}

func (r *PipelineModuleReconciler) pipelineModuleLabels(name string) map[string]string {
	return map[string]string{
		"ecida_type":  "module",
		"module_name": name,
	}
}

func (r *PipelineModuleReconciler) deploymentSpecPM(pipelineModule *ecidav1alpha.PipelineModule, moduleName types.NamespacedName) *appsv1.Deployment {
	fullName := r.pipelineModuleDepName(moduleName)
	image := pipelineModule.Spec.Image
	port := pipelineModule.Spec.Port
	replicas := int32(1) // Take care of listing the pods when changing this

	labels := r.pipelineModuleLabels(pipelineModule.Name)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: fullName.Name, Namespace: fullName.Namespace},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: image,
						Name:  pipelineModule.Name,
						Ports: []corev1.ContainerPort{{
							ContainerPort: int32(port),
						}},
					}},
				},
			},
		},
	}

	// Make the CR the owner of this deployment
	ctrl.SetControllerReference(pipelineModule, deployment, r.Scheme)

	return deployment
}

// SetupWithManager sets up the controller with the Manager.
func (r *PipelineModuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ecidav1alpha.PipelineModule{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
