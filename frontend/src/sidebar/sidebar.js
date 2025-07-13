export function createSidebar(showView) {
  const sidebar = document.createElement('div');
  sidebar.className = 'sidebar';

  const icons = [
    { emoji: '🛰️', id: 'network' },
    { emoji: '🛠️', id: 'tools' },
    { emoji: '📡', id: 'advanced-network-tools-section' },
    { emoji: '🔐', id: 'security' }, // optional: noch keine View
    { emoji: '📁', id: 'files' },
    { emoji: '⚙️', id: 'settings' }
  ];

  icons.forEach(icon => {
    const tile = document.createElement('div');
    tile.className = 'sidebar-tile';
    tile.innerText = icon.emoji;

    tile.onclick = () => showView(icon.id); // <<< Wichtig: nur "network", nicht "section-network"

    sidebar.appendChild(tile);
  });

  return sidebar;
}
