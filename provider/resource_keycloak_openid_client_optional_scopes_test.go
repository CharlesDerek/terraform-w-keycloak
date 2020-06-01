package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
	"regexp"
	"strings"
	"testing"
)

// All openid clients in Keycloak will automatically have these scopes listed as "optional client scopes".
func getPreAssignedOptionalClientScopes(t *testing.T) []string {
	keycloakVersionIsGreaterThanOrEqualTo6, err := keycloakVersionIsGreaterThanOrEqualTo(keycloakClient, getKeycloakVersion600())
	if err != nil {
		t.Fatal(err)
	}
	if keycloakVersionIsGreaterThanOrEqualTo6 {
		return []string{"address", "phone", "offline_access", "microprofile-jwt"}
	} else {
		return []string{"address", "phone", "offline_access"}
	}
}

func TestAccKeycloakOpenidClientOptionalScopes_basic(t *testing.T) {
	realm := "terraform-realm-" + acctest.RandString(10)
	client := "terraform-client-" + acctest.RandString(10)
	clientScope := "terraform-client-scope-" + acctest.RandString(10)

	clientScopes := append(getPreAssignedOptionalClientScopes(t), clientScope)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientOptionalScopes_basic(realm, client, clientScope),
				Check:  testAccCheckKeycloakOpenidClientHasOptionalScopes("keycloak_openid_client_optional_scopes.optional_scopes", clientScopes),
			},
			// we need a separate test step for destroy instead of using CheckDestroy because this resource is implicitly
			// destroyed at the end of each test via destroying clients
			{
				Config: testKeycloakOpenidClientOptionalScopes_noOptionalScopes(realm, client, clientScope),
				Check:  testAccCheckKeycloakOpenidClientHasNoOptionalScopes("keycloak_openid_client.client"),
			},
		},
	})
}

func TestAccKeycloakOpenidClientOptionalScopes_updateClientForceNew(t *testing.T) {
	realm := "terraform-realm-" + acctest.RandString(10)
	clientOne := "terraform-client-" + acctest.RandString(10)
	clientTwo := "terraform-client-" + acctest.RandString(10)
	clientScope := "terraform-client-scope-" + acctest.RandString(10)

	clientScopes := append(getPreAssignedOptionalClientScopes(t), clientScope)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientOptionalScopes_basic(realm, clientOne, clientScope),
				Check:  testAccCheckKeycloakOpenidClientHasOptionalScopes("keycloak_openid_client_optional_scopes.optional_scopes", clientScopes),
			},
			{
				Config: testKeycloakOpenidClientOptionalScopes_basic(realm, clientTwo, clientScope),
				Check:  testAccCheckKeycloakOpenidClientHasOptionalScopes("keycloak_openid_client_optional_scopes.optional_scopes", clientScopes),
			},
		},
	})
}

func TestAccKeycloakOpenidClientOptionalScopes_updateInPlace(t *testing.T) {
	realm := "terraform-realm-" + acctest.RandString(10)
	client := "terraform-client-" + acctest.RandString(10)
	clientScope := "terraform-client-scope-" + acctest.RandString(10)

	allClientScopes := append(getPreAssignedOptionalClientScopes(t), clientScope)

	clientScopeToRemove := allClientScopes[acctest.RandIntRange(0, 2)]
	var subsetOfClientScopes []string
	for _, cs := range allClientScopes {
		if cs != clientScopeToRemove {
			subsetOfClientScopes = append(subsetOfClientScopes, cs)
		}
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// init
			{
				Config: testKeycloakOpenidClientOptionalScopes_listOfScopes(realm, client, clientScope, allClientScopes),
				Check:  testAccCheckKeycloakOpenidClientHasOptionalScopes("keycloak_openid_client_optional_scopes.optional_scopes", allClientScopes),
			},
			// remove
			{
				Config: testKeycloakOpenidClientOptionalScopes_listOfScopes(realm, client, clientScope, subsetOfClientScopes),
				Check:  testAccCheckKeycloakOpenidClientHasOptionalScopes("keycloak_openid_client_optional_scopes.optional_scopes", subsetOfClientScopes),
			},
			// add
			{
				Config: testKeycloakOpenidClientOptionalScopes_listOfScopes(realm, client, clientScope, allClientScopes),
				Check:  testAccCheckKeycloakOpenidClientHasOptionalScopes("keycloak_openid_client_optional_scopes.optional_scopes", allClientScopes),
			},
		},
	})
}

