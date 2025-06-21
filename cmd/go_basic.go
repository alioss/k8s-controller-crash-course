package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var goBasicCmd = &cobra.Command{
	Use:   "go-basic",
	Short: "Run golang basic code with structured logging",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Starting go-basic command")
		
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
		log.Info().Msg("Displaying current users")
		k8s.GetUsers()

		// Add new user to struct
		log.Info().Str("new_user", "anonymous").Msg("Adding new user")
		k8s.AddNewUser("anonymous")

		// Print users one more time
		log.Info().
			Int("updated_users_count", len(k8s.Users)).
			Msg("Displaying updated users")
		k8s.GetUsers()
		
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

// GetUsers method using structured logging instead of fmt.Println
func (k8s Kubernetes) GetUsers() {
	for i, user := range k8s.Users {
		log.Debug().
			Int("index", i+1).
			Str("username", user).
			Msg("User found")
	}
}

// AddNewUser method with logging
func (k8s *Kubernetes) AddNewUser(user string) {
	k8s.Users = append(k8s.Users, user)
	log.Debug().
		Str("added_user", user).
		Int("total_users", len(k8s.Users)).
		Msg("User added successfully")
	}
