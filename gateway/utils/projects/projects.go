package projects

import (
	"errors"
	"sync"

	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/auth"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud"
	"github.com/spaceuptech/space-cloud/gateway/modules/crud/driver"
	"github.com/spaceuptech/space-cloud/gateway/modules/eventing"
	"github.com/spaceuptech/space-cloud/gateway/modules/filestore"
	"github.com/spaceuptech/space-cloud/gateway/modules/functions"
	"github.com/spaceuptech/space-cloud/gateway/modules/realtime"
	"github.com/spaceuptech/space-cloud/gateway/modules/schema"
	"github.com/spaceuptech/space-cloud/gateway/modules/userman"
	"github.com/spaceuptech/space-cloud/gateway/utils/admin"
	"github.com/spaceuptech/space-cloud/gateway/utils/graphql"
	"github.com/spaceuptech/space-cloud/gateway/utils/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/utils/metrics"
	"github.com/spaceuptech/space-cloud/gateway/utils/routing"
	"github.com/spaceuptech/space-cloud/gateway/utils/syncman"
)

// ProjectState holds the module state of a project
type ProjectState struct {
	//Config         *config.Project
	Auth           *auth.Module
	Crud           *crud.Module
	Schema         *schema.Schema
	UserManagement *userman.Module
	FileStore      *filestore.Module
	Functions      *functions.Module
	Realtime       *realtime.Module
	Eventing       *eventing.Module
	Graph          *graphql.Module
}

// Projects is the stub to manage the state of the various modules
type Projects struct {
	lock               sync.RWMutex
	nodeID             string
	removeProjectScope bool
	projects           map[string]*ProjectState
	h                  *driver.Handler

	// Global managers
	syncMan     *syncman.Manager
	adminMan    *admin.Manager
	metrics     *metrics.Module
	letsencrypt *letsencrypt.LetsEncrypt
	routing     *routing.Routing
}

// New creates a new Projects instance
func New(nodeID string, removeProjectScope bool, h *driver.Handler,
	adminMan *admin.Manager, syncMan *syncman.Manager, metrics *metrics.Module,
	le *letsencrypt.LetsEncrypt, r *routing.Routing) *Projects {
	return &Projects{nodeID: nodeID, removeProjectScope: removeProjectScope, projects: map[string]*ProjectState{}, h: h,
		syncMan: syncMan, adminMan: adminMan, metrics: metrics, letsencrypt: le, routing: r}
}

// LoadProject returns the state of the project specified
func (p *Projects) LoadProject(project string) (*ProjectState, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if state, p := p.projects[project]; p {
		return state, nil
	}

	return nil, errors.New("project not found in server state")
}

// DeleteProject deletes a single project
func (p *Projects) DeleteProject(project string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	// Remove config from global modules
	_ = p.letsencrypt.DeleteProjectDomains(project)
	p.routing.DeleteProjectRoutes(project)

	delete(p.projects, project)
}

// Iter iterates over all the projects and passes it in the provided function.
// Iteration stops if the function returns false
func (p *Projects) Iter(fn func(string, *ProjectState) bool) bool {
	p.lock.RLock()
	defer p.lock.RUnlock()

	for project, state := range p.projects {
		if !fn(project, state) {
			return false
		}
	}

	return true
}

// NewProject creates a new project with all modules in the default state.
// It will overwrite the existing project if any
func (p *Projects) NewProject(project string) (*ProjectState, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	// Create the fundamental modules
	c := crud.Init(p.h, p.adminMan)
	s := schema.Init(c, p.removeProjectScope)

	a := auth.Init(p.nodeID, c, s, p.removeProjectScope)
	a.SetMakeHttpRequest(p.syncMan.MakeHTTPRequest)

	fn := functions.Init(a, p.syncMan)

	f := filestore.Init(a)

	// Initialise the eventing module and set the crud module hooks
	e := eventing.New(a, c, s, fn, p.adminMan, p.syncMan, f)
	f.SetEventingModule(e)

	// Set hooks
	c.SetHooks(&model.CrudHooks{
		Create: e.HookDBCreateIntent,
		Update: e.HookDBUpdateIntent,
		Delete: e.HookDBDeleteIntent,
		Batch:  e.HookDBBatchIntent,
		Stage:  e.HookStage,
	}, p.metrics.AddDBOperation)

	rt, err := realtime.Init(p.nodeID, e, a, c, s, p.metrics, p.syncMan)
	if err != nil {
		return nil, err
	}

	u := userman.Init(c, a)
	graph := graphql.New(a, c, fn, s)

	state := &ProjectState{Crud: c, Schema: s, Functions: fn, Auth: a, UserManagement: u, FileStore: f, Realtime: rt,
		Eventing: e, Graph: graph}

	p.projects[project] = state

	return state, nil
}

// GetProjectIDs returns an array of project ids present in the project configuration
func (p *Projects) GetProjectIDs() []string {
	p.lock.RLock()
	defer p.lock.RUnlock()

	var ids []string
	for id := range p.projects {
		ids = append(ids, id)
	}

	return ids
}