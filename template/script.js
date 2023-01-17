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
        // Extract the compiler version from the server response
        const compilerVersionMatch = data.match(/Switched global version to (.*)/);
        const compilerVersion = compilerVersionMatch ? compilerVersionMatch[1] : "N/A";
        compilerVersionLabel.innerText = `Solidity Compiler Used For Analysis : ${compilerVersion}`;
        compilerVersionLabel.classList.add("CompilerVersion");

        // Extract the issues from the data
        let issues = data.match(/[^\r\n]+/g);
        let summary = "";
        let count = 0;
        //iterate over each issue and show the summary
        issues.forEach(function(issue) {
            if (issue.startsWith("Reentrancy") || issue.startsWith("Suicide") || issue.startsWith("Unchecked low-level calls")) {
                summary += issue + "<br>";
                count++;
            }
        });
        if (count > 0) {
          auditResultsLabel.innerHTML = `<h3>Critical Issues: <span class="red">${count}</span> </h3> ${summary}`;        } else {
          auditResultsLabel.innerHTML = "No critical issues found";
        }
        // Display the audit results on the page
        //auditResultsLabel.innerText = data;
      }).catch(error => {
          console.log(error);
    });
});
