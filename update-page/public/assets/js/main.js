// const API_URL = "http://localhost:8000/api"
// const API_URL = `http://${document.location.host}/api`
const API_URL = `/api`
const WS_URL = "/ws"

window.addEventListener('load', function () {
    console.log("Page is fully loaded.");

    loadMatches(`${API_URL}/matches`)
    connectWebsocket(`ws://${document.location.host}${WS_URL}`)

    // set footer text
    const dateFooter = document.querySelector('footer .js-date-container');
    const dateOptions = {weekday: 'long', year: 'numeric', month: 'long', day: 'numeric'};
    dateFooter.innerText = new Date().toLocaleDateString("en-US", dateOptions);
})

function connectWebsocket(url) {
    // Check if the browser supports WebSocket
    if (window['WebSocket']) {
        console.log("supports websockets");

        socket = new WebSocket(url);

        socket.onopen = function (event) {
            // TODO
            console.info("Socket: connected", event);
        };

        socket.onclose = function (event) {
            // TODO
            console.info("Socket: close", event);
        };

        socket.onmessage = function (event) {
            // TODO
            console.info("Socket: message", event);

            // parse websocket message as JSON
            const eventData = JSON.parse(event.data);

            // TODO Assign JSON data to new Event Object
            // Assign JSON data to new Event Object
            const eventObject = Object.assign(new Event, eventData)

            // Let router manage message
            routeEvent(eventObject)
        };
    } else {
        alert("Not supporting websockets");
    }
}

/**
 * routeEvent is a proxy function that routes
 * events into their correct Handler
 * based on the type field
 */
function routeEvent(event) {
    if (event.type === undefined) {
        alert("no 'type' field in event");
    }

    switch (event.type) {
        case "send_odds":
            const messageEvent = Object.assign(new MessageEvent, event.payload)
            break;
        default:
            alert("unsupported message type");
            break;
    }
}

/**
 * Event is used to wrap all messages Send and Received on the Websocket
 */
class Event {
    // Each Event needs a Type
    // The payload is not required
    constructor(type, payload) {
        this.type = type;
        this.payload = payload;
    }
}

/**
 * NewMessageEvent is messages coming from clients
 * */
class NewMessageEvent {
    constructor(message, from, sent) {
        this.message = message;
        this.sent = sent;
    }
}
