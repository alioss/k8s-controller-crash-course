package informer

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var deploymentInformer cache.SharedIndexInformer
var eventInformer cache.SharedIndexInformer

// Event represents a simplified Kubernetes event
type Event struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Reason    string `json:"reason"`
	Message   string `json:"message"`
	Type      string `json:"type"`
	Object    string `json:"object"`
	Timestamp string `json:"timestamp"`
}

// StartDeploymentInformer starts a shared informer for Deployments in the default namespace.
func StartDeploymentInformer(ctx context.Context, clientset *kubernetes.Clientset) {
	factory := informers.NewSharedInformerFactoryWithOptions(
		clientset,
		30*time.Second,
		informers.WithNamespace("default"),
		informers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.FieldSelector = fields.Everything().String()
		}),
	)
	deploymentInformer = factory.Apps().V1().Deployments().Informer()

	deploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			log.Info().Msgf("Deployment added: %s", getDeploymentName(obj))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			log.Info().Msgf("Deployment updated: %s", getDeploymentName(newObj))
		},
		DeleteFunc: func(obj interface{}) {
			log.Info().Msgf("Deployment deleted: %s", getDeploymentName(obj))
		},
	})

	log.Info().Msg("Starting deployment informer...")
	factory.Start(ctx.Done())
	for t, ok := range factory.WaitForCacheSync(ctx.Done()) {
		if !ok {
			log.Error().Msgf("Failed to sync informer for %v", t)
			os.Exit(1)
		}
	}
	log.Info().Msg("Deployment informer cache synced. Watching for events...")
	<-ctx.Done() // Block until context is cancelled
}

// StartEventInformer starts a shared informer for Events in the default namespace.
func StartEventInformer(ctx context.Context, clientset *kubernetes.Clientset) {
	factory := informers.NewSharedInformerFactoryWithOptions(
		clientset,
		30*time.Second,
		informers.WithNamespace("default"),
		informers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.FieldSelector = fields.Everything().String()
		}),
	)
	eventInformer = factory.Core().V1().Events().Informer()

	eventInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if event, ok := obj.(*corev1.Event); ok {
				log.Info().Msgf("📅 Event [%s]: %s/%s - %s: %s", 
					event.Type, event.InvolvedObject.Kind, event.InvolvedObject.Name, 
					event.Reason, event.Message)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			if event, ok := newObj.(*corev1.Event); ok {
				log.Debug().Msgf("📅 Event Updated [%s]: %s/%s - %s: %s", 
					event.Type, event.InvolvedObject.Kind, event.InvolvedObject.Name, 
					event.Reason, event.Message)
			}
		},
		DeleteFunc: func(obj interface{}) {
			log.Debug().Msg("📅 Event deleted")
		},
	})

	log.Info().Msg("Starting event informer...")
	factory.Start(ctx.Done())
	for t, ok := range factory.WaitForCacheSync(ctx.Done()) {
		if !ok {
			log.Error().Msgf("Failed to sync event informer for %v", t)
			os.Exit(1)
		}
	}
	log.Info().Msg("Event informer cache synced. Watching for events...")
	<-ctx.Done() // Block until context is cancelled
}

// GetDeploymentNames returns a slice of deployment names from the informer's cache.
func GetDeploymentNames() []string {
	var names []string
	if deploymentInformer == nil {
		return names
	}
	for _, obj := range deploymentInformer.GetStore().List() {
		if d, ok := obj.(*appsv1.Deployment); ok {
			names = append(names, d.Name)
		}
	}
	return names
}

// GetRecentEvents returns a slice of recent events from the informer's cache.
func GetRecentEvents() []Event {
	var events []Event
	if eventInformer == nil {
		return events
	}
	for _, obj := range eventInformer.GetStore().List() {
		if e, ok := obj.(*corev1.Event); ok {
			event := Event{
				Name:      e.Name,
				Namespace: e.Namespace,
				Reason:    e.Reason,
				Message:   e.Message,
				Type:      e.Type,
				Object:    e.InvolvedObject.Kind + "/" + e.InvolvedObject.Name,
				Timestamp: e.CreationTimestamp.Format("2006-01-02 15:04:05"),
			}
			events = append(events, event)
		}
	}
	return events
}

func getDeploymentName(obj any) string {
	if d, ok := obj.(metav1.Object); ok {
		return d.GetName()
	}
	return "unknown"
}
