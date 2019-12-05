package compexplorer

import (
	"fmt"
	"github.com/RayHuangCN/kube-jarvis/pkg/logger"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/cluster"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/cluster/custom/nodeexec"
	"golang.org/x/sync/errgroup"
	"strings"
	"sync"
)

// Bare get component information from cmd
type Bare struct {
	logger       logger.Logger
	cmdName      string
	nodes        []string
	nodeExecutor nodeexec.Executor
}

// NewBare create and int a StaticPods ComponentExecutor
func NewBare(logger logger.Logger, cmdName string, nodes []string, executor nodeexec.Executor) *Bare {
	return &Bare{
		logger:       logger,
		cmdName:      cmdName,
		nodes:        nodes,
		nodeExecutor: executor,
	}
}

// Component get cluster components
func (b *Bare) Component() ([]cluster.Component, error) {
	cmd := fmt.Sprintf("pgrep %s &&  cat /proc/`pgrep %s`/cmdline | xargs -0 | tr ' ' '\\n'", b.cmdName, b.cmdName)
	result := make([]cluster.Component, 0)
	lk := sync.Mutex{}
	conCtl := make(chan struct{}, 200)
	g := errgroup.Group{}

	for _, tempN := range b.nodes {
		n := tempN
		g.Go(func() error {
			conCtl <- struct{}{}
			defer func() { <-conCtl }()

			out, _, err := b.nodeExecutor.DoCmd(n, []string{
				"/bin/sh", "-c", cmd,
			})
			if err != nil {
				if !strings.Contains(err.Error(), "terminated with exit code") {
					b.logger.Errorf("do command on node %s failed :%v", n, err)
				}
				return err
			}

			cmp := cluster.Component{
				Name: b.cmdName,
				Node: n,
				Args: map[string]string{},
			}

			lines := strings.Split(out, "\n")
			for i, line := range lines {
				line = strings.TrimSpace(line)
				line = strings.TrimLeft(line, "-")
				if line == "" {
					continue
				}

				if i == 0 {
					cmp.IsRunning = true
					continue
				}

				spIndex := strings.IndexAny(line, "=")
				if spIndex == -1 {
					continue
				}

				k := line[0:spIndex]
				v := line[spIndex+1:]
				cmp.Args[strings.TrimSpace(k)] = strings.TrimSpace(v)
			}

			lk.Lock()
			result = append(result, cmp)
			lk.Unlock()
			return nil
		})
	}
	_ = g.Wait()

	return result, nil
}
