use dioxus::prelude::*;

use server::apis::oauth2::line_callback;

#[component]
pub fn LineCallBack(code: String) -> Element {
    let profile = use_server_future(move || {
        let code = code.clone();
        async { line_callback(code).await.unwrap() }
    })?;

    rsx! { "{profile().unwrap():?}" }
}
