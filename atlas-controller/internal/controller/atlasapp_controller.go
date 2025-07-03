/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	atlasv1 "atlas-controller/api/v1"
)

// AtlasAppReconciler reconciles a AtlasApp object
type AtlasAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=atlas.io,resources=atlasapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=atlas.io,resources=atlasapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=atlas.io,resources=atlasapps/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *AtlasAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// 1. Fetch the AtlasApp instance
	var atlasApp atlasv1.AtlasApp
	if err := r.Get(ctx, req.NamespacedName, &atlasApp); err != nil {
		if errors.IsNotFound(err) {
			log.Info("AtlasApp resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get AtlasApp")
		return ctrl.Result{}, err
	}

	log.Info("Reconciling AtlasApp", "environment", atlasApp.Spec.Environment, "version", atlasApp.Spec.Version)

	// 2. Check if approval is required for prod deployments
	if atlasApp.Spec.Environment == "prod" && atlasApp.Spec.RequireApproval && !atlasApp.Status.ApprovalRequired {
		return r.handleApprovalRequired(ctx, &atlasApp)
	}

	// 3. Create or update the deployment
	if err := r.reconcileDeployment(ctx, &atlasApp); err != nil {
		return r.updateStatus(ctx, &atlasApp, "Failed", false, err.Error())
	}

	// 4. Create or update the service
	if err := r.reconcileService(ctx, &atlasApp); err != nil {
		return r.updateStatus(ctx, &atlasApp, "Failed", false, err.Error())
	}

	// 5. Check deployment status
	ready, err := r.checkDeploymentStatus(ctx, &atlasApp)
	if err != nil {
		return r.updateStatus(ctx, &atlasApp, "Failed", false, err.Error())
	}

	if !ready {
		return r.updateStatus(ctx, &atlasApp, "Deploying", false, "Waiting for deployment to be ready")
	}

	// 6. Perform health check
	if atlasApp.Spec.HealthCheckPath != "" {
		healthy, err := r.performHealthCheck(ctx, &atlasApp)
		if err != nil {
			return r.updateStatus(ctx, &atlasApp, "Failed", false, fmt.Sprintf("Health check failed: %v", err))
		}
		if !healthy {
			return r.updateStatus(ctx, &atlasApp, "Unhealthy", false, "Health check failed")
		}
	}

	// 7. Update status to Ready
	if result, err := r.updateStatus(ctx, &atlasApp, "Ready", true, "Application is healthy and ready"); err != nil {
		return result, err
	}

	// 8. Handle auto-promotion
	if atlasApp.Spec.AutoPromote && atlasApp.Spec.NextEnvironment != "" {
		return r.handleAutoPromotion(ctx, &atlasApp)
	}

	return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
}

// reconcileDeployment creates or updates the deployment
func (r *AtlasAppReconciler) reconcileDeployment(ctx context.Context, atlasApp *atlasv1.AtlasApp) error {
	log := log.FromContext(ctx)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "atlas",
			Namespace: atlasApp.Namespace,
			Labels: map[string]string{
				"app":                          "atlas",
				"atlas.io/environment":         atlasApp.Spec.Environment,
				"atlas.io/version":             atlasApp.Spec.Version,
				"atlas.io/managed-by":          "atlas-controller",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &atlasApp.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "atlas",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":                  "atlas",
						"atlas.io/environment": atlasApp.Spec.Environment,
						"atlas.io/version":     atlasApp.Spec.Version,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "atlas",
							Image: fmt.Sprintf("nginx:%s", atlasApp.Spec.Version),
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "MIGRATION_ID",
									Value: fmt.Sprintf("%d", atlasApp.Spec.MigrationId),
								},
								{
									Name:  "ENVIRONMENT",
									Value: atlasApp.Spec.Environment,
								},
							},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt(80),
									},
								},
								InitialDelaySeconds: 30,
								PeriodSeconds:       10,
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/",
										Port: intstr.FromInt(80),
									},
								},
								InitialDelaySeconds: 5,
								PeriodSeconds:       5,
							},
						},
					},
				},
			},
		},
	}

	// Set AtlasApp as the owner of the Deployment
	if err := ctrl.SetControllerReference(atlasApp, deployment, r.Scheme); err != nil {
		return err
	}

	// Check if deployment already exists
	found := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.Create(ctx, deployment)
		if err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	// Update existing deployment if needed
	if found.Spec.Template.Spec.Containers[0].Image != deployment.Spec.Template.Spec.Containers[0].Image ||
		found.Spec.Template.Spec.Containers[0].Env[0].Value != deployment.Spec.Template.Spec.Containers[0].Env[0].Value {
		log.Info("Updating Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		found.Spec = deployment.Spec
		err = r.Update(ctx, found)
		if err != nil {
			return err
		}
	}

	return nil
}

// reconcileService creates or updates the service
func (r *AtlasAppReconciler) reconcileService(ctx context.Context, atlasApp *atlasv1.AtlasApp) error {
	log := log.FromContext(ctx)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "atlas",
			Namespace: atlasApp.Namespace,
			Labels: map[string]string{
				"app":                  "atlas",
				"atlas.io/environment": atlasApp.Spec.Environment,
				"atlas.io/managed-by":  "atlas-controller",
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "atlas",
			},
			Ports: []corev1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(80),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}

	// Set AtlasApp as the owner of the Service
	if err := ctrl.SetControllerReference(atlasApp, service, r.Scheme); err != nil {
		return err
	}

	// Check if service already exists
	found := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.Create(ctx, service)
		return err
	} else if err != nil {
		return err
	}

	return nil
}

