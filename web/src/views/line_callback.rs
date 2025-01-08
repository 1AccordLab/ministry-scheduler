use dioxus::prelude::*;

use server::apis::oauth2::line_callback;

#[component]
pub fn LineCallBack(code: String) -> Element {
    let profile = use_resource(move || {
        let code = code.clone();
        async { line_callback(code).await.unwrap() }
    });

    match profile() {
        Some(profile) => rsx! { "{profile:?}" },
        None => rsx! { "no profile fetched" },
    }
}
