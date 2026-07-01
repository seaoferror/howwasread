import { StyleSheet, Text, View } from "react-native";
import {
  useGetOfflineConversationDetail,
  useJoinOfflineConversation,
  useQuitOfflineConversation,
} from "@/hooks/useConversation";
import CustomButton from "@/components/CustomButton";
import { colors } from "@/constants";
import { openBrowserAsync } from "expo-web-browser";
import Toast from "react-native-toast-message";
import { router } from "expo-router";

interface OfflineConversationDetailProps {
  id: string;
}

export default function OfflineConversationDetail({
  id,
}: OfflineConversationDetailProps) {
  const { data } = useGetOfflineConversationDetail(id);
  const joinOfflineConversationMutation = useJoinOfflineConversation();
  const quitOfflineConversationMutation = useQuitOfflineConversation();

  const handleButtonPress = () => {
    if (data?.isParticipant) {
      quitOfflineConversationMutation.mutate(
        { conversationId: id },
        { onSuccess: () => router.push("/conversations") },
      );
      return;
    }
    joinOfflineConversationMutation.mutate(
      { conversationId: String(id) },
      {
        onSuccess: () => {
          Toast.show({
            type: "success",
            text1: "We invite you to the group chat room! Check it out!",
          });
        },
      },
    );
  };

  return (
    data && (
      <View>
        <View style={[styles.content, data.isModerator && styles.moderator]}>
          <Text style={styles.when}>
            {new Intl.DateTimeFormat("en-US", {
              weekday: "short",
              year: "numeric",
              month: "short",
              day: "numeric",
              hour: "2-digit",
              minute: "2-digit",
              hourCycle: "h12",
            })
              .format(new Date(data.time))
              .replace(/\sat\s/, " ")}
            {` For ${data.length}m`}
          </Text>
          {data.novel && <Text style={styles.detail}>Novel: {data.novel}</Text>}
          {data.shortStory && (
            <Text style={styles.detail}>Short story: {data.shortStory}</Text>
          )}
          {data.poem && <Text style={styles.detail}>Poem: {data.poem}</Text>}
          {data.play && <Text style={styles.detail}>Play: {data.play}</Text>}
          {data.film && <Text style={styles.detail}>Film: {data.film}</Text>}
          <Text style={styles.detail}>Written by: {data.writtenBy}</Text>
          {data.rule ? (
            <View>
              <Text style={styles.ruleHeader}>Rule</Text>{" "}
              <Text style={styles.detail}>{data.rule}</Text>
            </View>
          ) : (
            <Text style={styles.ruleHeader}>No rule</Text>
          )}
          {
            <Text style={styles.detail}>
              number of participants: {data.numberOfParticipants}
            </Text>
          }
          {data.location && (
            <Text style={styles.detail}>Location: {data.location}</Text>
          )}
          <Text
            style={[styles.detail, { color: colors.BLUE_500 }]}
            onPress={() => openBrowserAsync(data.mapsLink)}
          >
            Maps link
          </Text>
        </View>
        <CustomButton
          label={
            data.isParticipant
              ? "Quit offline conversation"
              : "Join offline conversation"
          }
          onPress={handleButtonPress}
          style={{ paddingHorizontal: 20 }}
          disabled={
            joinOfflineConversationMutation.isPending ||
            quitOfflineConversationMutation.isPending
          }
        />
      </View>
    )
  );
}
const styles = StyleSheet.create({
  container: { backgroundColor: colors.SAND_110 },
  content: {
    padding: 16,
    gap: 17,
  },
  when: {
    fontSize: 19,
    color: colors.BLACK,
    fontWeight: 500,
    marginVertical: 6,
  },
  detail: {
    fontSize: 17,
    fontWeight: 300,
  },
  ruleHeader: {
    marginTop: 6,
    fontSize: 18,
    fontWeight: 400,
  },
  moderator: {
    backgroundColor: colors.SAND_110,
  },
});
