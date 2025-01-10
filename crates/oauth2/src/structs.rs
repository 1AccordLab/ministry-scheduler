use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Profile {
    user_id: String,
    display_name: String,
    picture_url: String,
    status_message: String,
}
