package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
)

func resourceKeycloakLdapMsadUserAccountControlMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakLdapMsadUserAccountControlMapperCreate,
		Read:   resourceKeycloakLdapMsadUserAccountControlMapperRead,
		Update: resourceKeycloakLdapMsadUserAccountControlMapperUpdate,
		Delete: resourceKeycloakLdapMsadUserAccountControlMapperDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the mapper when displayed in the console.",
			},
			"realm_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The realm in which the ldap user federation provider exists.",
			},
			"ldap_user_federation_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ldap user federation provider to attach this mapper to.",
			},
			"ldap_password_policy_hints_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func getLdapMsadUserAccountControlMapperFromData(data *schema.ResourceData) *keycloak.LdapMsadUserAccountControlMapper {
	return &keycloak.LdapMsadUserAccountControlMapper{
		Id:                   data.Id(),
		Name:                 data.Get("name").(string),
		RealmId:              data.Get("realm_id").(string),
		LdapUserFederationId: data.Get("ldap_user_federation_id").(string),

		LdapPasswordPolicyHintsEnabled: data.Get("ldap_password_policy_hints_enabled").(bool),
	}
}

func setLdapMsadUserAccountControlMapperData(data *schema.ResourceData, ldapMsadUserAccountControlMapper *keycloak.LdapMsadUserAccountControlMapper) {
	data.SetId(ldapMsadUserAccountControlMapper.Id)

	data.Set("id", ldapMsadUserAccountControlMapper.Id)
	data.Set("name", ldapMsadUserAccountControlMapper.Name)
	data.Set("realm_id", ldapMsadUserAccountControlMapper.RealmId)
	data.Set("ldap_user_federation_id", ldapMsadUserAccountControlMapper.LdapUserFederationId)

	data.Set("ldap_password_policy_hints_enabled", ldapMsadUserAccountControlMapper.LdapPasswordPolicyHintsEnabled)
}

func resourceKeycloakLdapMsadUserAccountControlMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	ldapMsadUserAccountControlMapper := getLdapMsadUserAccountControlMapperFromData(data)

	err := keycloakClient.NewLdapMsadUserAccountControlMapper(ldapMsadUserAccountControlMapper)
	if err != nil {
		return err
	}

	setLdapMsadUserAccountControlMapperData(data, ldapMsadUserAccountControlMapper)

	return resourceKeycloakLdapMsadUserAccountControlMapperRead(data, meta)
}

func resourceKeycloakLdapMsadUserAccountControlMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	ldapMsadUserAccountControlMapper, err := keycloakClient.GetLdapMsadUserAccountControlMapper(realmId, id)
	if err != nil {
		return err
	}

	setLdapMsadUserAccountControlMapperData(data, ldapMsadUserAccountControlMapper)

	return nil
}

func resourceKeycloakLdapMsadUserAccountControlMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	ldapMsadUserAccountControlMapper := getLdapMsadUserAccountControlMapperFromData(data)

	err := keycloakClient.UpdateLdapMsadUserAccountControlMapper(ldapMsadUserAccountControlMapper)
	if err != nil {
		return err
	}

	setLdapMsadUserAccountControlMapperData(data, ldapMsadUserAccountControlMapper)

	return nil
}

func resourceKeycloakLdapMsadUserAccountControlMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteLdapMsadUserAccountControlMapper(realmId, id)
}
