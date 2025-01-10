use chrono::{Duration, Utc};
use dioxus::prelude::{server_context, ServerFnError};
use oauth2::{
    basic::BasicClient, AuthUrl, AuthorizationCode, ClientId, ClientSecret, CsrfToken, RedirectUrl,
    Scope, TokenResponse, TokenUrl,
};
use std::env;
use thiserror::Error;
use uuid::Uuid;

use crate::structs::Profile;

#[derive(Error, Debug)]
enum AuthError {
    #[error("failed to fetch token")]
    FetchTokenFailed,

    #[error("failed to fetch profile")]
    FetchProfileFailed,
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
    let token = client
        .exchange_code(AuthorizationCode::new(code))
        .request_async(oauth2::reqwest::async_http_client)
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

fn create_client() -> BasicClient {
    let client_id = ClientId::new(env::var("LINE_CHANNEL_ID").unwrap());
    let client_secret = ClientSecret::new(env::var("LINE_CHANNEL_SECRET").unwrap());
    let redirect_url = RedirectUrl::new(env::var("REDIRECT_URL").unwrap()).unwrap();

    BasicClient::new(
        client_id,
        Some(client_secret),
        AuthUrl::new(env::var("LINE_API_AUTHORIZE").unwrap()).unwrap(),
        Some(TokenUrl::new(env::var("LINE_API_TOKEN").unwrap()).unwrap()),
    )
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
