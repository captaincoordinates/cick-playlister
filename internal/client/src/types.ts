import { components } from "./generated/types";

type TrackInfo = components["schemas"]["TrackInfo"];

export interface Provider {
  get icon(): string;
  getHandlerData(input: string): HandlerData | undefined;
}

export enum HandlerType {
  Playlist = "playlist",
  Album = "album",
  Track = "track",
}

export interface HandlerData {
  provider: string;
  type: HandlerType;
  handle: (apiUrl: string) => Promise<TrackInfo[]>;
}

export enum FillRowResult {
  Success,
  NoFreeRow,
  Duplicate,
}