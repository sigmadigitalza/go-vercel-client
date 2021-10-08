package vercel_client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ProjectApi struct {
	c       *http.Client
	baseUrl *url.URL
}

// The Project CRUD

func (p *ProjectApi) CreateProject(ctx context.Context, options *CreateProjectOptions) (*Project, error) {
	rel := &url.URL{Path: "/v6/projects"}
	u := p.baseUrl.ResolveReference(rel)

	body := &CreateProjectRequest{
		Name:                        options.Name,
		BuildCommand:                options.BuildCommand,
		OutputDirectory:             options.OutputDirectory,
		CommandForIgnoringBuildStep: options.CommandForIgnoringBuildStep,
	}

	if options.Framework != "" {
		body.Framework = &options.Framework
	}

	if options.RepositoryType != "" && options.RepositoryName != "" {
		body.GitRepository = &GitRepositoryRequest{
			Type: options.RepositoryType,
			Repo: options.RepositoryName,
		}
	}

	if options.RootDirectory != "" {
		body.RootDirectory = &options.RootDirectory
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	resp, err := p.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, formatError(resp.Body)
	}

	return newProject(resp.Body)
}

func (p *ProjectApi) GetProject(ctx context.Context, name string) (*Project, error) {
	rel := &url.URL{Path: fmt.Sprintf("/v1/projects/%s", name)}
	u := p.baseUrl.ResolveReference(rel)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, formatError(resp.Body)
	}

	return newProject(resp.Body)
}

func (p *ProjectApi) UpdateProject(ctx context.Context, name string, project *Project) (*Project, error) {
	rel := &url.URL{Path: fmt.Sprintf("/v1/projects/%s", name)}
	u := p.baseUrl.ResolveReference(rel)

	body := &UpdateProjectRequest{
		BuildCommand:                nil,
		OutputDirectory:             nil,
		CommandForIgnoringBuildStep: project.CommandForIgnoringBuildStep,
	}

	if project.Framework != "" {
		body.Framework = &project.Framework
	}

	if project.BuildCommand != "" {
		body.BuildCommand = &project.BuildCommand
	}

	if project.OutputDirectory != "" {
		body.OutputDirectory = &project.OutputDirectory
	}

	if project.RootDirectory != "" {
		body.RootDirectory = &project.RootDirectory
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, u.String(), bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	resp, err := p.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, formatError(resp.Body)
	}

	return newProject(resp.Body)
}

func (p *ProjectApi) DeleteProject(ctx context.Context, name string) error {
	rel := &url.URL{Path: fmt.Sprintf("/v1/projects/%s", name)}
	u := p.baseUrl.ResolveReference(rel)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := p.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return formatError(resp.Body)
	}

	return nil
}

// The Project ENV CRUD

func (p *ProjectApi) CreateProjectEnv(ctx context.Context, id string, envType string, key string, value string, target []string) (*ProjectEnv, error) {
	rel := &url.URL{Path: fmt.Sprintf("/v7/projects/%s/env", id)}
	u := p.baseUrl.ResolveReference(rel)

	body := &CreateProjectEnvRequest{
		Type:   envType,
		Key:    key,
		Value:  value,
		Target: target,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	resp, err := p.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, formatError(resp.Body)
	}

	return newProjectEnv(resp.Body)
}

func (p *ProjectApi) GetProjectEnvs(ctx context.Context, id string, decrypt bool) ([]*ProjectEnv, error) {
	rel := &url.URL{Path: fmt.Sprintf("/v7/projects/%s/env", id)}
	u := p.baseUrl.ResolveReference(rel)

	q := u.Query()
	q.Add("decrypt", fmt.Sprintf("%t", decrypt))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, formatError(resp.Body)
	}

	var getProjectEnvsResponse GetProjectEnvsResponse
	err = json.NewDecoder(resp.Body).Decode(&getProjectEnvsResponse)
	if err != nil {
		return nil, err
	}

	return getProjectEnvsResponse.Envs, nil
}

func (p *ProjectApi) EditProjectEnv(ctx context.Context, id string, env *ProjectEnv) (*ProjectEnv, error) {
	rel := &url.URL{Path: fmt.Sprintf("/v7/projects/%s/env/%s", id, env.Id)}
	u := p.baseUrl.ResolveReference(rel)

	body := &CreateProjectEnvRequest{
		Type:   env.Type,
		Key:    env.Key,
		Value:  env.Value,
		Target: env.Target,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, u.String(), bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	resp, err := p.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, formatError(resp.Body)
	}

	return newProjectEnv(resp.Body)
}

func (p *ProjectApi) DeleteProjectEnv(ctx context.Context, id string, envId string) error {
	rel := &url.URL{Path: fmt.Sprintf("/v7/projects/%s/env/%s", id, envId)}
	u := p.baseUrl.ResolveReference(rel)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := p.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return formatError(resp.Body)
	}

	return nil
}

// The Project Domain CRUD

func (p *ProjectApi) AddDomain(ctx context.Context, id string, domain string, redirect string) ([]*Domain, error) {
	rel := &url.URL{Path: fmt.Sprintf("/v1/projects/%s/alias", id)}
	u := p.baseUrl.ResolveReference(rel)

	body := &CreateDomainRequest{
		Domain:   domain,
		Redirect: redirect,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	resp, err := p.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, formatError(resp.Body)
	}

	var domains []*Domain
	err = json.NewDecoder(resp.Body).Decode(&domains)
	if err != nil {
		return nil, err
	}

	return domains, nil
}

func (p *ProjectApi) UpdateDomain(ctx context.Context, id string, domain string, redirect string) ([]*Domain, error) {
	rel := &url.URL{Path: fmt.Sprintf("/v1/projects/%s/alias", id)}
	u := p.baseUrl.ResolveReference(rel)

	body := &CreateDomainRequest{
		Domain:   domain,
		Redirect: redirect,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, u.String(), bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	resp, err := p.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, formatError(resp.Body)
	}

	var domains []*Domain
	err = json.NewDecoder(resp.Body).Decode(&domains)
	if err != nil {
		return nil, err
	}

	return domains, nil
}

func (p *ProjectApi) DeleteDomain(ctx context.Context, id string, domain string) ([]*Domain, error) {
	rel := &url.URL{Path: fmt.Sprintf("/v1/projects/%s/alias", id)}
	u := p.baseUrl.ResolveReference(rel)

	q := u.Query()
	q.Add("domain", domain)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, formatError(resp.Body)
	}

	var domains []*Domain
	err = json.NewDecoder(resp.Body).Decode(&domains)
	if err != nil {
		return nil, err
	}

	return domains, nil
}

// Helper functions

func newProject(r io.Reader) (*Project, error) {
	var project Project
	err := json.NewDecoder(r).Decode(&project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func newProjectEnv(r io.Reader) (*ProjectEnv, error) {
	var projectEnv ProjectEnv
	err := json.NewDecoder(r).Decode(&projectEnv)
	if err != nil {
		return nil, err
	}

	return &projectEnv, nil
}

func formatError(r io.Reader) error {
	var e errorResponse
	err := json.NewDecoder(r).Decode(&e)
	if err != nil {
		return err
	}

	return errors.New(e.Error.Code)
}