func TestAccKeycloakOpenidClientOptionalScopes_validateClientDoesNotExist(t *testing.T) {
	realm := "terraform-realm-" + acctest.RandString(10)
	client := "terraform-client-" + acctest.RandString(10)
	clientScope := "terraform-client-scope-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenidClientOptionalScopes_validationNoClient(realm, client, clientScope),
				ExpectError: regexp.MustCompile("validation error: client with id .+ does not exist"),
			},
		},
	})
}

func TestAccKeycloakOpenidClientOptionalScopes_validateClientAccessType(t *testing.T) {
	realm := "terraform-realm-" + acctest.RandString(10)
	client := "terraform-client-" + acctest.RandString(10)
	clientScope := "terraform-client-scope-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakOpenidClientOptionalScopes_validationBearerOnlyClient(realm, client, clientScope),
				ExpectError: regexp.MustCompile("validation error: client with id .+ uses access type BEARER-ONLY which does not use scopes"),
			},
		},
	})
}

// if a optional client scope is manually detached from a client with optional scopes controlled by this resource, terraform should add it again
func TestAccKeycloakOpenidClientOptionalScopes_authoritativeAdd(t *testing.T) {
	realm := "terraform-realm-" + acctest.RandString(10)
	client := "terraform-client-" + acctest.RandString(10)
	clientScopes := append(getPreAssignedOptionalClientScopes(t),
		"terraform-client-scope-"+acctest.RandString(10),
		"terraform-client-scope-"+acctest.RandString(10),
		"terraform-client-scope-"+acctest.RandString(10),
	)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientOptionalScopes_multipleClientScopes(realm, client, clientScopes, clientScopes),
				Check:  testAccCheckKeycloakOpenidClientHasOptionalScopes("keycloak_openid_client_optional_scopes.optional_scopes", clientScopes),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					client, err := keycloakClient.GetOpenidClientByClientId(realm, client)
					if err != nil {
						t.Fatal(err)
					}

					clientToManuallyDetach := clientScopes[acctest.RandIntRange(0, len(clientScopes)-1)]
					err = keycloakClient.DetachOpenidClientOptionalScopes(realm, client.Id, []string{clientToManuallyDetach})
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakOpenidClientOptionalScopes_multipleClientScopes(realm, client, clientScopes, clientScopes),
				Check:  testAccCheckKeycloakOpenidClientHasOptionalScopes("keycloak_openid_client_optional_scopes.optional_scopes", clientScopes),
			},
		},
	})
}

// if an optional client scope is manually attached to a client with optional scopes controlled by this resource, terraform should detach it
func TestAccKeycloakOpenidClientOptionalScopes_authoritativeRemove(t *testing.T) {
	realm := "terraform-realm-" + acctest.RandString(10)
	client := "terraform-client-" + acctest.RandString(10)

	randomClientScopes := []string{
		"terraform-client-scope-" + acctest.RandString(10),
		"terraform-client-scope-" + acctest.RandString(10),
		"terraform-client-scope-" + acctest.RandString(10),
	}
	allClientScopes := append(getPreAssignedOptionalClientScopes(t), randomClientScopes...)

	clientToManuallyAttach := randomClientScopes[acctest.RandIntRange(0, len(randomClientScopes)-1)]
	var attachedClientScopes []string
	for _, clientScope := range allClientScopes {
		if clientScope != clientToManuallyAttach {
			attachedClientScopes = append(attachedClientScopes, clientScope)
		}
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientOptionalScopes_multipleClientScopes(realm, client, allClientScopes, attachedClientScopes),
				Check:  testAccCheckKeycloakOpenidClientHasOptionalScopes("keycloak_openid_client_optional_scopes.optional_scopes", attachedClientScopes),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					client, err := keycloakClient.GetOpenidClientByClientId(realm, client)
					if err != nil {
						t.Fatal(err)
					}

					err = keycloakClient.AttachOpenidClientOptionalScopes(realm, client.Id, []string{clientToManuallyAttach})
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakOpenidClientOptionalScopes_multipleClientScopes(realm, client, allClientScopes, attachedClientScopes),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakOpenidClientHasOptionalScopes("keycloak_openid_client_optional_scopes.optional_scopes", attachedClientScopes),
					testAccCheckKeycloakOpenidClientOptionalScopeIsNotAttached("keycloak_openid_client_optional_scopes.optional_scopes", clientToManuallyAttach),
				),
			},
		},
	})
}

