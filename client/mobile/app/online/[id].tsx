import { Platform, StyleSheet, Text, View } from "react-native";
import { router, useFocusEffect, useLocalSearchParams } from "expo-router";
import {
  colors,
  queryKey,
  SEAT_COORDINATES,
  SEAT_FILL_ORDER,
} from "@/constants";
import { useRef, useState } from "react";
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
import { getSecureAsync } from "@/util/storage";
import { SafeAreaView } from "react-native-safe-area-context";
import OnlineConversationRoomHeader from "@/components/conversation/OnlineConversationRoomHeader";
import { Feather, Ionicons } from "@expo/vector-icons";
import { useMyProfile } from "@/hooks/useMyProfile";
import { useActionSheet } from "@expo/react-native-action-sheet";
import { useSendLike } from "@/hooks/useChat";
import Toast from "react-native-toast-message";
import CustomButton from "@/components/CustomButton";
import queryClient from "@/api/queryClient";
import { useBanParticipant } from "@/hooks/useConversation";

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

export default function OnlineConversationScreen() {
  const { profile } = useMyProfile();
  const myId = Platform.OS === "ios" ? localDevId.ios : localDevId.android;

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
  const { showActionSheetWithOptions } = useActionSheet();
  const coordinates = SEAT_COORDINATES[Number(capacity)];
  const fillOrder = SEAT_FILL_ORDER[Number(capacity)];
  const sendLikeMutation = useSendLike();
  const banParticipantMutation = useBanParticipant();

  const [seatAssignments, setSeatAssignments] = useState<SeatAssignment[]>([]);
  const [mute, setMute] = useState<boolean>(false);
  const participantIds = useRef<string[]>([]);
  const participantNames = useRef<Record<string, string>>({});
  const participantMutes = useRef<Record<string, boolean>>({});
  const ws = useRef<WebSocket>(null);
  const peers = useRef<Record<string, RTCPeerConnection>>({});
  const localAudio = useRef<MediaStream>(null);
  const remoteAudios = useRef<Record<string, MediaStream>>({});

  const coordinateSeat = () => {
    const unique = [...new Set(participantIds.current)].sort((a, b) =>
      a.localeCompare(b),
    );
    const occupantBySeat: Record<number, string> = {};
    unique.forEach((id, idx) => {
      const seatIndex = fillOrder[idx];
      occupantBySeat[seatIndex] = id;
    });

    const nextAssignments: SeatAssignment[] = coordinates.map(
      (coordinate, idx) => ({
        id: occupantBySeat[idx],
        name: participantNames.current[occupantBySeat[idx]],
        mute: participantMutes.current[occupantBySeat[idx]],
        ...coordinate,
      }),
    );
    setSeatAssignments(nextAssignments);
  };

  const toggleAudio = () => {
    if (localAudio.current) {
      const audioTrack = localAudio.current.getAudioTracks()[0];
      audioTrack.enabled = !audioTrack.enabled;
      participantMutes.current[myId] = !participantMutes.current[myId];
      ws.current?.send(
        JSON.stringify({
          toId: Object.keys(peers.current),
          signal: { type: "mute" },
        }),
      );
      coordinateSeat();
      setMute(!mute);
    }
  };

  const handlePressParticipant = (id: string, name: string) => {
    if (isModerator) {
      showActionSheetWithOptions(
        {
          options: [`Send like to ${name}`, `Ban ${name}`, "Cancel"],
          destructiveButtonIndex: 1,
          cancelButtonIndex: 2,
        },
        (selectedIndex?: number) => {
          switch (selectedIndex) {
            case 0:
              sendLikeMutation.mutate(
                { toId: id },
                {
                  onSuccess: () => {
                    Toast.show({
                      type: "success",
                      text1: `${name} will know you sent like after this conversation`,
                    });
                  },
                },
              );
              break;
            case 1:
              banParticipantMutation.mutate(
                { conversationId: String(conversationId), banId: id },
                {
                  onSuccess: () => {
                    ws.current?.send(
                      JSON.stringify({
                        toId: id,
                        signal: { type: "ban" },
                      }),
                    );
                  },
                },
              );
              break;
            case 2:
              break;
          }
        },
      );
      return;
    }
    showActionSheetWithOptions(
      {
        options: [`Send like to ${name}`, "Cancel"],
        cancelButtonIndex: 1,
      },
      (selectedIndex?: number) => {
        switch (selectedIndex) {
          case 0:
            sendLikeMutation.mutate(
              { toId: id },
              {
                onSuccess: () => {
                  Toast.show({
                    type: "success",
                    text1: `${name} will know you sent like after this conversation`,
                  });
                },
              },
            );
            break;
          case 1:
            break;
        }
      },
    );
  };

  useFocusEffect(() => {
    const joinConversation = async () => {
      console.log("try to connect ws");
      ws.current = new WebSocket(
        `ws://${baseUrl.ios}:8080/online/conversation/join?id=${conversationId}`,
        undefined,
        {
          headers: {
            Authorization: `Bearer ${await getSecureAsync("accessToken")}`,
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
      participantIds.current = [myId];
      participantNames.current[myId] = profile.name;
      participantMutes.current[myId] = false;
      coordinateSeat();

      ws.current.onmessage = async (event) => {
        console.log("get message");

        const data: ConversationSignalResponse = JSON.parse(event.data);
        if (!data.signal) {
          const unique = [...new Set(data.fromIds)];
          if (unique.length >= Number(capacity)) {
            router.replace("/conversations");
            return;
          }
          participantIds.current = [...participantIds.current, ...unique];

          for (const fromId of unique) {
            participantMutes.current[fromId] = false;
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

            ws.current?.send(
              JSON.stringify({
                toIds: [fromId],
                signal: { type: "name-offer", name: profile.name },
              }),
            );
          }
          coordinateSeat();
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
        if (data.signal.type === "name-offer") {
          participantNames.current = {
            ...participantNames.current,
            [fromId]: data.signal.name,
          };
          ws.current?.send(
            JSON.stringify({
              toIds: [fromId],
              signal: { type: "name-answer", name: profile.name },
            }),
          );
          coordinateSeat();
          return;
        }
        if (data.signal.type === "name-answer") {
          participantNames.current = {
            ...participantNames.current,
            [fromId]: data.signal.name,
          };
          coordinateSeat();
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
          participantIds.current = participantIds.current.filter(
            (x) => x !== fromId,
          );
          delete participantNames.current[fromId];
          coordinateSeat();
        }
        if (data.signal.type === "ban") {
          await queryClient.invalidateQueries({
            queryKey: [
              queryKey.CONVERSATION,
              queryKey.GET_ONLINE_CONVERSATIONS,
            ],
          });
          router.replace("/conversations");
        }
        if (data.signal.type === "mute") {
          participantMutes.current[fromId] = !participantMutes.current[fromId];
          coordinateSeat();
        }
      };
    };
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
      localAudio.current?.release();
      for (const remoteAudio of Object.values(remoteAudios.current)) {
        remoteAudio.release();
      }
      ws.current?.close();
    };
  });

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
          {seatAssignments.map((seat, idx) => (
            <View
              key={idx}
              style={[
                styles.participantSeat,
                { left: `${seat.left}%`, top: `${seat.top}%` },
              ]}
            >
              {seat.id ? (
                <>
                  {seat.mute ? (
                    <Feather name="mic-off" size={10} color="black" />
                  ) : (
                    <></>
                  )}
                  <Ionicons
                    name="person-circle-outline"
                    size={24}
                    color="black"
                    onPress={
                      seat.id !== myId
                        ? () =>
                            handlePressParticipant(
                              seat.id ?? "",
                              seat.name ?? "",
                            )
                        : undefined
                    }
                  />
                </>
              ) : (
                <Feather name="circle" size={24} color="black" />
              )}
              <Text style={styles.participantId}>{seat.name ?? ""}</Text>
            </View>
          ))}
        </View>
        <View style={styles.controls}>
          {mute ? (
            <CustomButton label="turn off your mic" onPress={toggleAudio} />
          ) : (
            <CustomButton label="turn on your mic" onPress={toggleAudio} />
          )}
        </View>
      </View>
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
});
