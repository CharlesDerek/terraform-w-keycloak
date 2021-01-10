package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
	"regexp"
	"testing"
)

func TestAccKeycloakOpenIdScriptProtocolMapper_basicClient(t *testing.T) {
	t.Parallel()

	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-script-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_script_protocol_mapper.script_mapper_client"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdScriptProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdScriptProtocolMapper_basic_client(clientId, mapperName),
				Check:  testKeycloakOpenIdScriptProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdScriptProtocolMapper_basicClientScope(t *testing.T) {
	t.Parallel()

	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-script-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_script_protocol_mapper.script_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdScriptProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdScriptProtocolMapper_basic_clientScope(clientScopeId, mapperName),
				Check:  testKeycloakOpenIdScriptProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdScriptProtocolMapper_import(t *testing.T) {
	t.Parallel()

	clientId := "terraform-openid-client-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-script-mapper-" + acctest.RandString(5)

	clientResourceName := "keycloak_openid_script_protocol_mapper.script_mapper_client"
	clientScopeResourceName := "keycloak_openid_script_protocol_mapper.script_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdScriptProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdScriptProtocolMapper_import(clientId, clientScopeId, mapperName),
				Check: resource.ComposeTestCheckFunc(
					testKeycloakOpenIdScriptProtocolMapperExists(clientResourceName),
					testKeycloakOpenIdScriptProtocolMapperExists(clientScopeResourceName),
				),
			},
			{
				ResourceName:      clientResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getGenericProtocolMapperIdForClient(clientResourceName),
			},
			{
				ResourceName:      clientScopeResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getGenericProtocolMapperIdForClientScope(clientScopeResourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdScriptProtocolMapper_update(t *testing.T) {
	t.Parallel()

	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-script-mapper-" + acctest.RandString(5)

	attributeName := "claim-" + acctest.RandString(10)
	updatedAttributeName := "claim-update-" + acctest.RandString(10)
	resourceName := "keycloak_openid_script_protocol_mapper.script_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdScriptProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdScriptProtocolMapper_claim(clientId, mapperName, attributeName),
				Check:  testKeycloakOpenIdScriptProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdScriptProtocolMapper_claim(clientId, mapperName, updatedAttributeName),
				Check:  testKeycloakOpenIdScriptProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdScriptProtocolMapper_createAfterManualDestroy(t *testing.T) {
	t.Parallel()

	var mapper = &keycloak.OpenIdScriptProtocolMapper{}

	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-script-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_script_protocol_mapper.script_mapper_client"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdScriptProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdScriptProtocolMapper_basic_client(clientId, mapperName),
				Check:  testKeycloakOpenIdScriptProtocolMapperFetch(resourceName, mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteOpenIdScriptProtocolMapper(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)
					if err != nil {
						t.Error(err)
					}
				},
				Config: testKeycloakOpenIdScriptProtocolMapper_basic_client(clientId, mapperName),
				Check:  testKeycloakOpenIdScriptProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdScriptProtocolMapper_validateClaimValueType(t *testing.T) {
	t.Parallel()

	mapperName := "terraform-openid-connect-script-mapper-" + acctest.RandString(10)
	invalidClaimValueType := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdScriptProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdScriptProtocolMapper_claimValueType(mapperName, invalidClaimValueType),
				ExpectError: regexp.MustCompile("expected claim_value_type to be one of .+ got " + invalidClaimValueType),
			},
		},
	})
}

func TestAccKeycloakOpenIdScriptProtocolMapper_updateClientIdForceNew(t *testing.T) {
	t.Parallel()

	clientId := "terraform-client-" + acctest.RandString(10)
	updatedClientId := "terraform-client-update-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-script-mapper-" + acctest.RandString(5)

	attributeName := "claim-" + acctest.RandString(10)
	resourceName := "keycloak_openid_script_protocol_mapper.script_mapper"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdScriptProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdScriptProtocolMapper_claim(clientId, mapperName, attributeName),
				Check:  testKeycloakOpenIdScriptProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdScriptProtocolMapper_claim(updatedClientId, mapperName, attributeName),
				Check:  testKeycloakOpenIdScriptProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdScriptProtocolMapper_updateClientScopeForceNew(t *testing.T) {
	t.Parallel()

	mapperName := "terraform-openid-connect-script-mapper-" + acctest.RandString(5)
	clientScopeId := "terraform-client-" + acctest.RandString(10)
	newClientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	resourceName := "keycloak_openid_script_protocol_mapper.script_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccKeycloakOpenIdScriptProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdScriptProtocolMapper_basic_clientScope(clientScopeId, mapperName),
				Check:  testKeycloakOpenIdScriptProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdScriptProtocolMapper_basic_clientScope(newClientScopeId, mapperName),
				Check:  testKeycloakOpenIdScriptProtocolMapperExists(resourceName),
			},
		},
	})
}

