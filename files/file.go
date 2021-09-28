package files

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// InputParams get values from the parameters passed to the cli
func InputParams() ProcessInput {
	var input ProcessInput

	// define the flags accepted in the command line
	actorsFilename := flag.String("actors", "actors.csv", "actors filename to be processed")
	eventsFilename := flag.String("events", "events.csv", "events filename to be processed")
	reposFilename := flag.String("repos", "repos.csv", "repos filename to be processed")
	commitsFilename := flag.String("commits", "commits.csv", "commits filename to be processed")

	flag.Parse()

	input.ActorsFile = *actorsFilename
	input.EventsFile = *eventsFilename
	input.ReposFile = *reposFilename
	input.CommitsFile = *commitsFilename
	return input
}

// ProcessOrders function that process the files and prints out the top 10
// of commits, watched and actors
func ProcessOrders(actorFilename, eventFilename, repoFilename string) error {
	// open files
	actorFile, err := os.Open(actorFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read actor file content. Error: %s", err.Error())
		return fmt.Errorf("could not read actor file content. Error: %s", err.Error())
	}
	eventFile, err := os.Open(eventFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read event file content. Error: %s", err.Error())
		return fmt.Errorf("could not read event file content. Error: %s", err.Error())
	}
	repoFile, err := os.Open(repoFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read repo file content. Error: %s", err.Error())
		return fmt.Errorf("could not read repo file content. Error: %s", err.Error())
	}

	// get file data asynchronously
	actorData := make(chan map[string][]Actor)
	eventData := make(chan map[string][]Event)
	repoData := make(chan map[string][]Repo)

	go func() {
		eventData <- getEventData(eventFile)
	}()

	go func() {
		actorData <- getActorData(actorFile)
	}()

	go func() {
		repoData <- getRepoData(repoFile)
	}()

	actors := <-actorData
	repos := <-repoData
	events := <-eventData

	// filter data
	created := getFilteredActors(events["createevent"], events["pushevent"], actors)
	commits := getFilteredRepos(events["pushevent"], repos)
	watched := getFilteredRepos(events["watchevent"], repos)

	printActor("Top 10 active users", created)
	printRepo("Top 10 commits pushed", commits)
	printRepo("Top 10 events watched", watched)

	return nil
}

// printActor prints top 10 actors
func printActor(title string, actors []Actor) {
	fmt.Println()
	fmt.Println(title)
	fmt.Println("-----------------")
	for i := 0; i < 10; i ++ {
		fmt.Println("Username: ", actors[i].Username)
		fmt.Println("Created Count: ", actors[i].CreateCount)
		fmt.Println("Committed count: ", actors[i].CommitCount)
	}
}

// printRepo prints top 10 repos
func printRepo(title string, repos []Repo) {
	fmt.Println()
	fmt.Println(title)
	fmt.Println("-----------------")
	for i := 0; i < 10; i ++ {
		fmt.Println("Name: ", repos[i].Name)
		fmt.Println("Countount: ", repos[i].Count)
	}
}

// getFilteredActors get an array with filtered actors
func getFilteredActors(createEvent, pushEvent []Event, actors map[string][]Actor) []Actor {
	// filter first the created events
	topCreated := make(map[string]Actor)
	for _, event := range createEvent {
		created, ok := topCreated[event.ActorID]
		if !ok {
			created.CreateCount = 0
			created.ID = event.ActorID
			actor, _ := actors[event.ActorID]
			if actor != nil {
				created.Username = actor[0].Username
			}
		}
		created.CreateCount += 1
		topCreated[event.ActorID] = created
	}

	sortedCreated := make([]Actor, 0)
	for _, filteredCreated := range topCreated {
		sortedCreated = append(sortedCreated, filteredCreated)
	}

	// Sort by count, keeping original order or equal elements.
	sort.SliceStable(sortedCreated, func(i, j int) bool {
		return sortedCreated[i].CreateCount > sortedCreated[j].CreateCount
	})

	// filter then pushed events
	topPushed := make(map[string]Actor)
	for _, event := range pushEvent {
		pushed, ok := topPushed[event.ActorID]
		if !ok {
			pushed.CommitCount = 0
			pushed.ID = event.ActorID
			actor, _ := actors[event.ActorID]
			if actor != nil {
				pushed.Username = actor[0].Username
			}
		}
		pushed.CommitCount += 1
		topCreated[event.ActorID] = pushed
	}

	sortedPushed := make([]Actor, 0)
	for _, created := range sortedCreated {
		pushed, _ := topPushed[created.ID]
		created.CommitCount = pushed.CommitCount
		sortedPushed = append(sortedPushed, created)
	}

	// Sort by count, keeping original order or equal elements.
	sort.SliceStable(sortedPushed, func(i, j int) bool {
		return sortedPushed[i].CommitCount > sortedPushed[j].CommitCount &&
			sortedPushed[i].CreateCount > sortedPushed[j].CreateCount
	})

	return sortedPushed
}

