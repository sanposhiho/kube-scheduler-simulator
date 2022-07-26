package v1alpha1

import (
	"context"
	"errors"

	runtimeschema "k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"

	"golang.org/x/xerrors"
)

// Run runs the event.
// It returns boolean shows whether the Scenario should finish in this step.
func (s *ScenarioOperation) Run(ctx context.Context, cfg *rest.Config) (bool, error) {
	// TODO: validation webhook will reject when there are multiple non-nil operations in single Event.
	switch {
	case s.Create != nil:
		ope := s.Create
		gvk := ope.Object.GetObjectKind().GroupVersionKind()
		client, err := buildClient(gvk, cfg)
		_, err = client.Create(ctx, ope.Object, ope.CreateOptions)
		if err != nil {
			return true, xerrors.Errorf("run create operation: id: %s error: %w", s.ID, err)
		}
	case s.Patch != nil:
		ope := s.Patch
		gvk := ope.TypeMeta.GroupVersionKind()
		client, err := buildClient(gvk, cfg)
		_, err = client.Patch(ctx, ope.ObjectMeta.Name, ope.PatchType, []byte(ope.Patch), ope.PatchOptions)
		if err != nil {
			return true, xerrors.Errorf("run create operation: id: %s error: %w", s.ID, err)
		}
	case s.Delete != nil:
		ope := s.Delete
		gvk := ope.TypeMeta.GroupVersionKind()
		client, err := buildClient(gvk, cfg)
		err = client.Delete(ctx, ope.ObjectMeta.Name, ope.DeleteOptions)
		if err != nil {
			return true, xerrors.Errorf("run create operation: id: %s error: %w", s.ID, err)
		}
	case s.Done != nil:
		return true, nil
	}

	return true, ErrUnknownOperation
}

func (o *DoneOperation) run(id string, step ScenarioStep) (func(status *ScenarioStatus), error) {
	return func(status *ScenarioStatus) {
		status.Phase = ScenarioSucceeded
		status.ScenarioResult.Timeline[step] = append(status.ScenarioResult.Timeline[step], ScenarioTimelineEvent{
			ID:   id,
			Step: step,
		})
	}, nil
}

var ErrUnknownOperation = errors.New("")

func buildClient(gvk runtimeschema.GroupVersionKind, cfg *rest.Config) (dynamic.NamespaceableResourceInterface, error) {
	cli, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, xerrors.Errorf("build dynamic client: %w", err)
	}

	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, xerrors.Errorf("build discovery client: %w", err)
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, xerrors.Errorf("build mapping from RESTMapper: %w", err)
	}

	return cli.Resource(mapping.Resource), nil
}
