use dioxus::prelude::*;

use crate::views::Hero;

#[component]
pub fn Home() -> Element {
    rsx! {
        Hero {}
    }
}
