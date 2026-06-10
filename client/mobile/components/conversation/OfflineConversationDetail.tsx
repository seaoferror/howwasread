import { StyleSheet, Text, View } from "react-native";
import {
  useJoinOfflineConversation,
  useQuitOfflineConversation,
} from "@/hooks/useConversation";
import CustomButton from "@/components/CustomButton";
import { colors } from "@/constants";
import { openBrowserAsync } from "expo-web-browser";

interface OfflineConversationDetailProps {
  id: string;
  novel: string;
  poem: string;
  shortStory: string;
  play: string;
  film: string;
  writtenBy: string;
  rule: string;
  time: string;
  length: number;
  mapsLink: string;
  location: string;
  isModerator: boolean;
  isParticipant: boolean;
  numberOfParticipants: number;
}

export default function OfflineConversationDetail({
  id,
  novel,
  poem,
  shortStory,
  play,
  film,
  writtenBy,
  rule,
  time,
  length,
  mapsLink,
  location,
  isModerator,
  isParticipant,
  numberOfParticipants,
}: OfflineConversationDetailProps) {
  const joinOfflineConversationMutation = useJoinOfflineConversation();
  const quitOfflineConversationMutation = useQuitOfflineConversation();

  const handleButtonPress = () => {
    if (isParticipant) {
      quitOfflineConversationMutation.mutate({ conversationId: id });
      return;
    }
    joinOfflineConversationMutation.mutate({ conversationId: id });
  };

  return (
    <View>
      <View style={[styles.content, isModerator && styles.moderator]}>
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
            .format(new Date(time))
            .replace(/\sat\s/, " ")}
          {` For ${length}m`}
        </Text>
        {novel && <Text style={styles.detail}>Novel: {novel}</Text>}
        {shortStory && (
          <Text style={styles.detail}>Short story: {shortStory}</Text>
        )}
        {poem && <Text style={styles.detail}>Poem: {poem}</Text>}
        {play && <Text style={styles.detail}>Play: {play}</Text>}
        {film && <Text style={styles.detail}>Film: {film}</Text>}
        {writtenBy && (
          <Text style={styles.detail}>Written by: {writtenBy}</Text>
        )}
        {rule ? (
          <View>
            <Text style={styles.ruleHeader}>Rule</Text>{" "}
            <Text style={styles.detail}>{rule}</Text>
          </View>
        ) : (
          <Text style={styles.ruleHeader}>No rule</Text>
        )}
        {
          <Text style={styles.detail}>
            number of participants: {numberOfParticipants}
          </Text>
        }
        {location && <Text style={styles.detail}>Location: {location}</Text>}
        <Text
          style={[styles.detail, { color: colors.BLUE_500 }]}
          onPress={() => openBrowserAsync(mapsLink)}
        >
          Maps link
        </Text>
      </View>
      <CustomButton
        label={
          isParticipant
            ? "Quit offline conversation"
            : "Join offline conversation"
        }
        onPress={handleButtonPress}
      />
    </View>
  );
}
const styles = StyleSheet.create({
  container: { backgroundColor: colors.SAND_110 },
  content: {
    padding: 16,
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
