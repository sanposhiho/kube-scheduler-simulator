package workermanager

import (
	"errors"
	"fmt"
	"sync"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kube-scheduler-simulator/scenario/api/v1alpha1"
	"sigs.k8s.io/kube-scheduler-simulator/scenario/waitermanager"
	"sigs.k8s.io/kube-scheduler-simulator/scenario/worker"
)

type Manager struct {
	mu            sync.RWMutex
	queue         []string
	workers       map[string]*worker.ScenarioWorker
	client        client.Client
	waiterManager *waitermanager.Manager
}

func New(cli client.Client, wm *waitermanager.Manager) *Manager {
	// TODO: add exist scenario
	return &Manager{
		workers:       make(map[string]*worker.ScenarioWorker),
		client:        cli,
		waiterManager: wm,
	}
}

func (m *Manager) AddOrUpdateScenario(s *v1alpha1.Scenario) {
	if w, ok := m.workers[s.Name]; ok {
		w.HandleUpdate(s)
		return
	}

	m.workers[s.Name] = worker.New(s, m.client, m.waiterManager)
	m.queue = append(m.queue, s.Name)
}

func (m *Manager) RunNext() error {
	next, err := m.pop()
	if err != nil {
		return fmt.Errorf("pop next scenario worker from queue: %w", err)
	}

	go next.Run()

	return nil
}

func (m *Manager) pop() (*worker.ScenarioWorker, error) {
	if len(m.queue) == 0 {
		return nil, nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	first := m.queue[0]
	m.queue = m.queue[1:]
	w, ok := m.workers[first]
	if !ok {
		return nil, errWorkerNotFound
	}

	return w, nil
}

var errWorkerNotFound = errors.New("worker not found")
