const form = document.getElementById('uploadForm');
const selectedFile = document.getElementById("selectedFile");

form.addEventListener('submit', (event) => {
    event.preventDefault();
    const data = new FormData(form);
    fetch('/uploadfile', {
        method: 'POST',
        body: data
    }).then(response => response.json())
    .then(data => {
        console.log(data);
        // Do something with the data received from the server
    }).catch(error => {
        console.log(error);
    });
});
