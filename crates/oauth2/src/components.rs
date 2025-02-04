use dioxus::prelude::*;

use crate::server_fns::{line_auth, line_callback};

#[component]
pub fn LineLogin() -> Element {
    use_future(move || async {
        let auth_url = line_auth().await.unwrap();
        web_sys::window()
            .unwrap()
            .open_with_url_and_target(&auth_url, "_self")
            .unwrap();
    });

    rsx! {}
}

#[component]
pub fn LineCallBack(code: String, state: String) -> Element {
    let profile = use_server_future(move || {
        let code = code.clone();
        let state = state.clone();
        async { line_callback(code, state).await }
    })?;

    rsx! { "{profile().unwrap():?}" }
}
