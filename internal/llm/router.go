package llm

import (
	"agent-coach/internal/storage"
	"context"
	"fmt"
	"sync"

	"github.com/google/martian/log"
)

type Router struct {
	db              *storage.DB
	providers       map[string]Provider
	defaultProvider Provider
	mu              sync.RWMutex
}

func NewRouter(db *storage.DB) (*Router, error) {
	r := &Router{
		db:              db,
		providers:       make(map[string]Provider),
		defaultProvider: nil,
	}
	err := r.loadProvidersFromDB(context.Background())
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Router) SaveProviderConfig(ctx context.Context, config *LLMProviderConfig) error {
	query := `
		INSERT INTO llm_providers (name, provider, base_url, api_key, default_model, is_default, is_active)
		VALUES (:name, :provider, :base_url, :api_key, :default_model, :is_default, :is_active)
	`
	_, err := r.db.NamedExecContext(ctx, query, config)
	if err != nil {
		return err
	}

	return r.refreshProviders(ctx)
}

func (r *Router) Complete(ctx context.Context, req CompletionRequest, providers ...string) (*CompletionResponse, error) {
	provider := r.resolveProvider(providers...)
	if provider == nil {
		return nil, fmt.Errorf("no provider found")
	}
	return provider.Complete(ctx, req)
}

func (r *Router) loadProvidersFromDB(ctx context.Context) error {
	configs, err := r.getProviderConfigs(ctx)
	if err != nil {
		return err
	}
	for _, config := range configs {
		if !config.IsActive {
			continue
		}
		provider, err := r.createProviderFromConfig(config)
		if err != nil {
			log.Errorf("failed to create provider from config: %v", err)
			continue
		}
		r.providers[config.Name] = provider
		if config.IsDefault {
			r.defaultProvider = provider
		}
	}

	return nil
}

func (r *Router) getProviderConfigs(ctx context.Context) ([]*LLMProviderConfig, error) {
	query := `
		SELECT * FROM llm_providers ORDER BY is_default DESC
	`
	var configs []*LLMProviderConfig
	if err := r.db.SelectContext(ctx, &configs, query); err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *Router) createProviderFromConfig(config *LLMProviderConfig) (Provider, error) {
	switch config.Provider {
	case "openrouter":
		return NewOpenRouterProvider(config), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", config.Provider)
	}
}

func (r *Router) refreshProviders(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers = make(map[string]Provider)
	r.defaultProvider = nil
	return r.loadProvidersFromDB(ctx)
}

func (r *Router) resolveProvider(providers ...string) Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, name := range providers {
		if provider, ok := r.providers[name]; ok {
			if !provider.IsAvailable() {
				continue
			}
			return provider
		}
	}
	return r.defaultProvider
}
