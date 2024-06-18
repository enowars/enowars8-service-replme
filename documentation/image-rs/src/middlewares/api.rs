use actix_web::{
    body::EitherBody,
    dev::{self, Service, ServiceRequest, ServiceResponse, Transform},
    Error, HttpResponse,
};
use futures_util::future::LocalBoxFuture;
use std::future::{ready, Ready};

pub struct Api {
    pub api_key: String,
}

impl<S, B> Transform<S, ServiceRequest> for Api
where
    S: Service<ServiceRequest, Response = ServiceResponse<B>, Error = Error>,
    S::Future: 'static,
    B: 'static,
{
    type Response = ServiceResponse<EitherBody<B>>;
    type Error = Error;
    type InitError = ();
    type Transform = ApiMiddleware<S>;
    type Future = Ready<Result<Self::Transform, Self::InitError>>;

    fn new_transform(&self, service: S) -> Self::Future {
        ready(Ok(ApiMiddleware {
            service,
            api_key: self.api_key.clone(),
        }))
    }
}

pub struct ApiMiddleware<S> {
    service: S,
    pub api_key: String,
}

impl<S, B> Service<ServiceRequest> for ApiMiddleware<S>
where
    S: Service<ServiceRequest, Response = ServiceResponse<B>, Error = Error>,
    S::Future: 'static,
    B: 'static,
{
    type Response = ServiceResponse<EitherBody<B>>;
    type Error = Error;
    type Future = LocalBoxFuture<'static, Result<Self::Response, Self::Error>>;

    dev::forward_ready!(service);

    fn call(&self, request: ServiceRequest) -> Self::Future {
        let path = request.match_info().get("api_key");
        if let Some(api_key) = path {
            if self.api_key == api_key {
                let res = self.service.call(request);

                return Box::pin(async move { res.await.map(ServiceResponse::map_into_left_body) });
            }
        }
        let (request, _pl) = request.into_parts();
        let response = HttpResponse::Unauthorized()
            .body("Invalid API key")
            .map_into_right_body();

        Box::pin(async { Ok(ServiceResponse::new(request, response)) })
    }
}
