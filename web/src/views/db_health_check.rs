use dioxus::prelude::*;

use server::apis::db_health_check;

#[component]
pub fn DbHealthCheck() -> Element {
    let is_health = use_server_future(|| async { db_health_check().await.unwrap() })?;

    rsx! {
        p { "{is_health().unwrap()}" }
    }
}
