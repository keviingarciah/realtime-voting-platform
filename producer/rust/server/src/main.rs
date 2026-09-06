use rocket::response::status::BadRequest;
use rocket::serde::json::Json;
use rocket::config::SecretKey;
use rocket_cors::{AllowedOrigins, CorsOptions};
use std::convert::Infallible;
use serde_json;
use rocket::serde::Serialize; // Añade esta línea
use rdkafka::producer::{FutureProducer, FutureRecord};
use rdkafka::util::Timeout;

#[derive(rocket::serde::Deserialize, Clone, Serialize)] // Añade Serialize aquí
struct Data {
    name: String,
    album: String,
    year: String,
    rank: String,
}

async fn kafka_producer(payload: Json<Data>) -> Result<Json<String>, Infallible> {
    let payload_data = payload.into_inner(); // Obtén el Data de payload
    let payload_string = match serde_json::to_string(&payload_data) { // Serializa payload_data
        Ok(v) => v,
        Err(e) => {
            println!("Failed to serialize payload to JSON: {:?}", e);
            "Error".to_string() // return an empty string if there's an error
        }
    };

    // create producer
    let producer: FutureProducer = rdkafka::ClientConfig::new()
        .set("bootstrap.servers", "my-cluster-kafka-bootstrap:9092")
        .create()
        .expect("Producer creation error");


    // Produce a message to Kafka
    let topic = "so1-proyecto2";
    let delivery_status = producer.send(
        FutureRecord::<String, _>::to(&topic)
            .payload(&payload_string.as_bytes().to_vec()),
        Timeout::Never,
    ).await;
    
    match delivery_status {
        Ok((partition, offset)) => println!("Message delivered to partition {} at offset {}", partition, offset),
        Err((e, _)) => println!("Error producing message: {:?}", e),
    }
    
    // return the payload
    Ok(Json("Received".to_string()))
}

#[rocket::post("/data", data = "<data>")]
async fn receive_data(data: Json<Data>) -> Result<String, BadRequest<String>> {
    let received_data = data.into_inner();

    // Call kafka_producer
    let result = kafka_producer(Json(received_data.clone())).await;
    match result {
        Ok(_) => println!("Data sent to Kafka producer successfully"),
        Err(_) => println!("Failed to send data to Kafka producer"),
    }

    Ok("Data received and sent to Kafka producer".to_string())
}

#[rocket::main]
async fn main() {
    let secret_key = SecretKey::generate(); // Genera una nueva clave secreta

    // Configuración de opciones CORS
    let cors = CorsOptions::default()
        .allowed_origins(AllowedOrigins::all())
        .to_cors()
        .expect("failed to create CORS fairing");

    let config = rocket::Config {
        address: "0.0.0.0".parse().unwrap(),
        port: 8080,
        secret_key: secret_key.unwrap(), // Desempaqueta la clave secreta generada
        ..rocket::Config::default()
    };

    // Montar la aplicación Rocket con el middleware CORS
    rocket::custom(config)
        .attach(cors)
        .mount("/", rocket::routes![receive_data])
        .launch()
        .await
        .unwrap();
}