import { components } from "./generated/types";
import { FillRowResult, HandlerData, Provider } from "./types";
import { Spotify } from "./providers/spotify";
import { apiUrlBase } from "./constants";

type TrackInfo = components["schemas"]["TrackInfo"];
type TrackDisplay =  Omit<TrackInfo, "isSingle">;

class CickPlaylisterClient {

  public readonly anchorId: string = "cick-playlister-client-anchor";

  private readonly urlInputId: string = "cick-playlister-input";
  private readonly formSubmitId: string = "cick-playlister-submit";
  private readonly feedbackElementId: string = "cick-playlister-feedback";
  private readonly trackHashAttribute: string = "data-row-track-hash";
  private readonly trackHashEmptyValue: string = "empty";
  private readonly trackRowCounterAttribute: string = "data-row-counter";
  private readonly trackSingleValue: string = "Single";
  private readonly providers: {[index: string]: Provider} = {
    [Spotify.identifier]: new Spotify(),
  };
  private readonly boundEscapeKeyHandler: (event: KeyboardEvent) => void = this.escapeKeyHandler.bind(this);

  public show(): void {
    const providerIcons = Object.entries(this.providers).map(([identifier, provider]) => {
      return `
      <span class="modal-supported-provider">
        <img
          src="${apiUrlBase}/client/assets/${provider.icon}"
          alt="${identifier}"
          title="${identifier}"
          width="20px"
          height="20px"
          />
      </span>
      `
    });
    const anchor = this.anchor;
    anchor.innerHTML = `
      <div id="cick-playlister-modal">
        <div class="modal-container">
          <div class="modal-content">
            <div class="modal-title-container">
              <span class="modal-title">CICK Playlister</span>
              <span class="modal-close" onclick="window.cickPlaylisterClient.hide()">&times;</span>
            </div>
            <div class="modal-supported-providers-container">
              <span class="modal-supported-providers-title">
                Supported Services
              </span>
              ${providerIcons}
            </div>
            <form onsubmit="window.cickPlaylisterClient.processInput(); return false;">
              <div class="modal-input-container">
                <span class="modal-input-field-container">
                  <input
                    id="${this.urlInputId}"
                    class="modal-input-field"
                    type="text"
                    placeholder="Paste URL"
                    oninput="window.cickPlaylisterClient.urlChanged()"
                    />
                </span>
                <span>
                  <button
                    id="${this.formSubmitId}"
                    class="modal-submit-button"
                    type="submit" >
                    Fill
                  </button>
                </span>
              </div>
            </form>
            <div id="${this.feedbackElementId}"></div>
            <br />
            <hr />
            <br />
            <div id="cick-playlister-usage">
              Read more about the CICK Playlister <a href="https://github.com/captaincoordinates/cick-playlister/blob/main/USAGE.md" target="_blank">here</a>.
            </div>
          </div>
        </div>
      </div>
    `;
    anchor.style.display = "block";
    this.urlInput.focus();
    window.addEventListener("keydown", this.boundEscapeKeyHandler);
  }

  public hide(): void {
    const anchor = this.anchor;
    anchor.innerHTML = "";
    anchor.style.display = "none";
    window.removeEventListener("keydown", this.boundEscapeKeyHandler);
  }