func testAccKeycloakOpenIdScriptProtocolMapperDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for resourceName, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_openid_script_protocol_mapper" {
				continue
			}

			mapper, _ := getScriptMapperUsingState(state, resourceName)

			if mapper != nil {
				return fmt.Errorf("openid script protocol mapper with id %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testKeycloakOpenIdScriptProtocolMapperExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getScriptMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testKeycloakOpenIdScriptProtocolMapperFetch(resourceName string, mapper *keycloak.OpenIdScriptProtocolMapper) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		fetchedMapper, err := getScriptMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		mapper.Id = fetchedMapper.Id
		mapper.ClientId = fetchedMapper.ClientId
		mapper.ClientScopeId = fetchedMapper.ClientScopeId
		mapper.RealmId = fetchedMapper.RealmId

		return nil
	}
}

func getScriptMapperUsingState(state *terraform.State, resourceName string) (*keycloak.OpenIdScriptProtocolMapper, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]
	clientScopeId := rs.Primary.Attributes["client_scope_id"]

	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	return keycloakClient.GetOpenIdScriptProtocolMapper(realm, clientId, clientScopeId, id)
}

func testKeycloakOpenIdScriptProtocolMapper_basic_client(clientId, mapperName string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = data.keycloak_realm.realm.id
	client_id = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_script_protocol_mapper" "script_mapper_client" {
	name       = "%s"
	realm_id   = data.keycloak_realm.realm.id
	client_id  = keycloak_openid_client.openid_client.id
	script     = "exports = 'foo';"
	claim_name = "bar"
}`, testAccRealm.Realm, clientId, mapperName)
}

func testKeycloakOpenIdScriptProtocolMapper_basic_clientScope(clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = data.keycloak_realm.realm.id
}

resource "keycloak_openid_script_protocol_mapper" "script_mapper_client_scope" {
	name            = "%s"
	realm_id        = data.keycloak_realm.realm.id
	client_scope_id = keycloak_openid_client_scope.client_scope.id
	script          = "exports = 'foo';"
	claim_name      = "bar"
}`, testAccRealm.Realm, clientScopeId, mapperName)
}

func testKeycloakOpenIdScriptProtocolMapper_import(clientId, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id    = data.keycloak_realm.realm.id
	client_id   = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_script_protocol_mapper" "script_mapper_client" {
	name       = "%s"
	realm_id   = data.keycloak_realm.realm.id
	client_id  = keycloak_openid_client.openid_client.id
	script     = "exports = 'foo';"
	claim_name = "bar"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = data.keycloak_realm.realm.id
}

resource "keycloak_openid_script_protocol_mapper" "script_mapper_client_scope" {
	name            = "%s"
	realm_id        = data.keycloak_realm.realm.id
	client_scope_id = keycloak_openid_client_scope.client_scope.id
	script          = "exports = 'foo';"
	claim_name      = "bar"
}`, testAccRealm.Realm, clientId, mapperName, clientScopeId, mapperName)
}

func testKeycloakOpenIdScriptProtocolMapper_claim(clientId, mapperName, attributeName string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	realm_id  = data.keycloak_realm.realm.id
	client_id = "%s"

	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_script_protocol_mapper" "script_mapper" {
	name       = "%s"
	realm_id   = data.keycloak_realm.realm.id
	client_id  = keycloak_openid_client.openid_client.id
	script     = "exports = '%s';"
	claim_name = "bar"
}`, testAccRealm.Realm, clientId, mapperName, attributeName)
}

func testKeycloakOpenIdScriptProtocolMapper_claimValueType(mapperName, claimValueType string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_script_protocol_mapper" "script_mapper_validation" {
	name              = "%s"
	realm_id          = data.keycloak_realm.realm.id
	script            = "exports = 'foo';"
	claim_name        = "bar"
	claim_value_type  = "%s"
}`, testAccRealm.Realm, mapperName, claimValueType)
}
