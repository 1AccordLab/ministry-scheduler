use dioxus::prelude::*;

use server::apis::oauth2::line_auth;

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
