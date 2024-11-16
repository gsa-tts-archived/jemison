package main

import (
	"os"
	"testing"

	"github.com/GSA-TTS/jemison/internal/env"
	"github.com/stretchr/testify/assert"
)

// Tests need a backend
// docker compose -f backend.yaml up

func setup( /* t *testing.T */ ) func(t *testing.T) {
	os.Setenv("ENV", "LOCALHOST")
	env.InitGlobalEnv("entree")
	env.SetupLogging("entree")
	return func(t *testing.T) {
		// t.Log("teardown test case")
	}
}
func TestFetchWithHallPass(t *testing.T) {
	setup()
	which, _ := fetchOne(true, 1, "https", "digitalcorps.gsa.gov", "/")
	assert.Equal(t, HallPass, which)
}

func TestFetchWithoutHallPass(t *testing.T) {
	setup()
	which, _ := fetchOne(false, 1, "https", "digitalcorps.gsa.gov", "/")
	assert.Contains(t, []FetchStatus{DeadlineNotYetReached, DeadlinePassed}, which)

}

func TestUnknownUrlAtKnownHost(t *testing.T) {
	/*
		entree               | {"level":"debug","timestamp":"2024-11-16T15:00:32.046Z","caller":"entree/work.go:113","msg":"non-indexable content","pid":13,"host_id":1,"host":"digitalcorps.gsa.gov","path":"/blogs/from-silicon-valley-to-govtech-how-a-data-scientist-turned-their-passion-for-helping-people-into-a-career-in-public-service","hallpass":false}
	*/
	setup()
	which, _ := fetchOne(false, 1, "https", "digitalcorps.gsa.gov", "/blogs/from-silicon-valley-to-govtech-how-a-data-scientist-turned-their-passion-for-helping-people-into-a-career-in-public-service")
	assert.Equal(t, NotPreviouslySeen, which)
}
