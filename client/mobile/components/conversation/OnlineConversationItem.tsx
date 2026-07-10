import { Pressable, StyleSheet, Text, View } from "react-native";
import { colors } from "@/constants";
import { router } from "expo-router";
import { OnlineConversationFeedResponse } from "@/types/conversation";
import { useActionSheet } from "@expo/react-native-action-sheet";
import { reportUser } from "@/api/chat";
import { useBlockConversation } from "@/hooks/useConversation";
import { reportOnlineConversation } from "@/api/conversation";
import Toast from "react-native-toast-message";

interface OnlineConversationItemProps {
  conversation: OnlineConversationFeedResponse;
}

export default function OnlineConversationItem({
  conversation,
}: OnlineConversationItemProps) {
  console.log(conversation.id);
  const { showActionSheetWithOptions } = useActionSheet();
  const blockConversationMutation = useBlockConversation();

  const handleLongPress = () => {
    showActionSheetWithOptions(
      {
        options: ["Delete from feed", `Report and Delete from feed`, "Cancel"],
        destructiveButtonIndex: 1,
        cancelButtonIndex: 2,
      },
      async (selectedIndex?: number) => {
        switch (selectedIndex) {
          case 0:
            blockConversationMutation.mutate({
              id: conversation.id,
            });
            break;
          case 1:
            blockConversationMutation.mutate({
              id: conversation.id,
            });
            const reportPromises = [
              reportOnlineConversation({ id: conversation.id }),
              ...conversation.moderatorIds.map((mid) =>
                reportUser({ id: mid }),
              ),
            ];
            try {
              await Promise.all(reportPromises);
            } catch (e) {
              console.log(e);
            }
            Toast.show({
              type: "info",
              text1: "Success report",
              text2: "We will review this conversation, sorry for inconvenience.",
            });
        }
      },
    );
  };

  return (
    <Pressable
      style={styles.container}
      onPress={() =>
        router.push({
          pathname: `/online/[id]`,
          params: {
            id: conversation.id,
            novel: conversation.novel ?? "",
            shortStory: conversation.shortStory ?? "",
            poem: conversation.poem ?? "",
            play: conversation.play ?? "",
            film: conversation.film ?? "",
            by: conversation.by ?? "",
            rule: conversation.rule ?? "",
            capacity: conversation.capacity,
            when: conversation.when,
            length: conversation.length,
            isModerator: conversation.isModerator ? "true" : "",
          },
        })
      }
      onLongPress={() => handleLongPress()}
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
        {!!conversation.novel && (
          <Text style={styles.detail}>Novel: {conversation.novel}</Text>
        )}
        {!!conversation.shortStory && (
          <Text style={styles.detail}>
            Short story: {conversation.shortStory}
          </Text>
        )}
        {!!conversation.poem && (
          <Text style={styles.detail}>Poem: {conversation.poem}</Text>
        )}
        {!!conversation.play && (
          <Text style={styles.detail}>Play: {conversation.play}</Text>
        )}
        {!!conversation.film && (
          <Text style={styles.detail}>Film: {conversation.film}</Text>
        )}
        {!!conversation.by && (
          <Text style={styles.detail}>By: {conversation.by}</Text>
        )}
        {!!conversation.rule ? (
          <View>
            <Text style={styles.ruleHeader}>Rule</Text>
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
