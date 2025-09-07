import client/components/counter
import gleam/erlang/process
import gleam/http
import gleam/http/request.{type Request as HttpRequest}
import gleam/http/response.{type Response as HttpResponse}
import lustre/attribute.{attribute}
import lustre/element
import lustre/element/html
import lustre/server_component
import mist
import server/server_components/real_time_counter
import wisp
import wisp/wisp_mist

type Context {
  Context(static_dir: String, secret_key_base: String)
}

pub fn main() {
  wisp.configure_logger()

  let ctx =
    Context(
      static_dir: get_static_dir(),
      secret_key_base: wisp.random_string(64),
    )

  let assert Ok(_) =
    router(_, ctx)
    |> mist.new
    |> mist.port(3000)
    |> mist.start

  process.sleep_forever()
}

fn get_static_dir() -> String {
  let assert Ok(priv_dir) = wisp.priv_directory("server")
  priv_dir <> "/static"
}

fn router(req: HttpRequest(_), ctx: Context) -> HttpResponse(_) {
  case req.method, request.path_segments(req) {
    http.Get, ["ws", path] -> ws_router(req, ctx, path)
    _, _ -> http_router(req, ctx)
  }
}

fn ws_router(req: HttpRequest(_), ctx: Context, path: String) -> HttpResponse(_) {
  // ** Adds websocket APIs here **
  case path {
    "counter" -> real_time_counter.serve(req)
    _ -> http_router(req, ctx)
  }
}

fn http_router(req: HttpRequest(_), ctx: Context) -> HttpResponse(_) {
  let handle_request = fn(req: wisp.Request) -> wisp.Response {
    // ** Middlewares **
    use <- wisp.log_request(req)
    use <- wisp.rescue_crashes
    use req <- wisp.handle_head(req)
    use <- wisp.serve_static(req, under: "/static", from: ctx.static_dir)

    // ** Adds http APIs here **
    case req.method, request.path_segments(req) {
      http.Get, [] -> serve_index()
      _, _ -> wisp.not_found()
    }
  }

  // convert wisp handler to mist handler
  let mist_handler =
    handle_request
    |> wisp_mist.handler(ctx.secret_key_base)

  mist_handler(req)
}

fn serve_index() -> wisp.Response {
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

  html
  |> element.to_document_string()
  |> wisp.html_response(200)
}
