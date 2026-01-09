import ewe.{type Request, type Response}
import gleam/erlang/process
import gleam/http.{Get}
import gleam/http/request
import logging
import server/counter
import server/internal/utils.{
  log_request, send_html, send_not_found, serve_static,
}

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
    Get, ["counter"] -> counter.serve(req)
    _, _ -> send_not_found()
  }
}

fn serve_index(_req: Request) -> Response {
  send_html("<h1>Hello, world!</h1>")
}
