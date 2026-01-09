import ewe.{type Request, type Response, type ResponseBody, TextData}
import gleam/erlang/process
import gleam/http.{Get}
import gleam/http/request
import gleam/http/response
import gleam/int
import gleam/json.{type Json}
import gleam/list
import gleam/option.{None}
import gleam/result
import gleam/string
import logging
import lustre/attribute
import lustre/element.{type Element}
import lustre/element/html
import marceau
import shared.{view}

pub fn main() -> Nil {
  logging.configure()
  logging.set_level(logging.Info)

  let assert Ok(_) =
    ewe.new(handler)
    |> ewe.bind_all
    |> ewe.listening(port: 8000)
    |> ewe.start

  process.sleep_forever()
}

fn handler(req: Request) -> Response {
  use <- log_request(req)

  case req.method, request.path_segments(req) {
    Get, [] -> serve_index(req)
    Get, ["static", path] -> serve_static(req, path)
    Get, ["counter"] -> serve_counter(req)
    _, _ -> send_not_found()
  }
}

const init_counter_model = 5

fn load_lustre_app(init_model: Json, views: List(Element(a))) -> Element(a) {
  html.html([], [
    html.head([], [
      html.script(
        [attribute.type_("module"), attribute.src("/static/app.js")],
        "",
      ),
      html.script(
        [attribute.type_("application/json"), attribute.id("model")],
        { init_model |> json.to_string },
      ),
    ]),
    html.body([], [
      html.div([attribute.id("app")], views),
    ]),
  ])
}

fn serve_index(_req: Request) -> Response {
  send_html("<h1>Hello, world!</h1>")
}

fn serve_counter(_req: Request) -> Response {
  let init_model = json.int(init_counter_model)
  let views = [view(init_counter_model)]
  let html = load_lustre_app(init_model, views) |> element.to_string
  send_html(html)
}

fn send_not_found() -> Response {
  response.new(404)
  |> response.set_body(TextData("Not found"))
}

fn send_html(html: String) -> Response {
  response.new(200)
  |> response.set_header("Content-Type", "text/html; charset=utf-8")
  |> response.set_body(TextData(html))
}

fn send_file(file: ResponseBody, content_type: String) -> Response {
  response.new(200)
  |> response.set_header("Content-Type", content_type)
  |> response.set_body(file)
}

fn serve_static(req: Request, path: String) -> Response {
  let mime_type =
    req.path
    |> string.split(".")
    |> list.last
    |> result.unwrap("")
    |> marceau.extension_to_mime_type

  let content_type = case mime_type {
    "application/json" | "text/" <> _ -> mime_type <> "; charset=utf-8"
    _ -> mime_type
  }

  case ewe.file("priv/static/" <> path, offset: None, limit: None) {
    Ok(file) -> send_file(file, content_type)
    Error(err) -> {
      echo err
      send_not_found()
    }
  }
}

fn log_request(req: Request, handler: fn() -> Response) -> Response {
  let response = handler()

  [
    int.to_string(response.status),
    " ",
    string.uppercase(http.method_to_string(req.method)),
    " ",
    req.path,
  ]
  |> string.concat
  |> logging.log(logging.Info, _)

  response
}
