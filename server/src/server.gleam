import gleam/erlang/process
import gleam/http
import mist
import wisp.{type Request, type Response}
import wisp/wisp_mist

pub fn main() -> Nil {
  wisp.configure_logger()

  let secret_key_base = wisp.random_string(64)

  let assert Ok(_) =
    handle_request
    |> wisp_mist.handler(secret_key_base)
    |> mist.new
    |> mist.port(3000)
    |> mist.start

  process.sleep_forever()
}

fn handle_request(req: Request) -> Response {
  use req <- middlewares(req)

  case req.method, wisp.path_segments(req) {
    http.Get, ["hello"] -> {
      wisp.response(200) |> wisp.string_body("hello world!")
    }

    _, _ -> wisp.not_found()
  }
}

fn middlewares(req: Request, next: fn(Request) -> Response) -> Response {
  use <- wisp.log_request(req)
  use <- wisp.rescue_crashes
  use req <- wisp.handle_head(req)

  next(req)
}
