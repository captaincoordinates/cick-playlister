import { HandlerData, HandlerType, Provider } from '../types';
import { Common } from './common';

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
            handle: Common.getPlaylistFetcher(Spotify.identifier, matchResult[2]),
          };
        case HandlerType.Track:
          return {
            provider: Spotify.identifier,
            type: HandlerType.Track,
            handle: Common.getTrackFetcher(Spotify.identifier, matchResult[2])
          };
        default:
          return undefined;
      }
    } else {
      return undefined;
    }
  }
}
