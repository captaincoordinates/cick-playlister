window.cickPlaylisterClient = {
  anchorId: "cick-playlister-client-anchor",
  show: function() {
    var anchor = this._getAnchor();
    anchor.innerHTML = [
      '<div id="cick-playlister-modal">',
      '  <div class="modal-container">',
      '    <div class="modal-content">',
      '      <span onclick="window.cickPlaylisterClient.hide()" class="modal-close">&times;</span>',
      '      <br />',
      '      <input id="playlist-spotify-input" type="text" placeholder="Spotify Playlist ID" />',
      '      <button onclick="window.cickPlaylisterClient.playlist(\'spotify\');">Fill</button>',
      '    </div>',
      '  </div>',
      '</div>'
    ].join("");
    anchor.style.display = "block";
  },
  hide: function() {
    var anchor = this._getAnchor();
    anchor.innerHTML = "";
    anchor.style.display = "none";
  },
  playlist: function(provider) {
    if (this.providers.hasOwnProperty(provider)) {
      if (this.providers[provider].extractors.hasOwnProperty("playlist")) {
        var inputId = [
            "playlist",
            provider,
            "input"
          ].join("-");
        var playlistId = this.providers[provider].extractors.playlist(
          document.getElementById(inputId).value
        );
        if (playlistId) {
          fetch([
            this._apiUrlBase,
            provider,
            "playlist",
            playlistId
          ].join("/"))
            .then(response => {
              if (response.ok) {
                return response.json();
              } else {
                throw new Error("Unexpected API response for Playlist ID");
              }
            })
            .then(data => {
              data.tracks.forEach(track => {
                this._addTrack(track);
              });
            })
            .catch(err => {
              console.log(err);
              this._reportError();
            })
          ;
        } else {
          this._reportFeedback("Invalid Playlist URL", inputId);
        }
      }
    } else {
      this._reportError();
    }
  },
  providers: {
    spotify: {
      extractors: {
        playlist: function(playlistUrl) {
          var match = playlistUrl.match(/^https:\/\/open\.spotify\.com\/playlist\/([^\?]+)(\?.*)?/);
          if (match) {
            return match[1];
          }
        },
      },
    }
  },
  _addTrack: function(track) {
    /*
      album: "Penny Penguin"
      artist: "Raffi, Good Lovelies"
      isNew: false
      isSingle: true
      track: "Penny Penguin"
    */
    // still need a way to ensure that tracks are not repeated, e.g. base64-encoding an identifier
    console.log(track);
  },
  _getAnchor: function() {
    return document.getElementById(this.anchorId);
  },
  _apiUrlBase: "http://localhost:8123",
  _reportError: function() {
    alert("Error in CICK Playlister. Please report an issue at https://github.com/captaincoordinates/cick-playlister/issues");
  },
  _reportFeedback: function(message, inputId) {
    console.log(message);
  }
};
