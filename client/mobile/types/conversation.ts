export interface CreateOnlineConversationRequest {
  novel?: string;
  shortStory?: string;
  poem?: string;
  play?: string;
  film?: string;
  by?: string;
  rule?: string;
  capacity: number;
  when: string;
  length: string;
}

export interface OnlineConversationFeedResponse {
  id: string;
  novel?: string;
  shortStory?: string;
  poem?: string;
  play?: string;
  film?: string;
  by?: string;
  rule?: string;
  capacity: number;
  when: string;
  length: string;
  ongoing: boolean;
  isModerator: boolean;
  isRegistrant: boolean;
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
