package types

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// --- Errors ---

type ErrInvalidGenerateLinkRequest struct {
	message string
}

func NewErrInvalidGenerateLinkRequest(message string) *ErrInvalidGenerateLinkRequest {
	return &ErrInvalidGenerateLinkRequest{message: message}
}

func (e *ErrInvalidGenerateLinkRequest) Error() string {
	return fmt.Sprintf("generate link request is invalid - %s", e.message)
}

var (
	ErrInvalidAdminAuditRequest        = errors.New("admin audit request is invalid - if Query is not nil, then query Column must be author, action or type, and value must be given")
	ErrInvalidAdminUpdateFactorRequest = errors.New("admin update factor request is invalid - nothing to update")
	ErrInvalidTokenRequest             = errors.New("token request is invalid - grant_type must be either one of password, refresh_token, or pkce, email and password must be provided for grant_type=password, refresh_token must be provided for grant_type=refresh_token, auth_code and code_verifier must be provided for grant_type=pkce")
	ErrInvalidVerifyRequest            = errors.New("verify request is invalid - type, token and redirect_to must be provided, and email or phone must be provided to VerifyForUser")
	ErrInvalidProviderRequest          = errors.New("provider must be one of: github, apple, kakao, google, keycloak")
)

// --- Request/Response Types ---
type AuditQueryColumn string

const (
	AuditQueryColumnAuthor AuditQueryColumn = "author"
	AuditQueryColumnAction AuditQueryColumn = "action"
	AuditQueryColumnType   AuditQueryColumn = "type"
)

type AuditQuery struct {
	Column AuditQueryColumn
	Value  string
}

type AdminAuditRequest struct {
	// Query, if provided, is used to search the audit log.
	// Logs will be returned where the chosen column matches the value.
	Query *AuditQuery

	// Pagination
	Page    uint
	PerPage uint
}

