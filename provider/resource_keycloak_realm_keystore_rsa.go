package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
	"strings"
)

var (
	keycloakRealmKeystoreRsaAlgorithm = []string{"RS256", "RS384", "RS512", "PS256", "PS384", "PS512"}
)

func resourceKeycloakRealmKeystoreRsa() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakRealmKeystoreRsaCreate,
		Read:   resourceKeycloakRealmKeystoreRsaRead,
		Update: resourceKeycloakRealmKeystoreRsaUpdate,
		Delete: resourceKeycloakRealmKeystoreRsaDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakRealmKeystoreRsaImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of provider when linked in admin console.",
			},
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Set if the keys can be used for signing",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Set if the keys are enabled",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Priority for the provider",
			},
			"algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(keycloakRealmKeystoreRsaAlgorithm, false),
				Default:      "RS256",
				Description:  "Intended algorithm for the key",
			},
			"private_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Private RSA Key encoded in PEM format",
			},
			"certificate": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "X509 Certificate encoded in PEM format",
			},
		},
	}
}

func getRealmKeystoreRsaFromData(keycloakClient *keycloak.KeycloakClient, data *schema.ResourceData) (*keycloak.RealmKeystoreRsa, error) {
	mapper := &keycloak.RealmKeystoreRsa{
		Id:      data.Id(),
		Name:    data.Get("name").(string),
		RealmId: data.Get("realm_id").(string),

		Active:      data.Get("active").(bool),
		Enabled:     data.Get("enabled").(bool),
		Priority:    data.Get("priority").(int),
		Algorithm:   data.Get("algorithm").(string),
		PrivateKey:  data.Get("private_key").(string),
		Certificate: data.Get("certificate").(string),
	}
	_, err := keycloakClient.VersionIsGreaterThanOrEqualTo(keycloak.Version_11)
	if err != nil {
		return nil, err
	}

	return mapper, nil
}

func setRealmKeystoreRsaData(keycloakClient *keycloak.KeycloakClient, data *schema.ResourceData,
	realmKey *keycloak.RealmKeystoreRsa, readFunc bool) error {
	data.SetId(realmKey.Id)

	data.Set("name", realmKey.Name)
	data.Set("realm_id", realmKey.RealmId)

	data.Set("active", realmKey.Active)
	data.Set("enabled", realmKey.Enabled)
	data.Set("priority", realmKey.Priority)
	data.Set("algorithm", realmKey.Algorithm)
	if realmKey.PrivateKey != "**********" {
		data.Set("private_key", realmKey.PrivateKey)
		data.Set("certificate", realmKey.Certificate)
	}

	_, err := keycloakClient.VersionIsGreaterThanOrEqualTo(keycloak.Version_11)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakRealmKeystoreRsaCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeystoreRsaFromData(keycloakClient, data)
	if err != nil {
		return err
	}

	err = keycloakClient.NewRealmKeystoreRsa(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeystoreRsaData(keycloakClient, data, realmKey, false)
	if err != nil {
		return err
	}

	return resourceKeycloakRealmKeystoreRsaRead(data, meta)
}

func resourceKeycloakRealmKeystoreRsaRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	realmKey, err := keycloakClient.GetRealmKeystoreRsa(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	err = setRealmKeystoreRsaData(keycloakClient, data, realmKey, true)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakRealmKeystoreRsaUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmKey, err := getRealmKeystoreRsaFromData(keycloakClient, data)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateRealmKeystoreRsa(realmKey)
	if err != nil {
		return err
	}

	err = setRealmKeystoreRsaData(keycloakClient, data, realmKey, false)
	if err != nil {
		return err
	}

	return nil
}

func resourceKeycloakRealmKeystoreRsaDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteRealmKeystoreRsa(realmId, id)
}

func resourceKeycloakRealmKeystoreRsaImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{keystoreId}}")
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
