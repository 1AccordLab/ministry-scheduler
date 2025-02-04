use axum::{
    async_trait,
    extract::{FromRef, FromRequestParts, Query, State},
    http::{request::Parts, StatusCode},
    response::{IntoResponse, Redirect, Response},
    Json,
};
use axum_extra::extract::{cookie::Cookie, CookieJar};
use oauth2::{
    basic::{
        BasicClient, BasicErrorResponse, BasicRevocationErrorResponse,
        BasicTokenIntrospectionResponse, BasicTokenResponse,
    },
    AuthUrl, AuthorizationCode, Client, ClientId, ClientSecret, CsrfToken, EndpointNotSet,
    EndpointSet, PkceCodeChallenge, PkceCodeVerifier, RedirectUrl, Scope, StandardRevocableToken,
    TokenResponse, TokenUrl,
};
use serde::Deserialize;
use serde_json::json;
use std::{collections::HashMap, env, sync::Arc};
use thiserror::Error;
use tokio::sync::Mutex;
use uuid::Uuid;

use crate::models::Profile;

pub type SessionStore<SessionId = String> = Arc<Mutex<HashMap<SessionId, AuthState>>>;

pub struct AuthState {
    profile: Option<Profile>,
    csrf_token: CsrfToken,
    pkce_code_verifier: String,
}

#[derive(Error, Debug)]
pub enum AuthError {
    #[error("failed to fetch token")]
    FetchTokenFailed,

    #[error("failed to fetch profile")]
    FetchProfileFailed,

    #[error("no session retrieved from cookie")]
    NoSessionFromCookie,

    #[error("no session found in store")]
    NoSessionInStore,

    #[error("csrf token mismatch")]
    CsrfTokenMismatch,
}

#[derive(Deserialize)]
pub struct Oauth2CallbackRequest {
    code: String,
    state: String,
}

pub async fn get_profile(ProfileExtractor(profile): ProfileExtractor) -> Json<Profile> {
    Json(profile)
}

pub async fn line_auth(
    State(session_store): State<SessionStore>,
    jar: CookieJar,
) -> (CookieJar, Redirect) {
    let client = create_client();
    let session_id = Uuid::new_v4();
    let (pkce_code_challenge, pkce_code_verifier) = PkceCodeChallenge::new_random_sha256();

    let (auth_url, csrf_token) = client
        .authorize_url(CsrfToken::new_random)
        .set_pkce_challenge(pkce_code_challenge)
        .add_scope(Scope::new("profile".to_string()))
        .add_scope(Scope::new("openid".to_string()))
        .url();

    session_store.lock().await.insert(
        session_id.to_string(),
        AuthState {
            profile: None,
            csrf_token,
            pkce_code_verifier: pkce_code_verifier.into_secret(),
        },
    );

    (
        jar.add(Cookie::build(("session_id", session_id.to_string())).path("/")),
        Redirect::temporary(auth_url.as_str()),
    )
}

pub async fn line_callback(
    Query(params): Query<Oauth2CallbackRequest>,
    State(session_store): State<SessionStore>,
    jar: CookieJar,
) -> Result<Redirect, AuthError> {
    let session_id = jar
        .get("session_id")
        .ok_or(AuthError::NoSessionFromCookie)?
        .value();
    let mut auth_state = session_store.lock().await;
    let auth_state = auth_state
        .get_mut(session_id)
        .ok_or(AuthError::NoSessionInStore)?;

    if auth_state.csrf_token.secret().ne(&params.state) {
        return Err(AuthError::CsrfTokenMismatch);
    }

    let client = create_client();

    let http_client = reqwest::ClientBuilder::new()
        // Following redirects opens the client up to SSRF vulnerabilities.
        .redirect(reqwest::redirect::Policy::none())
        .build()
        .unwrap();

    let token = client
        .exchange_code(AuthorizationCode::new(params.code))
        .set_pkce_verifier(PkceCodeVerifier::new(auth_state.pkce_code_verifier.clone()))
        .request_async(&http_client)
        .await
        .map_err(|_| AuthError::FetchTokenFailed)?;

    let profile: Profile = reqwest::Client::new()
        .get(env::var("LINE_API_PROFILE").unwrap())
        .bearer_auth(token.access_token().secret())
        .send()
        .await
        .map_err(|_| AuthError::FetchProfileFailed)?
        .json()
        .await
        .map_err(|_| AuthError::FetchProfileFailed)?;

    auth_state.profile = Some(profile.clone());

    Ok(Redirect::temporary("/profile"))
}

fn create_client() -> LineOAuthClient {
    let client_id = ClientId::new(env::var("LINE_CHANNEL_ID").unwrap());
    let client_secret = ClientSecret::new(env::var("LINE_CHANNEL_SECRET").unwrap());
    let auth_url = AuthUrl::new(env::var("LINE_API_AUTHORIZE").unwrap()).unwrap();
    let token_url = TokenUrl::new(env::var("LINE_API_TOKEN").unwrap()).unwrap();
    let redirect_url = RedirectUrl::new(env::var("REDIRECT_URL").unwrap()).unwrap();

    BasicClient::new(client_id)
        .set_client_secret(client_secret)
        .set_auth_uri(auth_url)
        .set_token_uri(token_url)
        .set_redirect_uri(redirect_url)
}

impl IntoResponse for AuthError {
    fn into_response(self) -> Response {
        let status = match self {
            AuthError::FetchTokenFailed => StatusCode::INTERNAL_SERVER_ERROR,
            AuthError::FetchProfileFailed => StatusCode::INTERNAL_SERVER_ERROR,

            AuthError::NoSessionFromCookie => StatusCode::UNAUTHORIZED,
            AuthError::NoSessionInStore => StatusCode::UNAUTHORIZED,

            AuthError::CsrfTokenMismatch => StatusCode::FORBIDDEN,
        };

        let body = json!({
            "error": self.to_string(),
        });

        (status, Json(body)).into_response()
    }
}

pub struct AuthRedirect;

impl IntoResponse for AuthRedirect {
    fn into_response(self) -> Response {
        Redirect::temporary("/oauth2/line/login").into_response()
    }
}

pub struct ProfileExtractor(Profile);

#[async_trait]
impl<S: Send + Sync> FromRequestParts<S> for ProfileExtractor
where
    SessionStore: FromRef<S>,
{
    type Rejection = AuthRedirect;

    async fn from_request_parts(parts: &mut Parts, state: &S) -> Result<Self, Self::Rejection> {
        let jar = CookieJar::from_request_parts(parts, state)
            .await
            .map_err(|_| AuthRedirect)?;

        let session_id = jar.get("session_id").ok_or(AuthRedirect)?.value();

        let session_store = SessionStore::from_ref(state);

        let profile = session_store
            .lock()
            .await
            .get(session_id)
            .ok_or(AuthRedirect)?
            .profile
            .clone()
            .ok_or(AuthRedirect)?;

        Ok(ProfileExtractor(profile))
    }
}

// I think the oauth2 crate should have a way to make this easier
type LineOAuthClient<
    HasAuthUrl = EndpointSet,
    HasDeviceAuthUrl = EndpointNotSet,
    HasIntrospectionUrl = EndpointNotSet,
    HasRevocationUrl = EndpointNotSet,
    HasTokenUrl = EndpointSet,
> = Client<
    BasicErrorResponse,
    BasicTokenResponse,
    BasicTokenIntrospectionResponse,
    StandardRevocableToken,
    BasicRevocationErrorResponse,
    HasAuthUrl,
    HasDeviceAuthUrl,
    HasIntrospectionUrl,
    HasRevocationUrl,
    HasTokenUrl,
>;
