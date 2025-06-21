package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var goBasicCmd = &cobra.Command{
	Use:   "go-basic",
	Short: "Run golang basic code with configurable log levels",
	Run: func(cmd *cobra.Command, args []string) {
		log.Trace().Msg("TRACE: Starting go-basic command")
		log.Debug().Msg("DEBUG: Initializing Kubernetes struct")
		log.Info().Msg("INFO: Starting go-basic command")
		
		// Go basic code to run functions
		k8s := Kubernetes{
			Name:    "k8s-demo-cluster",
			Version: "1.31",
			Users:   []string{"alex", "den"},
			NodeNumber: func() int {
				return 10
			},
		}

		log.Info().
			Str("cluster_name", k8s.Name).
			Str("version", k8s.Version).
			Int("initial_users_count", len(k8s.Users)).
			Msg("Kubernetes cluster initialized")

		// Print users using structured logging
		log.Debug().Msg("DEBUG: About to display current users")
		log.Info().Msg("Displaying current users")
		k8s.GetUsers()

		// Add new user to struct
		log.Debug().Str("new_user", "anonymous").Msg("DEBUG: Preparing to add new user")
		log.Info().Str("new_user", "anonymous").Msg("Adding new user")
		k8s.AddNewUser("anonymous")

		// Print users one more time
		log.Info().
			Int("updated_users_count", len(k8s.Users)).
			Msg("Displaying updated users")
		k8s.GetUsers()
		
		log.Trace().Msg("TRACE: go-basic command completed")
		log.Info().Msg("go-basic command completed successfully")
	},
}

func init() {
	rootCmd.AddCommand(goBasicCmd)
}

// Kubernetes struct with JSON tags
type Kubernetes struct {
	Name       string     `json:"name"`
	Version    string     `json:"version"`
	Users      []string   `json:"users,omitempty"`
	NodeNumber func() int `json:"-"`
}

// GetUsers method using structured logging with different levels
func (k8s Kubernetes) GetUsers() {
	log.Debug().Int("total_users", len(k8s.Users)).Msg("DEBUG: Starting user enumeration")
	
	for i, user := range k8s.Users {
		log.Trace().
			Int("index", i+1).
			Str("username", user).
			Msg("TRACE: Processing user")
			
		log.Debug().
			Int("index", i+1).
			Str("username", user).
			Msg("DEBUG: User found")
			
		log.Info().
			Str("username", user).
			Msg("User")
	}
}

// AddNewUser method with logging at different levels
func (k8s *Kubernetes) AddNewUser(user string) {
	log.Trace().Str("user_to_add", user).Msg("TRACE: Entering AddNewUser method")
	
	oldCount := len(k8s.Users)
	k8s.Users = append(k8s.Users, user)
	
	log.Debug().
		Str("added_user", user).
		Int("old_count", oldCount).
		Int("new_count", len(k8s.Users)).
		Msg("DEBUG: User added to slice")
		
	log.Info().
		Str("added_user", user).
		Int("total_users", len(k8s.Users)).
		Msg("User added successfully")
}
