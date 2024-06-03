import { components } from "./generated/types";

type TrackInfo = components["schemas"]["TrackInfo"];

export interface Provider {
  getHandlerData(input: string): HandlerData | undefined;
}

export enum HandlerType {
  Playlist = "playlist",
  Track = "track",
}

export interface HandlerData {
  provider: string;
  type: HandlerType;
  handle: (apiUrl: string) => Promise<TrackInfo[]>;
}