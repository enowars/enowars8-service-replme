use std::env;

use actix_session::{config::PersistentSession, storage::CookieSessionStore, SessionMiddleware};
use actix_web::{
    cookie::{self, Key},
    middleware,
    web::{self},
    App, HttpServer,
};
use image_rs::handlers;
use image_rs::middlewares;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init_from_env(env_logger::Env::new().default_filter_or("info"));

    log::info!("starting HTTP server at http://0.0.0.0:3000");

    let session_key = Key::generate();

    HttpServer::new(move || {
        App::new()
            .wrap(middleware::Logger::default())
            .wrap(
                SessionMiddleware::builder(CookieSessionStore::default(), session_key.clone())
                    .cookie_secure(false)
                    .session_lifecycle(
                        PersistentSession::default()
                            .session_ttl(cookie::time::Duration::minutes(15)),
                    )
                    .build(),
            )
            .service(
                web::scope("/api/{api_key}")
                    .wrap(middlewares::api::Api {
                        api_key: env::var("API_KEY").unwrap(),
                    })
                    .route(
                        "/auth/register",
                        web::post().to(handlers::auth::post_register),
                    )
                    .route("/auth/login", web::post().to(handlers::auth::post_login))
                    .service(
                        web::scope("/term")
                            // .wrap(middlewares::auth::Auth)
                            .route("/", web::get().to(handlers::term::get_term_ws)),
                    ),
            )
    })
    .bind(("0.0.0.0", 3000))?
    .run()
    .await
}
