import client/components/counter
import gleam/erlang/process
import gleam/http
import lustre/attribute.{attribute}
import lustre/element
import lustre/element/html
import mist
import wisp.{type Request, type Response}
import wisp/wisp_mist

pub fn main() {
  wisp.configure_logger()

  let assert Ok(priv_dir) = wisp.priv_directory("server")
  let static_dir = priv_dir <> "/static"

  let secret_key_base = wisp.random_string(64)

  let assert Ok(_) =
    handle_request(_, static_dir)
    |> wisp_mist.handler(secret_key_base)
    |> mist.new
    |> mist.port(3000)
    |> mist.start

  process.sleep_forever()
}

fn handle_request(req: Request, static_dir: String) -> Response {
  use req <- middlewares(req, static_dir)

  case req.method, wisp.path_segments(req) {
    http.Get, [] -> serve_index()
    _, _ -> wisp.not_found()
  }
}

fn middlewares(
  req: Request,
  static_dir: String,
  next: fn(Request) -> Response,
) -> Response {
  use <- wisp.log_request(req)
  use <- wisp.rescue_crashes
  use req <- wisp.handle_head(req)
  use <- wisp.serve_static(req, under: "/static", from: static_dir)

  next(req)
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
      ]),
      html.body([], [html.div([], [counter.element()])]),
    ])

  html
  |> element.to_document_string_tree
  |> wisp.html_response(200)
}
