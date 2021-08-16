package vercel_client

type errorResponse struct {
	Error *errorContent `json:"error"`
}

type errorContent struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type GetProjectsResponse struct {
	Project []*Project `json:"project"`
}

type Project struct {
	Id                          string          `json:"id"`
	Name                        string          `json:"name"`
	Framework                   string          `json:"framework"`
	RootDirectory               string          `json:"rootDirectory"`
	NodeVersion                 string          `json:"nodeVersion"`
	AccountId                   string          `json:"accountId"`
	UpdatedAt                   int64           `json:"updatedAt"`
	CreatedAt                   int64           `json:"createdAt"`
	Alias                       []*Domain       `json:"alias"`
	Link                        *RepositoryLink `json:"link"`
	BuildCommand                string          `json:"buildCommand"`
	OutputDirectory             string          `json:"outputDirectory"`
	CommandForIgnoringBuildStep string          `json:"commandForIgnoringBuildStep,omitempty"`
}

type RepositoryLink struct {
	Type string `json:"type"`
	Repo string `json:"repo"`
	Org  string `json:"org"`
}

type CreateProjectOptions struct {
	Name                        string
	Framework                   string
	RepositoryType              string
	RepositoryName              string
	RootDirectory               string
	BuildCommand                string
	OutputDirectory             string
	CommandForIgnoringBuildStep string
}

type CreateProjectRequest struct {
	Name                        string                `json:"name"`
	Framework                   string                `json:"framework"`
	RootDirectory               *string               `json:"rootDirectory"`
	GitRepository               *GitRepositoryRequest `json:"gitRepository,omitempty"`
	BuildCommand                string                `json:"buildCommand"`
	OutputDirectory             string                `json:"outputDirectory"`
	CommandForIgnoringBuildStep string                `json:"commandForIgnoringBuildStep,omitempty"`
}

type GitRepositoryRequest struct {
	Type string `json:"type"`
	Repo string `json:"repo"`
}

type UpdateProjectRequest struct {
	Framework                   string  `json:"framework"`
	RootDirectory               *string `json:"rootDirectory"`
	BuildCommand                string  `json:"buildCommand"`
	OutputDirectory             string  `json:"outputDirectory"`
	CommandForIgnoringBuildStep string  `json:"commandForIgnoringBuildStep"`
}

type CreateProjectEnvRequest struct {
	Type   string   `json:"type"`
	Key    string   `json:"key"`
	Value  string   `json:"value"`
	Target []string `json:"target"`
}

type GetProjectEnvsResponse struct {
	Envs []*ProjectEnv `json:"envs"`
}

type ProjectEnv struct {
	Id     string   `json:"id"`
	Type   string   `json:"type"`
	Key    string   `json:"key"`
	Value  string   `json:"value"`
	Target []string `json:"target"`
}

type Domain struct {
	Domain   string `json:"domain"`
	Redirect string `json:"redirect"`
}

type CreateDomainRequest struct {
	Domain   string `json:"domain"`
	Redirect string `json:"redirect"`
}
