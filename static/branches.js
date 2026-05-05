'use strict';

function attachBranchSelector(card, project) {
  const branchBtn = document.createElement('button');
  branchBtn.className = 'btn-sm branch-btn';
  branchBtn.textContent = '...';
  branchBtn.title = 'Manage branches';
  branchBtn.addEventListener('click', async (e) => {
    e.stopPropagation();
    await toggleBranchDropdown(card, project, branchBtn);
  });
  card.querySelector('.project-header').appendChild(branchBtn);
}

async function toggleBranchDropdown(card, project, btn) {
  const existing = card.querySelector('.branch-dropdown');
  if (existing) {
    existing.remove();
    return;
  }

  const dropdown = document.createElement('div');
  dropdown.className = 'branch-dropdown';

  const branches = await fetchBranches(project.path);
  if (branches.error) {
    dropdown.innerHTML = '<div class="branch-error">Failed to load branches</div>';
  } else {
    const list = document.createElement('div');
    list.className = 'branch-list';
    for (const branch of branches) {
      const item = document.createElement('button');
      item.className = 'branch-item';
      item.textContent = branch;
      item.addEventListener('click', async () => {
        await checkoutBranch(project.path, branch);
        dropdown.remove();
      });
      list.appendChild(item);
    }
    dropdown.appendChild(list);

    const createSection = document.createElement('div');
    createSection.className = 'branch-create';
    createSection.innerHTML = `
      <input type="text" class="branch-name-input" placeholder="new-branch-name">
      <button class="btn-sm branch-create-btn">Create</button>
    `;
    createSection.querySelector('.branch-create-btn').addEventListener('click', async () => {
      const input = createSection.querySelector('.branch-name-input');
      const name = input.value.trim();
      if (name) {
        await createBranch(project.path, name);
        dropdown.remove();
      }
    });
    dropdown.appendChild(createSection);
  }

  card.appendChild(dropdown);
}

async function fetchBranches(repoPath) {
  try {
    const res = await fetch(`/api/branches?repo=${encodeURIComponent(repoPath)}`);
    const data = await res.json();
    return data.branches || [];
  } catch (err) {
    return { error: err.message };
  }
}

async function checkoutBranch(repoPath, branch) {
  try {
    const res = await fetch(`/api/branches?repo=${encodeURIComponent(repoPath)}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ action: 'checkout', branch })
    });
    const data = await res.json();
    if (!res.ok) {
      console.error('checkout failed:', data.error);
      return false;
    }
    return true;
  } catch (err) {
    console.error('checkout error:', err);
    return false;
  }
}

async function createBranch(repoPath, branch) {
  try {
    const res = await fetch(`/api/branches?repo=${encodeURIComponent(repoPath)}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ action: 'create', branch })
    });
    const data = await res.json();
    if (!res.ok) {
      console.error('create failed:', data.error);
      return false;
    }
    return true;
  } catch (err) {
    console.error('create error:', err);
    return false;
  }
}

function attachStashToggle(card, project) {
  const stashBtn = document.createElement('button');
  stashBtn.className = 'btn-sm stash-btn';
  stashBtn.textContent = 'Stash';
  stashBtn.title = 'View stash';
  stashBtn.addEventListener('click', async (e) => {
    e.stopPropagation();
    await toggleStashPanel(card, project, stashBtn);
  });
  card.querySelector('.project-header').appendChild(stashBtn);
}

async function toggleStashPanel(card, project, btn) {
  const existing = card.querySelector('.stash-panel');
  if (existing) {
    existing.remove();
    return;
  }

  const panel = document.createElement('div');
  panel.className = 'stash-panel';

  const stashes = await fetchStashes(project.path);
  if (stashes.error) {
    panel.innerHTML = '<div class="stash-error">Failed to load stashes</div>';
  } else if (stashes.length === 0) {
    panel.innerHTML = '<div class="stash-empty">No stashes</div>';
  } else {
    const list = document.createElement('div');
    list.className = 'stash-list';
    for (const stash of stashes) {
      const item = document.createElement('div');
      item.className = 'stash-item';
      item.innerHTML = `
        <span class="stash-msg">${esc(stash.message)}</span>
        <span class="stash-meta">${esc(stash.author)} · ${esc(stash.timestamp)}</span>
        <button class="btn-sm stash-apply-btn">Apply</button>
        <button class="btn-sm stash-drop-btn">Drop</button>
      `;
      item.querySelector('.stash-apply-btn').addEventListener('click', async () => {
        await applyStash(project.path, stash.index);
        panel.remove();
      });
      item.querySelector('.stash-drop-btn').addEventListener('click', async () => {
        await dropStash(project.path, stash.index);
        panel.remove();
      });
      list.appendChild(item);
    }
    panel.appendChild(list);
  }

  card.appendChild(panel);
}

async function fetchStashes(repoPath) {
  try {
    const res = await fetch(`/api/stash?repo=${encodeURIComponent(repoPath)}`);
    const data = await res.json();
    return data.stashes || [];
  } catch (err) {
    return { error: err.message };
  }
}

async function applyStash(repoPath, index) {
  try {
    const res = await fetch(`/api/stash?repo=${encodeURIComponent(repoPath)}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ action: 'apply', index })
    });
    return res.ok;
  } catch (err) {
    return false;
  }
}

async function dropStash(repoPath, index) {
  try {
    const res = await fetch(`/api/stash?repo=${encodeURIComponent(repoPath)}`, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ index })
    });
    return res.ok;
  } catch (err) {
    return false;
  }
}