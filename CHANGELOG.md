## 1.4.0 (January 28, 2019)

FEATURES:

* new resource: keycloak_saml_user_property_protocol_mapper ([#85](https://github.com/charlesderek/terraform-w-keycloak/pull/85))

## 1.3.0 (January 25, 2019)

FEATURES:

* new resource: keycloak_saml_user_attribute_protocol_mapper ([#84](https://github.com/charlesderek/terraform-w-keycloak/pull/84))

## 1.2.0 (January 24, 2019)

FEATURES:

* new resource: keycloak_saml_client ([#82](https://github.com/charlesderek/terraform-w-keycloak/pull/82))

IMPROVEMENTS:

* add validation for usernames to ensure they are always lowercase ([#83](https://github.com/charlesderek/terraform-w-keycloak/pull/83))

## 1.1.0 (January 7, 2019)

IMPROVEMENTS:

* openid_client: add web_origins attribute ([#81](https://github.com/charlesderek/terraform-w-keycloak/pull/81))
* user: add initial_password attribute ([#77](https://github.com/charlesderek/terraform-w-keycloak/pull/77))

BUG FIXES:

* ldap mappers: don't assume component fields are returned by Keycloak API ([#80](https://github.com/charlesderek/terraform-w-keycloak/pull/80))

## 1.0.0 (December 16, 2018)

Initial Release!

Docs: https://charlesderek.github.io/terraform-w-keycloak
