package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
)

func TestAccKeycloakOpenIdUserClientRoleProtocolMapper_basicClient(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-client-role-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_user_client_role_protocol_mapper.user_client_role_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserClientRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserClientRoleProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserClientRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserClientRoleProtocolMapper_basicClientScope(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-client-role-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_user_client_role_protocol_mapper.user_client_role_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserClientRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserClientRoleProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName),
				Check:  testKeycloakOpenIdUserClientRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserClientRoleProtocolMapper_import(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-openid-client-" + acctest.RandString(10)
	clientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-client-role-mapper-" + acctest.RandString(5)

	clientResourceName := "keycloak_openid_user_client_role_protocol_mapper.user_client_role_mapper_client"
	clientScopeResourceName := "keycloak_openid_user_client_role_protocol_mapper.user_client_role_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserClientRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserClientRoleProtocolMapper_import(realmName, clientId, clientScopeId, mapperName),
				Check: resource.ComposeTestCheckFunc(
					testKeycloakOpenIdUserClientRoleProtocolMapperExists(clientResourceName),
					testKeycloakOpenIdUserClientRoleProtocolMapperExists(clientScopeResourceName),
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

func TestAccKeycloakOpenIdUserClientRoleProtocolMapper_update(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-client-role-mapper-" + acctest.RandString(5)

	claimName := "claim-name-" + acctest.RandString(10)
	updatedClaimName := "claim-name-update-" + acctest.RandString(10)

	resourceName := "keycloak_openid_user_client_role_protocol_mapper.user_client_role_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserClientRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserClientRoleProtocolMapper_claim(realmName, clientId, mapperName, claimName),
				Check:  testKeycloakOpenIdUserClientRoleProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserClientRoleProtocolMapper_claim(realmName, clientId, mapperName, updatedClaimName),
				Check:  testKeycloakOpenIdUserClientRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserClientRoleProtocolMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.OpenIdUserClientRoleProtocolMapper{}

	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-client-role-mapper-" + acctest.RandString(5)

	resourceName := "keycloak_openid_user_client_role_protocol_mapper.user_client_role_mapper_client"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserClientRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserClientRoleProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserClientRoleProtocolMapperFetch(resourceName, mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteOpenIdUserClientRoleProtocolMapper(mapper.RealmId, mapper.ClientId, mapper.ClientScopeId, mapper.Id)
					if err != nil {
						t.Error(err)
					}
				},
				Config: testKeycloakOpenIdUserClientRoleProtocolMapper_basic_client(realmName, clientId, mapperName),
				Check:  testKeycloakOpenIdUserClientRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserClientRoleProtocolMapper_validateClaimValueType(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-client-role-mapper-" + acctest.RandString(10)
	invalidClaimValueType := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserClientRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenIdUserClientRoleProtocolMapper_validateClaimValueType(realmName, mapperName, invalidClaimValueType),
				ExpectError: regexp.MustCompile("expected claim_value_type to be one of .+ got " + invalidClaimValueType),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserClientRoleProtocolMapper_updateClientIdForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	updatedClientId := "terraform-client-update-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-client-role-mapper-" + acctest.RandString(5)

	claimName := "claim-name-" + acctest.RandString(10)
	resourceName := "keycloak_openid_user_client_role_protocol_mapper.user_client_role_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserClientRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserClientRoleProtocolMapper_claim(realmName, clientId, mapperName, claimName),
				Check:  testKeycloakOpenIdUserClientRoleProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserClientRoleProtocolMapper_claim(realmName, updatedClientId, mapperName, claimName),
				Check:  testKeycloakOpenIdUserClientRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserClientRoleProtocolMapper_updateClientScopeForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-client-role-mapper-" + acctest.RandString(5)
	clientScopeId := "terraform-client-" + acctest.RandString(10)
	newClientScopeId := "terraform-client-scope-" + acctest.RandString(10)
	resourceName := "keycloak_openid_user_client_role_protocol_mapper.user_client_role_mapper_client_scope"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserClientRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserClientRoleProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName),
				Check:  testKeycloakOpenIdUserClientRoleProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserClientRoleProtocolMapper_basic_clientScope(realmName, newClientScopeId, mapperName),
				Check:  testKeycloakOpenIdUserClientRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserClientRoleProtocolMapper_updateRealmIdForceNew(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	newRealmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-client-role-mapper-" + acctest.RandString(5)

	claimName := "claim-name-" + acctest.RandString(10)
	resourceName := "keycloak_openid_user_client_role_protocol_mapper.user_client_role_mapper"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserClientRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserClientRoleProtocolMapper_claim(realmName, clientId, mapperName, claimName),
				Check:  testKeycloakOpenIdUserClientRoleProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserClientRoleProtocolMapper_claim(newRealmName, clientId, mapperName, claimName),
				Check:  testKeycloakOpenIdUserClientRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserClientRoleProtocolMapper_clientAssignment(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	assignedClientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-client-role-mapper-" + acctest.RandString(5)
	resourceName := "keycloak_openid_user_client_role_protocol_mapper.user_client_role_mapper_validation"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserClientRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserClientRoleProtocolMapper_clientAssignment(realmName, clientId, assignedClientId, mapperName),
				Check:  testKeycloakOpenIdUserClientRoleProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserClientRoleProtocolMapper_clientAssignment(realmName, clientId, assignedClientId, mapperName),
				Check:  testKeycloakOpenIdUserClientRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakOpenIdUserClientRoleProtocolMapper_clientAssignmentRolePrefix(t *testing.T) {
	realmName := "terraform-realm-" + acctest.RandString(10)
	clientId := "terraform-client-" + acctest.RandString(10)
	assignedClientId := "terraform-client-" + acctest.RandString(10)
	mapperName := "terraform-openid-connect-user-client-role-mapper-" + acctest.RandString(5)
	rolePrefix := "role-prefix-" + acctest.RandString(10)
	resourceName := "keycloak_openid_user_client_role_protocol_mapper.user_client_role_mapper_validation"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccKeycloakOpenIdUserClientRoleProtocolMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenIdUserClientRoleProtocolMapper_clientAssignmentRolePrefix(realmName, clientId, assignedClientId, mapperName, rolePrefix),
				Check:  testKeycloakOpenIdUserClientRoleProtocolMapperExists(resourceName),
			},
			{
				Config: testKeycloakOpenIdUserClientRoleProtocolMapper_clientAssignmentRolePrefix(realmName, clientId, assignedClientId, mapperName, rolePrefix),
				Check:  testKeycloakOpenIdUserClientRoleProtocolMapperExists(resourceName),
			},
		},
	})
}

func testAccKeycloakOpenIdUserClientRoleProtocolMapperDestroy() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for resourceName, rs := range state.RootModule().Resources {
			if rs.Type != "keycloak_openid_user_client_role_protocol_mapper" {
				continue
			}

			mapper, _ := getUserClientRoleMapperUsingState(state, resourceName)

			if mapper != nil {
				return fmt.Errorf("openid user attribute protocol mapper with id %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testKeycloakOpenIdUserClientRoleProtocolMapperExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, err := getUserClientRoleMapperUsingState(state, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testKeycloakOpenIdUserClientRoleProtocolMapperFetch(resourceName string, mapper *keycloak.OpenIdUserClientRoleProtocolMapper) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		fetchedMapper, err := getUserClientRoleMapperUsingState(state, resourceName)
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

func getUserClientRoleMapperUsingState(state *terraform.State, resourceName string) (*keycloak.OpenIdUserClientRoleProtocolMapper, error) {
	rs, ok := state.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found in TF state: %s ", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	clientId := rs.Primary.Attributes["client_id"]
	clientScopeId := rs.Primary.Attributes["client_scope_id"]

	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	return keycloakClient.GetOpenIdUserClientRoleProtocolMapper(realm, clientId, clientScopeId, id)
}

func testKeycloakOpenIdUserClientRoleProtocolMapper_basic_client(realmName, clientId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"
	access_type = "BEARER-ONLY"
}
resource "keycloak_openid_user_client_role_protocol_mapper" "user_client_role_mapper_client" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_id        = "${keycloak_openid_client.openid_client.id}"
	claim_name       = "foo"
	claim_value_type = "String"
}`, realmName, clientId, mapperName)
}

func testKeycloakOpenIdUserClientRoleProtocolMapper_basic_clientScope(realmName, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}
resource "keycloak_openid_user_client_role_protocol_mapper" "user_client_role_mapper_client_scope" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_scope_id  = "${keycloak_openid_client_scope.client_scope.id}"
	claim_name       = "foo"
	claim_value_type = "String"
}`, realmName, clientScopeId, mapperName)
}

func testKeycloakOpenIdUserClientRoleProtocolMapper_claim(realmName, clientId, mapperName, claimName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "%s"
	access_type = "BEARER-ONLY"
}
resource "keycloak_openid_user_client_role_protocol_mapper" "user_client_role_mapper" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_id        = "${keycloak_openid_client.openid_client.id}"
	claim_name       = "%s"
	claim_value_type = "String"
}`, realmName, clientId, mapperName, claimName)
}

func testKeycloakOpenIdUserClientRoleProtocolMapper_import(realmName, clientId, clientScopeId, mapperName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_openid_client" "openid_client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"
	access_type = "BEARER-ONLY"
}
resource "keycloak_openid_user_client_role_protocol_mapper" "user_client_role_mapper_client" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_id        = "${keycloak_openid_client.openid_client.id}"
	claim_name       = "foo"
	claim_value_type = "String"
}
resource "keycloak_openid_client_scope" "client_scope" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}
resource "keycloak_openid_user_client_role_protocol_mapper" "user_client_role_mapper_client_scope" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_scope_id  = "${keycloak_openid_client_scope.client_scope.id}"
	claim_name       = "foo"
	claim_value_type = "String"
}`, realmName, clientId, mapperName, clientScopeId, mapperName)
}

