import ewe.{type Request, type Response}
import gleam/json
import lustre/element
import server/internal/utils.{lustre_app, send_html}
import shared/counter

const init_model = 5

pub fn serve(_req: Request) -> Response {
  let model = json.int(init_model)
  let views = [counter.view(init_model)]
  let html = lustre_app(model, views) |> element.to_string
  send_html(html)
}
