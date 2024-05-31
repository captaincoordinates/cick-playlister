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
      '      <form onsubmit="window.cickPlaylisterClient.processInput(); return false;">',
      '        <input id="cick-playlister-input" type="text" placeholder="Paste URL" />',
      '        <button type="submit">Lookup</button>',
      '      </form>',
      '    </div>',
      '  </div>',
      '</div>'
    ].join("");
    anchor.style.display = "block";
    this._getUrlInput().focus();
  },
  hide: function() {
    var anchor = this._getAnchor();
    anchor.innerHTML = "";
    anchor.style.display = "none";
  },
  addTracks: function() {
    console.log("addTracks executed");
  },
  processInput: function() {
    var url = this._getUrlInput().value;
    if (url.length == 0) {
      this._reportFeedback("URL is empty", "cick-playlister-input");
      return;
    }
    var handlerData = this._getHandlerData(url)
    if (!handlerData) {
      this._reportFeedback("URL type is not currently supported", "cick-playlister-input");
      return;
    }
    this._updateHandleTypeDisplay(handlerData.provider, handlerData.type);
    handlerData.handle(this._apiUrlBase)
      .then(tracks => {
        this._classifyTableRows();
        tracks.forEach(track => {
          track = {
            artist: track.artist,
            track: track.track,
            album: track.isSingle ? this._trackSingleValue : track.album,
            isNew: track.isNew,
          };
          var trackHash = this._trackUniqueIdentifier(track);
          var trackRowEls = document.querySelectorAll("[" + this._trackHashAttribute + "=\"" + trackHash + "\"]");
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
  },
  _getHandlerData: function(url) {
    for (var key in this._providers) {
      var handlerData = this._providers[key].getHandlerData(url);
      if (!!handlerData) {
        return handlerData;
      }
    };
  },
  _providers: {
    spotify: {
      _identifier: "spotify",
      _protocolAndDomain: "https://open.spotify.com",
      getHandlerData: function(input) {
        var matchResult = input.match(new RegExp("^" + this._protocolAndDomain + "/(playlist|track)/([^\?]+)"))
        if (!!matchResult) {
          switch (matchResult[1]) {
            case "playlist":
              return {
                provider: this._identifier,
                type: matchResult[1],
                handle: function(provider, playlistId) {
                  return async function(apiUrlBase) {
                    return fetch(apiUrlBase + "/" + provider + "/playlist/" + playlistId)
                      .then(response => {
                        if (response.ok) {
                          return response.json().then(data => {
                            return data.tracks;
                          });
                        } else {
                          throw new Error("Unexpected API response for Playlist ID");
                        }
                      })
                    ;
                  };
                }(this._identifier, matchResult[2])
              };
            case "track":
              return {
                provider: this._identifier,
                type: matchResult[1],
                handle: function(provider, trackId) {
                  return async function(apiUrlBase) {
                    return fetch(apiUrlBase + "/" + provider + "/track/" + trackId)
                      .then(response => {
                        if (response.ok) {
                          return response.json()
                            .then(data => {
                              return [data];
                            });
                        } else {
                          throw new Error("Unexpected API response for Track ID");
                        }
                      })
                    ;
                  };
                }(this._identifier, matchResult[2])
              };
            default:
              return undefined;
          }
        } else {
          return undefined;
        }
      },
      processors: {
        playlist: function(playlistUrl) {
          var match = playlistUrl.match(/^https:\/\/open\.spotify\.com\/playlist\/([^\?]+)(\?.*)?/);
          if (match) {
            return match[1];
          }
        },
      },
    }
  },
  _updateHandleTypeDisplay: function(provider, type) {
    console.log("inform user of " + provider + " " + type + " handling");
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
      var isNewInput = null;
      var rowCounter = null;
      var artistInputIdRegex = /^edit-tracks-(\d+)-artist$/
      Array.from(rowEl.getElementsByTagName("input")).forEach(inputEl => {
        match = inputEl.id.match(artistInputIdRegex);
        if (!!match) {
          rowCounter = match[1];
          artistInput = this._getArtistInput(rowCounter);
          trackInput = this._getTrackInput(rowCounter);
          albumInput = this._getAlbumInput(rowCounter);
          isNewInput = this._getIsNewInput(rowCounter);
        }
      });
      if (artistInput && trackInput && albumInput) {
        var hashValue = this._trackHashEmptyValue;
        if (artistInput.value && trackInput.value && albumInput.value) {
          hashValue = this._trackUniqueIdentifier({
            artist: artistInput.value,
            track: trackInput.value,
            album: albumInput.value,
            isNew: isNewInput.checked,
          });
        }
        rowEl.setAttribute(this._trackHashAttribute, hashValue);
        rowEl.setAttribute(this._trackRowCounterAttribute, rowCounter);
      }
    });
  },
  _fillRow: function(track) {
    var trackHash = this._trackUniqueIdentifier(track);
    var existingRowEls = document.querySelectorAll("[" + this._trackHashAttribute + "=\"" + trackHash + "\"]");
    if (existingRowEls.length != 0) {
      console.log("track already added, skipping " + trackHash.artist + ": " + trackHash.track + " (" + trackHash.album + ")");
      return;
    }
    var emptyRowEls = document.querySelectorAll("[" + this._trackHashAttribute + "=\"" + this._trackHashEmptyValue + "\"]");
    if (emptyRowEls.length === 0) {
      this._reportFeedback("no rows available");
      return;
    }
    var nextRow = emptyRowEls[0];
    var rowCounter = parseInt(nextRow.getAttribute(this._trackRowCounterAttribute), 10);
    this._getArtistInput(rowCounter).value = track.artist;
    this._getTrackInput(rowCounter).value = track.track;
    this._getAlbumInput(rowCounter).value = track.album;
    this._getIsNewInput(rowCounter).checked = track.isNew;
    nextRow.setAttribute(this._trackHashAttribute, trackHash);
  },
  _getUrlInput: function() {
    return document.getElementById("cick-playlister-input");
  },
  _getArtistInput: function(rowCounter) {
    return document.getElementById("edit-tracks-" + rowCounter + "-artist");
  },
  _getTrackInput: function(rowCounter) {
    return document.getElementById("edit-tracks-" + rowCounter + "-title");
  },
  _getAlbumInput: function(rowCounter) {
    return document.getElementById("edit-tracks-" + rowCounter + "-album");
  },
  _getIsNewInput: function(rowCounter) {
    return document.getElementById("edit-tracks-" + rowCounter + "-newtrack");
  },
  _trackUniqueIdentifier: function(track) {
    console.debug("hashing with '" + track.artist + "', '" + track.title + "', '" + track.album + "'");
    return btoa(
      encodeURIComponent(
        JSON.stringify(track)
      )
    ).replace(/=/g, "_"); // need to check if this is required
  },
};
