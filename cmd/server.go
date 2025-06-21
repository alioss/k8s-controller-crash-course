package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
)

var serverPort int

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a FastHTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		// Configure logging before starting server
		level := parseLogLevel(logLevel)
		configureLogger(level)
		
		handler := func(ctx *fasthttp.RequestCtx) {
			path := string(ctx.Path())
			method := string(ctx.Method())
			
			// Log each incoming request
			log.Debug().
				Str("method", method).
				Str("path", path).
				Str("user_agent", string(ctx.UserAgent())).
				Str("remote_addr", ctx.RemoteAddr().String()).
				Msg("Incoming request")
			
			// Simple routing based on path
			switch path {
			case "/health":
				handleHealth(ctx)
			case "/":
				handleRoot(ctx)
			default:
				handleNotFound(ctx)
			}
			
			// Log response
			log.Info().
				Str("method", method).
				Str("path", path).
				Int("status", ctx.Response.StatusCode()).
				Msg("Request completed")
		}
		addr := fmt.Sprintf(":%d", serverPort)
		log.Info().Msgf("Starting FastHTTP server on %s", addr)
		if err := fasthttp.ListenAndServe(addr, handler); err != nil {
			log.Error().Err(err).Msg("Error starting FastHTTP server")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVar(&serverPort, "port", 8080, "Port to run the server on")
}

// handleRoot handles requests to the root path
func handleRoot(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/plain")
	fmt.Fprintf(ctx, "Hello from FastHTTP!")
}

// handleHealth handles health check requests
func handleHealth(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(200)
	fmt.Fprintf(ctx, `{"status":"healthy","service":"k8s-controller-tutorial","port":%d}`, serverPort)
	log.Debug().Msg("Health check requested")
}

// handleNotFound handles 404 responses
func handleNotFound(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(404)
	fmt.Fprintf(ctx, `{"error":"Not Found","path":"%s","method":"%s"}`, ctx.Path(), ctx.Method())
	log.Warn().Str("path", string(ctx.Path())).Msg("404 Not Found")
}
