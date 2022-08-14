package worker

import (
	"context"

	"golang.org/x/xerrors"
	"sigs.k8s.io/kube-scheduler-simulator/scenario/utils"

	"k8s.io/client-go/rest"

	simulationv1alpha1 "sigs.k8s.io/kube-scheduler-simulator/scenario/api/v1alpha1"
)

// stepper will run all the operations defined in the single step.
type stepper struct {
	config *rest.Config
	// step is major step
	step       int32
	operations []*simulationv1alpha1.ScenarioOperation
}

// run runs all operations registered in s.step.
// It returns boolean shows whether the Scenario should finish in this step.
func (s *stepper) run(ctx context.Context) (bool, error) {
	finish, err := s.operate(ctx)
	if err != nil {
		return true, xerrors.Errorf("run .spec.operation: %w", err)
	}

	s.wait()

	return finish, nil
}

func (s *stepper) wait() error {
	// TODO: implement
	return nil
}

func (s *stepper) operate(ctx context.Context) (bool, error) {
	eg := utils.NewErrGroupWithSemaphore(ctx)

	// whether the Scenario should finish in this step.
	finishFlag := false
	for _, operation := range s.operations {
		if err := eg.Sem.Acquire(ctx, 1); err != nil {
			return false, xerrors.Errorf("acquire semaphore: %w", err)
		}
		operation := operation
		eg.Grp.Go(func() error {
			finish, err := operation.Run(ctx, s.config)
			if err != nil {
				return xerrors.Errorf("run operation: operation id: %v, step %v, error: %w", operation.ID, operation.MajorStep, err)
			}
			if finish {
				// update flag
				finishFlag = finish
			}
			return nil
		})
	}

	if err := eg.Grp.Wait(); err != nil {
		return true, err
	}

	return finishFlag, nil
}
