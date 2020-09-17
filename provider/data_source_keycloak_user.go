package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
)

func dataSourceKeycloakUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKeycloakUserRead,
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email_verified": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attributes": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"federated_identity": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceKeycloakUserRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmID := data.Get("realm_id").(string)
	username := data.Get("username").(string)

	user, err := keycloakClient.GetUserByUsername(realmID, username)
	if err != nil {
		return err
	}

	mapFromUserToData(data, user)

	return nil
}
