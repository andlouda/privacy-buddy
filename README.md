 Privacy-Buddy Netzwerk-Analyse-Tool


  Privacy-Buddy ist eine Desktop-Anwendung zur Analyse und Diagnose von Netzwerken. Sie bietet eine benutzerfreundliche Oberfläche, um detaillierte Informationen
  über Netzwerkschnittstellen, Verbindungen und Systemressourcen zu erhalten und gängige Netzwerk-Tools auszuführen.

** Das Tool ist eine private Spielwiese und wurde durch die kostenlose gemini-cli nun auch ein Spielzeug für's Vibe Coding.

  Voraussetzungen


  Stellen Sie sicher, dass die folgenden Abhängigkeiten auf Ihrem System installiert sind, bevor Sie beginnen:


   * Go: Version 1.23 oder höher
   * Node.js & npm: Erforderlich für die Frontend-Abhängigkeiten
   * Wails CLI: Das Toolset zum Bauen von Wails-Anwendungen. Installation via go install github.com/wailsapp/wails/v2/cmd/wails@latest

  Installation

   1. Repository klonen:

   1     git clone https://github.com/your-username/Privacy-Buddy.git
   2     cd Privacy-Buddy



   2. Anwendung im Entwicklungsmodus starten:
      Führen Sie den folgenden Befehl im Hauptverzeichnis des Projekts aus. Wails kümmert sich automatisch um die Installation der Go- und npm-Abhängigkeiten.

   1     wails dev

      Die Anwendung wird gestartet und lädt automatisch neu, wenn Sie Änderungen am Go-Code oder an den Frontend-Dateien vornehmen.


   3. Anwendung für die Produktion bauen:
      Um eine plattformspezifische, ausführbare Datei zu erstellen, verwenden Sie den Build-Befehl:

   1     wails build

      Die kompilierte Anwendung finden Sie anschließend im build/bin-Verzeichnis.

  Verwendung


  Starten Sie die Anwendung, um auf das Haupt-Dashboard zuzugreifen. Die Benutzeroberfläche ist in mehrere Sektionen unterteilt, die über die Seitenleiste
  erreichbar sind.

  Kernfunktionen


   * Netzwerk-Dashboard: Zeigt eine Echtzeit-Übersicht der aktiven Netzwerkverbindungen, der ein- und ausgehenden Datenraten sowie der öffentlichen
     IP-Adresse.
   * Schnittstellen-Details: Listet alle verfügbaren Netzwerkschnittstellen mit ihren Konfigurationen auf, einschließlich IP- und MAC-Adressen.
   * ARP-Cache: Zeigt den aktuellen ARP-Cache des Systems an, um IP-Adressen den entsprechenden MAC-Adressen zuzuordnen.
   * Netzwerk-Tools:
       * Ping: Senden Sie ICMP-Anfragen an eine Ziel-IP oder einen Hostnamen, um die Erreichbarkeit zu prüfen.
       * Traceroute: Verfolgen Sie die Route, die Pakete von Ihrem Computer zu einem Zielhost nehmen.
       * Paketmitschnitt: Erfassen und analysieren Sie den Netzwerkverkehr auf einer ausgewählten Schnittstelle.
   * Systeminformationen: Bietet einen Überblick über grundlegende Systemmetriken wie CPU- und Speicherauslastung.


  Beispiel: Ping ausführen


   1. Navigieren Sie zum Abschnitt "Netzwerk-Tools".
   2. Geben Sie eine Ziel-IP-Adresse (z. B. 8.8.8.8) oder einen Hostnamen (z. B. google.com) in das Ping-Formular ein.
   3. Klicken Sie auf "Start". Die Ergebnisse werden live in der Ausgabe-Konsole angezeigt.

  Beispielausgabe:



   1 PING 8.8.8.8 (8.8.8.8) 56(84) bytes of data.
   2 64 bytes from 8.8.8.8: icmp_seq=1 ttl=116 time=10.5 ms
   3 64 bytes from 8.8.8.8: icmp_seq=2 ttl=116 time=11.2 ms
   4 64 bytes from 8.8.8.8: icmp_seq=3 ttl=116 time=10.8 ms


  Lizenz


  In den Projektdateien ist keine explizite Lizenz angegeben. Bitte kontaktieren Sie den Autor für Informationen zur Nutzung und Weitergabe.

  Hinweise zur Weiterentwicklung


  Die Anwendung ist modular aufgebaut und nutzt plattformspezifische Implementierungen (erkennbar an Dateiendungen wie _windows.go, _linux.go, _darwin.go), um
   eine breite Kompatibilität zu gewährleisten. Beiträge zur Erweiterung der Funktionalität oder zur Verbesserung der plattformspezifischen Features sind
  willkommen.


