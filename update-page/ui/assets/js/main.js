const LIVE_MATCHES_URL = '/api/matches/live';

window.addEventListener('load', function () {
    console.log("Page is fully loaded.");

    startMatchFetching(LIVE_MATCHES_URL);

    // set footer text
    const dateFooter = document.querySelector('footer .js-date-container');
    const dateOptions = {weekday: 'long', year: 'numeric', month: 'long', day: 'numeric'};
    dateFooter.innerText = new Date().toLocaleDateString("en-US", dateOptions);
});
