"use dom";

import React, { useEffect, useRef, useState } from "react";
import { colors, SEAT_COORDINATES, SEAT_FILL_ORDER } from "@/constants";
import {
  ConversationSignalResponse,
  SeatAssignment,
} from "@/types/conversation";

export default function WebRTCRoom({
  conversationId,
  capacity,
  accessToken,
  myId,
  myName,
  apiUrl,
  onParticipantAction,
  onBanExit,
  onExit,
}: {
  conversationId: string;
  capacity: number;
  accessToken: string;
  myId: string;
  myName: string;
  apiUrl: string;
  onParticipantAction: (id: string, name: string) => Promise<string>;
  onBanExit: () => Promise<void>;
  onExit: () => Promise<void>;
}) {
  const [seatAssignments, setSeatAssignments] = useState<SeatAssignment[]>([]);
  const [mute, setMute] = useState<boolean>(false);
  const [isWebSocketOpen, setIsWebSocketOpen] = useState(false);

  // Track remote streams to bind them to HTML5 Audio elements
  const [remoteStreams, setRemoteStreams] = useState<
    Record<string, MediaStream>
  >({});

  const participantIds = useRef<string[]>([]);
  const participantNames = useRef<Record<string, string>>({});
  const participantMutes = useRef<Record<string, boolean>>({});
  const ws = useRef<WebSocket | null>(null);
  const peers = useRef<Record<string, RTCPeerConnection>>({});
  const localAudio = useRef<MediaStream | null>(null);

  const coordinates = SEAT_COORDINATES[capacity] || SEAT_COORDINATES[2];
  const fillOrder = SEAT_FILL_ORDER[capacity] || SEAT_FILL_ORDER[2];

  const coordinateSeat = () => {
    const unique = [...new Set(participantIds.current)].sort((a, b) =>
      a.localeCompare(b),
    );
    const occupantBySeat: Record<number, string> = {};
    unique.forEach((id, idx) => {
      occupantBySeat[fillOrder[idx]] = id;
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
    try {
      if (localAudio.current && localAudio.current.active) {
        const audioTracks = localAudio.current.getAudioTracks();
        if (
          ws.current &&
          ws.current.readyState === 1 &&
          audioTracks.length > 0
        ) {
          const audioTrack = audioTracks[0];
          audioTrack.enabled = !audioTrack.enabled;
          participantMutes.current[myId] = !participantMutes.current[myId];
          ws.current.send(
            JSON.stringify({
              toIds: Object.keys(peers.current),
              signal: { type: "mute" },
            }),
          );
          coordinateSeat();
          setMute(!mute);
        }
      }
    } catch (error) {
      console.error("Failed to toggle audio:", error);
    }
  };

  const handlePressParticipant = async (id: string, name: string) => {
    // Await the native action response before triggering local WebSocket state
    const action = await onParticipantAction(id, name);
    if (action === "ban" && ws.current) {
      ws.current.send(JSON.stringify({ toId: id, signal: { type: "ban" } }));
    }
  };

  useEffect(() => {
    const joinConversation = async () => {
      // Standard browser WebSockets don't support custom headers.
      // Token is passed via query string here. Ensure your backend reads this.
      ws.current = new WebSocket(
        `wss://${apiUrl}/onlineconversation/join?id=${conversationId}&token=${accessToken}`,
      );

      ws.current.onopen = () => setIsWebSocketOpen(true);

      // Replaces react-native-webrtc specific methods with standard window.navigator APIs
      localAudio.current = await navigator.mediaDevices.getUserMedia({
        audio: true,
        video: false,
      });

      participantIds.current = [myId];
      participantNames.current[myId] = myName;
      participantMutes.current[myId] = false;
      coordinateSeat();

      ws.current.onmessage = async (event) => {
        const data: ConversationSignalResponse = JSON.parse(event.data);

        if (!data.signal) {
          const unique = [...new Set(data.fromIds)];
          if (unique.length >= capacity) {
            onExit();
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
              const [track] = event.streams;
              // Save standard remote stream into state so it can be bound to the <audio> elements
              setRemoteStreams((prev) => ({
                ...prev,
                [fromId]: track || new MediaStream([event.track]),
              }));
            });

            const offer = await peers.current[fromId].createOffer({
              offerToReceiveAudio: true,
              offerToReceiveVideo: false,
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
                signal: { type: "name-offer", name: myName },
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
        if (
          data.signal.type === "name-offer" ||
          data.signal.type === "name-answer"
        ) {
          participantNames.current = {
            ...participantNames.current,
            [fromId]: data.signal.name,
          };
          if (data.signal.type === "name-offer") {
            ws.current?.send(
              JSON.stringify({
                toIds: [fromId],
                signal: { type: "name-answer", name: myName },
              }),
            );
          }
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
          peers.current[fromId]?.close();
          delete peers.current[fromId];

          setRemoteStreams((prev) => {
            const updated = { ...prev };
            delete updated[fromId];
            return updated;
          });

          participantIds.current = participantIds.current.filter(
            (x) => x !== fromId,
          );
          delete participantNames.current[fromId];
          coordinateSeat();
        }
        if (data.signal.type === "ban") onBanExit();
        if (data.signal.type === "mute") {
          participantMutes.current[fromId] = !participantMutes.current[fromId];
          coordinateSeat();
        }
      };
    };

    joinConversation();

    return () => {
      if (ws.current && ws.current.readyState === 1) {
        ws.current.send(
          JSON.stringify({
            toIds: Object.keys(peers.current),
            signal: { type: "leave" },
          }),
        );
      }
      for (const peer of Object.values(peers.current)) peer.close();
      if (localAudio.current)
        localAudio.current.getTracks().forEach((track) => track.stop());
      ws.current?.close();
    };
  }, []);

  return (
    <div style={styles.participantContainer}>
      <div style={styles.participantArea}>
        {seatAssignments.map((seat, idx) => (
          <div
            key={idx}
            style={{
              ...styles.participantSeat,
              left: `${seat.left}%`,
              top: `${seat.top}%`,
            }}
          >
            {seat.id ? (
              <div
                style={{
                  cursor: seat.id !== myId ? "pointer" : "default",
                  textAlign: "center",
                }}
                onClick={
                  seat.id !== myId
                    ? () =>
                        handlePressParticipant(seat.id ?? "", seat.name ?? "")
                    : undefined
                }
              >
                {/* Replaced Expo Vector Icons with standard Emojis/HTML for pure web compatibility */}
                <div style={{ fontSize: "24px" }}>👤</div>
                {seat.mute && (
                  <div style={{ fontSize: "10px", marginTop: "2px" }}>🔇</div>
                )}
              </div>
            ) : (
              <div style={{ fontSize: "24px", color: colors.BLACK }}>⚪</div>
            )}
            <span style={styles.participantId}>{seat.name ?? ""}</span>
          </div>
        ))}
      </div>

      <div style={styles.controls}>
        {!isWebSocketOpen ? (
          <button style={{ ...styles.button, opacity: 0.5 }} disabled>
            Connecting...
          </button>
        ) : mute ? (
          <button style={styles.button} onClick={toggleAudio}>
            Turn on your mic
          </button>
        ) : (
          <button style={styles.button} onClick={toggleAudio}>
            Turn off your mic
          </button>
        )}
      </div>

      {/* Hidden standard DOM audio tags vital for natively playing remote tracks */}
      {Object.entries(remoteStreams).map(([id, stream]) => (
        <audio
          key={id}
          autoPlay
          playsInline
          ref={(audioEl) => {
            if (audioEl && audioEl.srcObject !== stream) {
              audioEl.srcObject = stream;
            }
          }}
          style={{ display: "none" }}
        />
      ))}
    </div>
  );
}

// Replaced React Native StyleSheet with standard React Web inline styles
const styles: Record<string, React.CSSProperties> = {
  participantContainer: {
    display: "flex",
    flexDirection: "column",
    flex: 1,
    height: "100%",
    width: "100%",
  },
  participantArea: {
    display: "flex",
    flex: 1,
    position: "relative",
  },
  participantSeat: {
    position: "absolute",
    width: "96px",
    display: "flex",
    flexDirection: "column",
    alignItems: "center",
    // Translated React Native transform array into a standard CSS string
    transform: "translate(-48px, -12px)",
  },
  participantId: {
    marginTop: "6px",
    fontSize: "12px",
    color: colors.BLACK,
    fontFamily: "sans-serif",
  },
  controls: {
    display: "flex",
    flexDirection: "row",
    justifyContent: "center",
    alignItems: "center",
    paddingTop: "24px",
    paddingBottom: "24px",
    paddingLeft: "16px",
    paddingRight: "16px",
    gap: "12px",
  },
  button: {
    padding: "10px 16px",
    backgroundColor: colors.BLACK || "#000",
    color: "#fff",
    border: "none",
    borderRadius: "8px",
    cursor: "pointer",
    fontSize: "14px",
    fontWeight: "bold",
  },
};
