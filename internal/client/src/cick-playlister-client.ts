import { components } from "./generated/types";
import { HandlerData, Provider } from "./types";
import { Spotify } from "./providers/spotify";
import { apiUrlBase } from "./constants";

type TrackInfo = components["schemas"]["TrackInfo"];
type TrackDisplay =  Omit<TrackInfo, "isSingle">;

class CickPlaylisterClient {

  public readonly anchorId: string = "cick-playlister-client-anchor";

  private readonly urlInputId: string = "cick-playlister-input";
  private readonly trackHashAttribute: string = "data-row-track-hash";
  private readonly trackHashEmptyValue: string = "empty";
  private readonly trackRowCounterAttribute: string = "data-row-counter";
  private readonly trackSingleValue: string = "Single";
  private readonly providers: {[index: string]: Provider} = {
    [Spotify.identifier]: new Spotify(),
  };

  public show(): void {
    const anchor = this.anchor;
    anchor.innerHTML = `
      <div id="cick-playlister-modal">
        <div class="modal-container">
          <div class="modal-content">
            <span onclick="window.cickPlaylisterClient.hide()" class="modal-close">&times;</span>
            <br />
            <form onsubmit="window.cickPlaylisterClient.processInput(); return false;">
              <input id="${this.urlInputId}" type="text" placeholder="Paste URL" />
              <button type="submit">Fill</button>
            </form>
          </div>
        </div>
      </div>
    `;
    anchor.style.display = "block";
    this.urlInput.focus();
  }

  public hide(): void {
    const anchor = this.anchor;
    anchor.innerHTML = "";
    anchor.style.display = "none";
  }

  public processInput(): void {
    var url = this.urlInput.value;
    if (url.length == 0) {
      this.reportFeedback("URL is empty", this.urlInputId);
      return;
    }
    var handlerData = this.getHandlerData(url);
    if (!handlerData) {
      this.reportFeedback("URL type is not currently supported", this.urlInputId);
      return;
    }
    this.updateHandleTypeDisplay(handlerData.provider, handlerData.type);
    handlerData.handle(apiUrlBase)
      .then(tracks => {
        this.classifyTableRows();
        tracks.forEach(track => {
          track = {
            artist: track.artist,
            track: track.track,
            album: track.isSingle ? this.trackSingleValue : track.album,
            isNew: track.isNew,
            isSingle: track.isSingle,
          };
          var trackHash = this.trackUniqueIdentifier(track);
          var trackRowEls = document.querySelectorAll(`[${this.trackHashAttribute}="${trackHash}"]`);
          switch (trackRowEls.length) {
            case 0:
              this.fillRow(track);
              break;
            case 1:
              console.debug("skipping " + trackHash);
              break;
            default:
              console.log("unexpectedly found multiple rows with " + trackHash);
              this.reportError();
          }
        });
      })
      .catch(err => {
        console.log(err);
        this.reportError();
      })
    ;
  }

  private getHandlerData(url: string): HandlerData | undefined {
    for (var key in this.providers) {
      var handlerData = this.providers[key].getHandlerData(url);
      if (!!handlerData) {
        return handlerData;
      }
    };
  }

  private get anchor(): HTMLElement {
    return document.getElementById(this.anchorId)!;
  }

  private get urlInput(): HTMLInputElement {
    return document.getElementById(this.urlInputId) as HTMLInputElement;
  }

  private updateHandleTypeDisplay(provider: string, type: string): void {
    console.log(`inform user of ${provider} ${type} handling`);
  }

  private reportError(): void {
    alert("Error in CICK Playlister. Please report an issue at https://github.com/captaincoordinates/cick-playlister/issues");
  }

  private reportFeedback(message: string, inputId: string): void {
    console.log(message);
  }