# Privacy-Buddy

Privacy-Buddy ist eine Desktop-Anwendung, die mit [Wails](https://wails.io/) entwickelt wurde. Sie kombiniert ein Go-Backend mit einem modernen JavaScript-Frontend, das mit Vite erstellt wurde.

## Architektur

Die Anwendung folgt einer klassischen Frontend-Backend-Architektur, die durch Wails ermöglicht wird.

```mermaid
graph TD
    A[Frontend (JavaScript/Vite)] -- Aufrufe --> B{Wails Bridge};
    B -- Go-Methoden --> C[Go Backend];
    C -- Interaktion --> D[Betriebssystem];

    subgraph Frontend
        A
    end

    subgraph Backend
        C
    end
```

### Frontend

Das Frontend befindet sich im `frontend`-Verzeichnis und ist eine Single-Page-Application (SPA), die mit Vite entwickelt wurde. Die Benutzeroberfläche wird mit HTML, CSS und JavaScript erstellt. Die Kommunikation mit dem Backend erfolgt über die von Wails bereitgestellte Bridge, die es ermöglicht, Go-Funktionen direkt aus dem JavaScript-Code aufzurufen.

- **`frontend/src/main.js`**: Der Haupteinstiegspunkt für die Frontend-Anwendung.
- **`frontend/index.html`**: Die Haupt-HTML-Datei.
- **`frontend/wailsjs`**: Enthält die von Wails generierten JavaScript-Bindings für die Go-Methoden.

### Backend

Das Backend ist in Go geschrieben und befindet sich im `backend`-Verzeichnis. Es ist in verschiedene Dienste unterteilt, die für unterschiedliche Aufgaben zuständig sind.

- **`main.go`**: Der Haupteinstiegspunkt für die Go-Anwendung. Hier wird die Wails-Anwendung initialisiert und die verschiedenen Backend-Dienste werden an das Frontend gebunden.
- **`app.go`**: Enthält die Haupt-App-Struktur und einige grundlegende Methoden.
- **`backend/`**: Enthält die Kernlogik der Anwendung, unterteilt in verschiedene Pakete:
    - **`system_service.go`**: Stellt grundlegende Systeminformationen bereit.
    - **`setup.go`**: Verantwortlich für die Einrichtungslogik.
    - **`network/`**: Enthält Dienste im Zusammenhang mit Netzwerkinformationen.
    - **`platform/`**: Enthält plattformspezifischen Code für verschiedene Betriebssysteme (Linux, Windows, macOS).

## Live-Entwicklung

Um die Anwendung im Live-Entwicklungsmodus auszuführen, führe `wails dev` im Projektverzeichnis aus. Dadurch wird ein Vite-Entwicklungsserver gestartet, der ein sehr schnelles Hot-Reload deiner Frontend-Änderungen ermöglicht. Wenn du in einem Browser entwickeln und auf deine Go-Methoden zugreifen möchtest, gibt es auch einen Entwicklungsserver, der auf http://localhost:34115 läuft. Verbinde dich damit in deinem Browser, und du kannst deinen Go-Code aus den Entwicklertools aufrufen.

## Bauen

Um ein weiterverteilbares Produktionspaket zu erstellen, verwende `wails build`.