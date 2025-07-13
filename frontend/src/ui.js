import { initializeAdvancedNetworkTools } from './sections/advanced_network_tools.js';
import { RegisterDesktopEntry } from "../wailsjs/go/backend/SetupService";
import { GenerateReport, SaveReport } from "../wailsjs/go/report/ReportService";
import { Ping, Traceroute } from "../wailsjs/go/tools/NetworkToolsService";


export function createBootBox(id, icon, title, parentElement = null) {
  const wrapper = document.createElement('div');
  wrapper.id = id;
  wrapper.className = 'terminal-output box blink-box hidden';

  const header = document.createElement('div');
  header.className = 'box-header';
  header.innerHTML = `${icon} <strong>${title}</strong>`;
  wrapper.appendChild(header);

  const content = document.createElement('div');
  content.className = 'box-content';
  wrapper.appendChild(content);

  const targetParent = parentElement || document.querySelector('#content-view'); // Standard-Parent ist jetzt #content-view
  targetParent.appendChild(wrapper);
  return { wrapper, content };
}

export function appendBootLine(contentEl, text) {
  const line = document.createElement('div');
  line.className = 'line-appear';
  line.textContent = text;
  contentEl.appendChild(line);
  return Promise.resolve();
}

export async function createNetworkSection(localIP, publicIP) {
  const box = createBootBox('network-status-box', '🌐', 'Netzwerkstatus');
  await appendBootLine(box.content, `Lokale IP: ${localIP}`);
  await appendBootLine(box.content, `Öffentliche IP: ${publicIP}`);
  return box.wrapper;
}

export async function createSystemInfoSection({ hostname, os, arch, username, workingDir }) {
  const box = createBootBox('system-info-box', '💻', 'Systeminformationen');
  box.wrapper.id = 'system-info-section';
  await appendBootLine(box.content, `Hostname: ${hostname}`);
  await appendBootLine(box.content, `OS: ${os} (${arch})`);
  await appendBootLine(box.content, `Benutzer: ${username}`);
  await appendBootLine(box.content, `Arbeitsverzeichnis: ${workingDir}`);
  return box.wrapper;
}

export function createNetworkToolsSection() {
  const section = document.createElement('div');
  section.id = 'network-tools-section';
  section.className = 'terminal-output box blink-box';
  section.innerHTML = `
    <div class="box-header">🛠️ Netzwerk-Tools</div>
    <div class="box-content">
      <input type="text" id="hostInput" placeholder="Host oder IP-Adresse" class="input-field">
      <div class="button-group">
        <button id="pingBtn" class="btn">Ping</button>
        <button id="tracerouteBtn" class="btn">Traceroute</button>
      </div>
      <pre id="networkToolsOutput" class="output-field"></pre>
    </div>
  `;

  // Event Listener hinzufügen
  const hostInput = section.querySelector('#hostInput');
  const pingBtn = section.querySelector('#pingBtn');
  const tracerouteBtn = section.querySelector('#tracerouteBtn');
  const outputField = section.querySelector('#networkToolsOutput');

  pingBtn.onclick = async () => {
    const host = hostInput.value.trim();
    if (!host) {
      outputField.textContent = "Bitte geben Sie einen Host ein.";
      return;
    }
    outputField.textContent = `Pinging ${host}...\n`;
    try {
      const result = await Ping(host);
      if (result.error) {
        outputField.textContent += `Fehler: ${result.error}\n`;
      } else {
        outputField.textContent += `Host: ${result.host} (${result.ip})\n`;
        outputField.textContent += `Pakete gesendet: ${result.packets}, Verlust: ${result.loss}\n`;
        outputField.textContent += `Min/Avg/Max RTT: ${result.minRtt}/${result.avgRtt}/${result.maxRtt}\n`;
      }
    } catch (err) {
      outputField.textContent += `Ein unerwarteter Fehler ist aufgetreten: ${err}\n`;
      console.error("Ping-Fehler:", err);
    }
  };

  tracerouteBtn.onclick = async () => {
    const host = hostInput.value.trim();
    if (!host) {
      outputField.textContent = "Bitte geben Sie einen Host ein.";
      return;
    }
    outputField.textContent = `Tracerouting to ${host}...\n`;
    try {
      const hops = await Traceroute(host);
      if (hops && hops.length > 0) {
        hops.forEach(hop => {
          outputField.textContent += `${hop.n}. ${hop.host} (${hop.address}) - ${hop.rtt}\n`;
        });
      } else {
        outputField.textContent += "Keine Hops gefunden oder Fehler.\n";
      }
    } catch (err) {
      outputField.textContent += `Ein unerwarteter Fehler ist aufgetreten: ${err}\n`;
      console.error("Traceroute-Fehler:", err);
    }
  };

  return section;
}

