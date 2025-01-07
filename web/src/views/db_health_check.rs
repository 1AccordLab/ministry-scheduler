use dioxus::prelude::*;

use server::apis::db_health_check;

#[component]
pub fn DbHealthCheck() -> Element {
    let mut msg = use_signal(String::new);

    rsx! {
        div {
            button {
                onclick: move |_| async move {
                    let data = db_health_check().await.unwrap();
                    msg.set(data);
                },
                "click me to check db health"
            }

            if !msg().is_empty() {
                p { "{msg}" }
            }
        }
    }
}
