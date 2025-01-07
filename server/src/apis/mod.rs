use dioxus::prelude::*;

#[server]
pub async fn db_health_check() -> Result<String, ServerFnError> {
    let db = crate::db::get().await;
    let result = sqlx::query!("SELECT 'health!' AS text")
        .fetch_one(db)
        .await?
        .text
        .unwrap();
    Ok(result)
}
