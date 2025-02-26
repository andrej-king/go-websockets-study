let socket = null;
const EVENT_TYPE_SUBSCRIBE = 'subscribe';
const EVENT_TYPE_LIVE_ODDS = 'live_odds';

function websocketListener(subscriptions = [], cb) {
    // Do nothing if socket already connected
    if (socket) {
        return;
    }

    if (window['WebSocket']) {
        // Create a WebSocket connection to the server
        socket = new WebSocket('ws://localhost:8080/ws');
    } else {
        alert("Not supporting websockets");
    }

    // Handle connection open event
    socket.onopen = function (event) {
        // subscribe
        if (subscriptions.length > 0) {
            subscriptions.map(object => {
                sendEvent(EVENT_TYPE_SUBSCRIBE, {name: object.name})
            });
        }
    };

    // Handle incoming messages from the server
    socket.onmessage = function (event) {
        const receivedEvent = Object.assign(new Event, JSON.parse(event.data));
        const resolver = subscriptions.filter(item => item.name === receivedEvent.type);

        // Handle different types of events
        if (receivedEvent.type === EVENT_TYPE_LIVE_ODDS && resolver.length > 0) {
            resolver[0].cb(receivedEvent.payload);
        } else {
            alert(`Unexpected type: ${receivedEvent.type}`)
            throw new Error(`Unexpected type: ${receivedEvent.type}`);
        }
    };

    // Handle connection close event
    socket.onclose = function (event) {
        disconnect();
    };
}

function disconnect() {
    // Do nothing if no have socket
    if (!socket) {
        return;
    }

    socket.close();

    console.log('connection closed');
}

/**
 * send event
 *
 * eventName - the event name to send on
 * payload - the data payload
 * */
function sendEvent(eventName, payload) {
    // Do nothing if not connected
    if (!socket || socket.readyState !== WebSocket.OPEN) {
        console.log({socket, readyState: socket.readyState});

        return;
    }

    const event = new Event(eventName, payload);

    // Format as JSON and send
    socket.send(JSON.stringify(event));
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
