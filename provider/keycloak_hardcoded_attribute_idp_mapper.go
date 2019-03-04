package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
)

func resourceKeycloakHardcodedAttributeIdpMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"attribute_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "OIDC Claim",
		},
		"attribute_value": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "User Attribute",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakIdentityProviderMapperCreate("hardcoded-attribute-idp-mapper")
	genericMapperResource.Read = resourceKeycloakIdentityProviderMapperRead("hardcoded-attribute-idp-mapper")
	genericMapperResource.Update = resourceKeycloakIdentityProviderMapperUpdate("hardcoded-attribute-idp-mapper")
	return genericMapperResource
}

func getHardcodedAttributeIdpMapperFromData(data *schema.ResourceData) (*keycloak.IdentityProviderMapper, error) {
	rec, _ := getIdentityProviderMapperFromData(data)
	rec.IdentityProviderMapper = "hardcoded-attribute-idp-mapper"
	rec.Config = &keycloak.IdentityProviderMapperConfig{
		Attribute:      data.Get("attribute_name").(string),
		AttributeValue: data.Get("attribute_value").(string),
	}
	return rec, nil
}

func setHardcodedAttributeIdpMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	data.Set("attribute_name", identityProviderMapper.Config.Attribute)
	data.Set("attribute_value", identityProviderMapper.Config.AttributeValue)
	return nil
}
