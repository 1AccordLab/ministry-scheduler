use chrono::{Duration, Utc};
use dioxus::prelude::{server_context, ServerFnError};
use oauth2::{
    basic::{
        BasicClient, BasicErrorResponse, BasicRevocationErrorResponse,
        BasicTokenIntrospectionResponse, BasicTokenResponse,
    },
    AuthUrl, AuthorizationCode, Client, ClientId, ClientSecret, CsrfToken, EndpointNotSet,
    EndpointSet, PkceCodeChallenge, PkceCodeVerifier, RedirectUrl, Scope, StandardRevocableToken,
    TokenResponse, TokenUrl,
};
use std::{
    collections::HashMap,
    env,
    sync::{Arc, LazyLock},
};
use thiserror::Error;
use tokio::sync::Mutex;
use uuid::Uuid;

use crate::models::Profile;

type SessionStore<SessionId = String> = Arc<Mutex<HashMap<SessionId, AuthState>>>;
static SESSION_STORE: LazyLock<SessionStore> = LazyLock::new(Default::default);

struct AuthState {
    profile: Option<Profile>,
    csrf_token: CsrfToken,
    pkce_code_verifier: String,
}

#[derive(Error, Debug)]
enum AuthError {
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

// TODO: impl this (get `session_id:Profile` from redis)
pub fn get_current_user() -> Option<String> {
    get_current_session_id()
}

pub async fn line_auth() -> String {
    let client = create_client();
    let session_id = Uuid::new_v4();
    let (pkce_code_challenge, pkce_code_verifier) = PkceCodeChallenge::new_random_sha256();

    let (auth_url, csrf_token) = client
        .authorize_url(CsrfToken::new_random)
        .set_pkce_challenge(pkce_code_challenge)
        .add_scope(Scope::new("profile".to_string()))
        .add_scope(Scope::new("openid".to_string()))
        .url();

    SESSION_STORE.lock().await.insert(
        session_id.to_string(),
        AuthState {
            profile: None,
            csrf_token,
            pkce_code_verifier: pkce_code_verifier.into_secret(),
        },
    );

    set_cookie(session_id);
    auth_url.to_string()
}

pub async fn line_callback(code: String, state: String) -> Result<Profile, ServerFnError> {
    let session_id = get_current_session_id().ok_or(AuthError::NoSessionFromCookie)?;
    let mut auth_state = SESSION_STORE.lock().await;
    let auth_state = auth_state
        .get_mut(&session_id)
        .ok_or(AuthError::NoSessionInStore)?;

    if auth_state.csrf_token.secret().ne(&state) {
        return Err(ServerFnError::new(AuthError::CsrfTokenMismatch));
    }

    let client = create_client();

    let http_client = reqwest::ClientBuilder::new()
        // Following redirects opens the client up to SSRF vulnerabilities.
        .redirect(reqwest::redirect::Policy::none())
        .build()
        .unwrap();

    let token = client
        .exchange_code(AuthorizationCode::new(code))
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
        .await?;

    auth_state.profile = Some(profile.clone());

    Ok(profile)
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

fn set_cookie(session_id: Uuid) {
    let expiration_date = Utc::now() + Duration::days(30);
    let expires = expiration_date.to_rfc2822();

    server_context().response_parts_mut().headers.insert(
        "Set-Cookie",
        format!("session_id={session_id}; Path=/; HttpOnly; Secure; Expires={expires}")
            .parse()
            .unwrap(),
    );
}

fn get_current_session_id() -> Option<String> {
    server_context()
        .request_parts()
        .headers
        .get("Cookie")?
        .to_str()
        .ok()?
        .split(';')
        .map(str::trim)
        .find_map(|cookie| cookie.split('=').last())
        .map(String::from)
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
