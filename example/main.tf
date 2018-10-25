provider "keycloak" {
  client_id     = "terraform"
  client_secret = "d0e0f95-0f42-4a63-9b1f-94274655669e"
  url           = "http://localhost:8080"
}

resource "keycloak_realm" "test" {
  realm                = "test"
  enabled              = true
  display_name         = "foo"

  account_theme        = "base"

  access_code_lifespan = "30m"
}

resource "keycloak_openid_client" "test_client" {
  client_id           = "test-client"
  name                = "test-client"
  realm_id            = "${keycloak_realm.test.id}"
  description         = "a test client"

  access_type         = "CONFIDENTIAL"
  valid_redirect_uris = [
    "http://localhost:8080/callback"
  ]
}

resource "keycloak_openid_client_scope" "test_client_scope" {
  name                = "foo1"
  realm_id            = "${keycloak_realm.test.id}"

  description         = "test"
  consent_screen_text = "hello"
}

resource "keycloak_ldap_user_federation" "openldap" {
  name                    = "openldap"
  realm_id                = "master"

  enabled                 = true

  username_ldap_attribute = "cn"
  rdn_ldap_attribute      = "cn"
  uuid_ldap_attribute     = "entryDN"
  user_object_classes     = [
    "simpleSecurityObject",
    "organizationalRole"
  ]
  connection_url          = "ldap://openldap"
  users_dn                = "dc=example,dc=org"
  bind_dn                 = "cn=admin,dc=example,dc=org"
  bind_credential         = "admin"
}

resource "keycloak_ldap_user_attribute_mapper" "attr_mapper" {
  name                    = "test mapper"
  realm_id                = "${keycloak_ldap_user_federation.openldap.realm_id}"
  ldap_user_federation_id = "${keycloak_ldap_user_federation.openldap.id}"

  user_model_attribute    = "foo"
  ldap_attribute          = "bar"
}

resource "keycloak_ldap_group_mapper" "group_mapper" {
  name                           = "group mapper"
  realm_id                       = "${keycloak_ldap_user_federation.openldap.realm_id}"
  ldap_user_federation_id        = "${keycloak_ldap_user_federation.openldap.id}"

  ldap_groups_dn                 = "dc=example,dc=org"
  group_name_ldap_attribute      = "cn"
  group_object_classes           = [
    "groupOfNames"
  ]
  membership_attribute_type      = "DN"
  membership_ldap_attribute      = "member"
  membership_user_ldap_attribute = "cn"
  memberof_ldap_attribute        = "memberOf"
}

resource "keycloak_ldap_msad_user_account_control_mapper" "msad_uac_mapper" {
  name                    = "uac-mapper1"
  realm_id                = "${keycloak_ldap_user_federation.openldap.realm_id}"
  ldap_user_federation_id = "${keycloak_ldap_user_federation.openldap.id}"
}

resource "keycloak_ldap_full_name_mapper" "full_name_mapper" {
  name                     = "full-name-mapper"
  realm_id                 = "${keycloak_ldap_user_federation.openldap.realm_id}"
  ldap_user_federation_id  = "${keycloak_ldap_user_federation.openldap.id}"

  ldap_full_name_attribute = "cn"
  read_only                = true
}

resource "keycloak_custom_user_federation" "custom" {
  name        = "custom1"
  realm_id    = "master"
  provider_id = "custom"

  enabled     = true
}

resource "keycloak_openid_user_attribute_protocol_mapper" "map_user_attributes_client" {
  name            = "tf-test-open-id-user-attribute-protocol-mapper-client"
  realm_id        = "${keycloak_realm.test.id}"
  client_id       = "${keycloak_openid_client.test_client.id}"
  user_attribute  = "foo"
  claim_name      = "bar"
}

resource "keycloak_openid_user_attribute_protocol_mapper" "map_user_attributes_client_scope" {
  name            = "tf-test-open-id-user-attribute-protocol-mapper-client-scope"
  realm_id        = "${keycloak_realm.test.id}"
  client_scope_id = "${keycloak_openid_client_scope.test_client_scope.id}"
  user_attribute  = "foo2"
  claim_name      = "bar2"
}

resource "keycloak_openid_group_membership_protocol_mapper" "map_group_memberships_client" {
  name            = "tf-test-open-id-group-membership-protocol-mapper-client"
  realm_id        = "${keycloak_realm.test.id}"
  client_id       = "${keycloak_openid_client.test_client.id}"
  claim_name      = "bar"
}

resource "keycloak_openid_group_membership_protocol_mapper" "map_group_memberships_client_scope" {
  name            = "tf-test-open-id-group-membership-protocol-mapper-client-scope"
  realm_id        = "${keycloak_realm.test.id}"
  client_scope_id = "${keycloak_openid_client_scope.test_client_scope.id}"
  claim_name      = "bar2"
}
