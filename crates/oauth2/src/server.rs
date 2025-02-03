use chrono::{Duration, Utc};
use dioxus::prelude::{server_context, ServerFnError};
use oauth2::{
    basic::{
        BasicClient, BasicErrorResponse, BasicRevocationErrorResponse,
        BasicTokenIntrospectionResponse, BasicTokenResponse,
    },
    AuthUrl, AuthorizationCode, Client, ClientId, ClientSecret, CsrfToken, EndpointNotSet,
    EndpointSet, RedirectUrl, Scope, StandardRevocableToken, TokenResponse, TokenUrl,
};
use std::env;
use thiserror::Error;
use uuid::Uuid;

use crate::models::Profile;

#[derive(Error, Debug)]
enum AuthError {
    #[error("failed to fetch token")]
    FetchTokenFailed,

    #[error("failed to fetch profile")]
    FetchProfileFailed,
}

// TODO: impl this (get `session_id:Profile` from redis)
pub fn get_current_user() -> Option<String> {
    get_current_session_id()
}

pub fn line_auth() -> String {
    let client = create_client();
    let session_id = Uuid::new_v4();

    let (auth_url, _csrf_token) = client
        .authorize_url(CsrfToken::new_random)
        .add_scope(Scope::new("profile".to_string()))
        .add_scope(Scope::new("openid".to_string()))
        .url();

    set_cookie(session_id);
    auth_url.to_string()
}

pub async fn line_callback(code: String) -> Result<Profile, ServerFnError> {
    let client = create_client();

    let http_client = reqwest::ClientBuilder::new()
        // Following redirects opens the client up to SSRF vulnerabilities.
        .redirect(reqwest::redirect::Policy::none())
        .build()
        .unwrap();

    let token = client
        .exchange_code(AuthorizationCode::new(code))
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
