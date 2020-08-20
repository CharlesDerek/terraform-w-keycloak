# keycloak_saml_client

Allows for creating and managing Keycloak clients that use the SAML protocol.

Clients are entities that can use Keycloak for user authentication. Typically,
clients are applications that redirect users to Keycloak for authentication
in order to take advantage of Keycloak's user sessions for SSO.

### Example Usage

```hcl
resource "keycloak_realm" "realm" {
    realm   = "my-realm"
    enabled = true
}

resource "keycloak_saml_client" "saml_client" {
	realm_id                = keycloak_realm.realm.id
	client_id               = "test-saml-client"
	name                    = "test-saml-client"

	sign_documents          = false
	sign_assertions         = true
	include_authn_statement = true

	signing_certificate = "${file("saml-cert.pem")}"
	signing_private_key = "${file("saml-key.pem")}"
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm this client is attached to.
- `client_id` - (Required) The unique ID of this client, referenced in the URI during authentication and in issued tokens.
- `name` - (Optional) The display name of this client in the GUI.
- `enabled` - (Optional) When false, this client will not be able to initiate a login or obtain access tokens. Defaults to `true`.
- `description` - (Optional) The description of this client in the GUI.
- `include_authn_statement` - (Optional) When `true`, an `AuthnStatement` will be included in the SAML response.
- `sign_documents` - (Optional) When `true`, the SAML document will be signed by Keycloak using the realm's private key.
- `sign_assertions` - (Optional) When `true`, the SAML assertions will be signed by Keycloak using the realm's private key, and embedded within the SAML XML Auth response.
- `encrypt_assertions` - (Optional) When `true`, the SAML assertions will be encrypted by Keycloak using the client's public key.
- `client_signature_required` - (Optional) When `true`, Keycloak will expect that documents originating from a client will be signed using the certificate and/or key configured via `signing_certificate` and `signing_private_key`.
- `force_post_binding` - (Optional) When `true`, Keycloak will always respond to an authentication request via the SAML POST Binding.
- `front_channel_logout` - (Optional) When `true`, this client will require a browser redirect in order to perform a logout.
- `name_id_format` - (Optional) Sets the Name ID format for the subject.
- `force_name_id_format` - (Optional) Ignore requested NameID subject format and use the one defined in `name_id_format` instead.
- `signature_algorithm` - (Optional) The signature algorithm used to sign documents. Should be one of "RSA_SHA1", "RSA_SHA256", "RSA_SHA512", or "DSA_SHA1".
- `root_url` - (Optional) When specified, this value is prepended to all relative URLs.
- `valid_redirect_uris` - (Optional) When specified, Keycloak will use this list to validate given Assertion Consumer URLs specified in the authentication request.
- `base_url` - (Optional) When specified, this URL will be used whenever Keycloak needs to link to this client.
- `master_saml_processing_url` - (Optional) When specified, this URL will be used for all SAML requests.
- `encryption_certificate` - (Optional) If assertions for the client are encrypted, this certificate will be used for encryption.
- `signing_certificate` - (Optional) If documents or assertions from the client are signed, this certificate will be used to verify the signature.
- `signing_private_key` - (Optional) If documents or assertions from the client are signed, this private key will be used to verify the signature.
- `idp_initiated_sso_url_name` - (Optional) URL fragment name to reference client when you want to do IDP Initiated SSO.
- `idp_initiated_sso_relay_state` - (Optional) Relay state you want to send with SAML request when you want to do IDP Initiated SSO.
- `assertion_consumer_post_url` - (Optional) SAML POST Binding URL for the client's assertion consumer service (login responses).
- `assertion_consumer_redirect_url` - (Optional) SAML Redirect Binding URL for the client's assertion consumer service (login responses).
- `logout_service_post_binding_url` - (Optional) SAML POST Binding URL for the client's single logout service.
- `logout_service_redirect_binding_url` - (Optional) SAML Redirect Binding URL for the client's single logout service.
- `full_scope_allowed` - (Optional) - Allow to include all roles mappings in the access token

### Import

Clients can be imported using the format `{{realm_id}}/{{client_keycloak_id}}`, where `client_keycloak_id` is the unique ID that Keycloak
assigns to the client upon creation. This value can be found in the URI when editing this client in the GUI, and is typically a GUID.

Example:

```bash
$ terraform import keycloak_saml_client.saml_client my-realm/dcbc4c73-e478-4928-ae2e-d5e420223352
```
