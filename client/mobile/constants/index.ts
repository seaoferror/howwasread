import { SeatCoordinate } from "@/types/conversation";

export const colors = {
  WHITE: "#FFF",
  BLACK: "#000",
  GRAY_50: "#FCFCFC",
  GRAY_100: "#F5F5F5",
  GRAY_200: "#EEEEEE",
  GRAY_300: "#E0E0E0",
  GRAY_400: "#BDBDBD",
  GRAY_500: "#9E9E9E",
  GRAY_700: "#616161",
  GRAY_900: "#1D1E21",
  RED_100: "#FFDFDF",
  RED_500: "#FF5F5F",
  GREEN_100: "#388E3C",
  GREEN_200: "#2E7D32",
  GREEN_300: "#388E3C",
  GREEN_600: "#1B5E20",
  SAND_100: "#FFF8E1",
  SAND_110: "#FFF7E6",
  SAND_120: "#FFF6DA",
  SAND_150: "#FFF3CF",
  SAND_200: "#FFECB3",
  SAND_300: "#FFE082",
  ORANGE_50: "#FFF3E0",
  ORANGE_100: "#FFE0B2",
  ORANGE_150: "#FFD8A8",
  ORANGE_200: "#FFCC80",
  BLUE_500: "#2196F3"
};

export const queryKey = {
  AUTH: "auth",
  CONVERSATION: "conversation",
  GET_MY_ID: "getMyId",
  GET_ONLINE_CONVERSATIONS: "getOnlineConversations",
  GET_ONLINE_CONVERSATION_PRE_ASSIGNED_IDS:
    "getOnlineConversationPreAssignedIds",
  PROFILE: "profile",
  GET_MY_PROFILE: "getMyProfile",
  GET_RECENT_MESSAGES: "getRecentMessages",
  CHAT: "chat",
  GET_CHAT_ROOM_INFO: "getChatRoomInfo",
  GET_PROFILE: "getProfile",
  FIND_MESSAGES_BY_ROOM_ID: "findMessagesByRoomId",
  GET_SIGNED_URL: "getSignedURL",
  MAP_OFFLINE_CONVERSATIONS: "mapOfflineConversations",
  GET_OFFLINE_CONVERSATION_DETAIL: "getOfflineConversationDetail",
  CHECK_BLOCK: "checkBlock",
  GET_CHAT_PARTICIPANT_IDS: "getChatParticipantIds",
  REPORT_USER: "reportUser",
  BLOCKED_CONVERSATIONS: "blockedConversations",
  SEARCH_OFFLINE_CONVERSATIONS: "searchOfflineConversations"
};

export const time = {
  TEN_MINUTES: 10 * 60 * 1000,
};

export const SEAT_COORDINATES: Record<number, SeatCoordinate[]> = {
  2: [
    { left: 34, top: 50 },
    { left: 66, top: 50 },
  ],
  3: [
    { left: 50, top: 16 },
    { left: 24, top: 68 },
    { left: 76, top: 68 },
  ],
  4: [
    { left: 50, top: 14 },
    { left: 78, top: 50 },
    { left: 50, top: 86 },
    { left: 22, top: 50 },
  ],
  5: [
    { left: 50, top: 12 },
    { left: 74, top: 34 },
    { left: 64, top: 76 },
    { left: 36, top: 76 },
    { left: 26, top: 34 },
  ],
  6: [
    { left: 50, top: 12 },
    { left: 74, top: 30 },
    { left: 74, top: 70 },
    { left: 50, top: 88 },
    { left: 26, top: 70 },
    { left: 26, top: 30 },
  ],
};

export const SEAT_FILL_ORDER: Record<number, number[]> = {
  2: [0, 1],
  3: [0, 1, 2],
  4: [0, 2, 1, 3],
  5: [0, 3, 1, 4, 2],
  6: [0, 3, 1, 5, 2, 4],
};
