import gleam/int
import lustre
import lustre/effect
import lustre/element.{text}
import lustre/element/html.{button, div, p}
import lustre/event.{on_click}

pub fn main() -> Nil {
  let init_model = 0
  let app = lustre.application(init:, update:, view:)
  let assert Ok(_) = lustre.start(app, "#app", init_model)

  Nil
}

type Model =
  Int

type Msg {
  Incr
  Decr
}

fn init(_init_model: Model) -> #(Model, effect.Effect(Msg)) {
  #(0, effect.none())
}

fn update(model: Model, msg: Msg) -> #(Model, effect.Effect(Msg)) {
  case msg {
    Incr -> #(model + 1, effect.none())
    Decr -> #(model - 1, effect.none())
  }
}

fn view(model: Model) -> element.Element(Msg) {
  let count = int.to_string(model)

  div([], [
    button([on_click(Incr)], [text(" + ")]),
    p([], [text(count)]),
    button([on_click(Decr)], [text(" - ")]),
  ])
}