// this resource doesn't support import because it can be created even if the desired state already exists in keycloak
func TestAccKeycloakOpenidClientOptionalScopes_noImportNeeded(t *testing.T) {
	realm := "terraform-realm-" + acctest.RandString(10)
	client := "terraform-client-" + acctest.RandString(10)
	clientScope := "terraform-client-scope-" + acctest.RandString(10)

	clientScopes := append(getPreAssignedOptionalClientScopes(t), clientScope)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakOpenidClientOptionalScopes_noOptionalScopes(realm, client, clientScope),
				Check:  testAccCheckKeycloakOpenidClientOptionalScopeIsNotAttached("keycloak_openid_client.client", clientScope),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					openidClient, err := keycloakClient.GetOpenidClientByClientId(realm, client)
					if err != nil {
						t.Fatal(err)
					}

					err = keycloakClient.AttachOpenidClientOptionalScopes(realm, openidClient.Id, clientScopes)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakOpenidClientOptionalScopes_basic(realm, client, clientScope),
				Check:  testAccCheckKeycloakOpenidClientHasOptionalScopes("keycloak_openid_client_optional_scopes.optional_scopes", clientScopes),
			},
		},
	})
}

// by optional, keycloak clients have the optional scopes "address", "phone" and
// "offline_access" "microprofile-jwt" attached. if you create this resource with only one scope, it
// won't remove these two scopes, because the creation of a new resource should
// not result in anything destructive. thus, a following plan will not be empty,
// as terraform will think it needs to remove these scopes, which is okay to do
// during an update
func TestAccKeycloakOpenidClientOptionalScopes_profileAndEmailOptionalScopes(t *testing.T) {
	realm := "terraform-realm-" + acctest.RandString(10)
	client := "terraform-client-" + acctest.RandString(10)
	clientScope := "terraform-client-scope-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:             testKeycloakOpenidClientOptionalScopes_listOfScopes(realm, client, clientScope, []string{clientScope}),
				Check:              testAccCheckKeycloakOpenidClientHasOptionalScopes("keycloak_openid_client.client", append(getPreAssignedOptionalClientScopes(t), clientScope)),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Keycloak throws a 500 if you attempt to attach an optional scope that is already attached as a default scope
func TestAccKeycloakOpenidClientOptionalScopes_validateDuplicateScopeAssignment(t *testing.T) {
	realm := "terraform-realm-" + acctest.RandString(10)
	client := "terraform-client-" + acctest.RandString(10)
	clientScope := "terraform-client-scope-" + acctest.RandString(10)

	defaultClientScopes := append(preAssignedDefaultClientScopes, clientScope)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// attach default scopes, including the custom scope
			{
				Config: testKeycloakOpenidClientDefaultScopes_basic(realm, client, clientScope),
				Check:  testAccCheckKeycloakOpenidClientHasDefaultScopes("keycloak_openid_client_default_scopes.default_scopes", defaultClientScopes),
			},
			// attach optional scopes with the custom scope, expect an error since it is already in use
			{
				Config:      testKeycloakOpenidClientOptionalScopes_duplicateScopeAssignment(realm, client, clientScope),
				ExpectError: regexp.MustCompile("validation error: scope .+ is already attached to client as a default scope"),
			},
		},
	})
}

func getOptionalClientScopesFromState(resourceName string, s *terraform.State) ([]*keycloak.OpenidClientScope, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realm := rs.Primary.Attributes["realm_id"]

	var client string
	if strings.HasPrefix(resourceName, "keycloak_openid_client_optional_scopes") {
		client = rs.Primary.Attributes["client_id"]
	} else {
		client = rs.Primary.ID
	}

	keycloakOptionalClientScopes, err := keycloakClient.GetOpenidClientOptionalScopes(realm, client)
	if err != nil {
		return nil, err
	}

	return keycloakOptionalClientScopes, nil
}

