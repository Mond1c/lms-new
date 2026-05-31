package handler

import (
	"context"

	"connectrpc.com/connect"
	identityv1 "github.com/Mond1c/lms/gen/go/lms/v1"
	"github.com/Mond1c/lms/services/identity/internal/service"
)

func (h *Identity) LinkVCSIdentity(
	ctx context.Context,
	req *connect.Request[identityv1.LinkVCSIdentityRequest],
) (*connect.Response[identityv1.VCSIdentity], error) {
	m := req.Msg
	identity, err := h.vcsIdentities.Link(ctx, service.LinkVCSIdentityInput{
		UserID:         m.GetUserId(),
		Provider:       providerFromProto(m.GetProvider()),
		ExternalUserID: m.GetExternalUserId(),
		ExternalLogin:  m.GetExternalLogin(),
		AccessToken:    m.GetAccessToken(),
		RefreshToken:   m.GetRefreshToken(),
		ExpiresAt:      tsFromProto(m.GetExpiresAt()),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}
	return connect.NewResponse(vcsIdentityToProto(identity)), nil
}

func (h *Identity) UnlinkVCSIdentity(
	ctx context.Context,
	req *connect.Request[identityv1.UnlinkVCSIdentityRequest],
) (*connect.Response[identityv1.UnlinkVCSIdentityResponse], error) {
	if err := h.vcsIdentities.Unlink(ctx, req.Msg.GetUserId(), providerFromProto(req.Msg.GetProvider())); err != nil {
		return nil, toConnectErr(err)
	}
	return connect.NewResponse(&identityv1.UnlinkVCSIdentityResponse{}), nil
}

func (h *Identity) ListVCSIdentities(
	ctx context.Context,
	req *connect.Request[identityv1.ListVCSIdentitiesRequest],
) (*connect.Response[identityv1.ListVCSIdentitiesResponse], error) {
	identities, err := h.vcsIdentities.List(ctx, req.Msg.GetUserId())
	if err != nil {
		return nil, toConnectErr(err)
	}

	protoIdentities := make([]*identityv1.VCSIdentity, len(identities))
	for i, identity := range identities {
		protoIdentities[i] = vcsIdentityToProto(identity)
	}
	return connect.NewResponse(&identityv1.ListVCSIdentitiesResponse{
		Identities: protoIdentities,
	}), nil
}
