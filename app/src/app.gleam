import gleam/dynamic/decode
import gleam/json
import gleam/result
import lustre
import lustre/effect.{type Effect}
import lustre/element.{type Element} as lustre_element
import plinth/browser/document
import plinth/browser/element
import shared/counter

pub fn main() -> Nil {
  let init_model =
    Model(counter: {
      document.query_selector("#model")
      |> result.map(element.inner_text)
      |> result.unwrap("")
      |> json.parse(decode.int)
      |> result.unwrap(0)
    })

  let app = lustre.application(init:, update:, view:)
  let assert Ok(_) = lustre.start(app, "#app", init_model)

  Nil
}

type Model {
  Model(counter: counter.Model)
}

type Msg {
  Counter(counter.Msg)
}

fn init(init_model: Model) -> #(Model, Effect(Msg)) {
  #(init_model, effect.none())
}

fn update(model: Model, msg: Msg) -> #(Model, Effect(Msg)) {
  case msg {
    Counter(msg) -> {
      let #(model, effect) = counter.update(model.counter, msg)
      #(Model(counter: model), effect |> effect.map(Counter))
    }
  }
}

fn view(model: Model) -> Element(Msg) {
  counter.view(model.counter) |> lustre_element.map(Counter)
}
