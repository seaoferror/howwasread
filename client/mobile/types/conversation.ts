export interface CreateOnlineConversationRequest {
  novel?: string;
  shortStory?: string;
  poem?: string;
  play?: string;
  film?: string;
  writtenBy?: string;
  rule?: string;
  capacity: number;
  time: string;
  length: string;
}

export interface OnlineConversationDetailResponse {
  id: string;
  novel?: string;
  shortStory?: string;
  poem?: string;
  play?: string;
  film?: string;
  writtenBy?: string;
  rule?: string;
  capacity: number;
  time: string;
  length: string;
  canEnter: boolean;
  isRegistrant: boolean;
  moderatorIds: string[];
}

export interface BanParticipantRequest {
  conversationId: string;
  banId: string;
}

export interface ConversationSignalResponse {
  fromIds: string[];
  signal?: PeerSignal;
}

type PeerSignal =
  | { type: "offer" | "answer"; sdp: string }
  | { type: "candidate"; candidate: RTCIceCandidate }
  | { type: "leave" }
  | { type: "ban" }
  | { type: "mute" }
  | { type: "name-offer" | "name-answer"; name: string };

interface RTCIceCandidate {
  candidate?: string;
  sdpMLineIndex?: number | null;
  sdpMid?: string | null;
}

export type SeatCoordinate = { left: number; top: number };

export interface SeatAssignment {
  id?: string;
  name?: string;
  mute?: boolean;
  left: number;
  top: number;
}

export interface OfflineConversationMapResponse {
  id: string;
  writtenBy: string;
  lat: number;
  lng: number;
}

export interface OnlineConversationFeedResponse {
  id: string;
  novel: string;
  poem: string;
  shortStory: string;
  play: string;
  film: string;
  writtenBy: string;
  time: string;
}

export interface OfflineConversationSearchResponse {
  id: string;
  novel: string;
  poem: string;
  shortStory: string;
  play: string;
  film: string;
  writtenBy: string;
  time: string;
  lat: number;
  lng: number;
}

export interface CreateOfflineConversationRequest {
  novel?: string;
  shortStory?: string;
  poem?: string;
  play?: string;
  film?: string;
  writtenBy: string;
  rule?: string;
  time: string;
  length: number;
  mapsLink: string;
  location: string;
  city?: string;
  lat: number;
  lng: number;
  h3Res5: string;
  h3Res7: string;
}

export interface OfflineConversationDetailResponse {
  novel?: string | null;
  poem?: string | null;
  shortStory?: string | null;
  play?: string | null;
  film?: string | null;
  writtenBy: string;
  rule?: string | null;
  time: string;
  length: number;
  mapsLink: string;
  location: string;
  isModerator: boolean;
  isParticipant: boolean;
  numberOfParticipants: number;
  moderatorIds: string[];
}

export interface GetTurnResponse {
  uris: string[]
  username: string
  credential: string
}

