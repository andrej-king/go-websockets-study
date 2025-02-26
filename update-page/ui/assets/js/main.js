const LIVE_MATCHES_URL = '/api/matches/live';

window.addEventListener('load', function () {
    console.log("Page is fully loaded.");

    loadMatches(LIVE_MATCHES_URL);
    websocketListener([{name: EVENT_TYPE_LIVE_ODDS, cb: refreshTableOdd}]);

    // set footer text
    const dateFooter = document.querySelector('footer .js-date-container');
    const dateOptions = {weekday: 'long', year: 'numeric', month: 'long', day: 'numeric'};
    dateFooter.innerText = new Date().toLocaleDateString("en-US", dateOptions);
});
