use dioxus::prelude::{server, server_fn, ServerFnError};
use oauth2::{
    basic::BasicClient, AuthUrl, AuthorizationCode, ClientId, ClientSecret, CsrfToken, RedirectUrl,
    Scope, TokenResponse, TokenUrl,
};
use serde::{Deserialize, Serialize};
use std::env;

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Profile {
    user_id: String,
    display_name: String,
    picture_url: String,
    status_message: String,
}

#[server]
pub async fn line_auth() -> Result<String, ServerFnError> {
    let client = create_client();
    let (auth_url, _csrf_token) = client
        .authorize_url(CsrfToken::new_random)
        .add_scope(Scope::new("profile".to_string()))
        .add_scope(Scope::new("openid".to_string()))
        .url();
    Ok(auth_url.to_string())
}

#[server]
pub async fn line_callback(code: String) -> Result<Profile, ServerFnError> {
    let client = create_client();
    let token = client
        .exchange_code(AuthorizationCode::new(code))
        .request_async(oauth2::reqwest::async_http_client)
        .await?;
    let profile: Profile = reqwest::Client::new()
        .get("https://api.line.me/v2/profile")
        .bearer_auth(token.access_token().secret())
        .send()
        .await?
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
