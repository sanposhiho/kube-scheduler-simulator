package webhook

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"

	"sigs.k8s.io/kube-scheduler-simulator/scenario/waitermanager"

	definederr "sigs.k8s.io/kube-scheduler-simulator/scenario/errors"

	"sigs.k8s.io/kube-scheduler-simulator/scenario/utils"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kube-scheduler-simulator/scenario/api/v1alpha1"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/xerrors"
)

// AdmissionWebhookServer is server for simulator.
type AdmissionWebhookServer struct {
	e       *echo.Echo
	client  client.Client
	manager waitermanager.Manager
}

// NewSimulatorServer initialize AdmissionWebhookServer.
func NewSimulatorServer() *AdmissionWebhookServer {
	e := echo.New()
	e.Use(middleware.Logger())

	// initialize AdmissionWebhookServer.
	s := &AdmissionWebhookServer{e: e}
	s.e.Logger.SetLevel(log.INFO)

	s.e.GET("/admissionwebhook/validation", s.ValidationHandler)

	return s
}

// Start starts AdmissionWebhookServer.
func (s *AdmissionWebhookServer) Start(port int) (
	func(), // function for shutdown
	error,
) {
	e := s.e

	go func() {
		if err := e.Start(":" + strconv.Itoa(port)); err != nil && !xerrors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatalf("failed to start admission webhook server successfully: %v", err)
		}
	}()

	shutdownFn := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Warnf("failed to shutdown admission webhook server successfully: %v", err)
		}
	}

	return shutdownFn, nil
}

func (s *AdmissionWebhookServer) ValidationHandler(c echo.Context) error {
	ctx := c.Request().Context()
	running, err := utils.FetchRunningScenario(ctx, s.client)
	if err != nil {
		if errors.Is(err, definederr.ErrNoRunningScenario) {
			return nil
		}
		return err
	}

	if running.Status.StepStatus.Phase != v1alpha1.StepPhaseOperating && running.Status.StepStatus.Phase != v1alpha1.StepPhaseControllerRunning {
		return nil
	}

	controllers := sets.NewString()
	for _, enabled := range running.Spec.Controllers.PreparingControllers.Enabled {
		controllers.Insert(enabled.Name)
	}
	go func() {
		// TODO: set timeout
		ctx := context.Background()

		done, err := s.manager.Run(ctx, controllers)
		if err != nil {
			// TODO: log error
			// TODO: change scenario status to fail?
			return
		}
		if !done {
			return
		}

		running, err := utils.FetchRunningScenario(ctx, s.client)
		if err != nil {
			// TODO: log error
			// TODO: change scenario status to fail?
			return
		}

		running.Status.StepStatus.Phase = v1alpha1.StepPhaseOperatingCompleted

		if err := s.client.Update(ctx, running, nil); err != nil {
			// TODO: log error
			// TODO: change scenario status to fail?
			return
		}
	}()

	if running.Status.StepStatus.Phase == v1alpha1.StepPhaseControllerRunning {
		// The running simulated controller will be stopped.
		running.Status.StepStatus.Phase = v1alpha1.StepPhaseControllerPaused
		running.Status.StepStatus.Step.Minor++
		if err := s.client.Update(ctx, running, nil); err != nil {
			return err
		}
	}

	return nil
}
