package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
)

func resourceKeycloakOpenidClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientCreate,
		Read:   resourceKeycloakOpenidClientRead,
		Delete: resourceKeycloakOpenidClientDelete,
		Update: resourceKeycloakOpenidClientUpdate,
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func getOpenidClientFromData(data *schema.ResourceData) *keycloak.OpenidClient {
	return &keycloak.OpenidClient{
		Id:       data.Id(),
		ClientId: data.Get("client_id").(string),
		RealmId:  data.Get("realm_id").(string),
	}
}

func setOpenidClientData(data *schema.ResourceData, client *keycloak.OpenidClient) {
	data.SetId(client.Id)

	data.Set("client_id", client.ClientId)
	data.Set("realm_id", client.RealmId)
}

func resourceKeycloakOpenidClientCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	client := getOpenidClientFromData(data)

	err := keycloakClient.NewOpenidClient(client)
	if err != nil {
		return err
	}

	setOpenidClientData(data, client)

	return resourceKeycloakOpenidClientRead(data, meta)
}

func resourceKeycloakOpenidClientRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	client, err := keycloakClient.GetOpenidClient(realmId, id)
	if err != nil {
		return err
	}

	setOpenidClientData(data, client)

	return nil
}

func resourceKeycloakOpenidClientUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	client := getOpenidClientFromData(data)

	err := keycloakClient.UpdateOpenidClient(client)
	if err != nil {
		return err
	}

	setOpenidClientData(data, client)

	return nil
}

func resourceKeycloakOpenidClientDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteOpenidClient(realmId, id)
}