  public processInput(): void {
    this.clearFeedback();
    var url = this.urlInput.value;
    if (url.length == 0) {
      this.reportFeedback("URL is empty");
      return;
    }
    var handlerData = this.getHandlerData(url);
    if (!handlerData) {
      this.reportFeedback("URL type is not currently supported");
      return;
    }
    this.disableUrlInput();
    this.updateHandleTypeDisplay(handlerData.provider, handlerData.type);
    handlerData.handle(apiUrlBase)
      .then(tracks => {
        this.classifyTableRows();
        const counts = {
          success: 0,
          noFreeRow: 0,
          duplicate: 0,
        };
        if (tracks.length > 0) {
          tracks.forEach(track => {
            track = {
              artist: track.artist,
              track: track.track,
              album: track.isSingle ? this.trackSingleValue : track.album,
              isNew: track.isNew,
              isSingle: track.isSingle,
            };
            switch(this.fillRow(track)) {
              case FillRowResult.Success:
                counts.success++;
                break;
              case FillRowResult.Duplicate:
                counts.duplicate++;
                break;
              case FillRowResult.NoFreeRow:
                counts.noFreeRow++;
                break;
            }
          });
          if (counts.success === tracks.length) {
            this.reportFeedback(`${this.trackCountString(counts.success)} filled successfully`);
          } else {
            const feedbackParts = [
              `${this.trackCountString(counts.success)} filled`,
            ];
            if (counts.noFreeRow > 0) {
              feedbackParts.push(`${this.trackCountString(counts.noFreeRow)} skipped as all rows are filled`);
            }
            if (counts.duplicate > 0) {
              feedbackParts.push(`${this.trackCountString(counts.duplicate)} skipped as duplicate`);
            }
            this.reportFeedback(feedbackParts.join(", "));
          }
        } else {
          this.reportFeedback("No tracks found");
        }
        this.enableUrlInput();
      })
      .catch(err => {
        console.log(err);
        this.reportError();
        this.clearFeedback();
        this.enableUrlInput();
      })
    ;
  }

  public urlChanged(): void {
    this.clearFeedback();
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

  private get formSubmitButton(): HTMLButtonElement {
    return document.getElementById(this.formSubmitId) as HTMLButtonElement;
  }

  private updateHandleTypeDisplay(provider: string, type: string): void {
    this.reportFeedback(`Processing ${provider}: ${type}...`);
  }

  private reportError(): void {
    alert("Error in CICK Playlister. Please report an issue at https://github.com/captaincoordinates/cick-playlister/issues");
  }

  private reportFeedback(message: string): void {
    (document.getElementById(this.feedbackElementId) as HTMLElement).innerText = message;
  }

  public clearFeedback(): void {
    this.reportFeedback("");
  }

  private disableUrlInput(): void {
    this.urlInput.disabled = true;
    this.formSubmitButton.disabled = true;
  }

  private enableUrlInput(): void {
    this.urlInput.disabled = false;
    this.formSubmitButton.disabled = false;
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

  private fillRow(track: TrackInfo): FillRowResult {
    const trackHash = this.trackUniqueIdentifier(track);
    const existingRowEls = document.querySelectorAll(`[${this.trackHashAttribute}="${trackHash}"]`);
    if (existingRowEls.length != 0) {
      console.log(`skipping duplicate in fillRow '${track.artist}': '${track.track}'`);
      return FillRowResult.Duplicate;
    }
    const emptyRowEls = document.querySelectorAll(`[${this.trackHashAttribute}="${this.trackHashEmptyValue}"]`);
    if (emptyRowEls.length === 0) {
      console.log(`no rows available for '${track.artist}': '${track.track}'`);
      return FillRowResult.NoFreeRow;
    }
    const nextRow = emptyRowEls[0];
    const rowCounter = parseInt(nextRow.getAttribute(this.trackRowCounterAttribute)!, 10);
    this.getArtistInput(rowCounter).value = track.artist;
    this.getTrackInput(rowCounter).value = track.track;
    this.getAlbumInput(rowCounter).value = track.album;
    this.getIsNewInput(rowCounter).checked = track.isNew;
    nextRow.setAttribute(this.trackHashAttribute, trackHash);
    return FillRowResult.Success;
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

  private escapeKeyHandler(event: KeyboardEvent): void {
    if (event.key === "Escape") {
      this.hide();
    }
  }

  private trackCountString(count: number): string {
    return `${count} track${count === 1 ? "" : "s"}`;
  }
}

(<any>window).cickPlaylisterClient = new CickPlaylisterClient();
