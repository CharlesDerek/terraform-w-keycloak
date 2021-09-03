package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
)

func resourceKeycloakCustomIdentityProviderMapper() *schema.Resource {
	mapperSchema := map[string]*schema.Schema{
		"identity_provider_mapper": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "IDP Mapper Type",
		},
	}
	genericMapperResource := resourceKeycloakIdentityProviderMapper()
	genericMapperResource.Schema = mergeSchemas(genericMapperResource.Schema, mapperSchema)
	genericMapperResource.Create = resourceKeycloakIdentityProviderMapperCreate(getCustomIdentityProviderMapperFromData, setCustomIdentityProviderMapperData)
	genericMapperResource.Read = resourceKeycloakIdentityProviderMapperRead(setCustomIdentityProviderMapperData)
	genericMapperResource.Update = resourceKeycloakIdentityProviderMapperUpdate(getCustomIdentityProviderMapperFromData, setCustomIdentityProviderMapperData)
	return genericMapperResource
}

func getCustomIdentityProviderMapperFromData(data *schema.ResourceData, meta interface{}) (*keycloak.IdentityProviderMapper, error) {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	rec, _ := getIdentityProviderMapperFromData(data)
	identityProvider, err := keycloakClient.GetIdentityProvider(rec.Realm, rec.IdentityProviderAlias)
	if err != nil {
		return nil, handleNotFoundError(err, data)
	}

	identityProviderMapper := data.Get("identity_provider_mapper").(string)
	if strings.Contains(identityProviderMapper, "%s") {
		rec.IdentityProviderMapper = fmt.Sprintf(identityProviderMapper, identityProvider.ProviderId)
	} else {
		rec.IdentityProviderMapper = identityProviderMapper
	}

	return rec, nil
}

func setCustomIdentityProviderMapperData(data *schema.ResourceData, identityProviderMapper *keycloak.IdentityProviderMapper) error {
	setIdentityProviderMapperData(data, identityProviderMapper)
	setExtraConfigData(data, identityProviderMapper.Config.ExtraConfig)
	data.Set("identity_provider_mapper", identityProviderMapper.IdentityProviderMapper)

	return nil
}
