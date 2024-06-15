javascript:(function() {
    if (!window.location.href.match(/^(http(s)?:\/\/(www\.)?smithersradio.com|file:\/\/).*/)) {
        console.warn("url not supported by CICK Playlister: " + window.location.href);
        return;
    }
    function reportLoadError() {
        alert("Error loading CICK Playlister. Please report an issue at https://github.com/captaincoordinates/cick-playlister/issues");
    }
    if (window.cickPlaylisterClient) {
        window.cickPlaylisterClient.show();
        return;
    }
    var stylesheet = document.createElement("link");
    stylesheet.rel = "stylesheet";
    stylesheet.href = "http://localhost:8123/client/dist/cick-playlister-client.css";
    stylesheet.onerror = reportLoadError;
    stylesheet.onload = function() {
        /*
        Make script load conditional on stylesheet load so that styling is guaranteed
        to be applied before the script runs and adds DOM elements.
        */
        var script = document.createElement("script");
        script.src = "http://localhost:8123/client/dist/cick-playlister-client.js";
        script.onerror = reportLoadError;
        script.onload = function() {
            var anchor = document.createElement("div");
            anchor.id = window.cickPlaylisterClient.anchorId;
            document.body.appendChild(anchor);
            window.cickPlaylisterClient.show();
        };
        document.body.appendChild(script);
    };
    document.body.appendChild(stylesheet);
})();
