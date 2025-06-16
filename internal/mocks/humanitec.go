package mocks

import (
	"github.com/lil-yellow-flower/humctl-wrapper-demo/internal/humanitec"
)

// MockHumanitecClient is a mock implementation of the Humanitec client
type MockHumanitecClient struct {
	CreateAppFunc    func(name string, skipEnvCreation bool) (*humanitec.App, error)
	DeleteAppFunc    func(name string) error
	GetAppFunc       func(name string) (*humanitec.App, error)
	GetAppsFunc      func() ([]humanitec.App, error)
	UpdateAppFunc    func(name, newName string) (*humanitec.App, error)
}

func (m *MockHumanitecClient) CreateApp(name string, skipEnvCreation bool) (*humanitec.App, error) {
	if m.CreateAppFunc != nil {
		return m.CreateAppFunc(name, skipEnvCreation)
	}
	return nil, nil
}

func (m *MockHumanitecClient) DeleteApp(name string) error {
	if m.DeleteAppFunc != nil {
		return m.DeleteAppFunc(name)
	}
	return nil
}

func (m *MockHumanitecClient) GetApp(name string) (*humanitec.App, error) {
	if m.GetAppFunc != nil {
		return m.GetAppFunc(name)
	}
	return nil, nil
}

func (m *MockHumanitecClient) GetApps() ([]humanitec.App, error) {
	if m.GetAppsFunc != nil {
		return m.GetAppsFunc()
	}
	return nil, nil
}

func (m *MockHumanitecClient) UpdateApp(name, newName string) (*humanitec.App, error) {
	if m.UpdateAppFunc != nil {
		return m.UpdateAppFunc(name, newName)
	}
	return nil, nil
} 