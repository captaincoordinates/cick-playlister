import { components } from "../generated/types";
import { HandlerData, HandlerType, Provider } from '../types';

type TrackInfo = components["schemas"]["TrackInfo"];
type PlaylistInfo = components["schemas"]["PlaylistInfo"];

export class Spotify implements Provider {

  public static readonly identifier: string = "spotify";
  private readonly protocolAndDomain: string = "https://open.spotify.com";

  public getHandlerData(input: string): HandlerData | undefined {
    const matchResult = input.match(new RegExp("^" + this.protocolAndDomain + "/(playlist|track)/([^\?]+)"))
    if (!!matchResult) {
      switch (matchResult[1]) {
        case HandlerType.Playlist:
          return {
            provider: Spotify.identifier,
            type: HandlerType.Playlist,
            handle: function(provider, playlistId) {
              return async function(apiUrlBase) {
                return fetch(apiUrlBase + "/" + provider + "/playlist/" + playlistId)
                  .then(response => {
                    if (response.ok) {
                      return response.json()
                        .then((data: PlaylistInfo) => {
                          return data.tracks;
                        }
                      );
                    } else {
                      throw new Error("Unexpected API response for Playlist ID");
                    }
                  })
                ;
              };
            }(Spotify.identifier, matchResult[2])
          };
        case HandlerType.Track:
          return {
            provider: Spotify.identifier,
            type: HandlerType.Track,
            handle: function(provider, trackId) {
              return async function(apiUrlBase) {
                return fetch(apiUrlBase + "/" + provider + "/track/" + trackId)
                  .then(response => {
                    if (response.ok) {
                      return response.json()
                        .then((data: TrackInfo) => {
                          return [data];
                        });
                    } else {
                      throw new Error("Unexpected API response for Track ID");
                    }
                  })
                ;
              };
            }(Spotify.identifier, matchResult[2])
          };
        default:
          return undefined;
      }
    } else {
      return undefined;
    }
  }
}