  private classifyTableRows(): void {
    Array.from(document.getElementById("station-playlist-tracks-table")!.getElementsByTagName("tr")).forEach(rowEl => {
      let artistInput: HTMLInputElement | null = null;
      let trackInput: HTMLInputElement | null = null;
      let albumInput: HTMLInputElement | null = null;
      let isNewInput: HTMLInputElement | null = null;
      let rowCounter: number = -1;
      const artistInputIdRegex = /^edit-tracks-(\d+)-artist$/
      Array.from(rowEl.getElementsByTagName("input")).forEach(inputEl => {
        const match = inputEl.id.match(artistInputIdRegex);
        if (!!match) {
          rowCounter = parseInt(match[1], 10);
          artistInput = this.getArtistInput(rowCounter);
          trackInput = this.getTrackInput(rowCounter);
          albumInput = this.getAlbumInput(rowCounter);
          isNewInput = this.getIsNewInput(rowCounter);
        }
      });
      if (artistInput && trackInput && albumInput && isNewInput) {
        let hashValue = this.trackHashEmptyValue;
        if (
          (artistInput as HTMLInputElement).value &&
          (trackInput as HTMLInputElement).value &&
          (albumInput as HTMLInputElement).value
        ) {
          hashValue = this.trackUniqueIdentifier({
            artist: (artistInput as HTMLInputElement).value,
            track: (trackInput as HTMLInputElement).value,
            album: (albumInput as HTMLInputElement).value,
            isNew: (isNewInput as HTMLInputElement).checked,
          });
        }
        rowEl.setAttribute(this.trackHashAttribute, hashValue);
        rowEl.setAttribute(this.trackRowCounterAttribute, rowCounter.toString());
      }
    });
  }

  private fillRow(track: TrackInfo): void {
    const trackHash = this.trackUniqueIdentifier(track);
    const existingRowEls = document.querySelectorAll(`[${this.trackHashAttribute}="${trackHash}"]`);
    if (existingRowEls.length != 0) {
      console.log(`track already added, skipping ${track.artist}: ${track.track} (${track.album})`);
      return;
    }
    const emptyRowEls = document.querySelectorAll(`[${this.trackHashAttribute}="${this.trackHashEmptyValue}"]`);
    if (emptyRowEls.length === 0) {
      this.reportFeedback("no rows available", this.urlInputId);
      return;
    }
    const nextRow = emptyRowEls[0];
    const rowCounter = parseInt(nextRow.getAttribute(this.trackRowCounterAttribute)!, 10);
    this.getArtistInput(rowCounter).value = track.artist;
    this.getTrackInput(rowCounter).value = track.track;
    this.getAlbumInput(rowCounter).value = track.album;
    this.getIsNewInput(rowCounter).checked = track.isNew;
    nextRow.setAttribute(this.trackHashAttribute, trackHash);
  }

  private getArtistInput(rowCounter: number): HTMLInputElement {
    return document.getElementById("edit-tracks-" + rowCounter + "-artist") as HTMLInputElement;
  }

  private getTrackInput(rowCounter: number): HTMLInputElement {
    return document.getElementById("edit-tracks-" + rowCounter + "-title") as HTMLInputElement;
  }

  private getAlbumInput(rowCounter: number): HTMLInputElement {
    return document.getElementById("edit-tracks-" + rowCounter + "-album") as HTMLInputElement;
  }
  
  private getIsNewInput(rowCounter: number): HTMLInputElement {
    return document.getElementById("edit-tracks-" + rowCounter + "-newtrack") as HTMLInputElement;
  }

  private trackUniqueIdentifier(track: TrackDisplay): string {
    console.debug(`hashing with ${track.artist}: ${track.track} (${track.album}) ${track.isNew ? "[new]" : ""}`);
    return btoa(
      encodeURIComponent(
        JSON.stringify({
          artist: track.artist,
          track: track.track,
          album: track.album,
          isNew: track.isNew,
        })
      )
    ).replace(/=/g, "_"); // need to check if this is required
  }
}

(<any>window).cickPlaylisterClient = new CickPlaylisterClient();
