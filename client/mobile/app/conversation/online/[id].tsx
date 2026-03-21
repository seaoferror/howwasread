import { Platform, Pressable, StyleSheet, Text, View } from "react-native";
import { router, useFocusEffect, useLocalSearchParams } from "expo-router";
import { colors, SEAT_COORDINATES, SEAT_FILL_ORDER } from "@/constants";
import { useEffect, useRef, useState } from "react";
import {
  ConversationSignalResponse,
  SeatAssignment,
} from "@/types/conversation";
import {
  RTCPeerConnection,
  RTCIceCandidate,
  RTCSessionDescription,
  mediaDevices,
  MediaStream,
} from "react-native-webrtc";
import { baseUrl, localDevId } from "@/api/axios";
import { getSecureStore } from "@/util/secureStore";
import { SafeAreaView } from "react-native-safe-area-context";
import OnlineConversationRoomHeader from "@/components/conversation/OnlineConversationRoomHeader";
import { Feather, Ionicons } from "@expo/vector-icons";
import { useAuth } from "@/hooks/useAuth";
import { useGetOnlineConversationPreAssignedIds } from "@/hooks/useConversation";

//Prevent ts compiler kept trying to read dom definition of WebSocket written by Microsoft, which don't have header options
declare const WebSocket: {
  prototype: WebSocket;
  new (
    url: string,
    protocols?: string | string[] | null,
    options?: {
      headers?: { [header: string]: string };
      [key: string]: any;
    } | null,
  ): WebSocket;
};

