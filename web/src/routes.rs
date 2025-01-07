use dioxus::prelude::*;

use crate::{
    views::{Blog, DbHealthCheck, Home},
    Navbar,
};

#[derive(Debug, Clone, Routable, PartialEq)]
#[rustfmt::skip]
pub enum Route {
    #[layout(Navbar)]
    #[route("/")]
    Home {},

    #[route("/blog/:id")]
    Blog { id: i32 },

    #[route("/db_health_check")]
    DbHealthCheck {},
}
