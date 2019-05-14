package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
	"strings"
)

func resourceKeycloakOpenidClientAuthorizationResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientAuthorizationResourceCreate,
		Read:   resourceKeycloakOpenidClientAuthorizationResourceRead,
		Delete: resourceKeycloakOpenidClientAuthorizationResourceDelete,
		Update: resourceKeycloakOpenidClientAuthorizationResourceUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakOpenidClientAuthorizationResourceImport,
		},
		Schema: map[string]*schema.Schema{
			"resource_server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"uris": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"icon_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"owner_managed_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"scopes": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func getOpenidClientAuthorizationResourceFromData(data *schema.ResourceData) *keycloak.OpenidClientAuthorizationResource {
	var uris []string
	var scopes []keycloak.OpenidClientAuthorizationScope
	attributes := map[string][]string{}
	if v, ok := data.GetOk("uris"); ok {
		for _, uri := range v.(*schema.Set).List() {
			uris = append(uris, uri.(string))
		}
	}
	if v, ok := data.GetOk("scopes"); ok {
		for _, scope := range v.(*schema.Set).List() {
			scopes = append(scopes, keycloak.OpenidClientAuthorizationScope{
				Name: scope.(string),
			})
		}
	}
	if v, ok := data.GetOk("attributes"); ok {
		for key, value := range v.(map[string]interface{}) {
			attributes[key] = strings.Split(value.(string), ",")
		}
	}
	resource := keycloak.OpenidClientAuthorizationResource{
		Id:                 data.Id(),
		DisplayName:        data.Get("display_name").(string),
		Name:               data.Get("name").(string),
		IconUri:            data.Get("icon_uri").(string),
		OwnerManagedAccess: data.Get("owner_managed_access").(bool),
		Type:               data.Get("type").(string),
		ResourceServerId:   data.Get("resource_server_id").(string),
		RealmId:            data.Get("realm_id").(string),
		Uris:               uris,
		Scopes:             scopes,
		Attributes:         attributes,
	}
	return &resource
}

func setOpenidClientAuthorizationResourceData(data *schema.ResourceData, resource *keycloak.OpenidClientAuthorizationResource) {
	scopes := []string{}
	for _, scope := range resource.Scopes {
		scopes = append(scopes, scope.Name)
	}
	data.SetId(resource.Id)
	data.Set("resource_server_id", resource.ResourceServerId)
	data.Set("realm_id", resource.RealmId)
	data.Set("display_name", resource.DisplayName)
	data.Set("name", resource.Name)
	data.Set("icon_uri", resource.IconUri)
	data.Set("owner_managed_access", resource.OwnerManagedAccess)
	data.Set("type", resource.Type)
	data.Set("uris", resource.Uris)
	data.Set("attributes", resource.Attributes)
	data.Set("scopes", scopes)
}

func resourceKeycloakOpenidClientAuthorizationResourceCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	resource := getOpenidClientAuthorizationResourceFromData(data)

	err := keycloakClient.NewOpenidClientAuthorizationResource(resource)
	if err != nil {
		return err
	}

	setOpenidClientAuthorizationResourceData(data, resource)

	return resourceKeycloakOpenidClientAuthorizationResourceRead(data, meta)
}

func resourceKeycloakOpenidClientAuthorizationResourceRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	resourceServerId := data.Get("resource_server_id").(string)
	id := data.Id()

	resource, err := keycloakClient.GetOpenidClientAuthorizationResource(realmId, resourceServerId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setOpenidClientAuthorizationResourceData(data, resource)

	return nil
}

func resourceKeycloakOpenidClientAuthorizationResourceUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	resource := getOpenidClientAuthorizationResourceFromData(data)

	err := keycloakClient.UpdateOpenidClientAuthorizationResource(resource)
	if err != nil {
		return err
	}

	setOpenidClientAuthorizationResourceData(data, resource)

	return nil
}

func resourceKeycloakOpenidClientAuthorizationResourceDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	resourceServerId := data.Get("resource_server_id").(string)
	id := data.Id()

	return keycloakClient.DeleteOpenidClientAuthorizationResource(realmId, resourceServerId, id)
}

func resourceKeycloakOpenidClientAuthorizationResourceImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{resourceServerId}}/{{authorizationResourceId}}")
	}
	d.Set("realm_id", parts[0])
	d.Set("resource_server_id", parts[1])
	d.SetId(parts[3])

	return []*schema.ResourceData{d}, nil
}
