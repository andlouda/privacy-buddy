import './app.css';
import './style.css';
import './sections.css';
import './sidebar/sidebar.css'; // Sidebar CSS laden
import { runBootSequence } from './boot.js';
import { createSidebar } from './sidebar/sidebar.js'; // createSidebar importieren
import { showView } from './navigation.js'; // showView importieren

document.addEventListener('DOMContentLoaded', async () => {
  // Start the boot sequence
  const appInitialContext = await runBootSequence(); // Kontext abfangen

  // Sidebar einmalig rendern
  const sidebarContainer = document.getElementById('sidebar-container');
  sidebarContainer.appendChild(createSidebar((id) => showView(id, appInitialContext))); // Kontext an showView übergeben

  // Nach dem Boot, die Netzwerkansicht standardmäßig anzeigen
  showView('network', appInitialContext);
});