// checkDeploymentStatus checks if the deployment is ready
func (r *AtlasAppReconciler) checkDeploymentStatus(ctx context.Context, atlasApp *atlasv1.AtlasApp) (bool, error) {
	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: "atlas", Namespace: atlasApp.Namespace}, deployment)
	if err != nil {
		return false, err
	}

	// Update replica counts in status
	atlasApp.Status.ReadyReplicas = deployment.Status.ReadyReplicas
	atlasApp.Status.TotalReplicas = deployment.Status.Replicas

	// Check if all replicas are ready
	return deployment.Status.ReadyReplicas == *deployment.Spec.Replicas && deployment.Status.ReadyReplicas > 0, nil
}

// performHealthCheck performs application health check
func (r *AtlasAppReconciler) performHealthCheck(ctx context.Context, atlasApp *atlasv1.AtlasApp) (bool, error) {
	// TODO: Implement actual HTTP health check to the service
	// For now, we'll just simulate it
	log := log.FromContext(ctx)
	log.Info("Performing health check", "path", atlasApp.Spec.HealthCheckPath)
	
	// Simulate health check success
	return true, nil
}

// updateStatus updates the AtlasApp status with retry logic
func (r *AtlasAppReconciler) updateStatus(ctx context.Context, atlasApp *atlasv1.AtlasApp, phase string, ready bool, message string) (ctrl.Result, error) {
	// Use retry logic to handle conflicts
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Fetch the latest version
		latest := &atlasv1.AtlasApp{}
		if err := r.Get(ctx, client.ObjectKeyFromObject(atlasApp), latest); err != nil {
			return err
		}
		
		// Update status fields
		latest.Status.Phase = phase
		latest.Status.Ready = ready
		latest.Status.Message = message
		now := metav1.Now()
		latest.Status.LastUpdate = &now
		
		// Update the status
		return r.Status().Update(ctx, latest)
	})
	
	if retryErr != nil {
		return ctrl.Result{}, retryErr
	}

	if !ready {
		return ctrl.Result{RequeueAfter: time.Second * 30}, nil
	}

	return ctrl.Result{}, nil
}

// handleApprovalRequired sets the approval required status
func (r *AtlasAppReconciler) handleApprovalRequired(ctx context.Context, atlasApp *atlasv1.AtlasApp) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Production deployment requires approval", "app", atlasApp.Name)

	atlasApp.Status.ApprovalRequired = true
	atlasApp.Status.Phase = "PendingApproval"
	atlasApp.Status.Message = "Production deployment requires manual approval"

	if err := r.Status().Update(ctx, atlasApp); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
}

// handleAutoPromotion handles automatic promotion to next environment
func (r *AtlasAppReconciler) handleAutoPromotion(ctx context.Context, atlasApp *atlasv1.AtlasApp) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	
	// Check if promotion to prod requires approval
	if atlasApp.Spec.NextEnvironment == "prod" {
		log.Info("Promotion to prod requires approval", "current", atlasApp.Spec.Environment)
		atlasApp.Status.PromotionPending = true
		atlasApp.Status.Message = "Promotion to production requires manual approval"
		if err := r.Status().Update(ctx, atlasApp); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Create AtlasApp in next environment
	nextApp := &atlasv1.AtlasApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("atlas-%s", atlasApp.Spec.NextEnvironment),
			Namespace: atlasApp.Spec.NextEnvironment,
		},
		Spec: atlasv1.AtlasAppSpec{
			Environment:     atlasApp.Spec.NextEnvironment,
			Version:         atlasApp.Spec.Version,
			MigrationId:     atlasApp.Spec.MigrationId,
			Replicas:        atlasApp.Spec.Replicas,
			AutoPromote:     atlasApp.Spec.Environment != "stage", // Only auto-promote from dev to stage
			NextEnvironment: getNextEnvironment(atlasApp.Spec.NextEnvironment),
			RequireApproval: atlasApp.Spec.NextEnvironment == "prod",
			HealthCheckPath: atlasApp.Spec.HealthCheckPath,
		},
	}

	// Check if the next environment AtlasApp already exists
	existingApp := &atlasv1.AtlasApp{}
	err := r.Get(ctx, types.NamespacedName{Name: nextApp.Name, Namespace: nextApp.Namespace}, existingApp)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating AtlasApp in next environment", "environment", nextApp.Spec.Environment, "version", nextApp.Spec.Version)
		if err := r.Create(ctx, nextApp); err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		// Update existing app if version or migration changed
		if existingApp.Spec.Version != nextApp.Spec.Version || existingApp.Spec.MigrationId != nextApp.Spec.MigrationId {
			log.Info("Updating AtlasApp in next environment", "environment", nextApp.Spec.Environment, "version", nextApp.Spec.Version)
			existingApp.Spec.Version = nextApp.Spec.Version
			existingApp.Spec.MigrationId = nextApp.Spec.MigrationId
			if err := r.Update(ctx, existingApp); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// getNextEnvironment returns the next environment in the promotion chain
func getNextEnvironment(currentEnv string) string {
	switch currentEnv {
	case "dev":
		return "stage"
	case "stage":
		return "prod"
	default:
		return ""
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *AtlasAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&atlasv1.AtlasApp{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