export function createSettingsSection() {
  const section = document.createElement('div');
  section.id = 'settings-section';
  section.className = 'terminal-output box blink-box';
  section.innerHTML = `
    <div class="box-header">⚙️ Einstellungen</div>
    <div class="box-content">Einstellungen folgen...</div>
  `;

  const reportBtn = document.createElement('button');
  reportBtn.textContent = '📄 Bericht erstellen';
  reportBtn.className = 'btn';
  reportBtn.onclick = () => {
    GenerateReport().then(reportData => {
      window.runtime.SaveFileDialog({
        defaultFilename: 'privacy-buddy-report.json',
        defaultDirectory: '~/',
        title: 'Diagnosebericht speichern',
        filters: [{ displayName: 'JSON Files (*.json)', pattern: '*.json' }],
      }).then(filePath => {
        if (filePath) {
          SaveReport(reportData, filePath).then(() => {
            alert("Bericht wurde erfolgreich gespeichert.");
          }).catch(err => {
            console.error("Fehler beim Speichern des Berichts:", err);
            alert("Fehler beim Speichern des Berichts.");
          });
        }
      });
    });
  };

  const installBtn = document.createElement('button');
  installBtn.textContent = '🛠️ Installieren';
  installBtn.className = 'btn install-btn';
  installBtn.onclick = () => {
    RegisterDesktopEntry()
      .then(() => alert("✅ Eintrag wurde erstellt!"))
      .catch(err => {
        console.error("Install-Fehler:", err);
        alert("❌ Registrierung fehlgeschlagen.");
      });
  };

  const exitBtn = document.createElement('button');
  exitBtn.textContent = '🚪 Beenden';
  exitBtn.className = 'btn exit-btn';
  exitBtn.onclick = () => Quit();

  section.appendChild(reportBtn);
  section.appendChild(installBtn);
  section.appendChild(exitBtn);

  return section;
}

export function createFileSection() {
  const section = document.createElement('div');
  section.id = 'files-section';
  section.className = 'terminal-output box blink-box';
  section.innerHTML = `
    <div class="box-header">📁 Dateien</div>
    <div class="box-content">Bald verfügbar</div>
  `;
  return section;
}

export function createGreetingSection(username) {
  const section = document.createElement('div');
  section.className = 'terminal-output box blink-box';
  section.innerHTML = `<div class="box-content" style="text-align: center;">👋 Hallo, ${username}!</div>`;
  return section;
}

export async function createAdvancedNetworkToolsSection() {
  console.log("createAdvancedNetworkToolsSection: Funktion wird ausgeführt.");
  const section = document.createElement('div');
  section.id = 'advanced-network-tools-section';
  section.className = 'terminal-output box blink-box section';
  section.innerHTML = `<div class="box-header"> Erweiterte Netzwerk-Tools</div><div class="box-content">Lade Inhalt...</div>`;

  console.log("createAdvancedNetworkToolsSection: Starte fetch für HTML-Inhalt.");
  try {
    const r = await fetch('./src/sections/advanced_network_tools.html');
    if (!r.ok) {
        console.error(`Fehler beim Laden der HTML-Datei: ${r.status} ${r.statusText}`);
        throw new Error('Netzwerkantwort war nicht ok.');
    }
    const html = await r.text();

    console.log("createAdvancedNetworkToolsSection: HTML-Inhalt erfolgreich geladen.");
    section.innerHTML = html;
    console.log("createAdvancedNetworkToolsSection: HTML in DOM-Struktur eingefügt.");
    if (typeof initializeAdvancedNetworkTools === 'function') {
        setTimeout(() => {
          initializeAdvancedNetworkTools(section);
          console.log("createAdvancedNetworkToolsSection: Event-Handler initialisiert.");
        }, 0);
    } else {
        console.warn("initializeAdvancedNetworkTools ist keine Funktion. Event-Handler wurden nicht gesetzt.");
    }
  } catch (error) {
    console.error('Fehler während des Fetch-Vorgangs:', error);
    section.querySelector('#advanced-network-tools-container').innerHTML = `<p style="color: red;">Fehler beim Laden der Tools: ${error.message}</p>`;
  }

  console.log("createAdvancedNetworkToolsSection: Gebe das erstellte (jetzt gefüllte) Element zurück.");
  return section;
}
