package files

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_GetActorData(t *testing.T) {
	actorFile, err := os.Open("../actorstest.csv")
	if err != nil {
		t.Fail()
	}

	actors := getActorData(actorFile)
	assert.NotNil(t, actors, "actors should not be empty")
	assert.NotEmpty(t, actors, "actors should not be empty")
}

func Test_GetEventData(t *testing.T) {
	eventFile, err := os.Open("../eventstest.csv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read event file content. Error: %s", err.Error())
	}

	events := getEventData(eventFile)
	assert.NotNil(t, events, "events should not be empty")
	assert.NotEmpty(t, events, "events should not be empty")
}

func Test_GetRepoData(t *testing.T) {
	repoFile, err := os.Open("../repostest.csv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read repo file content. Error: %s", err.Error())
	}

	repos := getRepoData(repoFile)
	assert.NotNil(t, repos, "repos should not be empty")
	assert.NotEmpty(t, repos, "repos should not be empty")
}

func Test_GetFilteredActors(t *testing.T) {
	actorFile, err := os.Open("../actorstest.csv")
	if err != nil {
		t.Fail()
	}

	actors := getActorData(actorFile)
	assert.NotNil(t, actors, "actors should not be empty")
	assert.NotEmpty(t, actors, "actors should not be empty")

	eventFile, err := os.Open("../eventstest.csv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read event file content. Error: %s", err.Error())
	}

	events := getEventData(eventFile)
	assert.NotNil(t, events, "events should not be empty")
	assert.NotEmpty(t, events, "events should not be empty")

	filteredActors := getFilteredActors(events["createevent"], events["pushevent"], actors)
	assert.NotNil(t, filteredActors, "filteredActors should not be empty")
	assert.NotEmpty(t, filteredActors, "filteredActors should not be empty")
}

func Test_GetFilteredRepos(t *testing.T) {
	repoFile, err := os.Open("../repostest.csv")
	if err != nil {
		t.Fail()
	}

	repos := getRepoData(repoFile)
	assert.NotNil(t, repos, "repos should not be empty")
	assert.NotEmpty(t, repos, "repos should not be empty")

	eventFile, err := os.Open("../eventstest.csv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read event file content. Error: %s", err.Error())
	}

	events := getEventData(eventFile)
	assert.NotNil(t, events, "events should not be empty")
	assert.NotEmpty(t, events, "events should not be empty")

	filteredRepos := getFilteredRepos(events["pushevent"], repos)
	assert.NotNil(t, filteredRepos, "filteredRepos should not be empty")
	assert.NotEmpty(t, filteredRepos, "filteredRepos should not be empty")
}

func Test_ProcessOrders(t *testing.T) {
	err := ProcessOrders("../actors.csv", "../events.csv", "../repos.csv")
	assert.Nil(t, err, "error should be nil")
}

func Test_ProcessOrdersFail(t *testing.T) {
	err := ProcessOrders("actors.csv", "../events.csv", "../repos.csv")
	assert.NotNil(t, err, "error should not be nil")
	assert.Equal(t, "could not read actor file content. Error: open actors.csv: no such file or directory",
		err.Error())

	err = ProcessOrders("../actors.csv", "events.csv", "../repos.csv")
	assert.NotNil(t, err, "error should not be nil")
	assert.Equal(t, "could not read event file content. Error: open events.csv: no such file or directory",
		err.Error())

	err = ProcessOrders("../actors.csv", "../events.csv", "repos.csv")
	assert.NotNil(t, err, "error should not be nil")
	assert.Equal(t, "could not read repo file content. Error: open repos.csv: no such file or directory",
		err.Error())
}