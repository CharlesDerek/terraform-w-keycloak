package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
)

func dataSourceKeycloakRole() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeycloakRoleRead,
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceKeycloakRoleRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	clientId := data.Get("client_id").(string)
	roleName := data.Get("name").(string)

	role, err := keycloakClient.GetRoleByName(realmId, clientId, roleName)
	if err != nil {
		return err
	}

	mapFromRoleToData(data, role)

	return nil
}
