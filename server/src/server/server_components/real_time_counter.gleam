import client/components/counter
import gleam/erlang/process.{type Subject}
import gleam/http/request.{type Request as HttpRequest}
import gleam/http/response.{type Response as HttpResponse}
import gleam/json
import gleam/option.{Some}
import lustre
import lustre/element.{type Element}
import lustre/server_component.{
  type ClientMessage, client_message_to_json, runtime_message_decoder,
}
import mist.{type Next, type WebsocketConnection, type WebsocketMessage}

pub fn serve(req: HttpRequest(_)) -> HttpResponse(_) {
  mist.websocket(req, handler, on_init, on_close)
}

pub fn element() -> Element(_) {
  server_component.element([server_component.route("/ws/counter")], [])
}

type CounterSocket {
  CounterSocket(
    self: Subject(CounterMsg),
    component: lustre.Runtime(counter.Msg),
  )
}

type CounterMsg =
  ClientMessage(counter.Msg)

fn on_init(_) {
  // start the server component runtime on the server-side
  let assert Ok(component) =
    lustre.start_server_component(counter.component(), Nil)

  // client <-> websocket server(`process.subject`) <-> server component runtime
  let self = process.new_subject()

  // `handler()` method will handle messages received from client/server component runtime
  let selector =
    process.new_selector()
    |> process.select(self)

  // establish the connection between:
  // `component`: client
  // `self`: server component runtime
  server_component.register_subject(self)
  |> lustre.send(to: component)

  #(CounterSocket(self, component), Some(selector))
}

fn handler(
  socket: CounterSocket,
  msg: WebsocketMessage(CounterMsg),
  connection: WebsocketConnection,
) -> Next(CounterSocket, CounterMsg) {
  case msg {
    // the websocket server receives a message from the client first,
    // decode it, and pass to the server component runtime
    mist.Text(json) -> {
      echo "message received from client: " <> json

      // send decoded message to server component runtime
      let assert Ok(runtime_msg) = json.parse(json, runtime_message_decoder())
      lustre.send(socket.component, runtime_msg)

      mist.continue(socket)
    }

    mist.Binary(_) -> {
      mist.continue(socket)
    }

    // the server component sends a message to the websocket server first,
    // and forward the message to the client
    mist.Custom(client_msg) -> {
      let text_frame =
        client_message_to_json(client_msg)
        |> json.to_string
      let assert Ok(_) = mist.send_text_frame(connection, text_frame)
      echo "message sent to client: " <> text_frame

      mist.continue(socket)
    }

    mist.Closed | mist.Shutdown -> {
      server_component.deregister_subject(socket.self)
      |> lustre.send(to: socket.component)

      mist.stop()
    }
  }
}

fn on_close(socket: CounterSocket) {
  server_component.deregister_subject(socket.self)
  |> lustre.send(to: socket.component)
}
