# Web Desktop Javascript Gateway Interface (WDJGI) Dokumentation

## Was ist WDJGI?

WDJGI steht für Web Desktop Javascript Gateway Interface. Mit einfachen Worten können Sie Ihrem System mit JavaScript Funktionen hinzufügen :)

## Verwendung

1. Legen Sie Ihre JS- oder WDJGI-Datei im Verzeichnis web/* ab (z.B. ./app/Dummy/backend/test.js).
2. Laden Sie Ihr Skript, indem Sie einen AJAX-Request an ```/system/wdjgi/interface?script={yourfile}.js``` senden (z.B. /system/wdjgi/interface?script=Dummy/backend/test.js).
3. Warten Sie auf die Antwort des Skripts, indem Sie in Ihrem Skript `sendResp` aufrufen.

## Initialisierungsskript für Module

Um ein Modul ohne einen Hauptfunktionsaufruf in main.go zu initialisieren, können Sie ein "init.wdjgi"-Skript in Ihrem Modulroot unter ./app/meinModul erstellen, wobei "meinModul" der Name Ihres Moduls ist.

Um das Modul zu registrieren, rufen Sie die Funktion "registerModule" mit den JSON-Stringify-Modullaunch-Infos auf, wie im folgenden JavaScript-Objekt beschrieben.

```javascript
// Definieren Sie die launchInfo für das Modul
var moduleLaunchInfo = {
    Name: "NotepadA",
    Desc: "Der beste Code-Editor auf WDOS",
    Group: "Office",
    IconPath: "NotepadA/img/module_icon.png",
    Version: "1.2",
    StartDir: "NotepadA/index.html",
    SupportFW: true,
    LaunchFWDir: "NotepadA/index.html",
    SupportEmb: true,
    LaunchEmb: "NotepadA/embedded.html",
    InitFWSize: [1024, 768],
    InitEmbSize: [360, 200],
    SupportedExt: [".bat",".coffee",".cpp",".cs",".csp",".csv",".fs",".dockerfile",".go",".html",".ini",".java",".js",".lua",".mips",".md", ".sql",".txt",".php",".py",".ts",".xml",".yaml"]
}

// Modul registrieren
registerModule(JSON.stringify(moduleLaunchInfo));
```

Sie können auch die Datenbanktabelle in diesem Abschnitt des Codes erstellen. Zum Beispiel:

```javascript
// Erstellen Sie eine Datenbank für dieses Modul
newDBTableIfNotExists("meinModul")
```

## Anwendungsbeispiele

Sehen Sie sich web/UnitTest/backend/*.js für weitere Informationen zur Verwendung von WDJGI in Web-Apps an.

Für Subdienste finden Sie weitere Beispiele unter subservice/demo/wdjgi/.

### Zugriff von der Frontend

Um Serverfunktionen vom Frontend aus zu nutzen (z.B. beim Erstellen einer serverlosen Webanwendung auf Basis von WDOS), können Sie die Funktion `ao_module.js` aufrufen, um ein WDJGI-Skript im Verzeichnis ```./app``` auszuführen. Sie finden den `ao_module.js`-Wrapper im Verzeichnis ```./app/script/```.

Hier ist ein Beispiel aus dem Musikmodul zum Auflisten von Dateien in der Nähe der geöffneten Musikdatei.

./app/Music/embedded.html

```javascript
ao_module_agirun("Music/functions/getMeta.js", {
    file: encodeURIComponent(playingFileInfo.filepath)
}, function(data){
    songList = data;
    for (var i = 0; i < data.length; i++){
        // Hier etwas machen
    }
});
```

./app/Music/functions/getMeta.js

```javascript
// Helferfunktionen definieren
function bytesToSize(bytes) {
    var sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    if (bytes == 0) return '0 Byte';
    var i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)));
    return (bytes / Math.pow(1024, i)).toFixed(2) + ' ' + sizes[i];
}

// Hauptlogik
if (requirelib("filelib") == true){
    // Dateiname aus Parametern erhalten
    var openingFilePath = decodeURIComponent(file);
    var dirname = openingFilePath.split("/")
    dirname.pop()
    dirname = dirname.join("/");

    // Dateien in der Nähe scannen
    var nearbyFiles = filelib.aglob(dirname + "/*") 
    var audioFiles = [];
    var supportedFormats = [".mp3",".flac",".wav",".ogg",".aac",".webm",".mp4"];
    // Für jede Datei in der Nähe
    for (var i =0; i < nearbyFiles.length; i++){
        var thisFile = nearbyFiles[i];
        var ext = thisFile.split(".").pop();
        ext = "." + ext;
        // Überprüfen, ob die Dateierweiterung in der unterstützten Erweiterungsliste vorhanden ist
        for (var k = 0; k < supportedFormats.length; k++){
            if (filelib.isDir(nearbyFiles[i]) == false && supportedFormats[k] == ext){
                var fileExt = ext.substr(1);
                var fileName = thisFile.split("/").pop();
                var fileSize = filelib.filesize(thisFile);
                var humanReadableFileSize = bytesToSize(fileSize);

                var thisFileInfo = [];
                thisFileInfo.push(fileName);
                thisFileInfo.push(thisFile);
                thisFileInfo.push(fileExt);
                thisFileInfo.push(humanReadableFileSize);

                audioFiles.push(thisFileInfo);
                break;
            }
        }
    }
    sendJSONResp(JSON.stringify(audioFiles));
}
```

### Zugriff von Subservice-Backend

Es ist auch möglich, auf das WDJGI-Gateway von einem Subservice-Backend aus zuzugreifen. Sie können die `wd`-Bibliothek aus ```./subservice/demo/wd``` einbinden. Hier ist ein Beispiel, das aus einem Demo-Subservice stammt und auf Ihre Desktop-Dateiliste zugreift.

```go
package main
import (
    wd "Ihr/paket/name/wd"
)

var handler *wd.WDHandler

//...

func main(){
    // Andere Flags hier einfügen

    // Subservice-Pipeline und Flag-Parsing starten (Dieser Funktionsaufruf führt auch flag.parse() aus)
    handler = wd.HandleFlagParse(wd.ServiceInfo{
        Name: "Demo Subservice",
        Desc: "Ein einfacher Subservice-Code, der zeigt, wie Subservices in WDOS funktionieren",            
        Group: "Entwicklung",
        IconPath: "demo/demo.png",
    })

    //... 
    // Handler und

 andere Logik hier
    //...

    // Hier können Sie auf das WDJGI-Gateway zugreifen
    go func(){
        for{
            wd.ExecuteWDJGICall("Desktop/services/fileManager/main.js", nil, nil)
            time.Sleep(5 * time.Second) // oder einen anderen Intervall, den Sie mögen
        }
    }()
    
    //...
    // Rest des Codes
    //...
}
```

## Weiterführende Informationen

- WDJGI ist ein fortschrittliches Werkzeug, das es Benutzern ermöglicht, schnell und einfach neue Funktionen in WDOS hinzuzufügen.
- Bitte beachten Sie, dass WDJGI für fortgeschrittene Benutzer gedacht ist, die mit JavaScript und Go vertraut sind.
- Wenn Sie auf Probleme oder Fragen stoßen, wenden Sie sich bitte an die WDOS-Community. Wir sind hier, um zu helfen!
