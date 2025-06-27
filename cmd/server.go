package cmd

import (
	"fmt"
	"os"

	"github.com/buaazp/fasthttprouter"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	"github.com/yourusername/k8s-controller-tutorial/pkg/api"
	"github.com/yourusername/k8s-controller-tutorial/pkg/ctrl"
	"github.com/yourusername/k8s-controller-tutorial/pkg/informer"
	"github.com/yourusername/k8s-controller-tutorial/pkg/port"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"

	frontendv1alpha1 "github.com/yourusername/k8s-controller-tutorial/pkg/apis/frontend/v1alpha1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrlruntime "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"

	_ "github.com/yourusername/k8s-controller-tutorial/docs" // for swagger
)

var serverPort int
var serverKubeconfig string
var serverInCluster bool
var enableLeaderElection bool
var leaderElectionNamespace string
var metricsPort int
var portBaseURL string
var portClientID string
var portClientSecret string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a FastHTTP server and deployment informer",
	Run: func(cmd *cobra.Command, args []string) {
		level := parseLogLevel(logLevel)
		configureLogger(level)

		logf.SetLogger(zap.New(zap.UseDevMode(true)))
		logf.SetLogger(zerologr.New(&log.Logger))

		scheme := runtime.NewScheme()
		if err := clientgoscheme.AddToScheme(scheme); err != nil {
			log.Error().Err(err).Msg("Failed to add client-go scheme")
			os.Exit(1)
		}
		if err := frontendv1alpha1.AddToScheme(scheme); err != nil {
			log.Error().Err(err).Msg("Failed to add FrontendPage scheme")
			os.Exit(1)
		}
		mgr, err := ctrlruntime.NewManager(ctrlruntime.GetConfigOrDie(), manager.Options{
			Scheme:                  scheme,
			LeaderElection:          enableLeaderElection,
			LeaderElectionID:        "k8s-controller-tutorial-leader-election",
			LeaderElectionNamespace: leaderElectionNamespace,
			Metrics:                 server.Options{BindAddress: fmt.Sprintf(":%d", metricsPort)},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to create controller manager")
			os.Exit(1)
		}

		// Setup Port.io client if configured
		var portClient *port.PortClient
		if portBaseURL != "" && portClientID != "" && portClientSecret != "" {
			log.Info().Msg("Setting up Port.io integration...")
			portClient = port.NewPortClient(portBaseURL, portClientID, portClientSecret)
			if err := portClient.Authenticate(); err != nil {
				log.Warn().Err(err).Msg("Failed to authenticate with Port.io, continuing without integration")
				portClient = nil
			} else {
				log.Info().Msg("Successfully connected to Port.io")
			}
		} else {
			log.Info().Msg("Port.io integration not configured (use --port-* flags to enable)")
		}

		if err := ctrl.AddFrontendControllerWithPort(mgr, portClient); err != nil {
			log.Error().Err(err).Msg("Failed to add frontend controller")
			os.Exit(1)
		}

		// --- API ROUTER SETUP ---
		router := fasthttprouter.New()
		frontendAPI := &api.FrontendPageAPI{
			K8sClient: mgr.GetClient(),
			Namespace: "default", // or make configurable
		}
		router.GET("/api/frontendpages", frontendAPI.ListFrontendPages)
		router.POST("/api/frontendpages", frontendAPI.CreateFrontendPage)
		router.GET("/api/frontendpages/:name", frontendAPI.GetFrontendPage)
		router.PUT("/api/frontendpages/:name", frontendAPI.UpdateFrontendPage)
		router.DELETE("/api/frontendpages/:name", frontendAPI.DeleteFrontendPage)

		// Swagger JSON endpoint
		router.GET("/swagger/doc.json", func(ctx *fasthttp.RequestCtx) {
			ctx.SetContentType("application/json")
			swaggerJSON := `{"swagger":"2.0","info":{"title":"FrontendPage API","version":"1.0"},"paths":{"/api/frontendpages":{"get":{"summary":"List FrontendPages","responses":{"200":{"description":"OK"}}},"post":{"summary":"Create FrontendPage","responses":{"201":{"description":"Created"}}}},"/api/frontendpages/{name}":{"get":{"summary":"Get FrontendPage","responses":{"200":{"description":"OK"}}},"put":{"summary":"Update FrontendPage","responses":{"200":{"description":"OK"}}},"delete":{"summary":"Delete FrontendPage","responses":{"204":{"description":"No Content"}}}}}}`
			ctx.SetBodyString(swaggerJSON)
		})

		// Simple Swagger UI (basic HTML)
		router.GET("/swagger/index.html", func(ctx *fasthttp.RequestCtx) {
			ctx.SetContentType("text/html")
			swaggerHTML := `<!DOCTYPE html>
<html>
<head>
    <title>FrontendPage API</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui.css" />
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui-bundle.js"></script>
    <script>
    SwaggerUIBundle({
        url: '/swagger/doc.json',
        dom_id: '#swagger-ui',
        presets: [
            SwaggerUIBundle.presets.apis,
            SwaggerUIBundle.presets.standalone
        ]
    });
    </script>
</body>
</html>`
			ctx.SetBodyString(swaggerHTML)
		})

		// Legacy endpoint for deployments
		router.GET("/deployments", func(ctx *fasthttp.RequestCtx) {
			ctx.Response.Header.Set("Content-Type", "application/json")
			deployments := informer.GetDeploymentNames()
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
		})

		go func() {
			log.Info().Msg("Starting controller-runtime manager...")
			if err := mgr.Start(cmd.Context()); err != nil {
				log.Error().Err(err).Msg("Manager exited with error")
				os.Exit(1)
			}
		}()

		addr := fmt.Sprintf(":%d", serverPort)
		log.Info().Msgf("Starting FastHTTP server on %s", addr)
		if err := fasthttp.ListenAndServe(addr, router.Handler); err != nil {
			log.Error().Err(err).Msg("Error starting FastHTTP server")
			os.Exit(1)
		}
	},
}

func getServerKubeClient(kubeconfigPath string, inCluster bool) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if inCluster {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}
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
	serverCmd.Flags().BoolVar(&enableLeaderElection, "enable-leader-election", true, "Enable leader election for controller manager")
	serverCmd.Flags().StringVar(&leaderElectionNamespace, "leader-election-namespace", "default", "Namespace for leader election")
	serverCmd.Flags().IntVar(&metricsPort, "metrics-port", 8081, "Port for controller manager metrics")
	
	// Port.io integration flags
	serverCmd.Flags().StringVar(&portBaseURL, "port-base-url", "", "Port.io API base URL (e.g., https://api.getport.io)")
	serverCmd.Flags().StringVar(&portClientID, "port-client-id", "", "Port.io client ID")
	serverCmd.Flags().StringVar(&portClientSecret, "port-client-secret", "", "Port.io client secret")
}
