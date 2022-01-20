package provider

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
)

var (
	keycloakOpenidClientAccessTypes                          = []string{"CONFIDENTIAL", "PUBLIC", "BEARER-ONLY"}
	keycloakOpenidClientAuthorizationPolicyEnforcementMode   = []string{"ENFORCING", "PERMISSIVE", "DISABLED"}
	keycloakOpenidClientResourcePermissionDecisionStrategies = []string{"UNANIMOUS", "AFFIRMATIVE", "CONSENSUS"}
	keycloakOpenidClientPkceCodeChallengeMethod              = []string{"", "plain", "S256"}
	keycloakOpenidClientAuthenticatorTypes                   = []string{"client-secret", "client-jwt", "client-x509", "client-secret-jwt"}
)

func resourceKeycloakOpenidClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientCreate,
		Read:   resourceKeycloakOpenidClientRead,
		Delete: resourceKeycloakOpenidClientDelete,
		Update: resourceKeycloakOpenidClientUpdate,
		// This resource can be imported using {{realm}}/{{client_id}}. The Client ID is displayed in the GUI
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakOpenidClientImport,
		},
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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(keycloakOpenidClientAccessTypes, false),
			},
			"client_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"client_authenticator_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(keycloakOpenidClientAuthenticatorTypes, false),
				Default:      "client-secret",
			},
			"standard_flow_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"implicit_flow_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"direct_access_grants_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"service_accounts_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"valid_redirect_uris": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},
			"web_origins": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},
			"root_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"admin_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"base_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_account_user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pkce_code_challenge_method": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(keycloakOpenidClientPkceCodeChallengeMethod, false),
			},
			"access_token_lifespan": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_offline_session_idle_timeout": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_offline_session_max_lifespan": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_session_idle_timeout": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_session_max_lifespan": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"exclude_session_state_from_auth_response": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"resource_server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorization": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_enforcement_mode": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(keycloakOpenidClientAuthorizationPolicyEnforcementMode, false),
						},
						"decision_strategy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(keycloakOpenidClientResourcePermissionDecisionStrategies, false),
							Default:      "UNANIMOUS",
						},
						"allow_remote_resource_management": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"keep_defaults": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"full_scope_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"consent_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"authentication_flow_binding_overrides": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"browser_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"direct_grant_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"login_theme": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"use_refresh_tokens": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"backchannel_logout_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backchannel_logout_session_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"backchannel_logout_revoke_offline_sessions": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"extra_config": {
				Type:             schema.TypeMap,
				Optional:         true,
				ValidateDiagFunc: validateExtraConfig(reflect.ValueOf(&keycloak.OpenidClientAttributes{}).Elem()),
			},
			"oauth2_device_authorization_grant_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"oauth2_device_code_lifespan": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oauth2_device_polling_interval": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		CustomizeDiff: customdiff.ComputedIf("service_account_user_id", func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
			return d.HasChange("service_accounts_enabled")
		}),
	}
}

