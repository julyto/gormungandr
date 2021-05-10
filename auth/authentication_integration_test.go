package auth

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var dockerDB *sql.DB

func TestMain(m *testing.M) {
	flag.Parse() //required to get Short() from testing
	if testing.Short() {
		log.Warn("skipping test Docker in short mode.")
		os.Exit(m.Run())
	}
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "9.4", []string{"POSTGRES_HOST_AUTH_METHOD=trust",
		"POSTGRESQL_PASSWORD=secret", "POSTGRES_USER=jormun", "POSTGRES_DB=jormun"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	conStr := ""
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 30 * time.Second
	if err = pool.Retry(func() error {
		var errr error
		conStr = fmt.Sprintf("user=jormun password=secret host=localhost port=%s dbname=jormun sslmode=disable", resource.GetPort("5432/tcp"))
		dockerDB, errr = sql.Open("postgres", conStr)
		if errr != nil {
			return errr
		}
		return dockerDB.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	//loading fixture
	cmd := exec.Command("psql", conStr, "-f", "fixtures/jormun.sql")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Could not restore postres backup: %s", err)
	}
	//Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestRealAuthenticate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test Docker in short mode.")
	}
	t.Parallel()
	_, err := authenticate("thisIsNotAkey", time.Now(), dockerDB)
	assert.Equal(t, ErrAuthenticationFailed, err)

	user, err := authenticate("115aa17b-63d3-4a31-acd6-edebebd4d415", time.Now(), dockerDB)
	assert.Nil(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "test", user.AppName)
	assert.Equal(t, "with_free_instances", user.Type)
	assert.Equal(t, "115aa17b-63d3-4a31-acd6-edebebd4d415", user.Token)
	assert.Equal(t, "navitia.io", user.EndPointName)
	assert.Equal(t, 1, user.EndPointId)

	//fr-idf is in opendata
	ok, err := isAuthorized(user, "fr-idf", dockerDB)
	assert.Nil(t, err)
	assert.True(t, ok)

	//Transilien is private but we have the authorization to use it
	ok, err = isAuthorized(user, "transilien", dockerDB)
	assert.Nil(t, err)
	assert.True(t, ok)

	//sncf is private And we don't have any authorization on it
	ok, err = isAuthorized(user, "sncf", dockerDB)
	assert.Nil(t, err)
	assert.False(t, ok)
}
