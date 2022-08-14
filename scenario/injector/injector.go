package injector

import (
	"context"
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/fields"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kube-scheduler-simulator/scenario/api/v1alpha1"
)

var (
	clientCache client.Client

	ErrTooManyRunningScenario = errors.New("too many running scenario")
)

func CheckScenarioStepPhase(controllerName string, config *rest.Config) error {
	ctx := context.Background()
	cli, err := setUpClientOrDie(config)
	if err != nil {
		return err
	}

	list := v1alpha1.ScenarioList{}
	if err := cli.List(ctx, &list, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("status.phase", string(v1alpha1.ScenarioRunning)),
	}); err != nil {
		return err
	}

	if len(list.Items) > 1 {
		return ErrTooManyRunningScenario
	}
	if len(list.Items) == 0 {
		return nil
	}

	running := list.Items[0]

	for {
		// TODO: fix busy wait
		if running.Status.StepStatus.Phase != v1alpha1.ControllerRunning && running.Status.StepStatus.RunningSimulatedController != controllerName {
			continue
		}
		break
	}

	return nil
}

func setUpClientOrDie(config *rest.Config) (client.Client, error) {
	if clientCache != nil {
		return clientCache, nil
	}

	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("initialize discovery client: %w", err)
	}

	scheme, err := v1alpha1.SchemeBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("build scheme builder: %w", err)
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	cli, err := client.New(config, client.Options{
		Scheme: scheme,
		Mapper: mapper,
		Opts:   client.WarningHandlerOptions{},
	})
	if err != nil {
		return nil, fmt.Errorf("create a new client: %w", err)
	}

	clientCache = cli
	return cli, nil
}
