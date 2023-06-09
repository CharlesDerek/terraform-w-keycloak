package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
)

type OpenidClientAuthorizationUserPolicy struct {
	Id               string   `json:"id,omitempty"`
	RealmId          string   `json:"-"`
	ResourceServerId string   `json:"-"`
	Name             string   `json:"name"`
	DecisionStrategy string   `json:"decisionStrategy"`
	Logic            string   `json:"logic"`
	Type             string   `json:"type"`
	Users            []string `json:"users"`
	Description      string   `json:"description"`
}

func (keycloakClient *KeycloakClient) NewOpenidClientAuthorizationUserPolicy(ctx context.Context, policy *OpenidClientAuthorizationUserPolicy) error {
	body, _, err := keycloakClient.post(ctx, fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/user", policy.RealmId, policy.ResourceServerId), policy)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &policy)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) UpdateOpenidClientAuthorizationUserPolicy(ctx context.Context, policy *OpenidClientAuthorizationUserPolicy) error {
	err := keycloakClient.put(ctx, fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/user/%s", policy.RealmId, policy.ResourceServerId, policy.Id), policy)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) DeleteOpenidClientAuthorizationUserPolicy(ctx context.Context, realmId, resourceServerId, policyId string) error {
	return keycloakClient.delete(ctx, fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/user/%s", realmId, resourceServerId, policyId), nil)
}

func (keycloakClient *KeycloakClient) GetOpenidClientAuthorizationUserPolicy(ctx context.Context, realmId, resourceServerId, policyId string) (*OpenidClientAuthorizationUserPolicy, error) {

	policy := OpenidClientAuthorizationUserPolicy{
		Id:               policyId,
		ResourceServerId: resourceServerId,
		RealmId:          realmId,
	}
	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/user/%s", realmId, resourceServerId, policyId), &policy, nil)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}
