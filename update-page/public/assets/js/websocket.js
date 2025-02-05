class WS {
    get newClientPromise() {
        return new Promise((resolve, reject) => {
            let wsClient = new WebSocket(`ws://${document.location.host}/ws`);
            console.log(wsClient);

            wsClient.onopen = () => {
                console.log("websocket connected")
                resolve(wsClient);
            }

            wsClient.onerror = (error) => reject(error);
        });
    }

    get clientPromise() {
        if (!this.promise) {
            this.promise = this.newClientPromise;
        }

        return this.promise;
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
