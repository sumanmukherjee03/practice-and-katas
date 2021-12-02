package handlers

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/pusher/pusher-http-go"
	"github.com/robfig/cron/v3"
	"github.com/tsawler/vigilate/internal/channeldata"
	"github.com/tsawler/vigilate/internal/config"
	"github.com/tsawler/vigilate/internal/driver"
	"github.com/tsawler/vigilate/internal/helpers"
	"github.com/tsawler/vigilate/internal/repository/dbrepo"
)

var (
	testSession *scs.SessionManager
)

func TestMain(m *testing.M) {
	testSession = scs.New()
	testSession.Lifetime = 24 * time.Hour
	testSession.Cookie.Persist = true
	testSession.Cookie.SameSite = http.SameSiteLaxMode
	testSession.Cookie.Secure = false

	mailQueue := make(chan channeldata.MailJob, 5)

	// define application configuration
	a := config.AppConfig{
		DB:           &driver.DB{},
		Session:      testSession,
		InProduction: false,
		Domain:       "localhost",
		MailQueue:    mailQueue,
	}

	app = &a

	preferenceMap := make(map[string]string)

	app.PreferenceMap = preferenceMap

	// create pusher client
	wsClient := pusher.Client{
		AppID:  "1",
		Secret: "123abc",
		Key:    "abc123",
		Secure: false,
		Host:   "localhost:4001",
	}

	app.WsClient = wsClient

	monitorMap := make(map[int]cron.EntryID)
	app.MonitorMap = monitorMap

	localZone, _ := time.LoadLocation("Local")
	scheduler := cron.New(cron.WithLocation(localZone), cron.WithChain(
		cron.DelayIfStillRunning(cron.DefaultLogger),
		cron.Recover(cron.DefaultLogger),
	))

	app.Scheduler = scheduler

	// Use a dummy repo for testing purposes
	repo := NewTestHandlers(app)
	NewHandlers(repo, app)

	helpers.NewHelpers(app)
	// Fix the path of the views for tests
	helpers.SetViews("./../../views")
}

// gets the context from the req, adds session to it and returns the context
// You can add this returned context back into the request
func getCtx(req *http.Request) context.Context {
	ctx, err := testSession.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

// NewTestHandlers creates a new repository
func NewTestHandlers(a *config.AppConfig) *DBRepo {
	return &DBRepo{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}