func testKeycloakOpenIdUserClientRoleProtocolMapper_validateClaimValueType(realmName, mapperName, claimValueType string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}
resource "keycloak_openid_client" "openid_client" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "openid-client"
	access_type = "BEARER-ONLY"
}
resource "keycloak_openid_user_client_role_protocol_mapper" "user_client_role_mapper_validation" {
	name             = "%s"
	realm_id         = "${keycloak_realm.realm.id}"
	client_id        = "${keycloak_openid_client.openid_client.id}"
	claim_name      = "foo"
	claim_value_type = "%s"
}`, realmName, mapperName, claimValueType)
}

func testKeycloakOpenIdUserClientRoleProtocolMapper_clientAssignment(realmName, clientId, assignedClientId, mapperName string) string {
	return fmt.Sprintf(`
	resource "keycloak_realm" "realm" {
		realm = "%s"
	}
	
	resource "keycloak_openid_client" "openid_client" {
		realm_id  = "${keycloak_realm.realm.id}"
		client_id = "%s"
	
		access_type = "BEARER-ONLY"
	}
	resource "keycloak_openid_client" "openid_client_assigned" {
		realm_id  = "${keycloak_realm.realm.id}"
		client_id = "%s"
	
		access_type = "BEARER-ONLY"
	}
	
	resource "keycloak_openid_user_client_role_protocol_mapper" "user_client_role_mapper_validation" {
		name             = "%s"
		realm_id         = "${keycloak_realm.realm.id}"
		client_id        = "${keycloak_openid_client.openid_client.id}"
	
		claim_name      = "foo"
		claim_value_type = "String"
		client_id_for_role_mappings = "${keycloak_openid_client.openid_client_assigned.id}"
	}`, realmName, clientId, assignedClientId, mapperName)
}

func testKeycloakOpenIdUserClientRoleProtocolMapper_clientAssignmentRolePrefix(realmName, clientId, assignedClientId, mapperName, rolePrefix string) string {
	return fmt.Sprintf(`
	resource "keycloak_realm" "realm" {
		realm = "%s"
	}
	
	resource "keycloak_openid_client" "openid_client" {
		realm_id  = "${keycloak_realm.realm.id}"
		client_id = "%s"
	
		access_type = "BEARER-ONLY"
	}
	resource "keycloak_openid_client" "openid_client_assigned" {
		realm_id  = "${keycloak_realm.realm.id}"
		client_id = "%s"
	
		access_type = "BEARER-ONLY"
	}
	
	resource "keycloak_openid_user_client_role_protocol_mapper" "user_client_role_mapper_validation" {
		name             = "%s"
		realm_id         = "${keycloak_realm.realm.id}"
		client_id        = "${keycloak_openid_client.openid_client.id}"
	
		claim_name      = "foo"
		claim_value_type = "String"
		client_id_for_role_mappings = "${keycloak_openid_client.openid_client_assigned.id}"
		client_role_prefix= "%s"
	}`, realmName, clientId, assignedClientId, mapperName, rolePrefix)
}
