package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	"github.com/yourusername/k8s-controller-tutorial/pkg/ctrl"
	"github.com/yourusername/k8s-controller-tutorial/pkg/informer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrlruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var serverPort int
var serverKubeconfig string
var serverInCluster bool

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a FastHTTP server and deployment informer",
	Run: func(cmd *cobra.Command, args []string) {
		level := parseLogLevel(logLevel)
		configureLogger(level)
		clientset, err := getServerKubeClient(serverKubeconfig, serverInCluster)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create Kubernetes client")
			os.Exit(1)
		}
		ctx := context.Background()
		go informer.StartDeploymentInformer(ctx, clientset)
		go informer.StartEventInformer(ctx, clientset)

		// Start controller-runtime manager and controller
		config, err := getKubeConfig(serverKubeconfig, serverInCluster)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get kubeconfig for controller-runtime")
			os.Exit(1)
		}
		mgr, err := ctrlruntime.NewManager(config, manager.Options{})
		if err != nil {
			log.Error().Err(err).Msg("Failed to create controller-runtime manager")
			os.Exit(1)
		}
		if err := ctrl.AddDeploymentController(mgr); err != nil {
			log.Error().Err(err).Msg("Failed to add deployment controller")
			os.Exit(1)
		}
		go func() {
			log.Info().Msg("Starting controller-runtime manager...")
			if err := mgr.Start(cmd.Context()); err != nil {
				log.Error().Err(err).Msg("Manager exited with error")
				os.Exit(1)
			}
		}()

		handler := func(ctx *fasthttp.RequestCtx) {
			requestID := uuid.New().String()
			ctx.Response.Header.Set("X-Request-ID", requestID)
			logger := log.With().Str("request_id", requestID).Logger()
			switch string(ctx.Path()) {
			case "/deployments":
				logger.Info().Msg("Deployments request received")
				ctx.Response.Header.Set("Content-Type", "application/json")
				deployments := informer.GetDeploymentNames()
				logger.Info().Msgf("Deployments: %v", deployments)
				ctx.SetStatusCode(200)
				ctx.Write([]byte("["))
				for i, name := range deployments {
					ctx.WriteString("\"")
					ctx.WriteString(name)
					ctx.WriteString("\"")
					if i < len(deployments)-1 {
						ctx.WriteString(",")
					}
				}
				ctx.Write([]byte("]"))
				return
			case "/events":
				logger.Info().Msg("Events request received")
				ctx.Response.Header.Set("Content-Type", "application/json")
				events := informer.GetRecentEvents()
				logger.Info().Msgf("Found %d events", len(events))
				
				// Use proper JSON encoding
				response := map[string]interface{}{
					"events": events,
					"total":  len(events),
				}
				
				jsonData, err := json.Marshal(response)
				if err != nil {
					logger.Error().Err(err).Msg("Failed to marshal events to JSON")
					ctx.SetStatusCode(500)
					ctx.WriteString("{\"error\":\"Failed to encode events\"}")
					return
				}
				
				ctx.SetStatusCode(200)
				ctx.Write(jsonData)
				return
			default:
				logger.Info().Msg("Default request received")
				fmt.Fprintf(ctx, "Hello from FastHTTP!")
			}
		}
		addr := fmt.Sprintf(":%d", serverPort)
		log.Info().Msgf("Starting FastHTTP server on %s (version: %s)", addr, appVersion)
		if err := fasthttp.ListenAndServe(addr, handler); err != nil {
			log.Error().Err(err).Msg("Error starting FastHTTP server")
			os.Exit(1)
		}
	},
}

func getKubeConfig(kubeconfigPath string, inCluster bool) (*rest.Config, error) {
	if inCluster {
		return rest.InClusterConfig()
	}
	return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
}

func getServerKubeClient(kubeconfigPath string, inCluster bool) (*kubernetes.Clientset, error) {
	config, err := getKubeConfig(kubeconfigPath, inCluster)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVar(&serverPort, "port", 8080, "Port to run the server on")
	serverCmd.Flags().StringVar(&serverKubeconfig, "kubeconfig", "", "Path to the kubeconfig file")
	serverCmd.Flags().BoolVar(&serverInCluster, "in-cluster", false, "Use in-cluster Kubernetes config")
}
