import gleam/int
import lustre/element.{type Element, text}
import lustre/element/html.{button, div, p}
import lustre/event.{on_click}

pub type Model =
  Int

pub type Msg {
  Incr
  Decr
}

pub fn view(model: Model) -> Element(Msg) {
  let count = int.to_string(model)

  div([], [
    button([on_click(Incr)], [text(" + ")]),
    p([], [text(count)]),
    button([on_click(Decr)], [text(" - ")]),
  ])
}
