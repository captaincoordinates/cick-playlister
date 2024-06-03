import { components } from "../generated/types";

type TrackInfo = components["schemas"]["TrackInfo"];
type PlaylistInfo = components["schemas"]["PlaylistInfo"];

export class Common {

  public static getPlaylistFetcher(provider: string, playlistId: string): (apiUrl: string) => Promise<TrackInfo[]> {
    return async function(apiUrlBase: string) {
      return fetch(`${apiUrlBase}/${provider}/playlist/${playlistId}`)
        .then(async response => {
          if (response.ok) {
            return response.json()
              .then((data: PlaylistInfo) => {
                return data.tracks;
              });
          } else {
            throw new Error("Unexpected API response for Playlist ID");
          }
        })
      ;
    };
  }

  public static getTrackFetcher(provider: string, trackId: string): (apiUrl: string) => Promise<TrackInfo[]>{
    return async function(apiUrlBase: string) {
      return fetch(`${apiUrlBase}/${provider}/track/${trackId}`)
        .then(async response => {
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
  }
}
