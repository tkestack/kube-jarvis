/*
* Tencent is pleased to support the open source community by making TKEStack
* available.
*
* Copyright (C) 2012-2019 Tencent. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the “License”); you may not use
* this file except in compliance with the License. You may obtain a copy of the
* License at
*
* https://opensource.org/licenses/Apache-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
* WARRANTIES OF ANY KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations under the License.
 */
package plugins

import (
	"sync"
)

// Progress show the progress of cluster Initialization
type Progress struct {
	sync.Mutex
	// IsDone if true if will Step done
	IsDone bool
	// Steps shows the sub progress of every step
	Steps map[string]*ProgressStep
	// CurStep is the name of current ProgressStep in Steps
	CurStep string
	// Total is the max step value of Progress
	Total int
	// Current is finished step value of Progress
	Current int
	// Percent is the percentage of overall progress
	Percent  float64
	watchers []func(p *Progress)
}

// ProgressStep is one step of Initialization
type ProgressStep struct {
	// Title is the short describe of this step
	Title string
	// Total is the total value of this step
	Total int
	// Current is current value of this step
	// Percent is the percentage of overall this step
	Percent float64
	Current int
}

// NewProgress return a new Progress
func NewProgress() *Progress {
	return &Progress{
		Steps: make(map[string]*ProgressStep),
	}
}

// Clone return a new Progress with the same value of origin
func (p *Progress) Clone() *Progress {
	progress := &Progress{
		IsDone:  p.IsDone,
		CurStep: p.CurStep,
		Total:   p.Total,
		Current: p.Current,
		Steps:   map[string]*ProgressStep{},
	}

	for name, step := range p.Steps {
		s := &ProgressStep{
			Title:   step.Title,
			Total:   step.Total,
			Current: step.Current,
		}

		progress.Steps[name] = s
	}
	return progress
}

// CreateStep create and return an new InitialProgressStep from a InitialProgress
func (p *Progress) CreateStep(name string, title string, total int) {
	p.update(func() {
		if _, exist := p.Steps[name]; exist {
			return
		}

		step := &ProgressStep{
			Title: title,
			Total: total,
		}
		p.Total += total
		p.Steps[name] = step
	})
}

// SetCurStep change current step of Progress
func (p *Progress) SetCurStep(name string) {
	p.update(func() {
		p.CurStep = name
	})
}

// AddPercent add current progress  percent
func (p *Progress) AddStepPercent(name string, n int) {
	p.update(func() {
		step := p.Steps[name]
		step.Current += n
		step.Percent = float64(step.Current) / float64(step.Total) * 100
		p.Current += n
		p.Percent = float64(p.Current) / float64(p.Total) * 100
	})
}

// Done
func (p *Progress) Done() {
	p.update(func() {
		p.IsDone = true
	})
}

// AddProgressUpdatedWatcher add a watcher that will be called once progress updated
func (p *Progress) AddProgressUpdatedWatcher(f func(p *Progress)) {
	p.watchers = append(p.watchers, f)
}

func (p *Progress) updated() {
	for _, f := range p.watchers {
		f(p)
	}
}

func (p *Progress) update(f func()) {
	p.Lock()
	defer p.Unlock()
	f()
	p.updated()
}