func getOpenidClientFromData(data *schema.ResourceData) (*keycloak.OpenidClient, error) {
	validRedirectUris := make([]string, 0)
	webOrigins := make([]string, 0)

	rootUrlData, rootUrlOk := data.GetOkExists("root_url")
	validRedirectUrisData, validRedirectUrisOk := data.GetOk("valid_redirect_uris")
	webOriginsData, webOriginsOk := data.GetOk("web_origins")

	rootUrlString := rootUrlData.(string)

	if validRedirectUrisOk {
		for _, validRedirectUri := range validRedirectUrisData.(*schema.Set).List() {
			validRedirectUris = append(validRedirectUris, validRedirectUri.(string))
		}
	}

	if webOriginsOk {
		for _, webOrigin := range webOriginsData.(*schema.Set).List() {
			webOrigins = append(webOrigins, webOrigin.(string))
		}
	}

	// Keycloak uses the root URL for web origins if not specified otherwise
	if rootUrlOk && rootUrlString != "" {
		if !validRedirectUrisOk {
			return nil, errors.New("valid_redirect_uris is required when root_url is given1")
		}
		if !webOriginsOk {
			return nil, errors.New("web_origins is required when root_url is given")
		}
		if _, adminOk := data.GetOk("admin_url"); !adminOk {
			return nil, errors.New("admin_url is required when root_url is given")
		}
	}

	openidClient := &keycloak.OpenidClient{
		Id:                        data.Id(),
		ClientId:                  data.Get("client_id").(string),
		RealmId:                   data.Get("realm_id").(string),
		Name:                      data.Get("name").(string),
		Enabled:                   data.Get("enabled").(bool),
		Description:               data.Get("description").(string),
		ClientSecret:              data.Get("client_secret").(string),
		ClientAuthenticatorType:   data.Get("client_authenticator_type").(string),
		StandardFlowEnabled:       data.Get("standard_flow_enabled").(bool),
		ImplicitFlowEnabled:       data.Get("implicit_flow_enabled").(bool),
		DirectAccessGrantsEnabled: data.Get("direct_access_grants_enabled").(bool),
		ServiceAccountsEnabled:    data.Get("service_accounts_enabled").(bool),
		FullScopeAllowed:          data.Get("full_scope_allowed").(bool),
		Attributes: keycloak.OpenidClientAttributes{
			PkceCodeChallengeMethod:               data.Get("pkce_code_challenge_method").(string),
			ExcludeSessionStateFromAuthResponse:   keycloak.KeycloakBoolQuoted(data.Get("exclude_session_state_from_auth_response").(bool)),
			AccessTokenLifespan:                   data.Get("access_token_lifespan").(string),
			LoginTheme:                            data.Get("login_theme").(string),
			ClientOfflineSessionIdleTimeout:       data.Get("client_offline_session_idle_timeout").(string),
			ClientOfflineSessionMaxLifespan:       data.Get("client_offline_session_max_lifespan").(string),
			ClientSessionIdleTimeout:              data.Get("client_session_idle_timeout").(string),
			ClientSessionMaxLifespan:              data.Get("client_session_max_lifespan").(string),
			UseRefreshTokens:                      keycloak.KeycloakBoolQuoted(data.Get("use_refresh_tokens").(bool)),
			BackchannelLogoutUrl:                  data.Get("backchannel_logout_url").(string),
			BackchannelLogoutRevokeOfflineTokens:  keycloak.KeycloakBoolQuoted(data.Get("backchannel_logout_revoke_offline_sessions").(bool)),
			BackchannelLogoutSessionRequired:      keycloak.KeycloakBoolQuoted(data.Get("backchannel_logout_session_required").(bool)),
			ExtraConfig:                           getExtraConfigFromData(data),
			Oauth2DeviceAuthorizationGrantEnabled: keycloak.KeycloakBoolQuoted(data.Get("oauth2_device_authorization_grant_enabled").(bool)),
			Oauth2DeviceCodeLifespan:              data.Get("oauth2_device_code_lifespan").(string),
			Oauth2DevicePollingInterval:           data.Get("oauth2_device_polling_interval").(string),
		},
		ValidRedirectUris: validRedirectUris,
		WebOrigins:        webOrigins,
		AdminUrl:          data.Get("admin_url").(string),
		BaseUrl:           data.Get("base_url").(string),
		ConsentRequired:   data.Get("consent_required").(bool),
	}

	if rootUrlOk {
		openidClient.RootUrl = &rootUrlString
	}

	if !openidClient.ImplicitFlowEnabled && !openidClient.StandardFlowEnabled {
		if _, ok := data.GetOk("valid_redirect_uris"); ok {
			return nil, errors.New("valid_redirect_uris cannot be set when standard or implicit flow is not enabled")
		}
	}

	if !openidClient.ImplicitFlowEnabled && !openidClient.StandardFlowEnabled && !openidClient.DirectAccessGrantsEnabled {
		if _, ok := data.GetOk("web_origins"); ok {
			return nil, errors.New("web_origins cannot be set when standard or implicit flow is not enabled")
		}
	}

	// access type
	if accessType := data.Get("access_type").(string); accessType == "PUBLIC" {
		openidClient.PublicClient = true
	} else if accessType == "BEARER-ONLY" {
		openidClient.BearerOnly = true
	}

	if v, ok := data.GetOk("authorization"); ok {
		openidClient.AuthorizationServicesEnabled = true
		authorizationSettingsData := v.(*schema.Set).List()[0]
		authorizationSettings := authorizationSettingsData.(map[string]interface{})
		openidClient.AuthorizationSettings = &keycloak.OpenidClientAuthorizationSettings{
			PolicyEnforcementMode:         authorizationSettings["policy_enforcement_mode"].(string),
			DecisionStrategy:              authorizationSettings["decision_strategy"].(string),
			AllowRemoteResourceManagement: authorizationSettings["allow_remote_resource_management"].(bool),
			KeepDefaults:                  authorizationSettings["keep_defaults"].(bool),
		}
	} else {
		openidClient.AuthorizationServicesEnabled = false
	}

	if v, ok := data.GetOk("authentication_flow_binding_overrides"); ok {
		authenticationFlowBindingOverridesData := v.(*schema.Set).List()[0]
		authenticationFlowBindingOverrides := authenticationFlowBindingOverridesData.(map[string]interface{})
		openidClient.AuthenticationFlowBindingOverrides = keycloak.OpenidAuthenticationFlowBindingOverrides{
			BrowserId:     authenticationFlowBindingOverrides["browser_id"].(string),
			DirectGrantId: authenticationFlowBindingOverrides["direct_grant_id"].(string),
		}
	}

	return openidClient, nil
}

