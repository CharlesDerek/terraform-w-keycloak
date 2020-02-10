package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakAttributeImporterIdentityProviderMapper_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	userAttribute := "terraform-" + acctest.RandString(10)
	claimName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakAttributeImporterIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAttributeImporterIdentityProviderMapper_basic(realmName, alias, mapperName, userAttribute, claimName),
				Check:  testAccCheckKeycloakAttributeImporterIdentityProviderMapperExists("keycloak_attribute_importer_identity_provider_mapper.oidc"),
			},
		},
	})
}

func TestAccKeycloakAttributeImporterIdentityProviderMapper_createAfterManualDestroy(t *testing.T) {
	var mapper = &keycloak.IdentityProviderMapper{}

	realmName := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	userAttribute := "terraform-" + acctest.RandString(10)
	claimName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakAttributeImporterIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAttributeImporterIdentityProviderMapper_basic(realmName, alias, mapperName, userAttribute, claimName),
				Check:  testAccCheckKeycloakAttributeImporterIdentityProviderMapperFetch("keycloak_attribute_importer_identity_provider_mapper.oidc", mapper),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteIdentityProviderMapper(mapper.Realm, mapper.IdentityProviderAlias, mapper.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakAttributeImporterIdentityProviderMapper_basic(realmName, alias, mapperName, userAttribute, claimName),
				Check:  testAccCheckKeycloakAttributeImporterIdentityProviderMapperExists("keycloak_attribute_importer_identity_provider_mapper.oidc"),
			},
		},
	})
}

func TestAccKeycloakAttributeImporterIdentityProviderMapper_basicUpdateRealm(t *testing.T) {
	firstRealm := "terraform-" + acctest.RandString(10)
	secondRealm := "terraform-" + acctest.RandString(10)
	mapperName := "terraform-" + acctest.RandString(10)
	alias := "terraform-" + acctest.RandString(10)
	userAttribute := "terraform-" + acctest.RandString(10)
	claimName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakAttributeImporterIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAttributeImporterIdentityProviderMapper_basic(firstRealm, alias, mapperName, userAttribute, claimName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAttributeImporterIdentityProviderMapperExists("keycloak_attribute_importer_identity_provider_mapper.oidc"),
					resource.TestCheckResourceAttr("keycloak_attribute_importer_identity_provider_mapper.oidc", "realm", firstRealm),
				),
			},
			{
				Config: testKeycloakAttributeImporterIdentityProviderMapper_basic(secondRealm, alias, mapperName, userAttribute, claimName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAttributeImporterIdentityProviderMapperExists("keycloak_attribute_importer_identity_provider_mapper.oidc"),
					resource.TestCheckResourceAttr("keycloak_attribute_importer_identity_provider_mapper.oidc", "realm", secondRealm),
				),
			},
		},
	})
}

func TestAccKeycloakAttributeImporterIdentityProviderMapper_basicUpdateAll(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	identityProviderAliasName := "terraform-" + acctest.RandString(10)

	firstMapper := &keycloak.IdentityProviderMapper{
		Realm:                 realmName,
		IdentityProviderAlias: identityProviderAliasName,
		Name:                  acctest.RandString(10),
		Config: &keycloak.IdentityProviderMapperConfig{
			UserAttribute: acctest.RandString(10),
			Attribute:     acctest.RandString(10),
		},
	}

	secondMapper := &keycloak.IdentityProviderMapper{
		Realm:                 realmName,
		IdentityProviderAlias: identityProviderAliasName,
		Name:                  acctest.RandString(10),
		Config: &keycloak.IdentityProviderMapperConfig{
			UserAttribute: acctest.RandString(10),
			Attribute:     acctest.RandString(10),
		},
	}

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakAttributeImporterIdentityProviderMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAttributeImporterIdentityProviderMapper_basicFromInterface(firstMapper),
				Check:  testAccCheckKeycloakAttributeImporterIdentityProviderMapperExists("keycloak_attribute_importer_identity_provider_mapper.saml"),
			},
			{
				Config: testKeycloakAttributeImporterIdentityProviderMapper_basicFromInterface(secondMapper),
				Check:  testAccCheckKeycloakAttributeImporterIdentityProviderMapperExists("keycloak_attribute_importer_identity_provider_mapper.saml"),
			},
		},
	})
}

func testAccCheckKeycloakAttributeImporterIdentityProviderMapperExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getKeycloakAttributeImporterIdentityProviderMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakAttributeImporterIdentityProviderMapperFetch(resourceName string, mapper *keycloak.IdentityProviderMapper) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedMapper, err := getKeycloakAttributeImporterIdentityProviderMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		mapper.IdentityProviderAlias = fetchedMapper.IdentityProviderAlias
		mapper.Realm = fetchedMapper.Realm
		mapper.Id = fetchedMapper.Id

		return nil
	}
}

func testAccCheckKeycloakAttributeImporterIdentityProviderMapperDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_attribute_importer_identity_provider_mapper" {
				continue
			}

			realm := rs.Primary.Attributes["realm"]
			alias := rs.Primary.Attributes["identity_provider_alias"]
			id := rs.Primary.ID

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			mapper, _ := keycloakClient.GetIdentityProviderMapper(realm, alias, id)
			if mapper != nil {
				return fmt.Errorf("oidc config with id %s still exists", id)
			}
		}

		return nil
	}
}

func getKeycloakAttributeImporterIdentityProviderMapperFromState(s *terraform.State, resourceName string) (*keycloak.IdentityProviderMapper, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realm := rs.Primary.Attributes["realm"]
	alias := rs.Primary.Attributes["identity_provider_alias"]
	id := rs.Primary.ID

	mapper, err := keycloakClient.GetIdentityProviderMapper(realm, alias, id)
	if err != nil {
		return nil, fmt.Errorf("error getting identity provider mapper config with id %s: %s", id, err)
	}

	return mapper, nil
}

func testKeycloakAttributeImporterIdentityProviderMapper_basic(realm, alias, name, userAttribute, claimName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_oidc_identity_provider" "oidc" {
	realm             = "${keycloak_realm.realm.id}"
	alias             = "%s"
	authorization_url = "https://example.com/auth"
	token_url         = "https://example.com/token"
	client_id         = "example_id"
	client_secret     = "example_token"
}

resource keycloak_attribute_importer_identity_provider_mapper oidc {
	realm                   = "${keycloak_realm.realm.id}"
	name                    = "%s"
	identity_provider_alias = "${keycloak_oidc_identity_provider.oidc.alias}"
	user_attribute          = "%s"
	claim_name              = "%s"
}
	`, realm, alias, name, userAttribute, claimName)
}

func testKeycloakAttributeImporterIdentityProviderMapper_basicFromInterface(mapper *keycloak.IdentityProviderMapper) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_saml_identity_provider" "saml" {
	realm                      = "${keycloak_realm.realm.id}"
	alias                      = "%s"
	single_sign_on_service_url = "https://example.com/auth"
}

resource keycloak_attribute_importer_identity_provider_mapper saml {
	realm                   = "${keycloak_realm.realm.id}"
	name                    = "%s"
	identity_provider_alias = "${keycloak_saml_identity_provider.saml.alias}"
	attribute_name          = "%s"
	user_attribute          = "%s"
}
	`, mapper.Realm, mapper.IdentityProviderAlias, mapper.Name, mapper.Config.Attribute, mapper.Config.UserAttribute)
}
