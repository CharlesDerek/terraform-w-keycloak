package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type OpenidClientAuthorizationJSPolicy struct {
	Id               string `json:"id,omitempty"`
	RealmId          string `json:"-"`
	ResourceServerId string `json:"-"`
	Name             string `json:"name"`
	DecisionStrategy string `json:"decisionStrategy"`
	Logic            string `json:"logic"`
	Type             string `json:"type"`
	Code             string `json:"code"`
	Description      string `json:"description"`
}

func (keycloakClient *KeycloakClient) NewOpenidClientAuthorizationJSPolicy(ctx context.Context, policy *OpenidClientAuthorizationJSPolicy) error {
	var body []byte
	var err error
	if strings.HasSuffix(policy.Code, ".js") {
		body, _, err = keycloakClient.post(ctx, fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/%s", policy.RealmId, policy.ResourceServerId, policy.Code), policy)
	} else {
		body, _, err = keycloakClient.post(ctx, fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/js", policy.RealmId, policy.ResourceServerId), policy)
	}
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &policy)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) UpdateOpenidClientAuthorizationJSPolicy(ctx context.Context, policy *OpenidClientAuthorizationJSPolicy) error {
	err := keycloakClient.put(ctx, fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/js/%s", policy.RealmId, policy.ResourceServerId, policy.Id), policy)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) DeleteOpenidClientAuthorizationJSPolicy(ctx context.Context, realmId, resourceServerId, policyId string) error {
	return keycloakClient.delete(ctx, fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/js/%s", realmId, resourceServerId, policyId), nil)
}

func (keycloakClient *KeycloakClient) GetOpenidClientAuthorizationJSPolicy(ctx context.Context, realmId, resourceServerId, policyId string) (*OpenidClientAuthorizationJSPolicy, error) {

	policy := OpenidClientAuthorizationJSPolicy{
		Id:               policyId,
		ResourceServerId: resourceServerId,
		RealmId:          realmId,
	}
	err := keycloakClient.get(ctx, fmt.Sprintf("/realms/%s/clients/%s/authz/resource-server/policy/js/%s", realmId, resourceServerId, policyId), &policy, nil)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}
