document.addEventListener('DOMContentLoaded', function () {
    let form = document.getElementById("myForm");
    form.addEventListener("submit", getData);

    //sort button
    let sortButton = document.getElementById("sortButton");
    sortButton.addEventListener("click", function (event) {
        event.preventDefault();
        sortTable();
    });

    function getData(event) {
        event.preventDefault();
        let selectedOption = document.getElementById("dios").value;

        const xhr = new XMLHttpRequest();
        if (selectedOption == "Go Web API") {
            xhr.open("GET", "http://192.168.0.100:6001/api/v1/detection");
        } else if (selectedOption == "---") {
            document.getElementById("demo").innerHTML = "<p style='font-size:16px'>No data source selected.</p>";
            return;
        }
        xhr.send();
        xhr.responseType = "json";
        xhr.onload = () => {
            if (xhr.readyState == 4 && xhr.status == 200) {
                generateHtml(xhr.response)
            } else {
                console.log(`Error: ${xhr.status}`);
                document.getElementById("demo").innerHTML = "ERROR!"
            }
        };
    }

    function generateHtml(data) {
        innerHtml = '<table id="dataTable" class="styled-table" style="padding: 12px; margin-left:auto; margin-right:auto; width:64%; text-align:center; font-size:16px"><tr><th>Object ID</th><th>Duration (milliseconds)</th><th>Timestamp</th></tr>';

        data.forEach(obj => {
            let formattedTime = new Date(obj.whattime).toISOString();
            innerHtml += '<tr><td>' + obj.id + '</td><td>' + obj.duration + '</td><td>' + formattedTime + '</td></tr>';
        });

        innerHtml += '</table>';
        document.getElementById("demo").innerHTML = innerHtml;
    }

    function sortTable() {
        let selectedOption = document.getElementById("dios").value;
        if(selectedOption != "---"){
            let table = document.getElementById("dataTable");
            let tbody = table.getElementsByTagName("tbody")[0];
            let rows = Array.from(tbody.getElementsByTagName("tr"));

            rows.sort((a, b) => {
                let idA = parseInt(a.cells[0].textContent);
                let idB = parseInt(b.cells[0].textContent);
                return idB - idA;
            });

            while (tbody.firstChild) {
                tbody.removeChild(tbody.firstChild);
            }

            rows.forEach(row => {
                tbody.appendChild(row);
            });
        }
    }
});
