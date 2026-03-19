import { Pressable, StyleSheet, Text, View } from "react-native";
import { colors } from "@/constants";
import { router } from "expo-router";
import { ConversationFeedResponse } from "@/types/conversation";
import { useAuth } from "@/hooks/useAuth";

interface OnlineConversationItemProps {
  conversation: ConversationFeedResponse;
}

export default function OnlineConversationItem({
  conversation,
}: OnlineConversationItemProps) {
  const { id } = useAuth();
  console.log(conversation.id);

  return (
    <Pressable
      style={styles.container}
      onPress={() =>
        router.replace({
          pathname: `/conversation/online/[id]`,
          params: {
            id: conversation.id,
            novel: conversation.novel ?? "",
            shortStory: conversation.shortStory ?? "",
            poem: conversation.poem ?? "",
            play: conversation.play ?? "",
            film: conversation.film ?? "",
            rule: conversation.rule ?? "",
            by: conversation.by ?? "",
            when: conversation.when,
            length: conversation.length,
            isModerator: conversation.isModerator ? "true" : "",
          },
        })
      }
    >
      <View
        style={[
          styles.content,
          conversation.isRegistrant && styles.registrant,
          conversation.isModerator && styles.moderator,
        ]}
      >
        <Text style={styles.detail}>
          {conversation.ongoing && "🔴 ongoing"}
        </Text>
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
            .format(new Date(conversation.when))
            .replace(/\sat\s/, " ")}
          {` For ${conversation.length.replace("0s", "")}`}
        </Text>
        {conversation.novel && (
          <Text style={styles.detail}>Novel: {conversation.novel}</Text>
        )}
        {conversation.shortStory && (
          <Text style={styles.detail}>
            Short story: {conversation.shortStory}
          </Text>
        )}
        {conversation.poem && (
          <Text style={styles.detail}>Poem: {conversation.poem}</Text>
        )}
        {conversation.play && (
          <Text style={styles.detail}>Play: {conversation.play}</Text>
        )}
        {conversation.film && (
          <Text style={styles.detail}>Film: {conversation.film}</Text>
        )}
        {conversation.by && (
          <Text style={styles.detail}>By: {conversation.by}</Text>
        )}
        {conversation.rule ? (
          <View>
            <Text style={styles.ruleHeader}>Rule</Text>{" "}
            <Text style={styles.detail}>{conversation.rule}</Text>
          </View>
        ) : (
          <Text style={styles.ruleHeader}>No rule</Text>
        )}
      </View>
    </Pressable>
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
    backgroundColor: colors.ORANGE_150,
  },
  registrant: {
    backgroundColor: colors.ORANGE_100,
  },
});
