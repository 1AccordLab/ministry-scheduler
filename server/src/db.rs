use sqlx::PgPool;
use tokio::sync::OnceCell;

static DB: OnceCell<PgPool> = OnceCell::<PgPool>::const_new();

async fn init() -> PgPool {
    PgPool::connect("postgres://postgres:postgres@localhost:5432/postgres")
        .await
        .unwrap()
}

pub async fn get() -> &'static PgPool {
    DB.get_or_init(init).await
}
