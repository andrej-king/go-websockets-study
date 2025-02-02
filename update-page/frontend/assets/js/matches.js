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

        // const table = document.getElementById('matches')
        const tableBody = document.querySelector('#matches tbody')

        await response.json()
            .then(data => {
                console.log(data)

                let row, nameCell, firsOddCell, drawOddCell, secondOddCell, icon;

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
                    firsOddCell.setAttribute('class', 'firsOdd  text-center')
                    firsOddCell.appendChild(document.createTextNode(value.firstWinOdd.toFixed(2)))

                    icon = document.createElement('i');
                    icon.setAttribute('class', 'ms-2')
                    firsOddCell.appendChild(icon)

                    // draw odd cell
                    drawOddCell = document.createElement('td')
                    drawOddCell.setAttribute('class', 'drawOdd text-center')
                    drawOddCell.appendChild(document.createTextNode(value.drawOdd.toFixed(2)))

                    icon = document.createElement('i');
                    icon.setAttribute('class', 'ms-2')
                    secondOddCell.appendChild(icon)

                    // second odd cell
                    secondOddCell = document.createElement('td')
                    secondOddCell.setAttribute('class', 'secondOdd text-center')
                    secondOddCell.appendChild(document.createTextNode(value.secondWinOdd.toFixed(2)))

                    icon = document.createElement('i');
                    icon.setAttribute('class', 'ms-2')
                    secondOddCell.appendChild(icon)

                    row.appendChild(nameCell)
                    row.appendChild(firsOddCell)
                    row.appendChild(drawOddCell)
                    row.appendChild(secondOddCell)

                    tableBody.appendChild(row)
                }
            })
    } catch (error) {
        console.error(error.message)
    }
}
