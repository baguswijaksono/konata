document.getElementById('executeBtn').addEventListener('click', async () => {
  const command = document.getElementById('curlCommand').value;
  const workspace = document.getElementById('workspaceSelect').value; // Get the selected workspace
  const response = await fetch('/execute', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ command, workspace }), // Send workspace with the command
  }).then(res => res.json());

  document.getElementById('responseOutput').innerText = response.response;
  loadHistory();  // Reload history after execution
});


    document.getElementById('createWorkspaceBtn').addEventListener('click', async () => {
      const name = document.getElementById('workspaceName').value;
      const config = document.getElementById('workspaceConfig').value;
      const response = await fetch('/workspace', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name, config }),
      }).then(res => res.json());

      alert(response.message);
      loadWorkspaces();  // Reload workspaces after creation
    });

    async function loadHistory() {
      const history = await fetch('/history').then(res => res.json());
      const historyList = document.getElementById('historyList');
      historyList.innerHTML = '';
      history.forEach(item => {
        historyList.innerHTML += `<li>${item.Command} - ${item.Timestamp}</li>`;
      });
    }

    async function loadWorkspaces() {
      const workspaces = await fetch('/workspaces').then(res => res.json());
      const workspaceSelect = document.getElementById('workspaceSelect');
      workspaceSelect.innerHTML = '';
      workspaces.forEach(workspace => {
        workspaceSelect.innerHTML += `<option value="${workspace.Name}">${workspace.Name}</option>`;
      });
    }

    loadHistory();  // Initial load of history on page load
    loadWorkspaces();  // Initial load of workspaces on page load