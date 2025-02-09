function loadPageMatches() {
    // api
    loadMatches(`/api/matches/live`);
}

function wsInitPageRequest() {
    const receiveOddsReadyEvent = new Event("ready_to_receive_odds", {type: 'live'});

    window.wsSingleton.clientPromise
        .then(wsClient => wsClient.send(JSON.stringify(receiveOddsReadyEvent)))
        .catch(error => console.log(error));
}
