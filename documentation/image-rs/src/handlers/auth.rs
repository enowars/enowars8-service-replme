use actix_session::Session;
use actix_web::{web, Error, HttpResponse, HttpResponseBuilder};

use serde::{Deserialize, Serialize};

use crate::services::user;

#[derive(Debug, Serialize, Deserialize)]
pub struct RegisterRequest {
    username: String,
    password: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct LoginRequest {
    username: String,
    password: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct SuccessResponse {
    success: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct ErrorResponse {
    error: String,
}

pub async fn post_register(
    session: Session,
    register_request: web::Json<RegisterRequest>,
) -> Result<HttpResponse, Error> {
    match user::create(&register_request.username, &register_request.password) {
        Ok(result) => {
            let _ = session.insert("username", register_request.username.clone());
            Ok(
                HttpResponseBuilder::new(result.status).json(SuccessResponse {
                    success: result.message,
                }),
            )
        }
        Err(error) => Ok(HttpResponseBuilder::new(error.status).json(ErrorResponse {
            error: error.message,
        })),
    }
}

pub async fn post_login(
    session: Session,
    login_request: web::Json<LoginRequest>,
) -> Result<HttpResponse, Error> {
    match user::login(&login_request.username, &login_request.password) {
        Ok(result) => {
            let _ = session.insert("username", login_request.username.clone());
            Ok(
                HttpResponseBuilder::new(result.status).json(SuccessResponse {
                    success: result.message,
                }),
            )
        }
        Err(error) => Ok(HttpResponseBuilder::new(error.status).json(ErrorResponse {
            error: error.message,
        })),
    }
}
