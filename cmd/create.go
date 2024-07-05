/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/whites11/podcast-sync-server/internal/db"
	"github.com/whites11/podcast-sync-server/internal/models"
	"github.com/whites11/podcast-sync-server/internal/repository"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user",
	Long:  `Create a new user`,
	RunE:  runCreateUser,
}

const (
	flagUsername = "username"
)

func init() {
	usersCmd.AddCommand(createCmd)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	createCmd.Flags().StringP(flagUsername, "u", "", "Username")

}

func runCreateUser(cmd *cobra.Command, args []string) error {
	username, err := cmd.Flags().GetString(flagUsername)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	dbfactory, err := db.NewSqliteFactory()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	db, err := dbfactory.Build()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	usersRepository, err := repository.NewUsersRepository(db)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("Enter the password of the new user: ")
	var pwd string
	// Read the users input
	_, err = fmt.Scan(&pwd)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	user, err := usersRepository.CreateUser(models.User{
		Username: username,
		Password: pwd,
	})
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Printf("User %s created (id = %d)\n", username, user.ID)

	return nil
}
