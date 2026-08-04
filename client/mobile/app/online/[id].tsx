import { Pressable, StyleSheet } from "react-native";
import { router, useLocalSearchParams, useNavigation } from "expo-router";
import { useEffect, useState } from "react";
import { getKVStore, getSecureAsync } from "@/db/storage";
import { SafeAreaView } from "react-native-safe-area-context";
import OnlineConversationRoomHeader from "@/components/conversation/OnlineConversationRoomHeader";
import { Ionicons } from "@expo/vector-icons";
import { useGetMyProfile } from "@/hooks/useProfile";
import { useActionSheet } from "@expo/react-native-action-sheet";
import { useSendMessage } from "@/hooks/useChat";
import Toast from "react-native-toast-message";
import queryClient from "@/api/queryClient";
import {
  useBanParticipant,
  useGetOnlineConversationDetail,
} from "@/hooks/useConversation";
import { colors, queryKey } from "@/constants";
import WebRTCRoom from "@/components/conversation/WebRTCRoom";

export default function OnlineConversationScreen() {
  const { data: profile } = useGetMyProfile();
  const { id: conversationId, isPersonal } = useLocalSearchParams();
  const { data: detail } = useGetOnlineConversationDetail({
    id: String(conversationId),
    isPersonal: Boolean(isPersonal),
  });

  const { showActionSheetWithOptions } = useActionSheet();
  const sendMessageMutation = useSendMessage();
  const banParticipantMutation = useBanParticipant();
  const navigation = useNavigation();

  const [accessToken, setAccessToken] = useState("");

  useEffect(() => {
    getSecureAsync("accessToken").then((token) => setAccessToken(token || ""));
  }, []);

  useEffect(() => {
    if (!isPersonal) return;
    navigation.setOptions({
      title: "Voice call",
      headerShown: true,
      headerLeft: () => (
        <Pressable onPress={() => router.back()}>
          <Ionicons name="chevron-back" size={28} color="black" />
        </Pressable>
      ),
    });
  }, [isPersonal]);

  // Pass an async function to the DOM component so it can trigger native UI
  const handleParticipantAction = async (
    id: string,
    name: string,
  ): Promise<string> => {
    return new Promise((resolve) => {
      const isModerator =
        detail?.moderatorIds.some(
          (m) => m === (profile?.id ?? getKVStore("myId")),
        ) ?? false;

      const options = isModerator
        ? [`Send like to ${name}`, `Ban ${name}`, "Cancel"]
        : [`Send like to ${name}`, "Cancel"];
      const destructiveButtonIndex = isModerator ? 1 : undefined;
      const cancelButtonIndex = isModerator ? 2 : 1;

      showActionSheetWithOptions(
        { options, destructiveButtonIndex, cancelButtonIndex },
        (selectedIndex?: number) => {
          if (selectedIndex === 0) {
            sendMessageMutation.mutate(
              {
                toIdType: "personal",
                toId: id,
                contentType: "text",
                contents: ["👍"],
              },
              {
                onSuccess: () =>
                  Toast.show({
                    type: "success",
                    text1: `You sent a like to ${name}`,
                  }),
              },
            );
            resolve("liked");
          } else if (isModerator && selectedIndex === 1) {
            banParticipantMutation.mutate(
              { conversationId: String(conversationId), banId: id },
              { onSuccess: () => resolve("ban") },
            );
          } else {
            resolve("cancelled");
          }
        },
      );
    });
  };

  const handleBanExit = async () => {
    await queryClient.invalidateQueries({
      queryKey: [queryKey.CONVERSATION, queryKey.GET_ONLINE_CONVERSATIONS],
    });
    router.replace("/conversations");
  };

  const handleExit = async () => {
    router.replace("/conversations");
  };

  if (!detail || !profile || !accessToken) return null;

  return (
    <SafeAreaView style={styles.container}>
      {!isPersonal && (
        <OnlineConversationRoomHeader
          novel={String(detail.novel)}
          shortStory={String(detail.shortStory)}
          poem={String(detail.poem)}
          play={String(detail.play)}
          film={String(detail.film)}
          writtenBy={String(detail.writtenBy)}
          rule={String(detail.rule)}
          time={String(detail.time)}
          length={String(detail.length)}
        />
      )}
      <WebRTCRoogm
        conversationId={String(conversationId)}
        capacity={Number(detail?.capacity ?? 2)}
        accessToken={accessToken}
        myId={profile.id ?? getKVStore("myId")}
        myName={profile.name ?? getKVStore("myName")}
        apiUrl={process.env.EXPO_PUBLIC_API_URL || ""}
        onParticipantAction={handleParticipantAction}
        onBanExit={handleBanExit}
        onExit={handleExit}
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: colors.SAND_110 },
});