func testAccCheckKeycloakOpenidClientHasOptionalScopes(resourceName string, tfOptionalClientScopes []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		keycloakOptionalClientScopes, err := getOptionalClientScopesFromState(resourceName, s)
		if err != nil {
			return err
		}

		for _, tfOptionalClientScope := range tfOptionalClientScopes {
			found := false

			for _, keycloakOptionalScope := range keycloakOptionalClientScopes {
				if keycloakOptionalScope.Name == tfOptionalClientScope {
					found = true

					break
				}
			}

			if !found {
				return fmt.Errorf("optional scope %s is not assigned to client", tfOptionalClientScope)
			}
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientHasNoOptionalScopes(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		keycloakOptionalClientScopes, err := getOptionalClientScopesFromState(resourceName, s)
		if err != nil {
			return err
		}

		if numberOfOptionalScopes := len(keycloakOptionalClientScopes); numberOfOptionalScopes != 0 {
			return fmt.Errorf("expected client to have no assigned optional scopes, but it has %d", numberOfOptionalScopes)
		}

		return nil
	}
}

func testAccCheckKeycloakOpenidClientOptionalScopeIsNotAttached(resourceName, clientScope string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		keycloakOptionalClientScopes, err := getOptionalClientScopesFromState(resourceName, s)
		if err != nil {
			return err
		}

		for _, keycloakOptionalClientScope := range keycloakOptionalClientScopes {
			if keycloakOptionalClientScope.Name == clientScope {
				return fmt.Errorf("expected client scope with name %s to not be attached to client", clientScope)
			}
		}

		return nil
	}
}

func testKeycloakOpenidClientOptionalScopes_basic(realm, client, clientScope string) string {
	keycloakVersionIsHigherOrEqualTo6, _ := keycloakVersionIsGreaterThanOrEqualTo(keycloakClient, getKeycloakVersion600())
	if keycloakVersionIsHigherOrEqualTo6 {
		return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "PUBLIC"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name        = "%s"
	realm_id    = "${keycloak_realm.realm.id}"

	description = "test description"
}

resource "keycloak_openid_client_optional_scopes" "optional_scopes" {
	realm_id       = "${keycloak_realm.realm.id}"
	client_id      = "${keycloak_openid_client.client.id}"
	optional_scopes = [
		"address",
		"phone",
		"offline_access",
		"microprofile-jwt",
		"${keycloak_openid_client_scope.client_scope.name}"
	]
}
	`, realm, client, clientScope)
	} else {
		return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "PUBLIC"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name        = "%s"
	realm_id    = "${keycloak_realm.realm.id}"

	description = "test description"
}

resource "keycloak_openid_client_optional_scopes" "optional_scopes" {
	realm_id       = "${keycloak_realm.realm.id}"
	client_id      = "${keycloak_openid_client.client.id}"
	optional_scopes = [
		"address",
		"phone",
		"offline_access",
		"${keycloak_openid_client_scope.client_scope.name}"
	]
}
	`, realm, client, clientScope)
	}
}

func testKeycloakOpenidClientOptionalScopes_noOptionalScopes(realm, client, clientScope string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "PUBLIC"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name        = "%s"
	realm_id    = "${keycloak_realm.realm.id}"

	description = "test description"
}
	`, realm, client, clientScope)
}

func testKeycloakOpenidClientOptionalScopes_listOfScopes(realm, client, clientScope string, listOfOptionalScopes []string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "PUBLIC"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name        = "%s"
	realm_id    = "${keycloak_realm.realm.id}"

	description = "test description"
}

resource "keycloak_openid_client_optional_scopes" "optional_scopes" {
	realm_id       = "${keycloak_realm.realm.id}"
	client_id      = "${keycloak_openid_client.client.id}"
	optional_scopes = %s

	depends_on = ["keycloak_openid_client_scope.client_scope"]
}
	`, realm, client, clientScope, arrayOfStringsForTerraformResource(listOfOptionalScopes))
}

func testKeycloakOpenidClientOptionalScopes_validationNoClient(realm, client, clientScope string) string {
	keycloakVersionIsHigherOrEqualTo6, _ := keycloakVersionIsGreaterThanOrEqualTo(keycloakClient, getKeycloakVersion600())
	if keycloakVersionIsHigherOrEqualTo6 {
		return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name        = "%s"
	realm_id    = "${keycloak_realm.realm.id}"

	description = "test description"
}

resource "keycloak_openid_client_optional_scopes" "optional_scopes" {
	realm_id       = "${keycloak_realm.realm.id}"
	client_id      = "%s"
	optional_scopes = [
		"address",
		"phone",
		"offline_access",
		"microprofile-jwt",
		"${keycloak_openid_client_scope.client_scope.name}"
	]
}
	`, realm, clientScope, client)
	} else {
		return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name        = "%s"
	realm_id    = "${keycloak_realm.realm.id}"

	description = "test description"
}

resource "keycloak_openid_client_optional_scopes" "optional_scopes" {
	realm_id       = "${keycloak_realm.realm.id}"
	client_id      = "%s"
	optional_scopes = [
		"address",
		"phone",
		"offline_access",
		"${keycloak_openid_client_scope.client_scope.name}"
	]
}
	`, realm, clientScope, client)
	}
}

