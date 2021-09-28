package files

type ProcessInput struct {
	ActorsFile  string
	EventsFile  string
	CommitsFile string
	ReposFile   string
}

type Actor struct {
	ID          string
	Username    string
	CreateCount int
	CommitCount int
}

type Event struct {
	ID      string
	Type    string
	ActorID string
	RepoID  string
}

type Repo struct {
	ID    string
	Name  string
	Count int
}
