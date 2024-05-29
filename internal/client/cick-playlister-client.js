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
              this._classifyTableRows();
              data.tracks.forEach(track => {
                var trackHash = this._hashEntry(
                  track.artist,
                  track.track,
                  track.album
                );
                var trackRowEls = document.querySelectorAll([
                  "[",
                  this._trackHashAttribute,
                  "=\"",
                  trackHash,
                  "\"]"
                ].join(""));
                switch (trackRowEls.length) {
                  case 0:
                    this._fillRow(track);
                    break;
                  case 1:
                    console.debug("skipping " + trackHash);
                    break;
                  default:
                    console.log("unexpectedly found multiple rows with " + trackHash);
                    this._reportError();
                }
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
  _getAnchor: function() {
    return document.getElementById(this.anchorId);
  },
  _apiUrlBase: "http://localhost:8123",
  _reportError: function() {
    alert("Error in CICK Playlister. Please report an issue at https://github.com/captaincoordinates/cick-playlister/issues");
  },
  _reportFeedback: function(message, inputId) {
    console.log(message);
  },
  _trackHashAttribute: "data-row-track-hash",
  _trackHashEmptyValue: "empty",
  _trackRowCounterAttribute: "data-row-counter",
  _trackSingleValue: "Single",
  _classifyTableRows: function() {
    Array.from(document.getElementById("station-playlist-tracks-table").getElementsByTagName("tr")).forEach(rowEl => {
      var artistInput = null;
      var trackInput = null;
      var albumInput = null;
      var rowCounter = null;
      var artistInputIdRegex = /^edit-tracks-(\d+)-artist$/
      Array.from(rowEl.getElementsByTagName("input")).forEach(inputEl => {
        match = inputEl.id.match(artistInputIdRegex);
        if (!!match) {
          rowCounter = match[1];
          artistInput = this._getArtistInput(rowCounter);
          trackInput = this._getTrackInput(rowCounter);
          albumInput = this._getAlbumInput(rowCounter);
        }
      });
      if (artistInput && trackInput && albumInput) {
        var hashValue = this._trackHashEmptyValue;
        if (artistInput.value && trackInput.value && albumInput.value) {
          hashValue = this._hashEntry(
            artistInput.value,
            trackInput.value,
            albumInput.value
          );
        }
        rowEl.setAttribute(this._trackHashAttribute, hashValue);
        rowEl.setAttribute(this._trackRowCounterAttribute, rowCounter);
      }
    });
  },
  _fillRow: function(track) {
    var emptyRowEls = document.querySelectorAll([
      "[",
      this._trackHashAttribute,
      "=\"",
      this._trackHashEmptyValue,
      "\"]"
    ].join(""));
    if (emptyRowEls.length === 0) {
      this._reportFeedback("no rows available");
      return;
    }
    var nextRow = emptyRowEls[0];
    var rowCounter = parseInt(nextRow.getAttribute(this._trackRowCounterAttribute), 10);
    this._getArtistInput(rowCounter).value = track.artist;
    this._getTrackInput(rowCounter).value = track.track;
    this._getAlbumInput(rowCounter).value = (track.isSingle ? this._trackSingleValue : track.album);
    this._getIsNewInput(rowCounter).checked = track.isNew;
    nextRow.setAttribute(this._trackHashAttribute, this._hashEntry(
      track.artist,
      track.track,
      track.album
    ));
  },
  _getArtistInput(rowCounter) {
    return document.getElementById("edit-tracks-" + rowCounter + "-artist");
  },
  _getTrackInput(rowCounter) {
    return document.getElementById("edit-tracks-" + rowCounter + "-title");
  },
  _getAlbumInput(rowCounter) {
    return document.getElementById("edit-tracks-" + rowCounter + "-album");
  },
  _getIsNewInput(rowCounter) {
    return document.getElementById("edit-tracks-" + rowCounter + "-newtrack");
  },
  _hashEntry: function(artist, title, album) {
    return btoa(
      encodeURIComponent(
        JSON.stringify({
          artist: artist,
          title: title,
          album: album
        })
      )
    );
  },
};
