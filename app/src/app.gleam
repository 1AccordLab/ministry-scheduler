import gleam/dynamic/decode
import gleam/json
import gleam/result
import lustre
import lustre/effect
import plinth/browser/document
import plinth/browser/element as browser_element
import shared.{type Model, type Msg, Decr, Incr, view}

pub fn main() -> Nil {
  let init_model = case
    document.query_selector("#model")
    |> result.map(browser_element.inner_text)
    |> result.unwrap("")
    |> json.parse(decode.int)
  {
    Ok(count) -> count
    Error(_) -> 0
  }

  let app = lustre.application(init:, update:, view:)
  let assert Ok(_) = lustre.start(app, "#app", init_model)

  Nil
}

fn init(init_model: Model) -> #(Model, effect.Effect(Msg)) {
  #(init_model, effect.none())
}

fn update(model: Model, msg: Msg) -> #(Model, effect.Effect(Msg)) {
  case msg {
    Incr -> #(model + 1, effect.none())
    Decr -> #(model - 1, effect.none())
  }
}
