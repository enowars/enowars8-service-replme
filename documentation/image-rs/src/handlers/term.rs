use serde::{Deserialize, Serialize};
use std::{env, os::unix::process::CommandExt, process::Command};

use actix_session::Session;
use actix_web::{
    web::{self},
    Error, HttpRequest, HttpResponse, HttpResponseBuilder,
};
use actix_web_actors::ws;

use crate::services::{term::Websocket, user::get_user_data};

#[derive(Debug, Serialize, Deserialize)]
pub struct ErrorResponse {
    error: String,
}

pub async fn get_term_ws(
    request: HttpRequest,
    stream: web::Payload,
) -> Result<HttpResponse, Error> {
    let command = Command::new("/etc/profiles/per-user/gandalf/bin/zsh");
    ws::start(Websocket::new(command), &request, stream)
}

// pub async fn get_term_ws(
//     request: HttpRequest,
//     session: Session,
//     stream: web::Payload,
// ) -> Result<HttpResponse, Error> {
//     match session.get::<String>("username") {
//         Ok(Some(username)) if !username.is_empty() => match get_user_data(&username) {
//             Ok(user_data) => match env::var("PATH") {
//                 Ok(path) => {
//                     let mut command = Command::new(user_data.shell);
//                     command.env_clear();
//                     command.env("TERM", "xterm");
//                     command.env("PATH", path);
//                     command.gid(user_data.gid);
//                     command.uid(user_data.uid);
//                     ws::start(Websocket::new(command), &request, stream)
//                 }
//                 Err(_) => Ok(HttpResponse::InternalServerError().json(ErrorResponse {
//                     error: "InternalServerError".to_owned(),
//                 })),
//             },
//             Err(error) => Ok(HttpResponseBuilder::new(error.status).json(ErrorResponse {
//                 error: error.message,
//             })),
//         },
//         _ => Ok(HttpResponse::Unauthorized().into()),
//     }
// }