// getFilteredRepos get an array with filtered repos
func getFilteredRepos(events []Event, repos map[string][]Repo) []Repo {
	topRepos := make(map[string]Repo)

	for _, event := range events {
		repo, ok := topRepos[event.RepoID]
		if !ok {
			repo.Count = 0
			repo.ID = event.RepoID
		}
		repo.Count += 1
		topRepos[event.RepoID] = repo
	}

	sortedRepos := make([]Repo, 0)

	for _, filteredRepo := range topRepos {
		repo, _ := repos[filteredRepo.ID]
		if repo != nil {
			filteredRepo.Name = repo[0].Name
		}
		sortedRepos = append(sortedRepos, filteredRepo)
	}

	// Sort by count, keeping original order or equal elements.
	sort.SliceStable(sortedRepos, func(i, j int) bool {
		return sortedRepos[i].Count > sortedRepos[j].Count
	})

	return sortedRepos
}

// getActorData gets data from file and return it as a go struct
func getActorData(actorFile *os.File) map[string][]Actor {
	reader := bufio.NewReader(actorFile)
	actors := make(map[string][]Actor)

	reader.ReadBytes('\n')
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "could not read file content. Error: %s", err.Error())
		}
		data := strings.Split(string(line), ",")
		if len(data) != 2 {
			fmt.Fprintf(os.Stderr, "data incorrectly formatted. %v", data)
		}
		actor := Actor{
			ID:       data[0],
			Username: strings.Replace(strings.ToLower(data[1]), "\n", "", -1),
		}

		uniqueActors, ok := actors[actor.ID]
		if !ok {
			uniqueActors = make([]Actor, 0)
		}
		uniqueActors = append(uniqueActors, actor)
		actors[actor.ID] = uniqueActors
	}
	return actors
}

// getEventData gets data from file and return it as a go struct
func getEventData(eventFile *os.File) map[string][]Event {
	reader := bufio.NewReader(eventFile)
	events := make(map[string][]Event)

	reader.ReadBytes('\n')
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "could not read file content. Error: %s", err.Error())
		}
		data := strings.Split(string(line), ",")

		// if data is not in the format expected, print out to std error
		if len(data) != 4 {
			fmt.Fprintf(os.Stderr, "data incorrectly formatted. %v", data)
		}
		event := Event{
			ID:      data[0],
			Type:    strings.ToLower(data[1]),
			ActorID: data[2],
			RepoID:  strings.Replace(strings.ToLower(data[3]), "\n", "", -1),
		}
		uniqueEvent, ok := events[event.Type]
		if !ok {
			uniqueEvent = make([]Event, 0)
		}
		uniqueEvent = append(uniqueEvent, event)
		events[event.Type] = uniqueEvent
	}
	return events
}

// getActorData gets data from file and return it as a go struct
func getRepoData(repoFile *os.File) map[string][]Repo {
	reader := bufio.NewReader(repoFile)
	repos := make(map[string][]Repo)

	reader.ReadBytes('\n')
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "could not read file content. Error: %s", err.Error())
		}
		data := strings.Split(string(line), ",")
		if len(data) != 2 {
			fmt.Fprintf(os.Stderr, "data incorrectly formatted. %v", data)
		}
		repo := Repo{
			ID:   data[0],
			Name: strings.Replace(strings.ToLower(data[1]), "\n", "", -1),
		}
		uniqueRepo, ok := repos[repo.ID]
		if !ok {
			uniqueRepo = make([]Repo, 0)
		}
		uniqueRepo = append(uniqueRepo, repo)
		repos[repo.ID] = uniqueRepo
	}
	return repos
}
