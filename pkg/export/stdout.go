package export

import (
	"context"

	"github.com/RayHuangCN/Jarvis/pkg/evaluate"

	"github.com/RayHuangCN/Jarvis/pkg/diagnose"
	"github.com/RayHuangCN/Jarvis/pkg/logger"
)

type Stdout struct {
	Logger logger.Logger
}

func NewStdout() *Stdout {
	return &Stdout{
		Logger: logger.NewLogger(),
	}
}

func (s *Stdout) CoordinateBegin(ctx context.Context) error {
	s.Logger.Infof("--------------- Jarvis [1.0] -----------------")
	return nil
}

func (s *Stdout) DiagnosticBegin(ctx context.Context, dia diagnose.Diagnostic) error {
	s.Logger.Infof("[Diagnostic begin] %s", dia.Name())
	return nil
}

func (s *Stdout) DiagnosticFinish(ctx context.Context, dia diagnose.Diagnostic) error {
	s.Logger.Infof("[Diagnostic finish] %s", dia.Name())
	return nil
}

func (s *Stdout) DiagnosticResult(ctx context.Context, result *diagnose.Result) error {
	if result.Error != nil {
		s.Logger.Errorf(result.Error.Error())
	} else {
		s.Logger.Infof("[%s] [%s] [%s] [%s] [%s]", result.Level, result.Name, result.ObjName, result.Desc, result.Proposal)
	}
	return nil
}

func (s *Stdout) EvaluationBegin(ctx context.Context, eva evaluate.Evaluator) error {
	s.Logger.Infof("[Evaluate begin] %s", eva.Name())
	return nil
}

func (s *Stdout) EvaluationFinish(ctx context.Context, eva evaluate.Evaluator) error {
	s.Logger.Infof("[Evaluate finish] %s", eva.Name())
	return nil
}

func (s *Stdout) EvaluationResult(ctx context.Context, result *evaluate.Result) error {
	s.Logger.Infof("[%s] [%s]", result.Name, result.Desc)
	return nil
}

func (s *Stdout) CoordinateFinish(ctx context.Context) error {
	s.Logger.Infof("--------------- Finish -----------------")
	return nil
}
