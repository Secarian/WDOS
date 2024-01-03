# WDOS 0.0.1

Dies ist die Go-Implementierung der WDOS Web-Desktop-Umgebung, die für Linux entwickelt wurde, aber aus irgendeinem Grund auch auf Windows und Mac OS funktioniert.

Diese README-Datei richtet sich ausschließlich an Entwickler. Wenn Sie ein normaler Benutzer sind, konsultieren Sie bitte die README-Datei außerhalb des /src-Ordners.

## Entwicklungsnotizen

- Beginnen Sie jedes Modul mit der Funktion {ModuleName}Init(), z. B. ```WiFiInit()```.
- Platzieren Sie Ihre Funktion im Modul (wenn möglich) und rufen Sie sie im Hauptprogramm auf.
- Ändern Sie die Reihenfolge in der startup() Funktion nicht, es sei denn, es ist notwendig.
- Bei Unsicherheiten fügen Sie Startflags hinzu (und verwenden Sie Startflags, um experimentelle Funktionen beim Start zu deaktivieren).

## Überschreiben von Vendor-Ressourcen

Wenn Sie Vendor-bezogene Ressourcen in WDOS 0.0.1 oder höher überschreiben möchten, erstellen Sie einen Ordner im Systemstamm mit dem Namen ```vendor-res``` und platzieren Sie die Ersatzdateien hier. Hier ist eine Liste der unterstützten Ersatzressourcendateien:

| Dateiname        | Empfohlenes Format | Verwendung              |
| ---------------- | ------------------ | ----------------------- |
| auth_bg.jpg      | 2938 x 1653 px     | Anmeldehintergrund      |
| auth_icon.png    | 5900 x 1180 px     | Logo der Authentifizierungsseite |
| vendor_icon.png  | 1560 x 600 px      | Marken-Symbol des Anbieters |

(Wird erweitert)

## Dateisystem-Virtualisierung und Abstraktionsebenen

Das WDOS-System enthält sowohl die Virtualisierungsebene als auch die Abstraktionsebene. Der einfachste Weg zu überprüfen, unter welcher Ebene Ihr Pfad liegt, besteht darin, den Startordnernamen zu betrachten.

| Pfadstruktur                                 | Beispiel Pfad                                        | Ebene                                            |
| -------------------------------------------- | ---------------------------------------------------- | ------------------------------------------------ |
| {vroot_id}:/{subpath}                        | user:/Desktop/meinedatei.txt                         | Dateisystem-Virtualisierungsebene (höchste Ebene) |
| fsh (*Dateisystem-Handler*) + subpath (String)| fsh (localfs) + /dateien/benutzer/alan/Desktop/meinedatei.txt | Dateisystem-Abstraktionsebene                     |
| {physische_position}/{subpath}                | /home/wd/wdos/dateien/benutzer/Desktop/meinedatei.txt| Physische (Festplatten-) Ebene                    |

Seit WDOS v2.000 haben wir der (schon komplexen) Dateisystem-Handler-Infrastruktur eine Dateisystem-Abstraktion (fsa oder manchmal als fshAbs abgekürzt, für "Dateisystem-Handler unterliegende Dateisystemabstraktion") hinzugefügt. Es gibt zwei Arten von fsh, die derzeit von der WDOS-Dateisystem-Abstraktionsebene unterstützt werden.

## WDOS JavaScript Gateway Interface / Plugin Loader

Die WDOS AJGI / AGI-Schnittstelle bietet eine JavaScript-programmierbare Schnittstelle für WDOS-Benutzer zur Erstellung von Plugins für das System. Um das Modul zu initialisieren, können Sie eine "init.agi"-Datei im Webverzeichnis des Moduls (auch als Modul-Stammverzeichnis bezeichnet) platzieren. Weitere Details finden Sie in der [AJGI-Dokumentation](AJGI Dokumentation.md).

AGI-Skripte können in verschiedenen Bereichen und mit verschiedenen Berechtigungen ausgeführt werden.

| Bereich                                        | Verwendbare Funktionen                                                      |
| ---------------------------------------------- | -------------------------------------------------------------------------- |
| WebApp-Startskript (init.agi)                  | Systemfunktionen und Registrierungen                                         |
| In der WebApp enthaltene Skripte                | Systemfunktionen und Benutzerfunktionen                                       |
| Andere (Web-Root / Serverless / Zeitplaner)     | Systemfunktionen, Benutzerfunktionen (mit Skript-Registrierungsbereich) und Serverless |

## Unterlogik und Konfiguration von Subdiensten

Um andere binär basierte Webserver in die Subdienstschnittstelle zu integrieren, erstellen Sie einen Ordner innerhalb von "./subservice/ihres_dienstes", in dem Ihre binäre ausführbare Datei den gleichen Namen wie das enthaltende Verzeichnis haben sollte. Wenn Sie beispielsweise ein Modul haben, das eine Web-UI namens "demo.exe" bereitstellt, sollten Sie die demo.exe in "./subservice/demo/demo.exe" platzieren.

Im Falle einer Linux-Umgebung wird zuerst überprüft, ob das Modul über apt-get installiert ist, indem das "which"-Programm verwendet wird. (Wenn Sie busybox haben, sollte es integriert sein). Wenn das Paket in der apt-Liste nicht gefunden wird, wird die Binärdatei des Programms unter dem Subdienstverzeichnis gesucht.

Bitte befolgen Sie die Namenskonvention, die in der build.sh-Vorlage angegeben ist. Zum Beispiel sucht die entsprechende Plattform nach dem entsprechenden binären ausführbaren Dateinamen:

```
demo_linux_amd64    => Linux AMD64
demo_linux_arm      => Linux ARMv6l / v7l
demo_linux_arm64    => Linux ARM64
demo_macOS_amd64    => MacOS AMD64 
```

