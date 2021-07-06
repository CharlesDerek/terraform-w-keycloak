package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
	"strings"
)

func resourceKeycloakUserRoles() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakUserRolesReconcile,
		Read:   resourceKeycloakUserRolesRead,
		Update: resourceKeycloakUserRolesReconcile,
		Delete: resourceKeycloakUserRolesDelete,
		// This resource can be imported using {{realm}}/{{userId}}.
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakUserRolesImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Required: true,
			},
			"exhaustive": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func userRolesId(realmId, userId string) string {
	return fmt.Sprintf("%s/%s", realmId, userId)
}

func addRolesToUser(keycloakClient *keycloak.KeycloakClient, clientRolesToAdd map[string][]*keycloak.Role, realmRolesToAdd []*keycloak.Role, user *keycloak.User) error {
	if len(realmRolesToAdd) != 0 {
		err := keycloakClient.AddRealmRolesToUser(user.RealmId, user.Id, realmRolesToAdd)
		if err != nil {
			return err
		}
	}

	for k, roles := range clientRolesToAdd {
		if len(roles) != 0 {
			err := keycloakClient.AddClientRolesToUser(user.RealmId, user.Id, k, roles)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func removeRolesFromUser(keycloakClient *keycloak.KeycloakClient, clientRolesToRemove map[string][]*keycloak.Role, realmRolesToRemove []*keycloak.Role, user *keycloak.User) error {
	if len(realmRolesToRemove) != 0 {
		err := keycloakClient.RemoveRealmRolesFromUser(user.RealmId, user.Id, realmRolesToRemove)
		if err != nil {
			return err
		}
	}

	for k, roles := range clientRolesToRemove {
		if len(roles) != 0 {
			err := keycloakClient.RemoveClientRolesFromUser(user.RealmId, user.Id, k, roles)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceKeycloakUserRolesReconcile(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	userId := data.Get("user_id").(string)
	roleIds := interfaceSliceToStringSlice(data.Get("role_ids").(*schema.Set).List())
	exhaustive := data.Get("exhaustive").(bool)

	user, err := keycloakClient.GetUser(realmId, userId)
	if err != nil {
		return err
	}

	if data.HasChange("role_ids") && !data.IsNewResource() {
		o, n := data.GetChange("role_ids")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := interfaceSliceToStringSlice(os.Difference(ns).List())

		tfRolesToRemove, err := getExtendedRoleMapping(keycloakClient, realmId, remove)
		if err != nil {
			return err
		}

		if err = removeRolesFromUser(keycloakClient, tfRolesToRemove.clientRoles, tfRolesToRemove.realmRoles, user); err != nil {
			return err
		}
	}

	tfRoles, err := getExtendedRoleMapping(keycloakClient, realmId, roleIds)
	if err != nil {
		return err
	}

	// get the list of currently assigned roles. Due to default realm and client roles
	// (e.g. roles of the account client) this is probably not empty upon resource creation
	roleMappings, err := keycloakClient.GetUserRoleMappings(realmId, userId)

	// sort into roles we need to add and roles we need to remove
	updates := calculateRoleMappingUpdates(tfRoles, intoRoleMapping(roleMappings))

	// add roles
	err = addRolesToUser(keycloakClient, updates.clientRolesToAdd, updates.realmRolesToAdd, user)
	if err != nil {
		return err
	}

	// remove roles if exhaustive (authoritative)
	if exhaustive {
		err = removeRolesFromUser(keycloakClient, updates.clientRolesToRemove, updates.realmRolesToRemove, user)
		if err != nil {
			return err
		}
	}

	data.SetId(userRolesId(realmId, userId))
	return resourceKeycloakUserRolesRead(data, meta)
}

func resourceKeycloakUserRolesRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	userId := data.Get("user_id").(string)
	sortedRoleIds := interfaceSliceToStringSlice(data.Get("role_ids").(*schema.Set).List())
	exhaustive := data.Get("exhaustive").(bool)

	// check if user exists, remove from state if not found
	if _, err := keycloakClient.GetUser(realmId, userId); err != nil {
		return handleNotFoundError(err, data)
	}

	roles, err := keycloakClient.GetUserRoleMappings(realmId, userId)
	if err != nil {
		return err
	}

	var roleIds []string

	for _, realmRole := range roles.RealmMappings {
		if exhaustive || stringSliceContains(sortedRoleIds, realmRole.Id) {
			roleIds = append(roleIds, realmRole.Id)
		}
	}

	for _, clientRoleMapping := range roles.ClientMappings {
		for _, clientRole := range clientRoleMapping.Mappings {
			if exhaustive || stringSliceContains(sortedRoleIds, clientRole.Id) {
				roleIds = append(roleIds, clientRole.Id)
			}
		}
	}

	data.Set("role_ids", roleIds)
	data.SetId(userRolesId(realmId, userId))

	return nil
}

func resourceKeycloakUserRolesDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	userId := data.Get("user_id").(string)

	user, err := keycloakClient.GetUser(realmId, userId)

	roleIds := interfaceSliceToStringSlice(data.Get("role_ids").(*schema.Set).List())
	rolesToRemove, err := getExtendedRoleMapping(keycloakClient, realmId, roleIds)
	if err != nil {
		return err
	}

	err = removeRolesFromUser(keycloakClient, rolesToRemove.clientRoles, rolesToRemove.realmRoles, user)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakUserRolesImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import format: {{realm}}/{{userId}}.")
	}

	realmId := parts[0]
	userId := parts[1]

	d.Set("realm_id", realmId)
	d.Set("user_id", userId)
	d.Set("exhaustive", true)

	d.SetId(userRolesId(realmId, userId))

	return []*schema.ResourceData{d}, resourceKeycloakUserRolesRead(d, meta)
}
