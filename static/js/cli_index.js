document.getElementById('executeBtn').addEventListener('click', async () => {
  const command = document.getElementById('curlCommand').value;
  const workspaceId = document.getElementById('workspaceSelect').value; // Get the selected workspace ID
  
  try {
    const response = await fetch('/execute', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ command, workspaceId }), // Send workspace ID with the command
    });
    
    const jsonResponse = await response.json();
    document.getElementById('responseOutput').innerText = jsonResponse.response;
    loadHistory(workspaceId);  // Reload history for the selected workspace after execution

  } catch (error) {
    console.error('Error executing command:', error);
    document.getElementById('responseOutput').innerText = 'Error executing command.';
  }
});

// Load command history based on selected workspace ID
async function loadHistory(workspaceId) {
  try {
    const response = await fetch(`/history/${encodeURIComponent(workspaceId)}`); // Send workspace ID as a query parameter
    const history = await response.json();
    const historyList = document.getElementById('historyList');
    historyList.innerHTML = '';
    history.forEach(item => {
      historyList.innerHTML += `<li>${item.Command} - ${item.Timestamp}</li>`;
    });
  } catch (error) {
    console.error('Error loading history:', error);
    document.getElementById('historyList').innerHTML = 'Error loading history.';
  }
}

// Load workspace list
async function loadWorkspaces() {
  try {
    const response = await fetch('/workspaces');
    const workspaces = await response.json();
    
    const workspaceSelect = document.getElementById('workspaceSelect');
    workspaceSelect.innerHTML = '';  // Clear existing options
    
    if (workspaces.length > 0) {
      workspaces.forEach(workspace => {
        workspaceSelect.innerHTML += `<option value="${workspace.ID}">${workspace.Name}</option>`;
      });
    } else {
      workspaceSelect.innerHTML = '<option disabled>No workspaces available</option>';
    }

    // Load history for the first workspace if available
    if (workspaces.length > 0) {
      loadHistory(workspaces[0].ID);
    }

  } catch (error) {
    console.error('Error loading workspaces:', error);
    const workspaceSelect = document.getElementById('workspaceSelect');
    workspaceSelect.innerHTML = '<option disabled>Error loading workspaces</option>';
  }
}

// Load history when workspace selection changes
document.getElementById('workspaceSelect').addEventListener('change', (event) => {
  const workspaceId = event.target.value;
  loadHistory(workspaceId);
});

// Load history and workspaces on page load
window.onload = () => {
  loadWorkspaces();
};
