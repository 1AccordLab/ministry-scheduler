import gleam/int
import lustre
import lustre/effect.{type Effect}
import lustre/element.{type Element}
import lustre/element/html
import lustre/event

const component_name = "lustre-counter"

pub fn register() -> Result(Nil, lustre.Error) {
  let component = lustre.component(init, update, view, [])
  lustre.register(component, component_name)
}

pub fn element() -> Element(Msg) {
  element.element(component_name, [], [])
}

pub type Model {
  Model(count: Int)
}

pub type Msg {
  Increment
  Decrement
}

fn init(_) -> #(Model, Effect(Msg)) {
  #(Model(count: 0), effect.none())
}

fn update(model: Model, msg: Msg) -> #(Model, Effect(Msg)) {
  case msg {
    Increment -> #(Model(count: model.count + 1), effect.none())
    Decrement -> #(Model(count: model.count - 1), effect.none())
  }
}

fn view(model: Model) -> Element(Msg) {
  let count = int.to_string(model.count)

  html.div([], [
    html.button([event.on_click(Decrement)], [element.text("-")]),
    element.text(count),
    html.button([event.on_click(Increment)], [element.text("+")]),
  ])
}
