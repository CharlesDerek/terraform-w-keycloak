provider "keycloak" {
  client_id     = "terraform"
  client_secret = "d0e0f95-0f42-4a63-9b1f-94274655669e"
  url           = "http://localhost:8080"
}

resource "keycloak_realm" "test" {
  realm                          = "test"
  enabled                        = true
  display_name                   = "foo"

  registration_allowed           = false
  registration_email_as_username = false
  edit_username_allowed          = false
  reset_password_allowed         = false
  remember_me                    = false
  verify_email                   = false
  login_with_email_allowed       = false
  duplicate_emails_allowed       = false
}

resource "keycloak_client" "test-client" {
  client_id = "test-client"
  realm_id  = "${keycloak_realm.test.id}"
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
