import { ScrollView, StyleSheet, Text, View } from "react-native";
import { useLocalSearchParams } from "expo-router";
import { useGetChatParticipants } from "@/hooks/useChat";
import OfflineConversationDetail from "@/components/conversation/OfflineConversationDetail";
import { getKVStore, setKVStore } from "@/db/storage";
import MemberItem from "@/components/chat/MemberItem";
import { useEffect } from "react";
import { z } from "zod";
import { SafeAreaView } from "react-native-safe-area-context";

const participantSchema = z.array(
  z.object({
    id: z.string(),
  }),
);

export default function ChatRoomDetailScreen() {
  const { id: roomId } = useLocalSearchParams();
  const { data: participants } = useGetChatParticipants(String(roomId));

  useEffect(() => {
    if (!participants) {
      return;
    }
    setKVStore("participants" + String(roomId), JSON.stringify(participants));
  }, [participants, roomId]);

  const local = participantSchema.safeParse(
    JSON.parse(getKVStore("participants" + String(roomId)) || "[]"),
  );
  if (!local.success) {
    return null;
  }

  const ps: { id: string }[] = participants ? participants : local.data;

  return (
    <SafeAreaView style={styles.container} edges={["top"]}>
      <ScrollView style={{marginBottom: 100}}>
        <View style={{gap:20}}>
          <View style={styles.box}>
            <OfflineConversationDetail id={String(roomId)} />
          </View>

          <View style={styles.box}>
            <Text style={styles.listTitle}>Member list</Text>
            <View style={styles.memberContainer}>
              {ps.map((p, idx) => (
                <MemberItem key={idx} id={p.id} />
              ))}
            </View>
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    marginTop: 50,
  },
  box: {
    padding: 16,
    marginHorizontal: 16,
    backgroundColor: "#ffffff",
    borderRadius: 8,
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 2,
  },
  listTitle: {
    fontSize: 16,
    fontWeight: "bold",
    marginBottom: 12,
  },
  memberContainer: {
    gap: 8,
  },
});
