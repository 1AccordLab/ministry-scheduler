use dioxus::prelude::{server, server_fn, ServerFnError};

#[cfg(feature = "server")]
use crate::server;

use crate::models::Profile;

#[server]
pub async fn line_auth() -> Result<String, ServerFnError> {
    let auth_url = server::line_auth().await;
    Ok(auth_url)
}

#[server]
pub async fn line_callback(code: String, state: String) -> Result<Profile, ServerFnError> {
    let profile = server::line_callback(code, state).await?;
    Ok(profile)
}
