const form = document.getElementById('uploadForm');
const selectedFile = document.getElementById("selectedFile");
const compilerVersionLabel = document.getElementById("CompilerVersion");
const auditResultsLabel = document.getElementById("AuditResults");

form.addEventListener('submit', (event) => {
    event.preventDefault();
    const data = new FormData(form);
    fetch('/uploadfile', {
        method: 'POST',
        body: data
    }).then(response => response.json())
    .then(data => {
        
        console.log(data);

      }).catch(error => {
          console.log(error);
    });
});