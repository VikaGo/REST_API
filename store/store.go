package store

import (
	context "context"
	"github.com/jmoiron/sqlx"
	"log"
	"time"

	"github.com/VikaGo/REST_API/config"
	"github.com/VikaGo/REST_API/logger"
	"github.com/VikaGo/REST_API/store/pg"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// Store contains all repositories
type Store struct {
	Pg   *sqlx.DB // for KeepAlivePg (see below)
	User UserRepo
}

// New creates new store
func New(ctx context.Context) (*Store, error) {
	cfg := config.Get()

	// Connect to PostgreSQL using sqlx
	pgDB, err := sqlx.Connect("postgres", cfg.PgURL)
	if err != nil {
		return nil, errors.Wrap(err, "sqlx.Connect failed")
	}

	err = pgDB.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	// Run PostgreSQL migrations
	log.Println("Running PostgreSQL migrations...")
	if err := runPgMigrations(pgDB); err != nil {
		return nil, errors.Wrap(err, "runPgMigrations failed")
	}

	var store Store

	// Init PostgreSQL repositories
	store.Pg = pgDB
	go store.KeepAlivePg()
	store.User = pg.NewUserRepo(pgDB)
	return &store, nil
}

// KeepAlivePollPeriod is a Pg keepalive check time period
const KeepAlivePollPeriod = 3

// KeepAlivePg makes sure PostgreSQL is alive and reconnects if needed
func (store *Store) KeepAlivePg() {
	logger := logger.Get()
	var err error
	for {
		// Check if PostgreSQL is alive every 3 seconds
		time.Sleep(time.Second * KeepAlivePollPeriod)
		lostConnect := false
		if store.Pg == nil {
			lostConnect = true
		} else if err = store.Pg.Ping(); err != nil {
			lostConnect = true
		}
		if !lostConnect {
			continue
		}
		logger.Debug().Msg("[store.KeepAlivePg] Lost PostgreSQL connection. Restoring...")
		store.Pg, err = sqlx.Connect("postgres", config.Get().PgURL)
		if err != nil {
			logger.Err(err)
			continue
		}
		logger.Debug().Msg("[store.KeepAlivePg] PostgreSQL reconnected")
	}
}
