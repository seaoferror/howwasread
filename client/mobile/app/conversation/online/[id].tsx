import { Platform, Pressable, StyleSheet, Text, View } from "react-native";
import { useFocusEffect, useLocalSearchParams } from "expo-router";
import { colors } from "@/constants";
import { useRef, useState } from "react";
import { ConversationSignalResponse } from "@/types/conversation";
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

//Prevent ts compiler kept trying to read dom definition of WebSocket written by Microsoft, which don't have header options
declare let WebSocket: {
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
  const {
    id: conversationId,
    novel,
    shortStory,
    poem,
    play,
    film,
    by,
    rule,
    when,
    length,
    isModerator,
  } = useLocalSearchParams();

  const [mute, setMute] = useState(true);
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
        for (const fromId of data.fromIds) {
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
        <View style={styles.content} />
        <View style={styles.controls}>
          <Pressable style={styles.muteButton} onPress={toggleAudio}>
            <Text style={styles.muteText}>{mute ? "Unmute" : "Mute"}</Text>
          </Pressable>
        </View>
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
  controls: {
    flexDirection: "row",
    justifyContent: "space-around",
    alignItems: "center",
    paddingVertical: 24,
    paddingHorizontal: 16,
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
