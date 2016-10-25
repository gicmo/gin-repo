package wire

type RepoAccessQuery struct {
	User string
	Path string
}

type RepoAccessInfo struct {
	Path string
	Push bool
}

type CreateRepo struct {
	Name        string
	Description string
}

type Repo struct {
	Name        string
	Owner       string
	Description string
	Public      bool
	Head        string
}

type Branch struct {
	Name   string
	Commit string
}
