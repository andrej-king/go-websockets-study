// const API_URL = "http://localhost:8000/api"
// const API_URL = `http://${document.location.host}/api`
const API_URL = `/api`
const WS_URL = "/ws"

window.addEventListener('load', function () {
    console.log("Page is fully loaded.");

    loadMatches(`${API_URL}/matches`)

    // set footer text
    const dateFooter = document.querySelector('footer .js-date-container');
    const dateOptions = {weekday: 'long', year: 'numeric', month: 'long', day: 'numeric'};
    dateFooter.innerText = new Date().toLocaleDateString("en-US", dateOptions);
})
