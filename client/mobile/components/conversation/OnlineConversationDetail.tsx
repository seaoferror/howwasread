import { StyleSheet, Text, View } from "react-native";
import {
  useBlockConversation,
  useDeregisterOnlineConversation,
  useGetOnlineConversationDetail,
  useRegisterOnlineConversation,
} from "@/hooks/useConversation";
import CustomButton from "@/components/CustomButton";
import { colors } from "@/constants";
import Toast from "react-native-toast-message";
import { useActionSheet } from "@expo/react-native-action-sheet";
import { reportUser } from "@/api/chat";
import { reportOnlineConversation } from "@/api/conversation";

interface OnlineConversationDetailProps {
  id: string;
}

export default function OnlineConversationDetail({
  id,
}: OnlineConversationDetailProps) {
  const { data } = useGetOnlineConversationDetail(id);
  const blockConversationMutation = useBlockConversation();
  const { showActionSheetWithOptions } = useActionSheet();
  const registerOnlineConversationMutation = useRegisterOnlineConversation();
  const deregisterOnlineConversationMutation =
    useDeregisterOnlineConversation();

  const handleMore = () => {
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
              id: id,
            });
            break;
          case 1:
            blockConversationMutation.mutate({
              id: id,
            });
            const reportPromises = [
              reportOnlineConversation({ id }),
              ...(data?.moderatorIds.map((mid) => reportUser({ id: mid })) ||
                []),
            ];
            try {
              await Promise.all(reportPromises);
            } catch (e) {
              console.log(e);
            }
            Toast.show({
              type: "info",
              text1: "Success report",
              text2:
                "We will review this conversation, sorry for inconvenience.",
            });
        }
      },
    );
  };

  const handleEnter = () => {};

  return (
    data && (
      <View>
        <View style={styles.box}>
          <View style={[styles.content]}>
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
                .format(new Date(data.when))
                .replace(/\sat\s/, " ")}
              {` For ${data.length}m`}
            </Text>
            {data.novel && (
              <Text style={styles.detail}>Novel: {data.novel}</Text>
            )}
            {data.shortStory && (
              <Text style={styles.detail}>Short story: {data.shortStory}</Text>
            )}
            {data.poem && <Text style={styles.detail}>Poem: {data.poem}</Text>}
            {data.play && <Text style={styles.detail}>Play: {data.play}</Text>}
            {data.film && <Text style={styles.detail}>Film: {data.film}</Text>}
            <Text style={styles.detail}>Written by: {data.by}</Text>
            {data.rule ? (
              <View>
                <Text style={styles.ruleHeader}>Rule</Text>{" "}
                <Text style={styles.detail}>{data.rule}</Text>
              </View>
            ) : (
              <Text style={styles.ruleHeader}>No rule</Text>
            )}
          </View>
          <View style={styles.buttonRow}>
            <CustomButton
              label={!data.isRegistrant ? "Register" : "Cancel registration"}
              onPress={
                !data.isRegistrant
                  ? () => registerOnlineConversationMutation.mutate({id: id})
                  : () => deregisterOnlineConversationMutation.mutate({id: id})
              }
              disabled={
                registerOnlineConversationMutation.isPending ||
                deregisterOnlineConversationMutation.isPending
              }
            />
            <CustomButton
              label={
                data.canEnter
                  ? "Enter"
                  : data.isRegistrant
                    ? "Enter available 15 minutes before conversation start"
                    : "Non-registrant enter open 10 minutes after conversation start"
              }
              onPress={() => handleEnter()}
              disabled={!data.canEnter}
            />
          </View>
        </View>
      </View>
    )
  );
}
const styles = StyleSheet.create({
  container: { backgroundColor: colors.SAND_110 },
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
  footer: {
    marginVertical: 50,
    alignItems: "center",
  },
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
  reportPressed: {
    opacity: 0.6,
  },
  buttonRow: {
    flexDirection: "row",
  },
});
