import { GetNetworkInterfaces } from '../../wailsjs/go/tools/NetworkToolsService';
import { GetPublicIPInfo } from '../../wailsjs/go/network/PublicIPService';
import { StartPacketCapture, StopPacketCapture, GetCaptureTemplates, SaveCaptureTemplate } from '../../wailsjs/go/tools/AdvancedNetworkToolsService';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';

// Hilfsfunktion, um die Interfaces-Tabelle zu generieren
function generateInterfaceTable(interfaces, externalInterfaceName) {
  if (!interfaces || interfaces.length === 0) {
    return '<p>No network interfaces found.</p>';
  }

  const headers = ['Name', 'Description', 'MAC Address', 'IP Addresses', 'Status'];
  let table = '<table><thead><tr>';
  headers.forEach(header => {
    table += `<th>${header}</th>`;
  });
  table += '</tr></thead><tbody>';

  interfaces.forEach(iface => {
    const name = iface.displayName || iface.name || '';
    const description = iface.description || '';

    const isUplink =
      externalInterfaceName &&
      (
        name.toLowerCase().includes(externalInterfaceName.toLowerCase()) ||
        description.toLowerCase().includes(externalInterfaceName.toLowerCase())
      );

    const isUp = iface.isUp;

    const statusText = isUplink
      ? '⚡ Uplink'
      : isUp
        ? '🟢 Up'
        : '🔴 Down';

    const rowClass = isUplink ? 'class="uplink-interface"' : (isUp ? 'class="up-interface"' : 'class="down-interface"');

    table += `<tr ${rowClass}>`;
    table += `<td>${name}</td>`;
    table += `<td>${description}</td>`;
    table += `<td>${iface.hardwareAddr || '-'}</td>`;
    table += `<td>${(iface.addrs || []).join('<br>') || '-'}</td>`;
    table += `<td>${statusText}</td>`;
    table += '</tr>';
  });

  table += '</tbody></table>';
  return table;
}


// Befüllt das Dropdown-Menü mit den verfügbaren Netzwerkschnittstellen.
function populateInterfaceDropdown(selectElement, interfaces) {
  selectElement.innerHTML = '';
  if (!interfaces || interfaces.length === 0) {
    const option = document.createElement('option');
    option.value = '';
    option.textContent = 'No interfaces available';
    option.disabled = true;
    selectElement.appendChild(option);
  } else {
    const defaultOption = document.createElement('option');
    defaultOption.value = '';
    defaultOption.textContent = 'Select an interface';
    selectElement.appendChild(defaultOption);
    interfaces.forEach(iface => {
      const option = document.createElement('option');
      option.value = iface.name; // Use the pcap name as the value
      option.textContent = iface.displayName || iface.name; // Display the user-friendly name
      selectElement.appendChild(option);
    });
  }
}

function populateTemplateDropdown(selectElement, templates, bpfFilterInput, captureDurationInput) {
  selectElement.innerHTML = '<option value="">Select a template</option>';
  templates.forEach(template => {
    const option = document.createElement('option');
    option.value = template.Name;
    option.textContent = template.Name + (template.Description ? ` (${template.Description})` : '');
    option.dataset.bpfFilter = template.BPFFilter;
    option.dataset.duration = template.Duration || ''; // Assuming templates might have a default duration
    selectElement.appendChild(option);
  });

  selectElement.addEventListener('change', () => {
    const selectedOption = selectElement.options[selectElement.selectedIndex];
    if (selectedOption.value) {
      bpfFilterInput.value = selectedOption.dataset.bpfFilter;
      if (selectedOption.dataset.duration) {
        captureDurationInput.value = selectedOption.dataset.duration;
      }
    } else {
      // Clear if "Select a template" is chosen
      bpfFilterInput.value = '';
      captureDurationInput.value = '10'; // Reset to default
    }
  });
}

