package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
)

func resourceKeycloakOidcGoogleIdentityProvider() *schema.Resource {
	oidcGoogleSchema := map[string]*schema.Schema{
		"alias": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The alias uniquely identifies an identity provider and it is also used to build the redirect uri. In case of google this is computed and always google",
		},
		"display_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Not used by this provider, Will be implicitly Google",
		},
		"provider_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "google",
			Description: "provider id, is always google, unless you have a extended custom implementation",
		},
		"client_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Client ID.",
		},
		"client_secret": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Client Secret.",
		},
		"hosted_domain": { //hostedDomain
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Set 'hd' query parameter when logging in with Google. Google will list accounts only for this domain. Keycloak validates that the returned identity token has a claim for this domain. When '*' is entered, any hosted account can be used.",
		},
		"use_user_ip_param": { //userIp
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Set 'userIp' query parameter when invoking on Google's User Info service.  This will use the user's ip address.  Useful if Google is throttling access to the User Info service.",
		},
		"request_refresh_token": { //offlineAccess
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Set 'access_type' query parameter to 'offline' when redirecting to google authorization endpoint, to get a refresh token back. Useful if planning to use Token Exchange to retrieve Google token to access Google APIs when the user is not at the browser.",
		},
		"default_scopes": { //defaultScope
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "openid profile email",
			Description: "The scopes to be sent when asking for authorization. See the documentation for possible values, separator and default value'. Default: 'openid profile email'",
		},
		"accepts_prompt_none_forward_from_client": { // acceptsPromptNoneForwardFromClient
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "This is just used together with Identity Provider Authenticator or when kc_idp_hint points to this identity provider. In case that client sends a request with prompt=none and user is not yet authenticated, the error will not be directly returned to client, but the request with prompt=none will be forwarded to this identity provider.",
		},
		"disable_user_info": { //disableUserInfo
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Disable usage of User Info service to obtain additional user information?  Default is to use this OIDC service.",
		},
		"hide_on_login_page": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Hide On Login Page.",
		},
		"extra_config": {
			Type:     schema.TypeMap,
			Optional: true,
		},
	}
	oidcResource := resourceKeycloakIdentityProvider()
	oidcResource.Schema = mergeSchemas(oidcResource.Schema, oidcGoogleSchema)
	oidcResource.Create = resourceKeycloakIdentityProviderCreate(getOidcGoogleIdentityProviderFromData, setOidcGoogleIdentityProviderData)
	oidcResource.Read = resourceKeycloakIdentityProviderRead(setOidcGoogleIdentityProviderData)
	oidcResource.Update = resourceKeycloakIdentityProviderUpdate(getOidcGoogleIdentityProviderFromData, setOidcGoogleIdentityProviderData)
	return oidcResource
}

func getOidcGoogleIdentityProviderFromData(data *schema.ResourceData) (*keycloak.IdentityProvider, error) {
	rec, _ := getIdentityProviderFromData(data)
	rec.ProviderId = data.Get("provider_id").(string)
	rec.Alias = "google"

	extraConfig := map[string]interface{}{}
	if v, ok := data.GetOk("extra_config"); ok {
		for key, value := range v.(map[string]interface{}) {
			extraConfig[key] = value
		}
	}

	rec.Config = &keycloak.IdentityProviderConfig{
		ClientId:                    data.Get("client_id").(string),
		ClientSecret:                data.Get("client_secret").(string),
		HideOnLoginPage:             keycloak.KeycloakBoolQuoted(data.Get("hide_on_login_page").(bool)),
		HostedDomain:                data.Get("hosted_domain").(string),
		UserIp:                      keycloak.KeycloakBoolQuoted(data.Get("use_user_ip_param").(bool)),
		OfflineAccess:               keycloak.KeycloakBoolQuoted(data.Get("request_refresh_token").(bool)),
		DefaultScope:                data.Get("default_scopes").(string),
		AcceptsPromptNoneForwFrmClt: keycloak.KeycloakBoolQuoted(data.Get("accepts_prompt_none_forward_from_client").(bool)),
		ExtraConfig:                 extraConfig,
		UseJwksUrl:                  keycloak.KeycloakBoolQuoted(true),
		DisableUserInfo:             keycloak.KeycloakBoolQuoted(data.Get("disable_user_info").(bool)),
	}

	return rec, nil
}

func setOidcGoogleIdentityProviderData(data *schema.ResourceData, identityProvider *keycloak.IdentityProvider) error {
	setIdentityProviderData(data, identityProvider)
	data.Set("provider_id", identityProvider.ProviderId)
	data.Set("client_id", identityProvider.Config.ClientId)
	data.Set("hide_on_login_page", identityProvider.Config.HideOnLoginPage)
	data.Set("hosted_domain", identityProvider.Config.HostedDomain)
	data.Set("use_user_ip_param", identityProvider.Config.UserIp)
	data.Set("request_refresh_token", identityProvider.Config.OfflineAccess)
	data.Set("default_scopes", identityProvider.Config.DefaultScope)
	data.Set("accepts_prompt_none_forward_from_client", identityProvider.Config.AcceptsPromptNoneForwFrmClt)
	data.Set("extra_config", identityProvider.Config.ExtraConfig)
	data.Set("use_jwks_url", identityProvider.Config.UseJwksUrl)
	data.Set("disable_user_info", identityProvider.Config.DisableUserInfo)
	return nil
}
