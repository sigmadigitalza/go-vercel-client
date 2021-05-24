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
	Id          string `json:"id"`
	Name        string `json:"name"`
	Framework   string `json:"framework"`
	NodeVersion string `json:"nodeVersion"`
	AccountId   string `json:"accountId"`
	UpdatedAt   int64  `json:"updatedAt"`
	CreatedAt   int64  `json:"createdAt"`
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Framework   string `json:"framework"`
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