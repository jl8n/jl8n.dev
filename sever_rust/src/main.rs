use axum::{routing::get, Router, http::StatusCode};
use std::net::SocketAddr;
use dotenv::dotenv;
use serde::{Serialize, Deserialize};
use serde_json::json;
use tower_http::cors::CorsLayer;


#[derive(Serialize, Deserialize)]
struct TrackInfo {
    artist_name: String,
    release_name: String,
    track_name: String,
}


#[tokio::main]
async fn main() {
    dotenv().ok();


    // Route all requests on "/" endpoint to anonymous handler.
    //
    // A handler is an async function which returns something that implements
    // `axum::response::IntoResponse`.

    // A closure or a function can be used as handler.

    let app = Router::new()
        .route("/", get(handler))
        .route("/external", get(external_handler))
        .layer(CorsLayer::permissive());


    // Address that server will bind to.
    let addr = SocketAddr::from(([127, 0, 0, 1], 3000));

    // Use `hyper::server::Server` which is re-exported through `axum::Server` to serve the app.
    axum::Server::bind(&addr)
        // Hyper server takes a make service.
        .serve(app.into_make_service())
        .await
        .unwrap();
}

async fn handler() -> &'static str {
    "Hello, world!"
}


async fn external_handler() -> Result<String, StatusCode> {
    let listenbrainz_api_token = std::env::var("LISTENBRAINZ_API_TOKEN")
        .expect("LISTENBRAINZ_API_TOKEN must be set.");
    let listenbrainz_user = std::env::var("LISTENBRAINZ_USER")
        .expect("LISTENBRAINZ_USER must be set.");
    let url = format!("https://api.listenbrainz.org/1/user/{}/playing-now", listenbrainz_user);

    let resp = reqwest::Client::new()
        .get(&url)
        .header(reqwest::header::AUTHORIZATION, format!("Token {}", listenbrainz_api_token))
        .send()
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let json_value = resp.json::<serde_json::Value>()
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let track_metadata = &json_value["payload"]["listens"][0]["track_metadata"];
    let artist_name = track_metadata["artist_name"].as_str().ok_or(StatusCode::INTERNAL_SERVER_ERROR)?;
    let release_name = track_metadata["release_name"].as_str().ok_or(StatusCode::INTERNAL_SERVER_ERROR)?;
    let track_name = track_metadata["track_name"].as_str().ok_or(StatusCode::INTERNAL_SERVER_ERROR)?;

    let response = json!({
        "artist": artist_name,
        "album": release_name,
        "track": track_name,
    });

    println!("{}", "GET /music");

    Ok(response.to_string())
}