type AuditLogEntry struct {
	ID        uuid.UUID              `json:"id" db:"id"`
	Payload   map[string]interface{} `json:"payload" db:"payload"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
	IPAddress string                 `json:"ip_address" db:"ip_address"`
}

type AdminAuditResponse struct {
	Logs []AuditLogEntry

	// Pagination
	TotalCount int
	TotalPages uint
	NextPage   uint
}

type LinkType string

const (
	LinkTypeSignup             LinkType = "signup"
	LinkTypeMagicLink          LinkType = "magiclink"
	LinkTypeRecovery           LinkType = "recovery"
	LinkTypeInvite             LinkType = "invite"
	LinkTypeEmailChangeCurrent LinkType = "email_change_current"
	LinkTypeEmailChangeNew     LinkType = "email_change_new"
)

type AdminGenerateLinkRequest struct {
	Type       LinkType               `json:"type"`
	Email      string                 `json:"email"`
	NewEmail   string                 `json:"new_email"`
	Password   string                 `json:"password"`
	Data       map[string]interface{} `json:"data"`
	RedirectTo string                 `json:"redirect_to"`
}

type AdminGenerateLinkResponse struct {
	ActionLink       string   `json:"action_link"`
	EmailOTP         string   `json:"email_otp"`
	HashedToken      string   `json:"hashed_token"`
	RedirectTo       string   `json:"redirect_to"`
	VerificationType LinkType `json:"verification_type"`

	User
}

type AdminCreateUserRequest struct {
	Aud          string                 `json:"aud,omitempty"`
	Role         string                 `json:"role,omitempty"`
	Email        string                 `json:"email,omitempty"`
	Phone        string                 `json:"phone,omitempty"`
	Password     *string                `json:"password,omitempty"` // Only if type = signup
	EmailConfirm bool                   `json:"email_confirm,omitempty"`
	PhoneConfirm bool                   `json:"phone_confirm,omitempty"`
	UserMetadata map[string]interface{} `json:"user_metadata,omitempty"`
	AppMetadata  map[string]interface{} `json:"app_metadata,omitempty"`
	BanDuration  time.Duration          `json:"ban_duration,omitempty"` // Cannot be "none" when creating a user, so just set it or leave it empty
}

type AdminCreateUserResponse struct {
	User
}

type AdminListUsersRequest struct {
	Page    *int
	PerPage *int
}

type AdminListUsersResponse struct {
	Users []User `json:"users"`
}

type AdminGetUserRequest struct {
	UserID uuid.UUID
}

type AdminGetUserResponse struct {
	User
}

type AdminUpdateUserRequest struct {
	UserID uuid.UUID `json:"-"`

	Aud          string                 `json:"aud,omitempty"`
	Role         string                 `json:"role,omitempty"`
	Email        string                 `json:"email,omitempty"`
	Phone        string                 `json:"phone,omitempty"`
	Password     string                 `json:"password,omitempty"`
	EmailConfirm bool                   `json:"email_confirm,omitempty"`
	PhoneConfirm bool                   `json:"phone_confirm,omitempty"`
	UserMetadata map[string]interface{} `json:"user_metadata,omitempty"`
	AppMetadata  map[string]interface{} `json:"app_metadata,omitempty"`
	BanDuration  *BanDuration           `json:"ban_duration,omitempty"`
}

type AdminUpdateUserResponse struct {
	User
}

type AdminDeleteUserRequest struct {
	UserID uuid.UUID
}

type AdminListUserFactorsRequest struct {
	UserID uuid.UUID
}

type AdminListUserFactorsResponse struct {
	Factors []Factor
}

type AdminUpdateUserFactorRequest struct {
	UserID   uuid.UUID `json:"-"`
	FactorID uuid.UUID `json:"-"`

	FriendlyName string `json:"friendly_name,omitempty"`
}

type AdminUpdateUserFactorResponse struct {
	Factor
}

type AdminDeleteUserFactorRequest struct {
	UserID   uuid.UUID
	FactorID uuid.UUID
}

type SAMLAttribute struct {
	Name    string      `json:"name,omitempty"`
	Names   []string    `json:"names,omitempty"`
	Default interface{} `json:"default,omitempty"`
}

type SAMLAttributeMapping struct {
	Keys map[string]SAMLAttribute `json:"keys,omitempty"`
}

type SAMLProvider struct {
	EntityID    string  `json:"entity_id"`
	MetadataXML string  `json:"metadata_xml,omitempty"`
	MetadataURL *string `json:"metadata_url,omitempty"`

	AttributeMapping SAMLAttributeMapping `json:"attribute_mapping,omitempty"`
}

type SSODomain struct {
	Domain string `db:"domain" json:"domain"`
}

type SSOProvider struct {
	ID           uuid.UUID    `json:"id"`
	ResourceID   *string      `json:"resource_id,omitempty"`
	SAMLProvider SAMLProvider `json:"saml,omitempty"`
	SSODomains   []SSODomain  `json:"domains"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type AdminListSSOProvidersResponse struct {
	Providers []SSOProvider `json:"items"`
}

type AdminCreateSSOProviderRequest struct {
	ResourceID       string               `json:"resource_id"`
	Type             string               `json:"type"`
	MetadataURL      string               `json:"metadata_url"`
	MetadataXML      string               `json:"metadata_xml"`
	Domains          []string             `json:"domains"`
	AttributeMapping SAMLAttributeMapping `json:"attribute_mapping"`
}

type AdminCreateSSOProviderResponse struct {
	SSOProvider
}

type AdminGetSSOProviderRequest struct {
	ProviderID uuid.UUID
}

type AdminGetSSOProviderResponse struct {
	SSOProvider
}

type AdminUpdateSSOProviderRequest struct {
	ProviderID uuid.UUID `json:"-"`

	ResourceID       string               `json:"resource_id"`
	Type             string               `json:"type"`
	MetadataURL      string               `json:"metadata_url"`
	MetadataXML      string               `json:"metadata_xml"`
	Domains          []string             `json:"domains"`
	AttributeMapping SAMLAttributeMapping `json:"attribute_mapping"`
}

type AdminUpdateSSOProviderResponse struct {
	SSOProvider
}

type AdminDeleteSSOProviderRequest struct {
	ProviderID uuid.UUID
}

type AdminDeleteSSOProviderResponse struct {
	SSOProvider
}

type Provider string
type FlowType string

const (
	ProviderApple     Provider = "apple"
	ProviderAzure     Provider = "azure"
	ProviderBitbucket Provider = "bitbucket"
	ProviderDiscord   Provider = "discord"
	ProviderGitHub    Provider = "github"
	ProviderGitLab    Provider = "gitlab"
	ProviderGoogle    Provider = "google"
	ProviderKeycloak  Provider = "keycloak"
	ProviderLinkedin  Provider = "linkedin"
	ProviderFacebook  Provider = "facebook"
	ProviderNotion    Provider = "notion"
	ProviderSpotify   Provider = "spotify"
	ProviderSlack     Provider = "slack"
	ProviderTwitch    Provider = "twitch"
	ProviderTwitter   Provider = "twitter"
	ProviderWorkOS    Provider = "workos"
	ProviderZoom      Provider = "zoom"
)

const (
	FlowImplicit FlowType = "implicit"
	FlowPKCE     FlowType = "pkce"
)

type AuthorizeRequest struct {
	Provider   Provider
	RedirectTo string
	FlowType   FlowType
	Scopes     string
}

type AuthorizeResponse struct {
	AuthorizationURL string
	Verifier         string
}

// adapted from https://go-review.googlesource.com/c/oauth2/+/463979/9/pkce.go#64
type PKCEParams struct {
	Challenge       string
	ChallengeMethod string
	Verifier        string
}

type FactorType string

const FactorTypeTOTP FactorType = "totp"

type EnrollFactorRequest struct {
	FriendlyName string     `json:"friendly_name"`
	FactorType   FactorType `json:"factor_type"`
	Issuer       string     `json:"issuer"`
}

type TOTPObject struct {
	QRCode string `json:"qr_code"`
	Secret string `json:"secret"`
	URI    string `json:"uri"`
}

type EnrollFactorResponse struct {
	ID   uuid.UUID  `json:"id"`
	Type FactorType `json:"type"`
	TOTP TOTPObject `json:"totp,omitempty"`
}

type ChallengeFactorRequest struct {
	FactorID uuid.UUID `json:"factor_id"`
}

type ChallengeFactorResponse struct {
	ID        uuid.UUID `json:"id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type VerifyFactorRequest struct {
	FactorID uuid.UUID `json:"-"`

	ChallengeID uuid.UUID `json:"challenge_id"`
	Code        string    `json:"code"`
}

type VerifyFactorResponse struct {
	Session
}

type UnenrollFactorRequest struct {
	FactorID uuid.UUID
}

type UnenrollFactorResponse struct {
	ID uuid.UUID `json:"id"`
}

type HealthCheckResponse struct {
	Version     string `json:"version"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type InviteRequest struct {
	Email      string                 `json:"email"`
	Data       map[string]interface{} `json:"data"`
	RedirectTo string                 `json:"-"`
}

type InviteResponse struct {
	User
}

// DEPRECATED: Use /otp with Email and CreateUser=true instead of /magiclink.
type MagiclinkRequest struct {
	Email string `json:"email"`

	// Provide Captcha token if enabled.
	SecurityEmbed
}

type OTPRequest struct {
	Email           string                 `json:"email"`
	Phone           string                 `json:"phone"`
	CreateUser      bool                   `json:"create_user"`
	Data            map[string]interface{} `json:"data"`
	RedirectTo      string                 `json:"-"`
	EmailRedirectTo string                 `json:"-"`

	// Provide Captcha token if enabled.
	SecurityEmbed
}

type RecoverRequest struct {
	Email      string `json:"email"`
	RedirectTo string `json:"-"`

	// Provide Captcha token if enabled.
	SecurityEmbed
}

type ResendRequest struct {
	Email           string `json:"email,omitempty"`
	Phone           string `json:"phone,omitempty"`
	EmailRedirectTo string `json:"-"`
}

type ExternalProviders struct {
	Apple     bool `json:"apple"`
	Azure     bool `json:"azure"`
	Bitbucket bool `json:"bitbucket"`
	Discord   bool `json:"discord"`
	Email     bool `json:"email"`
	Facebook  bool `json:"facebook"`
	GitHub    bool `json:"github"`
	GitLab    bool `json:"gitlab"`
	Google    bool `json:"google"`
	Keycloak  bool `json:"keycloak"`
	Linkedin  bool `json:"linkedin"`
	Notion    bool `json:"notion"`
	Phone     bool `json:"phone"`
	SAML      bool `json:"saml"`
	Slack     bool `json:"slack"`
	Spotify   bool `json:"spotify"`
	Twitch    bool `json:"twitch"`
	Twitter   bool `json:"twitter"`
	WorkOS    bool `json:"workos"`
	Zoom      bool `json:"zoom"`
}

type SettingsResponse struct {
	DisableSignup     bool              `json:"disable_signup"`
	Autoconfirm       bool              `json:"autoconfirm"`
	MailerAutoconfirm bool              `json:"mailer_autoconfirm"`
	PhoneAutoconfirm  bool              `json:"phone_autoconfirm"`
	SmsProvider       string            `json:"sms_provider"`
	MFAEnabled        bool              `json:"mfa_enabled"`
	External          ExternalProviders `json:"external"`
}

type SignupRequest struct {
	Email    string                 `json:"email,omitempty"`
	Phone    string                 `json:"phone,omitempty"`
	Password string                 `json:"password,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`

	// Provide Captcha token if enabled.
	SecurityEmbed
}

type SignupResponse struct {
	// Response if autoconfirm is off
	User

	// Response if autoconfirm is on
	Session
}

type SSORequest struct {
	// Use either ProviderID or Domain.
	ProviderID       uuid.UUID `json:"provider_id"`
	Domain           string    `json:"domain"`
	RedirectTo       string    `json:"redirect_to"`
	SkipHTTPRedirect bool      `json:"skip_http_redirect"`

	// Provide Captcha token if enabled.
	SecurityEmbed
}

type SSOResponse struct {
	// Returned only if SkipHTTPRedirect was set in request.
	URL string `json:"url"`

	// Returned otherwise.
	HTTPResponse *http.Response `json:"-"`
}

type TokenRequest struct {
	GrantType string `json:"-"`

	// Email or Phone, and Password, are required if GrantType is 'password'.
	// They must not be provided if GrantType is 'refresh_token'.
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Password string `json:"password,omitempty"`

	// RefreshToken is required if GrantType is 'refresh_token'.
	// It must not be provided if GrantType is 'password'.
	RefreshToken string `json:"refresh_token,omitempty"`

	// Code and CodeVerifier are required if GrantType is 'pkce'.
	Code         string `json:"auth_code,omitempty"`
	CodeVerifier string `json:"code_verifier,omitempty"`

	// Provider and IdToken are required if GrantType is 'id_token'.
	Provider    string `json:"provider,omitempty"`
	IdToken     string `json:"id_token,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	Nonce       string `json:"nonce,omitempty"`

	// Provide Captcha token if enabled. Not required if GrantType is 'refresh_token'.
	SecurityEmbed
}

type TokenResponse struct {
	Session
}

type UserResponse struct {
	User
}

type UpdateUserRequest struct {
	Email    string                 `json:"email,omitempty"`
	Password *string                `json:"password,omitempty"`
	Nonce    string                 `json:"nonce,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
	AppData  map[string]interface{} `json:"app_metadata,omitempty"`
	Phone    string                 `json:"phone,omitempty"`
}

type UpdateUserResponse struct {
	User
}

type VerificationType string

const (
	VerificationTypeSignup      = "signup"
	VerificationTypeRecovery    = "recovery"
	VerificationTypeInvite      = "invite"
	VerificationTypeMagiclink   = "magiclink"
	VerificationTypeEmailChange = "email_change"
	VerificationTypeSMS         = "sms"
	VerificationTypePhoneChange = "phone_change"
)

type VerifyRequest struct {
	Type       VerificationType
	Token      string
	RedirectTo string
}

type VerifyResponse struct {
	URL string

	// The fields below are returned only for a successful response.
	AccessToken  string
	TokenType    string
	ExpiresIn    int
	RefreshToken string
	Type         VerificationType

	// The fields below are returned if there was an error verifying.
	Error            string
	ErrorCode        string
	ErrorDescription string
}

type VerifyForUserRequest struct {
	Type       VerificationType `json:"type"`
	Token      string           `json:"token"`
	RedirectTo string           `json:"redirect_to"`
	Email      string           `json:"email"`
	Phone      string           `json:"phone"`

	// Provide Captcha token if enabled.
	// Not required for server version >= v2.30.1
	SecurityEmbed
}

type VerifyForUserResponse struct {
	Session
}

type ProviderType string

const (
	ProviderTypeOAuth2 ProviderType = "oauth2"
	ProviderTypeOIDC   ProviderType = "oidc"
)

type OIDCDiscovery struct {
	Issuer                 string   `json:"issuer"`
	AuthorizationEndpoint  string   `json:"authorization_endpoint"`
	TokenEndpoint          string   `json:"token_endpoint"`
	UserinfoEndpoint       string   `json:"userinfo_endpoint,omitempty"`
	JwksURI                string   `json:"jwks_uri"`
	ScopesSupported        []string `json:"scopes_supported,omitempty"`
	ResponseTypesSupported []string `json:"response_types_supported,omitempty"`
	GrantTypesSupported    []string `json:"grant_types_supported,omitempty"`
	SubjectTypesSupported  []string `json:"subject_types_supported,omitempty"`
}

// CustomOAuthProvider is a custom OAuth2 or OIDC provider configuration.
//
// The client secret is never returned by the server.
type CustomOAuthProvider struct {
	ID           uuid.UUID               `json:"id"`
	ProviderType ProviderType `json:"provider_type"`

	Identifier            string                 `json:"identifier"`
	Name                  string                 `json:"name"`
	ClientID              string                 `json:"client_id"`
	AcceptableClientIDs   []string               `json:"acceptable_client_ids"`
	Scopes                []string               `json:"scopes"`
	PKCEEnabled           bool                   `json:"pkce_enabled"`
	AttributeMapping      map[string]interface{} `json:"attribute_mapping"`
	CustomClaimsAllowlist []string               `json:"custom_claims_allowlist"`
	AuthorizationParams   map[string]interface{} `json:"authorization_params"`
	Enabled               bool                   `json:"enabled"`
	EmailOptional         bool                   `json:"email_optional"`

	// OIDC only
	Issuer            *string                `json:"issuer,omitempty"`
	DiscoveryURL      *string                `json:"discovery_url,omitempty"`
	SkipNonceCheck    bool                   `json:"skip_nonce_check"`
	CachedDiscovery *OIDCDiscovery `json:"discovery_document,omitempty"`

	// OAuth2 only
	AuthorizationURL *string `json:"authorization_url,omitempty"`
	TokenURL         *string `json:"token_url,omitempty"`
	UserinfoURL      *string `json:"userinfo_url,omitempty"`
	JwksURI          *string `json:"jwks_uri,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AdminListCustomOAuthProvidersRequest struct {
	ProviderType ProviderType `json:"-"`
}

type AdminListCustomOAuthProvidersResponse struct {
	Providers []CustomOAuthProvider `json:"providers"`
}

type AdminCreateCustomOAuthProviderRequest struct {
	ProviderType ProviderType `json:"provider_type"`

	// Identifier must be prefixed with "custom:", e.g. "custom:acme".
	Identifier            string                 `json:"identifier"`
	Name                  string                 `json:"name"`
	ClientID              string                 `json:"client_id"`
	ClientSecret          string                 `json:"client_secret"`
	AcceptableClientIDs   []string               `json:"acceptable_client_ids,omitempty"`
	Scopes                []string               `json:"scopes,omitempty"`
	PKCEEnabled           *bool                  `json:"pkce_enabled,omitempty"`
	AttributeMapping      map[string]interface{} `json:"attribute_mapping,omitempty"`
	CustomClaimsAllowlist []string               `json:"custom_claims_allowlist,omitempty"`
	AuthorizationParams   map[string]interface{} `json:"authorization_params,omitempty"`
	Enabled               *bool                  `json:"enabled,omitempty"`
	EmailOptional         *bool                  `json:"email_optional,omitempty"`

	// OIDC only
	Issuer         string  `json:"issuer,omitempty"`
	DiscoveryURL   *string `json:"discovery_url,omitempty"`
	SkipNonceCheck *bool   `json:"skip_nonce_check,omitempty"`

	// OAuth2 only
	AuthorizationURL string  `json:"authorization_url,omitempty"`
	TokenURL         string  `json:"token_url,omitempty"`
	UserinfoURL      string  `json:"userinfo_url,omitempty"`
	JwksURI          *string `json:"jwks_uri,omitempty"`
}

type AdminCreateCustomOAuthProviderResponse struct {
	CustomOAuthProvider
}

type AdminGetCustomOAuthProviderRequest struct {
	Identifier string `json:"-"`
}

type AdminGetCustomOAuthProviderResponse struct {
	CustomOAuthProvider
}

type AdminUpdateCustomOAuthProviderRequest struct {
	// Identifier of the provider to update, including the "custom:" prefix.
	// Sent in the URL, not the body.
	Identifier string `json:"-"`

	Name                  string                 `json:"name,omitempty"`
	ClientID              string                 `json:"client_id,omitempty"`
	ClientSecret          string                 `json:"client_secret,omitempty"`
	AcceptableClientIDs   []string               `json:"acceptable_client_ids,omitempty"`
	Scopes                []string               `json:"scopes,omitempty"`
	PKCEEnabled           *bool                  `json:"pkce_enabled,omitempty"`
	AttributeMapping      map[string]interface{} `json:"attribute_mapping,omitempty"`
	CustomClaimsAllowlist []string               `json:"custom_claims_allowlist,omitempty"`
	AuthorizationParams   map[string]interface{} `json:"authorization_params,omitempty"`
	Enabled               *bool                  `json:"enabled,omitempty"`
	EmailOptional         *bool                  `json:"email_optional,omitempty"`

	// OIDC only
	Issuer         string  `json:"issuer,omitempty"`
	DiscoveryURL   *string `json:"discovery_url,omitempty"`
	SkipNonceCheck *bool   `json:"skip_nonce_check,omitempty"`

	// OAuth2 only
	AuthorizationURL string  `json:"authorization_url,omitempty"`
	TokenURL         string  `json:"token_url,omitempty"`
	UserinfoURL      string  `json:"userinfo_url,omitempty"`
	JwksURI          *string `json:"jwks_uri,omitempty"`
}

type AdminUpdateCustomOAuthProviderResponse struct {
	CustomOAuthProvider
}

type AdminDeleteCustomOAuthProviderRequest struct {
	Identifier string `json:"-"`
}