### Startflags

Während des Starts des Subdienstes werden zwei Arten von Parametern übergeben. Hier sind einige Beispiele:

```
demo.exe -info
demo.exe -port 12810 -rpt "http://localhost:8080/api/ajgi/interface"
```

Im Falle des Erhalts des "info"-Flags sollte das Programm die JSON-Zeichenfolge mit den korrekten Modulinformationen gemäß der unten stehenden Struktur ausgeben.

```
// Struktur zur Speicherung von Modulinformationen
type serviceInfo struct {
    Name string                // Name dieses Moduls, z. B. "Audio"
    Desc string                // Beschreibung dieses Moduls
    Group string               // Gruppe des Moduls, z. B. "System" / "Medien" usw.
    IconPath string            // Pfad zum Modul-Icon-Bild, z. B. "Audio/img/funktions_icon.png"
    Version string             // Version des Moduls. Format: [0-

9]*.[0-9][0-9].[0-9]
    StartDir string            // Standardstartverzeichnis, z. B. "Audio/index.html"
    SupportFW bool             // Unterstützung von floatWindow. Wenn ja, wird das floatWindow-Verzeichnis geladen
    LaunchFWDir string         // Dieser Link wird anstelle von 'StartDir' gestartet, wenn im FW-Modus
    SupportEmb bool            // Unterstützung des eingebetteten Modus
    LaunchEmb string           // Dieser Link wird anstelle von StartDir / Fw gestartet, wenn eine Datei mit diesem Modul geöffnet wird
    InitFWSize []int           // Floatwindow-Initialgröße. [0] => Breite, [1] => Höhe
    InitEmbSize []int          // Initialgröße im eingebetteten Modus. [0] => Breite, [1] => Höhe
    SupportedExt []string      // Unterstützte Dateierweiterungen, z. B. ".mp3", ".flac", ".wav"
}

// Beispiel für die Verwendung beim Erhalten des -info-Flags
infoObject := serviceInfo{
    Name: "Demo Subdienst",
    Desc: "Ein einfacher Subdienstcode, der zeigt, wie Subdienste in WDOS funktionieren",
    Group: "Entwicklung",
    IconPath: "demo/icon.png",
    Version: "0.0.1",
    StartDir: "demo/home.html",
    SupportFW: true,
    LaunchFWDir: "demo/home.html",
    SupportEmb: true,
    LaunchEmb: "demo/eingebettet.html",
    InitFWSize: []int{720, 480},
    InitEmbSize: []int{720, 480},
    SupportedExt: []string{".txt", ".md"},
}

jsonString, _ := json.Marshal(infoObject)
fmt.Println(string(jsonString))
os.Exit(0)
```

Beim Erhalt des Port-Flags sollte das Programm die Web-UI am angegebenen Port starten. Hier ist ein Beispiel für die Implementierung einer solchen Funktionalität.

```
var port = flag.String("port", ":80", "Der Standard-Listening-Endpunkt für diesen Subdienst")
flag.Parse()
err := http.ListenAndServe(*port, nil)
if err != nil {
    log.Fatal(err)
}
```

### Subdienst-Ausführungseinstellungen

Standardmäßig erstellt die Subdienstroutine einen Reverse-Proxy mit integriertem URL-Umschreiben, der Ihre Web-UI startet, die von der binären ausführbaren Datei ausgeführt wird. Wenn Sie keine Reverse-Proxy-Verbindung benötigen, ein benutzerdefiniertes Startskript oder sonstiges wünschen, können Sie die folgenden Einstellungsdateien verwenden.

```
.noproxy        => Starten Sie keinen Proxy zum angegebenen Port
.startscript    => Senden Sie den Startparameter an die "start.bat" oder "start.sh" Datei anstelle der binären ausführbaren Datei
.disabled       => Diesen Subdienst beim Starten nicht laden. Der Benutzer kann ihn jedoch über die Einstellungsschnittstelle aktivieren
```

Hier ist ein Beispiel für eine "start.bat"-Datei, die zum Integrieren von Syncthing in das WDOS-Online-System mit einer ".startscript"-Datei neben der syncthing.exe-Datei verwendet wird.

```
if not exist ".\config" mkdir ".\config"
syncthing.exe -home=".\config" -no-browser -gui-address=127.0.0.1%2
```

## Systemd-Unterstützung

Um systemd in Ihrem Host zu aktivieren, der das WDOS-Online-System unterstützt, erstellen Sie ein Bash-Skript in Ihrem WDOS-Online-Stammverzeichnis mit dem Namen "start.sh" und füllen Sie es mit Ihren bevorzugten Startparametern. Der einfachste ist wie folgt:

```
#/bin/bash
sudo ./wdos_online_linux_amd64
```

Danach können Sie eine neue Datei namens "wdos.service" in /etc/systemd/system mit dem folgenden Inhalt erstellen (angenommen, Ihr WDOS-Online-Stammverzeichnis befindet sich unter /home/pi/wdos):

```
[Unit]
Description=WDOS Cloud Desktop Service.

[Service]
Type=simple
WorkingDirectory=/home/pi/wdos/
ExecStart=/bin/bash /home/pi/wdos/start.sh

Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Schließlich, um den Dienst zu aktivieren, verwenden Sie die folgenden systemd-Befehle:

```
# Skript während des Startvorgangs aktivieren
sudo systemctl enable wdos.service

# Den Dienst jetzt starten
sudo systemctl start wdos.service

# Den Status des Dienstes anzeigen
systemctl status wdos.service

# Den Dienst deaktivieren, wenn Sie ihn beim Start nicht mehr starten möchten
sudo systemctl disable wd-online.service
```