export default function OnlineConversationRoomScreen() {
  const { id: authId } = useAuth();
  const fallbackId =
    (Platform.OS === "ios" ? localDevId.ios : localDevId.android) ??
    "local-user";
  const myId = authId ?? fallbackId;

  const {
    id: conversationId,
    novel,
    shortStory,
    poem,
    play,
    film,
    by,
    rule,
    capacity,
    when,
    length,
    isModerator,
  } = useLocalSearchParams();

  const [mute, setMute] = useState(true);
  const [roomMemberIds, setRoomMemberIds] = useState<string[]>([myId]);
  const [seatAssignments, setSeatAssignments] = useState<SeatAssignment[]>([]);
  const ws = useRef<WebSocket>(null);
  const peers = useRef<Record<string, RTCPeerConnection>>({});
  const localAudio = useRef<MediaStream>(null);
  const remoteAudios = useRef<Record<string, MediaStream>>({});

  const joinConversation = async () => {
    console.log("try to connect ws");
    ws.current = new WebSocket(
      `ws://${baseUrl.ios}:8080/online/conversation/join?id=${conversationId}`,
      undefined,
      {
        headers: {
          Authorization: `Bearer ${await getSecureStore("accessToken")}`,
          "X-User-Id": `${Platform.OS === "ios" ? localDevId.ios : localDevId.android}`,
        },
      },
    );
    console.log(
      `ws://${baseUrl.ios}:8080/online/conversation/join?id=${conversationId}`,
    );
    localAudio.current = await mediaDevices.getUserMedia({
      audio: true,
      video: false,
    });

    ws.current.onmessage = async (event) => {
      console.log("get message");

      const data: ConversationSignalResponse = JSON.parse(event.data);
      if (!data.signal) {
        const unique = [...new Set(data.fromIds)];
        if (unique.length >= Number(capacity)) {
          router.replace("/conversation/online");
          return;
        }
        setRoomMemberIds([...roomMemberIds, ...unique]);
        for (const fromId of unique) {
          peers.current[fromId] = new RTCPeerConnection({
            iceServers: [{ urls: "stun:stun.l.google.com:19302" }],
          });

          localAudio.current?.getTracks().forEach((track) => {
            if (localAudio.current) {
              peers.current[fromId].addTrack(track, localAudio.current);
            }
          });

          peers.current[fromId].addEventListener(
            "icecandidate",
            (event: any) => {
              if (event.candidate) {
                ws.current?.send(
                  JSON.stringify({
                    toIds: [fromId],
                    signal: { type: "candidate", candidate: event.candidate },
                  }),
                );
              }
            },
          );

          peers.current[fromId].addEventListener("track", (event: any) => {
            if (!remoteAudios.current[fromId]) {
              remoteAudios.current[fromId] = new MediaStream();
            }
            if (event.track) {
              remoteAudios.current[fromId].addTrack(event.track);
            }
          });

          const offer = await peers.current[fromId].createOffer({
            offerToReceiveAudio: true,
            offerToReceiveVideo: false,
            voiceActivityDetection: true,
          });
          await peers.current[fromId].setLocalDescription(offer);
          ws.current?.send(
            JSON.stringify({
              toIds: [fromId],
              signal: peers.current[fromId].localDescription,
            }),
          );
        }
        //TODO: may be setPending false and make clean up possible in this point
        return;
      }
      const fromId = data.fromIds[0];
      if (data.signal.type === "offer") {
        const offer = new RTCSessionDescription(data.signal);
        await peers.current[fromId].setRemoteDescription(offer);
        const answer = await peers.current[fromId].createAnswer();
        await peers.current[fromId].setLocalDescription(answer);
        ws.current?.send(
          JSON.stringify({
            toIds: [fromId],
            signal: peers.current[fromId].localDescription,
          }),
        );
        return;
      }
      if (data.signal.type === "answer") {
        const offerDescription = new RTCSessionDescription(data.signal);
        await peers.current[fromId].setRemoteDescription(offerDescription);
        return;
      }
      if (data.signal.type === "candidate") {
        const iceCandidate = new RTCIceCandidate(data.signal.candidate);
        await peers.current[fromId].addIceCandidate(iceCandidate);
        return;
      }
      if (data.signal.type === "leave") {
        peers.current[fromId].close();
        delete peers.current[fromId];
        remoteAudios.current[fromId].release();
        delete remoteAudios.current[fromId];
        setRoomMemberIds(roomMemberIds.filter((x) => x !== fromId));
      }
    };
  };

  const toggleAudio = () => {
    if (localAudio.current) {
      const audioTrack = localAudio.current.getAudioTracks()[0];
      audioTrack.enabled = !audioTrack.enabled;
      setMute(!mute);
    }
  };

  useFocusEffect(() => {
    joinConversation();

    return () => {
      ws.current?.send(
        JSON.stringify({
          toIds: Object.keys(peers.current),
          signal: { type: "leave" },
        }),
      );
      for (const peer of Object.values(peers.current)) {
        peer.close();
      }
      peers.current = {};
      localAudio.current?.release();
      localAudio.current = null;
      for (const remoteAudio of Object.values(remoteAudios.current)) {
        remoteAudio.release();
      }
      remoteAudios.current = {};
      ws.current?.close();
      ws.current = null;
    };
  });

  useEffect(() => {
    const coordinates = SEAT_COORDINATES[Number(capacity)];
    const fillOrder = SEAT_FILL_ORDER[Number(capacity)];

    const unique = [...new Set(roomMemberIds)];
    const occupantBySeat: Record<number, string> = {};
    unique.forEach((id, i) => {
      const seatIndex = fillOrder[i];
      occupantBySeat[seatIndex] = id;
    });

    const nextAssignments: SeatAssignment[] = coordinates.map(
      (coordinate, i) => ({
        id: occupantBySeat[i],
        ...coordinate,
      }),
    );

    setSeatAssignments(nextAssignments);
  }, [capacity, roomMemberIds]);

  return (
    <SafeAreaView style={styles.container}>
      <OnlineConversationRoomHeader
        novel={String(novel)}
        shortStory={String(shortStory)}
        poem={String(poem)}
        play={String(play)}
        film={String(film)}
        by={String(by)}
        rule={String(rule)}
        when={String(when)}
        length={String(length)}
      />
      <View style={styles.participantContainer}>
        <View style={styles.participantArea}>
          {seatAssignments.map((seat, index) => (
            <View
              key={index}
              style={[
                styles.participantSeat,
                { left: `${seat.left}%`, top: `${seat.top}%` },
              ]}
            >
              {seat.id ? (
                <Ionicons
                  name="person-circle-outline"
                  size={24}
                  color="black"
                />
              ) : (
                <Feather name="circle" size={24} color="black" />
              )}
              <Text style={styles.participantId}>{seat.id ?? ""}</Text>
            </View>
          ))}
        </View>
        <View style={styles.controls}></View>
      </View>
      <View style={styles.chatContainer}></View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.SAND_110,
  },
  content: {
    flex: 1,
    padding: 16,
  },
  title: {
    fontSize: 18,
    color: colors.BLACK,
    fontWeight: 600,
    marginVertical: 8,
  },
  description: {
    fontSize: 16,
    color: colors.BLACK,
    marginBottom: 14,
  },
  participantContainer: {
    flex: 1,
  },
  participantArea: {
    flex: 1,
    position: "relative",
  },
  participantSeat: {
    position: "absolute",
    width: 96,
    alignItems: "center",
    transform: [{ translateX: -48 }, { translateY: -12 }],
  },
  participantId: {
    marginTop: 6,
    fontSize: 12,
    color: colors.BLACK,
  },
  controls: {
    flexDirection: "row",
    justifyContent: "center",
    alignItems: "center",
    paddingVertical: 24,
    paddingHorizontal: 16,
    gap: 12,
  },
  muteButton: {
    backgroundColor: colors.SAND_110,
    borderWidth: 1,
    borderColor: colors.BLACK,
    borderRadius: 8,
    paddingVertical: 12,
    paddingHorizontal: 24,
  },
  muteText: {
    color: colors.BLACK,
    fontSize: 16,
    fontWeight: "600",
  },
  chatContainer: {},
});