func setOpenidClientData(keycloakClient *keycloak.KeycloakClient, data *schema.ResourceData, client *keycloak.OpenidClient) error {
	var serviceAccountUserId string
	if client.ServiceAccountsEnabled {
		serviceAccountUser, err := keycloakClient.GetOpenidClientServiceAccountUserId(client.RealmId, client.Id)
		if err != nil {
			return err
		}
		serviceAccountUserId = serviceAccountUser.Id
	}
	data.SetId(client.Id)
	data.Set("client_id", client.ClientId)
	data.Set("realm_id", client.RealmId)
	data.Set("name", client.Name)
	data.Set("enabled", client.Enabled)
	data.Set("description", client.Description)
	data.Set("client_secret", client.ClientSecret)
	data.Set("client_authenticator_type", client.ClientAuthenticatorType)
	data.Set("standard_flow_enabled", client.StandardFlowEnabled)
	data.Set("implicit_flow_enabled", client.ImplicitFlowEnabled)
	data.Set("direct_access_grants_enabled", client.DirectAccessGrantsEnabled)
	data.Set("service_accounts_enabled", client.ServiceAccountsEnabled)
	data.Set("valid_redirect_uris", client.ValidRedirectUris)
	data.Set("web_origins", client.WebOrigins)
	data.Set("admin_url", client.AdminUrl)
	data.Set("base_url", client.BaseUrl)
	data.Set("root_url", &client.RootUrl)
	data.Set("full_scope_allowed", client.FullScopeAllowed)
	data.Set("consent_required", client.ConsentRequired)

	data.Set("access_token_lifespan", client.Attributes.AccessTokenLifespan)
	data.Set("login_theme", client.Attributes.LoginTheme)
	data.Set("use_refresh_tokens", client.Attributes.UseRefreshTokens)
	data.Set("oauth2_device_authorization_grant_enabled", client.Attributes.Oauth2DeviceAuthorizationGrantEnabled)
	data.Set("oauth2_device_code_lifespan", client.Attributes.Oauth2DeviceCodeLifespan)
	data.Set("oauth2_device_polling_interval", client.Attributes.Oauth2DevicePollingInterval)
	data.Set("client_offline_session_idle_timeout", client.Attributes.ClientOfflineSessionIdleTimeout)
	data.Set("client_offline_session_max_lifespan", client.Attributes.ClientOfflineSessionMaxLifespan)
	data.Set("client_session_idle_timeout", client.Attributes.ClientSessionIdleTimeout)
	data.Set("client_session_max_lifespan", client.Attributes.ClientSessionMaxLifespan)
	data.Set("backchannel_logout_url", client.Attributes.BackchannelLogoutUrl)
	data.Set("backchannel_logout_revoke_offline_sessions", client.Attributes.BackchannelLogoutRevokeOfflineTokens)
	data.Set("backchannel_logout_session_required", client.Attributes.BackchannelLogoutSessionRequired)
	setExtraConfigData(data, client.Attributes.ExtraConfig)

	if client.AuthorizationServicesEnabled {
		data.Set("resource_server_id", client.Id)
	}

	if client.ServiceAccountsEnabled {
		data.Set("service_account_user_id", serviceAccountUserId)
	} else {
		data.Set("service_account_user_id", "")
	}

	// access type
	if client.PublicClient {
		data.Set("access_type", "PUBLIC")
	} else if client.BearerOnly {
		data.Set("access_type", "BEARER-ONLY")
	} else {
		data.Set("access_type", "CONFIDENTIAL")
	}

	if (keycloak.OpenidAuthenticationFlowBindingOverrides{}) == client.AuthenticationFlowBindingOverrides {
		data.Set("authentication_flow_binding_overrides", nil)
	} else {
		authenticationFlowBindingOverridesSettings := make(map[string]interface{})
		authenticationFlowBindingOverridesSettings["browser_id"] = client.AuthenticationFlowBindingOverrides.BrowserId
		authenticationFlowBindingOverridesSettings["direct_grant_id"] = client.AuthenticationFlowBindingOverrides.DirectGrantId
		data.Set("authentication_flow_binding_overrides", []interface{}{authenticationFlowBindingOverridesSettings})
	}

	return nil
}

func resourceKeycloakOpenidClientCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	client, err := getOpenidClientFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.ValidateOpenidClient(client)
	if err != nil {
		return err
	}

	err = keycloakClient.NewOpenidClient(client)
	if err != nil {
		return err
	}

	err = setOpenidClientData(keycloakClient, data, client)
	if err != nil {
		return err
	}

	return resourceKeycloakOpenidClientRead(data, meta)
}

func resourceKeycloakOpenidClientRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	client, err := keycloakClient.GetOpenidClient(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	err = setOpenidClientData(keycloakClient, data, client)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakOpenidClientUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	client, err := getOpenidClientFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.ValidateOpenidClient(client)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateOpenidClient(client)
	if err != nil {
		return err
	}

	err = setOpenidClientData(keycloakClient, data, client)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakOpenidClientDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteOpenidClient(realmId, id)
}

func resourceKeycloakOpenidClientImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{openidClientId}}")
	}
	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, resourceKeycloakOpenidClientRead(d, meta)
}
