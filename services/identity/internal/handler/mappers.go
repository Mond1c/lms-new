package handler

import (
	"time"

	identityv1 "github.com/Mond1c/lms/gen/go/lms/v1"
	"github.com/Mond1c/lms/services/identity/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func auditToProto(created, updated time.Time) *identityv1.AuditFields {
	return &identityv1.AuditFields{
		CreatedAt: timestamppb.New(created),
		UpdatedAt: timestamppb.New(updated),
	}
}

func tsToProto(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func tsFromProto(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime()
	return &t
}

func providerFromProto(p *identityv1.ProviderRef) domain.ProviderRef {
	if p == nil {
		return domain.ProviderRef{}
	}
	return domain.ProviderRef{Kind: int32(p.GetKind()), Instance: p.GetInstance()}
}

func providerToProto(p domain.ProviderRef) *identityv1.ProviderRef {
	return &identityv1.ProviderRef{
		Kind:     identityv1.ProviderKind(p.Kind),
		Instance: p.Instance,
	}
}

func roleFromProto(r identityv1.Role) domain.Role {
	switch r {
	case identityv1.Role_ROLE_STUDENT:
		return domain.RoleStudent
	case identityv1.Role_ROLE_INSTRUCTOR:
		return domain.RoleInstructor
	case identityv1.Role_ROLE_ADMIN:
		return domain.RoleAdmin
	default:
		return ""
	}
}

func roleToProto(r domain.Role) identityv1.Role {
	switch r {
	case domain.RoleStudent:
		return identityv1.Role_ROLE_STUDENT
	case domain.RoleInstructor:
		return identityv1.Role_ROLE_INSTRUCTOR
	case domain.RoleAdmin:
		return identityv1.Role_ROLE_ADMIN
	default:
		return identityv1.Role_ROLE_UNSPECIFIED
	}
}

func runnerFromProto(r identityv1.RunnerKind) domain.RunnerKind {
	switch r {
	case identityv1.RunnerKind_RUNNER_KIND_EXTERNAL_CI:
		return domain.RunnerExternalCI
	case identityv1.RunnerKind_RUNNER_KIND_SELF_HOSTED:
		return domain.RunnerSelfHosted
	default:
		return "" // service applies the default
	}
}

func runnerToProto(r domain.RunnerKind) identityv1.RunnerKind {
	switch r {
	case domain.RunnerExternalCI:
		return identityv1.RunnerKind_RUNNER_KIND_EXTERNAL_CI
	case domain.RunnerSelfHosted:
		return identityv1.RunnerKind_RUNNER_KIND_SELF_HOSTED
	default:
		return identityv1.RunnerKind_RUNNER_KIND_UNSPECIFIED
	}
}

func provisioningToProto(s domain.ProvisioningState) identityv1.ProvisioningState {
	switch s {
	case domain.ProvisioningPending:
		return identityv1.ProvisioningState_PROVISIONING_STATE_PENDING
	case domain.ProvisioningReady:
		return identityv1.ProvisioningState_PROVISIONING_STATE_READY
	case domain.ProvisioningFailed:
		return identityv1.ProvisioningState_PROVISIONING_STATE_FAILED
	default:
		return identityv1.ProvisioningState_PROVISIONING_STATE_UNSPECIFIED
	}
}

func gradingPolicyFromProto(p *identityv1.GradingPolicy) domain.GradingPolicy {
	if p == nil {
		return domain.GradingPolicy{}
	}
	return domain.GradingPolicy{
		WeightTests:       p.GetWeightTests(),
		WeightQuality:     p.GetWeightQuality(),
		DefenceMultiplier: p.GetDefenceMultiplier(),
		CustomFormula:     p.GetCustomFormula(),
	}
}

func gradingPolicyToProto(p domain.GradingPolicy) *identityv1.GradingPolicy {
	return &identityv1.GradingPolicy{
		WeightTests:       p.WeightTests,
		WeightQuality:     p.WeightQuality,
		DefenceMultiplier: p.DefenceMultiplier,
		CustomFormula:     p.CustomFormula,
	}
}

func courseToProto(c *domain.Course) *identityv1.Course {
	out := &identityv1.Course{
		Id:           c.ID,
		Code:         c.Code,
		Title:        c.Title,
		Description:  c.Description,
		InstructorId: c.InstructorID,
		Audit:        auditToProto(c.CreatedAt, c.UpdatedAt),
	}
	if c.VCS != nil {
		out.Vcs = &identityv1.VCSBinding{
			Provider:       providerToProto(c.VCS.Provider),
			TargetOrg:      c.VCS.TargetOrg,
			StudentTeam:    c.VCS.StudentTeam,
			ReviewerTeam:   c.VCS.ReviewerTeam,
			ReviewerLogins: c.VCS.ReviewerLogins,
		}
	}
	return out
}

func enrollmentToProto(e *domain.Enrollment) *identityv1.Enrollment {
	return &identityv1.Enrollment{
		Id:         e.ID,
		UserId:     e.UserID,
		CourseId:   e.CourseID,
		Role:       roleToProto(e.Role),
		EnrolledAt: timestamppb.New(e.EnrolledAt),
	}
}

func assignmentToProto(a *domain.Assignment) *identityv1.Assignment {
	return &identityv1.Assignment{
		Id:                      a.ID,
		CourseId:                a.CourseID,
		Slug:                    a.Slug,
		Title:                   a.Title,
		DescriptionMarkdown:     a.DescriptionMarkdown,
		Deadline:                tsToProto(a.Deadline),
		HardDeadline:            tsToProto(a.HardDeadline),
		MaxScore:                a.MaxScore,
		TemplateRepo:            a.TemplateRepo,
		RepoNamingPattern:       a.RepoNamingPattern,
		AutoRequestReviewOnPass: a.AutoRequestReviewOnPass,
		RequiresDefense:         a.RequiresDefense,
		GradingPolicy:           gradingPolicyToProto(a.GradingPolicy),
		Runner:                  runnerToProto(a.Runner),
		Audit:                   auditToProto(a.CreatedAt, a.UpdatedAt),
	}
}

func vcsIdentityToProto(vi *domain.VCSIdentity) *identityv1.VCSIdentity {
	return &identityv1.VCSIdentity{
		Provider:       providerToProto(vi.Provider),
		ExternalUserId: vi.ExternalUserID,
		ExternalLogin:  vi.ExternalLogin,
		LinkedAt:       timestamppb.New(vi.LinkedAt),
		TokenValid:     vi.TokenValid,
	}
}

func studentRepoToProto(sr *domain.StudentRepo) *identityv1.StudentRepo {
	return &identityv1.StudentRepo{
		Id:            sr.ID,
		UserId:        sr.UserID,
		AssignmentId:  sr.AssignmentID,
		Provider:      providerToProto(sr.Provider),
		FullName:      sr.FullName,
		ExternalId:    sr.ExternalID,
		State:         provisioningToProto(sr.State),
		CloneUrlHttps: sr.CloneURLHTTPS,
		CloneUrlSsh:   sr.CloneURLSSH,
		Audit:         auditToProto(sr.CreatedAt, sr.UpdatedAt),
	}
}
