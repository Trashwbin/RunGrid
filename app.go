package main

import (
	"context"

	"rungrid/backend/domain"
	"rungrid/backend/service"
	"rungrid/backend/storage"
	"rungrid/backend/storage/memory"
)

// App struct
type App struct {
	ctx     context.Context
	items   *service.ItemService
	groups  *service.GroupService
	closeFn func() error
}

// NewApp creates a new App application struct
func NewApp() *App {
	itemRepo := memory.NewItemRepository()
	groupRepo := memory.NewGroupRepository()

	return &App{
		items:  service.NewItemService(itemRepo),
		groups: service.NewGroupService(groupRepo),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// shutdown is called when the app is terminating.
func (a *App) shutdown(ctx context.Context) {
	if a.closeFn != nil {
		_ = a.closeFn()
	}
}

func (a *App) ListItems(groupID string, query string) ([]domain.Item, error) {
	return a.items.List(a.context(), storage.ItemFilter{GroupID: groupID, Query: query})
}

func (a *App) CreateItem(input domain.ItemInput) (domain.Item, error) {
	return a.items.Create(a.context(), input)
}

func (a *App) UpdateItem(input domain.ItemUpdate) (domain.Item, error) {
	return a.items.Update(a.context(), input)
}

func (a *App) DeleteItem(id string) error {
	return a.items.Delete(a.context(), id)
}

func (a *App) RecordLaunch(id string) (domain.Item, error) {
	return a.items.RecordLaunch(a.context(), id)
}

func (a *App) ListGroups() ([]domain.Group, error) {
	return a.groups.List(a.context())
}

func (a *App) CreateGroup(input domain.GroupInput) (domain.Group, error) {
	return a.groups.Create(a.context(), input)
}

func (a *App) UpdateGroup(input domain.Group) (domain.Group, error) {
	return a.groups.Update(a.context(), input)
}

func (a *App) DeleteGroup(id string) error {
	return a.groups.Delete(a.context(), id)
}

func (a *App) context() context.Context {
	if a.ctx != nil {
		return a.ctx
	}
	return context.Background()
}
