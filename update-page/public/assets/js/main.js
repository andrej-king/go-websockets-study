window.addEventListener('load', function () {
    console.log("Page is fully loaded.");

    // set footer text
    const dateFooter = document.querySelector('footer .js-date-container');
    const dateOptions = {weekday: 'long', year: 'numeric', month: 'long', day: 'numeric'};
    dateFooter.innerText = new Date().toLocaleDateString("en-US", dateOptions);
});
