import client/components/counter
import ewe.{type Request, type Response}
import gleam/erlang/process
import gleam/http
import gleam/http/request
import gleam/http/response
import gleam/int
import gleam/list
import gleam/option.{None}
import gleam/result
import gleam/string
import logging
import lustre/attribute.{attribute}
import lustre/element
import lustre/element/html
import lustre/server_component
import marceau
import server/server_components/real_time_counter

pub fn main() {
  logging.configure()
  logging.set_level(logging.Info)

  let assert Ok(_) =
    ewe.new(handler)
    |> ewe.bind_all()
    |> ewe.listening(port: 3000)
    |> ewe.start

  process.sleep_forever()
}

fn handler(req: Request) -> Response {
  use <- log_request(req)

  case req.method, request.path_segments(req) {
    http.Get, [] -> serve_index()
    http.Get, ["static", path] -> serve_static(req, path)
    http.Get, ["ws", "counter"] -> real_time_counter.serve(req)
    _, _ -> send_not_found()
  }
}

fn serve_static(req: Request, path: String) -> Response {
  let file_ext =
    req.path
    |> string.split(".")
    |> list.last
    |> result.unwrap("")

  let mime_type = marceau.extension_to_mime_type(file_ext)

  let content_type = case mime_type {
    "application/json" | "text/" <> _ -> mime_type <> "; charset=utf-8"
    _ -> mime_type
  }

  case ewe.file("priv/static/" <> path, offset: None, limit: None) {
    Ok(file) ->
      response.new(200)
      |> response.set_header("Content-Type", content_type)
      |> response.set_body(file)

    Error(_) -> send_not_found()
  }
}

fn serve_index() -> Response {
  let html =
    html.html([], [
      html.head([], [
        html.meta([attribute("charset", "utf-8")]),
        html.meta([
          attribute("content", "width=device-width, initial-scale=1"),
          attribute("name", "viewport"),
        ]),
        html.script(
          [attribute.src("/static/client.min.mjs"), attribute.type_("module")],
          "",
        ),
        server_component.script(),
      ]),
      html.body([], [
        html.div([], [html.p([], [html.text("Counter")]), counter.element()]),
        html.div([], [
          html.p([], [html.text("Real Time Counter")]),
          real_time_counter.element(),
        ]),
      ]),
    ])
    |> element.to_document_string

  send_html(html)
}

fn send_html(body: String) -> Response {
  response.new(200)
  |> response.set_header("Content-Type", "text/html; charset=utf-8")
  |> response.set_body(ewe.TextData(body))
}

fn send_not_found() -> Response {
  response.new(404)
  |> response.set_header("Content-Type", "text/plain; charset=utf-8")
  |> response.set_body(ewe.TextData("404 Not Found"))
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
