import { createNetworkSection, createSettingsSection, createFileSection, createGreetingSection, createSystemInfoSection, createNetworkToolsSection, createAdvancedNetworkToolsSection } from './ui.js';

function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' });
}

export async function showView(viewName, context) {
  const contentView = document.querySelector('#content-view'); // Ziel ist der neue Inhaltsbereich
  if (!contentView) {
    console.error("Fehler: #content-view Element nicht gefunden.");
    return;
  }
  // Alle vorhandenen Sektionen ausblenden
  document.querySelectorAll('.section').forEach(sec => sec.classList.remove('active'));

  contentView.innerHTML = ''; // Nur den Inhaltsbereich leeren

  let sectionElement;

  if (viewName === 'network') {
    const greetingSection = createGreetingSection(context.sys.username);
    contentView.appendChild(greetingSection);

    const networkSection = await createNetworkSection(context.localIP, context.publicIP);
    contentView.appendChild(networkSection);

    const systemInfoSection = await createSystemInfoSection(context.sys);
    contentView.appendChild(systemInfoSection);
    sectionElement = networkSection; // Eine der Sektionen als "aktive" Sektion markieren
  } else if (viewName === 'settings') {
    sectionElement = createSettingsSection();
    contentView.appendChild(sectionElement);
  } else if (viewName === 'files') {
    sectionElement = createFileSection();
    contentView.appendChild(sectionElement);
  } else if (viewName === 'tools') {
    sectionElement = createNetworkToolsSection();
    contentView.appendChild(sectionElement);
  } else if (viewName === 'advanced-network-tools-section') {
    console.log("Navigiere zu advanced-network-tools.");
    console.log("Navigiere zu 'advanced-network-tools' und rufe createAdvancedNetworkToolsSection auf.");
    sectionElement = await createAdvancedNetworkToolsSection();
    contentView.appendChild(sectionElement);
  }

  // Die neu hinzugefügte Sektion als aktiv markieren
  if (sectionElement) {
    sectionElement.classList.add('active');
  }
  scrollToTop();
}
