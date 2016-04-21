package rocket

type PushData struct {
	Repository    string `json:"repository"`
	RepositoryUrl string `json:"repositoryurl"`
	Message       string `json:"message"`
	CommitUrl     string `json:"commiturl"`
	UserName      string `json:"username"`
}
