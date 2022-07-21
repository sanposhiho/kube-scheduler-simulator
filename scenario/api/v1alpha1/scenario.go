package v1alpha1

import (
	"errors"

	runtimeschema "k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"

	"golang.org/x/xerrors"
)

//// Run runs the event.
//// It returns the func to update ScenarioResult and error.
//func (e *ScenarioOperation) Run(ctx context.Context, cfg *rest.Config) (StatusUpdateFn, error) {
//	// TODO: validation webhook will reject when there are multiple non-nil operations in single Event.
//	switch {
//	case e.CreateOperation != nil:
//		ope := e.CreateOperation
//		gvk := ope.Object.GetObjectKind().GroupVersionKind()
//		client, err := buildClient(gvk, cfg)
//		_, err = client.Create(ctx, ope.Object, ope.CreateOptions)
//		if err != nil {
//			return nil, xerrors.Errorf("run create operation: id: %s error: %w", e.ID, err)
//		}
//	case e.PatchOperation != nil:
//	case e.DeleteOperation != nil:
//	case e.DoneOperation != nil:
//		return e.DoneOperation.run(e.Step)
//	}
//
//	return true, nil
//}

func (o *DoneOperation) run(id string, step ScenarioStep) (func(status *ScenarioStatus), error) {
	return func(status *ScenarioStatus) {
		status.Phase = ScenarioSucceeded
		status.ScenarioResult.Timeline[step] = append(status.ScenarioResult.Timeline[step], ScenarioTimelineEvent{
			ID:             id,
			Step:           step,
			Create:         nil,
			Patch:          nil,
			Delete:         nil,
			Done:           nil,
			PodScheduled:   nil,
			PodUnscheduled: nil,
			PodPreempted:   nil,
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
