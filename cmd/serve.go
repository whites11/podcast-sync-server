package cmd

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/gorilla/securecookie"
	"github.com/spf13/cobra"

	"github.com/whites11/podcast-sync-server/internal/db"
	"github.com/whites11/podcast-sync-server/internal/dependencies"
	"github.com/whites11/podcast-sync-server/internal/repository"
	"github.com/whites11/podcast-sync-server/internal/routes"
	"github.com/whites11/podcast-sync-server/internal/settings"
)

const (
	flagPort = "port"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start listening",
	Long:  `Start listening`,
	RunE:  run,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().Uint16P(flagPort, "p", 3000, "Port to listen to")
}

func run(cmd *cobra.Command, args []string) error {
	port, err := cmd.Flags().GetUint16(flagPort)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	var deps *dependencies.Dependencies
	{
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

		devicesRepository, err := repository.NewDevicesRepository(db)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		subscriptionsRepository, err := repository.NewSubscriptionsRepository(db)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		episodesActionsRepository, err := repository.NewEpisodeActionsRepository(db)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		secretsStorage, err := settings.NewDatabaseStorage(db)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		secretsProvider, err := settings.New(secretsStorage)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		// Hash keys should be at least 32 bytes long
		hashKey, err := secretsProvider.GetOrGenerate("cookies-hash-key", 32)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
		blockKey, err := secretsProvider.GetOrGenerate("cookies-block-key", 32)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		secureCookie := securecookie.New([]byte(hashKey), []byte(blockKey))

		deps = dependencies.New(devicesRepository, episodesActionsRepository, subscriptionsRepository, usersRepository, secureCookie)
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	routes.NewProbesRouter().Setup(r)

	r.Route("/api", func(r chi.Router) {
		routes.NewV2Router(r, deps).Setup()
	})

	fmt.Printf("Listening on 0.0.0.0:%d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
