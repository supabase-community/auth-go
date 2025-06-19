package auth

import (
	"net/http"

	"github.com/supabase-community/auth-go/types"
)

// Create a new client using auth.New, then you can call the methods below.
//
// Some methods require bearer token authentication. To set the bearer token,
// use the WithToken(token) method.
type Client interface {

	// Options:

	// By default, the client will use the supabase project reference and assume
	// you are connecting to the Auth server as part of a supabase project.
	// To connect to a Auth server hosted elsewhere, you can specify a custom
	// URL using this method.
	//
	// It returns a copy of the client, so only requests made with the returned
	// copy will use the new URL.
	//
	// This method does not validate the URL you pass in.
	WithCustomAuthURL(url string) Client
	// WithToken sets the access token to pass when making admin requests that
	// require token authentication.
	//
	// It returns a copy of the client, so only requests made with the returned
	// copy will use the new token.
	//
	// The token can be your service role if running on a server.
	// REMEMBER TO KEEP YOUR SERVICE ROLE TOKEN SECRET!!!
	//
	// A user token can also be used to make requests on behalf of a user. This is
	// usually preferable to using a service role token.
	WithToken(token string) Client
	// WithClient allows you to pass in your own HTTP client.
	//
	// It returns a copy of the client, so only requests made with the returned
	// copy will use the new HTTP client.
	WithClient(client http.Client) Client

	// Endpoints:

	// GET /admin/audit
	//
	// Get audit logs.
	//
	// May optionally specify a query to use for filtering the audit logs. The
	// column and value must be specified if using a query.
	//
	// The result may also be paginated. By default, 50 results will be returned
	// per request. This can be configured with PerPage in the request. The response
	// will include the total number of results, as well as the total number of pages
	// and, if not already on the last page, the next page number.
	AdminAudit(req types.AdminAuditRequest) (*types.AdminAuditResponse, error)

	// POST /admin/generate_link
	//
	// Returns the corresponding email action link based on the type specified.
	// Among other things, the response also contains the query params of the action
	// link as separate JSON fields for convenience (along with the email OTP from
	// which the corresponding token is generated).
	//
	// Requires admin token.
	AdminGenerateLink(req types.AdminGenerateLinkRequest) (*types.AdminGenerateLinkResponse, error)

	// GET /admin/sso/providers
	//
	// Get a list of all SAML SSO Identity Providers in the system.
	AdminListSSOProviders() (*types.AdminListSSOProvidersResponse, error)
	// POST /admin/sso/providers
	//
	// Create a new SAML SSO Identity Provider.
	AdminCreateSSOProvider(req types.AdminCreateSSOProviderRequest) (*types.AdminCreateSSOProviderResponse, error)
	// GET /admin/sso/providers/{idp_id}
	//
	// Get a SAML SSO Identity Provider by ID.
	AdminGetSSOProvider(req types.AdminGetSSOProviderRequest) (*types.AdminGetSSOProviderResponse, error)
	// PUT /admin/sso/providers/{idp_id}
	//
	// Update a SAML SSO Identity Provider by ID.
	AdminUpdateSSOProvider(req types.AdminUpdateSSOProviderRequest) (*types.AdminUpdateSSOProviderResponse, error)
	// DELETE /admin/sso/providers/{idp_id}
	//
	// Delete a SAML SSO Identity Provider by ID.
	AdminDeleteSSOProvider(req types.AdminDeleteSSOProviderRequest) (*types.AdminDeleteSSOProviderResponse, error)

	// POST /admin/users
	//
	// Creates the user based on the user_id specified.
	//
	// Requires admin token.
	AdminCreateUser(req types.AdminCreateUserRequest) (*types.AdminCreateUserResponse, error)
	// GET /admin/users
	//
	// Get a list of users.
	//
	// Requires admin token.
	AdminListUsers(req types.AdminListUsersRequest) (*types.AdminListUsersResponse, error)
	// GET /admin/users/{user_id}
	//
	// Get a user by their user_id.
	AdminGetUser(req types.AdminGetUserRequest) (*types.AdminGetUserResponse, error)
	// PUT /admin/users/{user_id}
	//
	// Update a user by their user_id.
	AdminUpdateUser(req types.AdminUpdateUserRequest) (*types.AdminUpdateUserResponse, error)
	// DELETE /admin/users/{user_id}
	//
	// Delete a user by their user_id.
	AdminDeleteUser(req types.AdminDeleteUserRequest) error

	// GET /admin/users/{user_id}/factors
	//
	// Get a list of factors for a user.
	AdminListUserFactors(req types.AdminListUserFactorsRequest) (*types.AdminListUserFactorsResponse, error)
	// PUT /admin/users/{user_id}/factors/{factor_id}
	//
	// Update a factor for a user.
	AdminUpdateUserFactor(req types.AdminUpdateUserFactorRequest) (*types.AdminUpdateUserFactorResponse, error)
	// DELETE /admin/users/{user_id}/factors/{factor_id}
	//
	// Delete a factor for a user.
	AdminDeleteUserFactor(req types.AdminDeleteUserFactorRequest) error

	// GET /authorize
	//
	// Get access_token from external oauth provider.
	//
	// Scopes are optional additional scopes depending on the provider (email and
	// name are requested by default).
	//
	// If successful, the server returns a redirect response. This method will not
	// follow the redirect, but instead returns the URL the client was told to
	// redirect to.
	Authorize(req types.AuthorizeRequest) (*types.AuthorizeResponse, error)

	// POST /factors
	//
	// Enroll a new factor.
	EnrollFactor(req types.EnrollFactorRequest) (*types.EnrollFactorResponse, error)
	// POST /factors/{factor_id}/challenge
	//
	// Challenge a factor.
	ChallengeFactor(req types.ChallengeFactorRequest) (*types.ChallengeFactorResponse, error)
	// POST /factors/{factor_id}/verify
	//
	// Verify the challenge for an enrolled factor.
	VerifyFactor(req types.VerifyFactorRequest) (*types.VerifyFactorResponse, error)
	// DELETE /factors/{factor_id}
	//
	// Unenroll an enrolled factor.
	UnenrollFactor(req types.UnenrollFactorRequest) (*types.UnenrollFactorResponse, error)

	// GET/POST /callback
	//
	// Callback endpoint for external oauth providers to redirect to.
	//
	// There is no meaningful implementation of this as a client method, so it is
	// not included here.

	// GET /health
	//
	// Check the health of the Auth server.
	HealthCheck() (*types.HealthCheckResponse, error)

	// POST /invite
	//
	// Invites a new user with an email.
	//
	// Requires service_role or admin token.
	Invite(req types.InviteRequest) (*types.InviteResponse, error)

	// POST /logout
	//
	// Logout a user (Requires authentication).
	//
	// This will revoke all refresh tokens for the user. Remember that the JWT
	// tokens will still be valid for stateless auth until they expires.
	Logout() error

	// POST /magiclink
	//
	// DEPRECATED: Use /otp with Email and CreateUser=true instead of /magiclink.
	//
	// Magic Link. Will deliver a link (e.g.
	// /verify?type=magiclink&token=abcdefghijklmno) to the user based on email
	// address which they can use to redeem an access_token.
	//
	// By default Magic Links can only be sent once every 60 seconds.
	Magiclink(req types.MagiclinkRequest) error
	// POST /otp
	// One-Time-Password. Will deliver a magiclink or SMS OTP to the user depending
	// on whether the request contains an email or phone key.
	//
	// If CreateUser is true, the user will be automatically signed up if the user
	// doesn't exist.
	OTP(req types.OTPRequest) error

	// GET /reauthenticate
	//
	// Sends a nonce to the user's email (preferred) or phone. This endpoint
	// requires the user to be logged in / authenticated first. The user needs to
	// have either an email or phone number for the nonce to be sent successfully.
	Reauthenticate() error

	// POST /recover
	//
	// Password recovery. Will deliver a password recovery mail to the user based
	// on email address.
	//
	// By default recovery links can only be sent once every 60 seconds.
	Recover(req types.RecoverRequest) error

	// GET /settings
	//
	// Returns the publicly available settings for this auth instance.
	GetSettings() (*types.SettingsResponse, error)

	// POST /signup
	//
	// Register a new user with an email and password.
	Signup(req types.SignupRequest) (*types.SignupResponse, error)

	// Sign in with email and password
	//
	// This is a convenience method that calls Token with the password grant type
	SignInWithEmailPassword(email, password string) (*types.TokenResponse, error)
	// Sign in with phone and password
	//
	// This is a convenience method that calls Token with the password grant type
	SignInWithPhonePassword(phone, password string) (*types.TokenResponse, error)
	// Sign in with id token
	//
	// This is a convenience method that calls Token with the id_token grant type
	SignInWithIdToken(provider, idToken, nonce, accessToken, captchaToken string) (*types.TokenResponse, error)
	// Sign in with refresh token
	//
	// This is a convenience method that calls Token with the refresh_token grant type
	RefreshToken(refreshToken string) (*types.TokenResponse, error)
	// POST /token
	//
	// This is an OAuth2 endpoint that currently implements the password, id_token,
	//  and refresh_token grant types
	Token(req types.TokenRequest) (*types.TokenResponse, error)

	// GET /user
	//
	// Get the JSON object for the logged in user (requires authentication)
	GetUser() (*types.UserResponse, error)
	// PUT /user
	//
	// Update a user (Requires authentication). Apart from changing email/password,
	// this method can be used to set custom user data. Changing the email will
	// result in a magiclink being sent out.
	UpdateUser(req types.UpdateUserRequest) (*types.UpdateUserResponse, error)

	// GET /verify
	//
	// Verify a registration or a password recovery. Type can be signup or recovery
	// or magiclink or invite and the token is a token returned from either /signup
	// or /recover or /magiclink.
	//
	// The server returns a redirect response. This method will not follow the
	// redirect, but instead returns the URL the client was told to redirect to,
	// as well as parsing the parameters from the URL fragment.
	//
	// NOTE: This endpoint may return a nil error, but the Response can contain
	// error details extracted from the returned URL. Please check that the Error,
	// ErrorCode and/or ErrorDescription fields of the response are empty.
	Verify(req types.VerifyRequest) (*types.VerifyResponse, error)
	// POST /verify
	//
	// Verify a registration or a password recovery. Type can be signup or recovery
	// or magiclink or invite and the token is a token returned from either /signup
	// or /recover or /magiclink.
	//
	// This differs from GET /verify as it requires an email or phone to be given,
	// which is used to verify the token associated to the user. It also returns a
	// JSON response rather than a redirect.
	VerifyForUser(req types.VerifyForUserRequest) (*types.VerifyForUserResponse, error)

	// GET /sso/saml/metadata
	//
	// Get the SAML metadata for the configured SAML provider.
	//
	// If successful, the server returns an XML response. Making sense of this is
	// outside the scope of this client, so it is simply returned as []byte.
	SAMLMetadata() ([]byte, error)
	// POST /sso/saml/acs
	//
	// Implements the main Assertion Consumer Service endpoint behavior.
	//
	// This client does not provide a typed endpoint for SAML ACS. This method is
	// provided for convenience and will simply POST your HTTP request to the
	// endpoint and return the response.
	//
	// For required parameters, see the SAML spec or the Auth implementation
	// of this endpoint.
	//
	// The server may issue redirects. Using the default HTTP client, this method
	// will follow those redirects and return the final HTTP response. Should you
	// prefer the client not to follow redirects, you can provide a custom HTTP
	// client using WithClient(). See the example below.
	//
	// Example:
	//	c := http.Client{
	//		CheckRedirect: func(req *http.Request, via []*http.Request) error {
	//			return http.ErrUseLastResponse
	//		},
	//	}
	SAMLACS(req *http.Request) (*http.Response, error)

	// POST /sso
	//
	// Initiate an SSO session with the given provider.
	//
	// If successful, the server returns a redirect to the provider's authorization
	// URL. The client will follow it and return the final HTTP response.
	//
	// Auth allows you to skip following the redirect by setting SkipHTTPRedirect
	// on the request struct. In this case, the URL to redirect to will be returned
	// in the response.
	SSO(req types.SSORequest) (*types.SSOResponse, error)
}
