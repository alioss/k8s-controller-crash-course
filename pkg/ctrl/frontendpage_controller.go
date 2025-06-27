package ctrl

import (
	context "context"
	"reflect"
	"time"
	"fmt"
	"github.com/rs/zerolog/log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	frontendv1alpha1 "github.com/yourusername/k8s-controller-tutorial/pkg/apis/frontend/v1alpha1"
	"github.com/yourusername/k8s-controller-tutorial/pkg/port"
)

type FrontendPageReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	PortClient *port.PortClient
}

func buildConfigMap(page *frontendv1alpha1.FrontendPage) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      page.Name,
			Namespace: page.Namespace,
		},
		Data: map[string]string{
			"contents": page.Spec.Contents,
		},
	}
}

func buildDeployment(page *frontendv1alpha1.FrontendPage) *appsv1.Deployment {
	replicas := int32(page.Spec.Replicas)
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      page.Name,
			Namespace: page.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": page.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": page.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "frontend",
						Image: page.Spec.Image,
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "contents",
							MountPath: "/data",
						}},
					}},
					Volumes: []corev1.Volume{{
						Name: "contents",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: page.Name,
								},
							},
						},
					}},
				},
			},
		},
	}
}

func (r *FrontendPageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var page frontendv1alpha1.FrontendPage
	err := r.Get(ctx, req.NamespacedName, &page)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			// FrontendPage deleted: clean up resources
			log.Info().Msgf("FrontendPage deleted: %s %s", req.Name, req.Namespace)
			
			// Clean up in Port.io if configured
			if r.PortClient != nil {
				identifier := fmt.Sprintf("%s-%s", req.Name, req.Namespace)
				if err := r.PortClient.DeleteEntity("frontendpage", identifier); err != nil {
					log.Warn().Err(err).Msg("Failed to delete entity from Port.io")
				}
			}
			
			var cm corev1.ConfigMap
			cm.Name = req.Name
			cm.Namespace = req.Namespace
			_ = r.Delete(ctx, &cm) // ignore errors if not found
			var dep appsv1.Deployment
			dep.Name = req.Name
			dep.Namespace = req.Namespace
			_ = r.Delete(ctx, &dep)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// 1. Ensure ConfigMap exists and is up to date
	cm := buildConfigMap(&page)
	if err := ctrl.SetControllerReference(&page, cm, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	log.Info().Msgf("Reconciling ConfigMap for FrontendPage: %s %s",cm.Name, cm.Namespace)
	var existingCM corev1.ConfigMap
	cmErr := r.Get(ctx, req.NamespacedName, &existingCM)
	if cmErr != nil && errors.IsNotFound(cmErr) {
		if err := r.Create(ctx, cm); err != nil && !errors.IsAlreadyExists(err) {
			return ctrl.Result{}, err
		}
	} else if cmErr == nil && !reflect.DeepEqual(existingCM.Data, cm.Data) {
		existingCM.Data = cm.Data
		if err := r.Update(ctx, &existingCM); err != nil {
			return ctrl.Result{}, err
		}
	}

	// 2. Ensure Deployment exists and is up to date
	log.Info().Msgf("Reconciling Deployment: %s/%s", req.Namespace, req.Name)
	dep := buildDeployment(&page)
	if err := ctrl.SetControllerReference(&page, dep, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	log.Info().Msgf("Reconciling Deployment for FrontendPage: %s %s", dep.Name, dep.Namespace)
	var existingDep appsv1.Deployment
	depErr := r.Get(ctx, req.NamespacedName, &existingDep)
	if depErr != nil && errors.IsNotFound(depErr) {
		if err := r.Create(ctx, dep); err != nil && !errors.IsAlreadyExists(err) {
			return ctrl.Result{}, err
		}
	} else if depErr == nil {
		updated := false
		if *existingDep.Spec.Replicas != *dep.Spec.Replicas {
			existingDep.Spec.Replicas = dep.Spec.Replicas
			updated = true
		}
		if existingDep.Spec.Template.Spec.Containers[0].Image != dep.Spec.Template.Spec.Containers[0].Image {
			existingDep.Spec.Template.Spec.Containers[0].Image = dep.Spec.Template.Spec.Containers[0].Image
			updated = true
		}
		if updated {
			if err := r.Update(ctx, &existingDep); err != nil {
				if errors.IsConflict(err) {
					// Requeue to try again with the latest version
					return ctrl.Result{Requeue: true}, nil
				}
				return ctrl.Result{}, err
			}
		}
	}
	// 3. Sync with Port.io if configured
	if r.PortClient != nil {
		if err := r.syncToPort(ctx, &page); err != nil {
			log.Warn().Err(err).Msg("Failed to sync to Port.io")
			// Don't fail the reconciliation if Port sync fails
		}
	}

	return ctrl.Result{}, nil
}

// syncToPort syncs the FrontendPage state to Port.io
func (r *FrontendPageReconciler) syncToPort(ctx context.Context, page *frontendv1alpha1.FrontendPage) error {
	// Get deployment status
	var dep appsv1.Deployment
	depErr := r.Get(ctx, client.ObjectKey{Name: page.Name, Namespace: page.Namespace}, &dep)
	
	status := "Creating"
	if depErr == nil {
		if dep.Status.ReadyReplicas == dep.Status.Replicas && dep.Status.Replicas > 0 {
			status = "Ready"
		} else if dep.Status.Conditions != nil {
			for _, condition := range dep.Status.Conditions {
				if condition.Type == appsv1.DeploymentProgressing && condition.Status == corev1.ConditionFalse {
					status = "Failed"
					break
				}
			}
		}
	}

	// Create Port entity
	entity := port.FrontendPageEntity{
		Identifier: fmt.Sprintf("%s-%s", page.Name, page.Namespace),
		Title:      page.Name,
		Blueprint:  "frontendpage",
		Properties: map[string]interface{}{
			"contents":  page.Spec.Contents,
			"image":     page.Spec.Image,
			"replicas":  page.Spec.Replicas,
			"namespace": page.Namespace,
			"status":    status,
			"createdAt": page.CreationTimestamp.Format(time.RFC3339),
			"updatedAt": time.Now().Format(time.RFC3339),
		},
	}

	// Add URL if we can determine it (placeholder for now)
	if status == "Ready" {
		entity.Properties["url"] = fmt.Sprintf("http://%s.%s.svc.cluster.local", page.Name, page.Namespace)
	}

	return r.PortClient.CreateOrUpdateEntity(entity)
}

func AddFrontendController(mgr manager.Manager) error {
	return AddFrontendControllerWithPort(mgr, nil)
}

func AddFrontendControllerWithPort(mgr manager.Manager, portClient *port.PortClient) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&frontendv1alpha1.FrontendPage{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.ConfigMap{}).
		Complete(&FrontendPageReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		 PortClient: portClient,
	})
}
