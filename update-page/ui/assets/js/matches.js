const TABLE_RENDERED_ID = 'render-completed';

/**
 * @param {string} url API url for load matches
 *
 * @return void
 */
async function loadMatches(url) {
    try {
        const response = await fetch(url)
        if (!response.ok) {
            throw new Error(`Response status ${response.status}`);
        }

        const tableBody = document.querySelector('#matches tbody');

        await response.json()
            .then(data => {
                tableBody.getAttribute('data') === TABLE_RENDERED_ID
                    ? refreshTableOdd(data, tableBody)
                    : buildTableOdd(data, tableBody);
            })
    } catch (error) {
        console.error(error.message)
    }
}

function startMatchFetching(url, interval = 5000) {
    async function loop() {
        await loadMatches(url);
        setTimeout(loop, interval);
    }

    loop(url, interval);
}


function refreshTableOdd(data, tableBody) {
    let row, oddElement;

    Object.keys(data).forEach((key) => {
        row = document.getElementById(key);

        oddElement = row.getElementsByClassName('firsOdd')[0];
        compareAndReplaceOdd(oddElement, data[key].firstWinOdd);

        oddElement = row.getElementsByClassName('secondOdd')[0];
        compareAndReplaceOdd(oddElement, data[key].secondWinOdd);

        oddElement = row.getElementsByClassName('drawOdd')[0];
        compareAndReplaceOdd(oddElement, data[key].drawOdd);
    });
}

function compareAndReplaceOdd(oddElement, fetchedValue) {
    const currentValue = parseFloat(oddElement.innerHTML);

    if (fetchedValue !== currentValue) {
        oddElement.innerHTML = fetchedValue.toFixed(2);

        if (fetchedValue > currentValue) {
            oddElement.classList.remove('odd-down');
            oddElement.classList.add('odd-up');
        } else {
            oddElement.classList.remove('odd-up');
            oddElement.classList.add('odd-down');
        }

        setTimeout(() => {
            oddElement.classList.remove('odd-up');
            oddElement.classList.remove('odd-down');
        }, 3000);

        //console.log({currentValue, fetchedValue, newResult: oddElement.innerHTML});
    }
}

function buildTableOdd(data, tableBody) {
    let row, nameCell, firsOddCell, drawOddCell, secondOddCell;

    row = document.createElement('tr');
    nameCell = document.createElement('th');
    firsOddCell = document.createElement('th');
    firsOddCell.setAttribute('class', 'text-center')
    firsOddCell.appendChild(document.createTextNode("1"))

    drawOddCell = document.createElement('th');
    drawOddCell.setAttribute('class', 'text-center')
    drawOddCell.appendChild(document.createTextNode("X"))

    secondOddCell = document.createElement('th');
    secondOddCell.setAttribute('class', 'text-center')
    secondOddCell.appendChild(document.createTextNode("2"))

    row.appendChild(nameCell)
    row.appendChild(firsOddCell)
    row.appendChild(drawOddCell)
    row.appendChild(secondOddCell)

    tableBody.appendChild(row)

    for (const [key, value] of Object.entries(data)) {
        row = document.createElement('tr');
        row.setAttribute('id', key);

        nameCell = document.createElement('td');
        nameCell.setAttribute('class', 'name')
        nameCell.appendChild(document.createTextNode(value.name))

        // first odd cell
        firsOddCell = document.createElement('td')
        firsOddCell.setAttribute('class', 'firsOdd  text-center odd')
        firsOddCell.appendChild(document.createTextNode(value.firstWinOdd.toFixed(2)))

        // draw odd cell
        drawOddCell = document.createElement('td')
        drawOddCell.setAttribute('class', 'drawOdd text-center odd')
        drawOddCell.appendChild(document.createTextNode(value.drawOdd.toFixed(2)))

        // second odd cell
        secondOddCell = document.createElement('td')
        secondOddCell.setAttribute('class', 'secondOdd text-center odd')
        secondOddCell.appendChild(document.createTextNode(value.secondWinOdd.toFixed(2)))

        row.appendChild(nameCell)
        row.appendChild(firsOddCell)
        row.appendChild(drawOddCell)
        row.appendChild(secondOddCell)

        tableBody.appendChild(row)
    }

    tableBody.setAttribute('data', TABLE_RENDERED_ID);
}
