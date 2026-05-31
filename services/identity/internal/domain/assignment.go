package domain

import "time"

// RunnerKind selects where an assignment's tests run.
type RunnerKind string

const (
	RunnerExternalCI RunnerKind = "external_ci"
	RunnerSelfHosted RunnerKind = "self_hosted"
)

func (r RunnerKind) Valid() bool {
	switch r {
	case RunnerExternalCI, RunnerSelfHosted:
		return true
	}
	return false
}

// GradingPolicy configures how a submission's final grade is composed.
// See docs/architecture.md §4.
type GradingPolicy struct {
	WeightTests       float64
	WeightQuality     float64
	DefenceMultiplier bool
	CustomFormula     string
}

type Assignment struct {
	ID                      string
	CourseID                string
	Slug                    string
	Title                   string
	DescriptionMarkdown     string
	Deadline                *time.Time
	HardDeadline            *time.Time
	MaxScore                int32
	TemplateRepo            string
	RepoNamingPattern       string
	AutoRequestReviewOnPass bool
	RequiresDefense         bool
	GradingPolicy           GradingPolicy
	Runner                  RunnerKind
	CreatedAt               time.Time
	UpdatedAt               time.Time
}