func testKeycloakOpenidClientOptionalScopes_validationBearerOnlyClient(realm, client, clientScope string) string {

	keycloakVersionIsHigherOrEqualTo6, _ := keycloakVersionIsGreaterThanOrEqualTo(keycloakClient, getKeycloakVersion600())
	if keycloakVersionIsHigherOrEqualTo6 {
		return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name        = "%s"
	realm_id    = "${keycloak_realm.realm.id}"

	description = "test description"
}

resource "keycloak_openid_client_optional_scopes" "optional_scopes" {
	realm_id       = "${keycloak_realm.realm.id}"
	client_id      = "${keycloak_openid_client.client.id}"
	optional_scopes = [
		"address",
		"phone",
		"offline_access",
		"microprofile-jwt",
		"${keycloak_openid_client_scope.client_scope.name}"
	]
}
	`, realm, client, clientScope)
	} else {
		return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "BEARER-ONLY"
}

resource "keycloak_openid_client_scope" "client_scope" {
	name        = "%s"
	realm_id    = "${keycloak_realm.realm.id}"

	description = "test description"
}

resource "keycloak_openid_client_optional_scopes" "optional_scopes" {
	realm_id       = "${keycloak_realm.realm.id}"
	client_id      = "${keycloak_openid_client.client.id}"
	optional_scopes = [
		"address",
		"phone",
		"offline_access",
		"${keycloak_openid_client_scope.client_scope.name}"
	]
}
	`, realm, client, clientScope)
	}
}

func testKeycloakOpenidClientOptionalScopes_multipleClientScopes(realm, client string, allClientScopes, attachedClientScopes []string) string {
	var clientScopeResources strings.Builder
	for _, clientScope := range allClientScopes {
		if strings.HasPrefix(clientScope, "terraform") {
			clientScopeResources.WriteString(fmt.Sprintf(`
resource "keycloak_openid_client_scope" "client_scope_%s" {
	name        = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
}
		`, clientScope, clientScope))
		}
	}

	var attachedClientScopesInterpolated []string
	for _, attachedClientScope := range attachedClientScopes {
		if strings.HasPrefix(attachedClientScope, "terraform") {
			attachedClientScopesInterpolated = append(attachedClientScopesInterpolated, fmt.Sprintf("${keycloak_openid_client_scope.client_scope_%s.name}", attachedClientScope))
		} else {
			attachedClientScopesInterpolated = append(attachedClientScopesInterpolated, attachedClientScope)
		}
	}

	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "PUBLIC"
}

%s

resource "keycloak_openid_client_optional_scopes" "optional_scopes" {
	realm_id       = "${keycloak_realm.realm.id}"
	client_id      = "${keycloak_openid_client.client.id}"
	optional_scopes = %s
}
	`, realm, client, clientScopeResources.String(), arrayOfStringsForTerraformResource(attachedClientScopesInterpolated))
}

func testKeycloakOpenidClientOptionalScopes_duplicateScopeAssignment(realm, client, clientScope string) string {
	keycloakVersionIsHigherOrEqualTo6, _ := keycloakVersionIsGreaterThanOrEqualTo(keycloakClient, getKeycloakVersion600())
	if keycloakVersionIsHigherOrEqualTo6 {
		return fmt.Sprintf(`
%s

resource "keycloak_openid_client_optional_scopes" "optional_scopes" {
	realm_id       = "${keycloak_realm.realm.id}"
	client_id      = "${keycloak_openid_client.client.id}"
	optional_scopes = [
		"address",
		"phone",
		"offline_access",
		"microprofile-jwt",
		"${keycloak_openid_client_scope.client_scope.name}"
	]
}
	`, testKeycloakOpenidClientDefaultScopes_basic(realm, client, clientScope))
	} else {
		return fmt.Sprintf(`
%s

resource "keycloak_openid_client_optional_scopes" "optional_scopes" {
	realm_id       = "${keycloak_realm.realm.id}"
	client_id      = "${keycloak_openid_client.client.id}"
	optional_scopes = [
		"address",
		"phone",
		"offline_access",
		"${keycloak_openid_client_scope.client_scope.name}"
	]
}
	`, testKeycloakOpenidClientDefaultScopes_basic(realm, client, clientScope))
	}
}
