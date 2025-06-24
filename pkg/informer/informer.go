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
var nodeInformer cache.SharedIndexInformer

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

// StartNodeInformer starts a shared informer for Nodes (cluster-wide)
func StartNodeInformer(ctx context.Context, clientset *kubernetes.Clientset) {
	factory := informers.NewSharedInformerFactory(clientset, 30*time.Second)
	nodeInformer = factory.Core().V1().Nodes().Informer()

	nodeInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			log.Info().Msgf("Node added: %s", getNodeName(obj))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			log.Info().Msgf("Node updated: %s", getNodeName(newObj))
		},
		DeleteFunc: func(obj interface{}) {
			log.Info().Msgf("Node deleted: %s", getNodeName(obj))
		},
	})

	log.Info().Msg("Starting node informer...")
	factory.Start(ctx.Done())
	for t, ok := range factory.WaitForCacheSync(ctx.Done()) {
		if !ok {
			log.Error().Msgf("Failed to sync node informer for %v", t)
			os.Exit(1)
		}
	}
	log.Info().Msg("Node informer cache synced. Watching for node events...")
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

// GetNodeNames returns a slice of node names from the informer's cache.
func GetNodeNames() []string {
	var names []string
	if nodeInformer == nil {
		return names
	}
	for _, obj := range nodeInformer.GetStore().List() {
		if n, ok := obj.(*corev1.Node); ok {
			names = append(names, n.Name)
		}
	}
	return names
}

func getDeploymentName(obj any) string {
	if d, ok := obj.(metav1.Object); ok {
		return d.GetName()
	}
	return "unknown"
}

func getNodeName(obj any) string {
	if n, ok := obj.(metav1.Object); ok {
		return n.GetName()
	}
	return "unknown"
}