export function initializeAdvancedNetworkTools(sectionElement) {
  console.log('Initializing advanced network tools...');

  const getInterfacesBtn = sectionElement.querySelector('#get-interfaces-btn');
  const interfacesOutput = sectionElement.querySelector('#interfaces-output');
  const interfaceSelect = sectionElement.querySelector('#interface-select');

  // Packet Capture elements
  const captureInterfaceInput = sectionElement.querySelector('#capture-interface'); // This is now redundant, but keeping for reference
  const bpfFilterInput = sectionElement.querySelector('#bpf-filter');
  const captureDurationInput = sectionElement.querySelector('#capture-duration');
  const startCaptureBtn = sectionElement.querySelector('#start-capture-btn');
  const stopCaptureBtn = sectionElement.querySelector('#stop-capture-btn');
  const packetCaptureOutput = sectionElement.querySelector('#packet-capture-output');
  const templateSelect = sectionElement.querySelector('#template-select');

  // New elements for saving templates
  const newTemplateNameInput = sectionElement.querySelector('#new-template-name');
  const newTemplateDescriptionInput = sectionElement.querySelector('#new-template-description');
  const saveTemplateBtn = sectionElement.querySelector('#save-template-btn');
  const saveTemplateOutput = sectionElement.querySelector('#save-template-output');

  let packetCaptureListener = null; // To store the listener for later removal

  if (getInterfacesBtn && interfacesOutput && interfaceSelect &&
      startCaptureBtn && stopCaptureBtn && packetCaptureOutput && templateSelect &&
      newTemplateNameInput && newTemplateDescriptionInput && saveTemplateBtn && saveTemplateOutput) {

    // Initial state for capture buttons
    stopCaptureBtn.disabled = true;

    // Load and populate templates on initialization
    GetCaptureTemplates().then(templates => {
      populateTemplateDropdown(templateSelect, templates, bpfFilterInput, captureDurationInput);
    }).catch(error => {
      console.error('Error loading capture templates:', error);
      // Optionally display an error to the user
    });

    getInterfacesBtn.addEventListener('click', () => {
      interfacesOutput.textContent = 'Loading interfaces...';
      interfaceSelect.innerHTML = '<option>Loading...</option>';
      interfaceSelect.disabled = true;

      Promise.all([
        GetNetworkInterfaces(),
        GetPublicIPInfo()
      ])
      .then(([interfaces, publicIpInfo]) => {
        const externalInterfaceName = publicIpInfo ? publicIpInfo.interfaceName : '';
        
        interfacesOutput.innerHTML = generateInterfaceTable(interfaces, externalInterfaceName);
        populateInterfaceDropdown(interfaceSelect, interfaces);
        interfaceSelect.disabled = false;

        if (externalInterfaceName) {
          const matchingIface = interfaces.find(iface =>
            (iface.displayName && iface.displayName.toLowerCase().includes(externalInterfaceName.toLowerCase())) ||
            (iface.description && iface.description.toLowerCase().includes(externalInterfaceName.toLowerCase()))
          );

          if (matchingIface) {
            interfaceSelect.value = matchingIface.name; // Set correct value (PCAP name)
          } else {
            console.warn('No matching interface found for externalInterfaceName:', externalInterfaceName);
          }
        }
      })
      .catch(error => {
        console.error('Error getting network info:', error);
        interfacesOutput.textContent = `Error loading interfaces: ${error}`;
        populateInterfaceDropdown(interfaceSelect, []);
        interfaceSelect.disabled = false;
      });
    });

    // Start Capture Button Listener
    startCaptureBtn.addEventListener('click', async () => {
      const selectedInterface = interfaceSelect.value; // Use selected value from dropdown
      const bpfFilter = bpfFilterInput.value;
      const duration = parseInt(captureDurationInput.value, 10);

      if (!selectedInterface) {
        packetCaptureOutput.textContent = 'Please select an interface to capture.';
        return;
      }
      if (isNaN(duration) || duration <= 0) {
        packetCaptureOutput.textContent = 'Please enter a valid capture duration (seconds).';
        return;
      }

      packetCaptureOutput.textContent = `Starting packet capture on ${selectedInterface} for ${duration} seconds...`;
      packetCaptureOutput.innerHTML = ''; // Clear previous output

      startCaptureBtn.disabled = true;
      stopCaptureBtn.disabled = false;

      // Register event listener for incoming packets
      packetCaptureListener = EventsOn('packetCaptureEvent', (packet) => {
        const packetInfo = `[${packet.Timestamp}] ${packet.Source} -> ${packet.Destination} (${packet.Protocol}) Len: ${packet.Length} Summary: ${packet.Summary}`;
        const p = document.createElement('p');
        p.textContent = packetInfo;
        packetCaptureOutput.appendChild(p);
        // Scroll to bottom
        packetCaptureOutput.scrollTop = packetCaptureOutput.scrollHeight;
      });

      // Register event listener for capture stopped
      EventsOn('packetCaptureStopped', (message) => {
        packetCaptureOutput.textContent += `\nCapture stopped: ${message}`;
        startCaptureBtn.disabled = false;
        stopCaptureBtn.disabled = true;
        if (packetCaptureListener) {
          EventsOff('packetCaptureEvent'); // Unregister the packet listener
          packetCaptureListener = null;
        }
      });

      try {
        await StartPacketCapture(selectedInterface, bpfFilter, duration);
        packetCaptureOutput.textContent += '\nCapture started successfully. Waiting for packets...';
      } catch (error) {
        console.error('Error starting packet capture:', error);
        packetCaptureOutput.textContent = `Error starting capture: ${error}`;
        startCaptureBtn.disabled = false;
        stopCaptureBtn.disabled = true;
        if (packetCaptureListener) {
          EventsOff('packetCaptureEvent');
          packetCaptureListener = null;
        }
      }
    });

    // Stop Capture Button Listener
    stopCaptureBtn.addEventListener('click', async () => {
      packetCaptureOutput.textContent += '\nStopping capture...';
      stopCaptureBtn.disabled = true;
      startCaptureBtn.disabled = false; // Re-enable start button immediately

      try {
        await StopPacketCapture();
        // The 'packetCaptureStopped' event will handle final UI updates
      } catch (error) {
        console.error('Error stopping packet capture:', error);
        packetCaptureOutput.textContent += `\nError stopping capture: ${error}`;
      }
    });

    // Save Template Button Listener
    saveTemplateBtn.addEventListener('click', async () => {
      const templateName = newTemplateNameInput.value.trim();
      const templateDescription = newTemplateDescriptionInput.value.trim();
      const bpfFilter = bpfFilterInput.value.trim();
      const duration = parseInt(captureDurationInput.value, 10);

      if (!templateName) {
        saveTemplateOutput.textContent = 'Template Name cannot be empty.';
        return;
      }
      if (!bpfFilter) {
        saveTemplateOutput.textContent = 'BPF Filter cannot be empty.';
        return;
      }

      const newTemplate = {
        Name: templateName,
        Description: templateDescription,
        BPFFilter: bpfFilter,
        Duration: duration, // Include duration in the template
      };

      try {
        await SaveCaptureTemplate(newTemplate);
        saveTemplateOutput.textContent = `Template '${templateName}' saved successfully!`;
        newTemplateNameInput.value = '';
        newTemplateDescriptionInput.value = '';
        // Reload templates to update the dropdown
        const updatedTemplates = await GetCaptureTemplates();
        populateTemplateDropdown(templateSelect, updatedTemplates, bpfFilterInput, captureDurationInput);
      } catch (error) {
        console.error('Error saving template:', error);
        saveTemplateOutput.textContent = `Error saving template: ${error}`;
      }
    });

  } else {
    console.error('Could not find all required elements for advanced network tools.');
  }
}