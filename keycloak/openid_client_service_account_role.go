package keycloak

import (
	"fmt"
)

type OpenidClientServiceAccountRole struct {
	Id                   string `json:"id"`
	RealmId              string `json:"-"`
	ServiceAccountUserId string `json:"-"`
	Name                 string `json:"name,omitempty"`
	ClientRole           bool   `json:"clientRole"`
	Composite            bool   `json:"composite"`
	ContainerId          string `json:"containerId"`
	Description          string `json:"description"`
}

func (keycloakClient *KeycloakClient) NewOpenidClientServiceAccountRole(serviceAccountRole *OpenidClientServiceAccountRole) error {
	role, err := keycloakClient.GetClientRoleByName(serviceAccountRole.RealmId, serviceAccountRole.ContainerId, serviceAccountRole.Name)
	if err != nil {
		return err
	}
	serviceAccountRole.Id = role.Id
	serviceAccountRoles := []OpenidClientServiceAccountRole{*serviceAccountRole}
	_, _, err = keycloakClient.post(fmt.Sprintf("/realms/%s/users/%s/role-mappings/clients/%s", serviceAccountRole.RealmId, serviceAccountRole.ServiceAccountUserId, serviceAccountRole.ContainerId), serviceAccountRoles)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) DeleteOpenidClientServiceAccountRole(realm, serviceAccountUserId, clientId, roleId string) error {
	serviceAccountRole, err := keycloakClient.GetOpenidClientServiceAccountRole(realm, serviceAccountUserId, clientId, roleId)
	if err != nil {
		return err
	}
	serviceAccountRoles := []OpenidClientServiceAccountRole{*serviceAccountRole}
	err = keycloakClient.delete(fmt.Sprintf("/realms/%s/users/%s/role-mappings/clients/%s", realm, serviceAccountUserId, clientId), &serviceAccountRoles)
	if err != nil {
		return err
	}
	return nil
}

func (keycloakClient *KeycloakClient) GetOpenidClientServiceAccountRole(realm, serviceAccountUserId, clientId, roleId string) (*OpenidClientServiceAccountRole, error) {
	serviceAccountRoles := []OpenidClientServiceAccountRole{
		{
			Id:                   roleId,
			RealmId:              realm,
			ContainerId:          clientId,
			ServiceAccountUserId: serviceAccountUserId,
		},
	}
	err := keycloakClient.get(fmt.Sprintf("/realms/%s/users/%s/role-mappings/clients/%s", realm, serviceAccountUserId, clientId), &serviceAccountRoles, nil)
	if err != nil {
		return nil, err
	}
	for _, serviceAccountRole := range serviceAccountRoles {
		if serviceAccountRole.Id == roleId {
			return &serviceAccountRole, nil
		}
	}
	return &OpenidClientServiceAccountRole{}, nil
}
