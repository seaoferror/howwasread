import { Pressable, StyleSheet, Text, View } from "react-native";
import { useGetProfile } from "@/hooks/useProfile";
import { getKVStore } from "@/db/storage";
import { colors } from "@/constants";
import { router } from "expo-router";

export default function MemberItem({ id }: { id: string }) {
  const { data } = useGetProfile(id);

  const displayName = data?.name || getKVStore(id);

  return (
    <Pressable
      // 2. Add visual feedback when the item is pressed
      style={({ pressed }) => [styles.container, pressed && styles.pressed]}
      onPress={() =>
        router.push({ pathname: "/chat/[id]", params: { id: id } })
      }
    >
      <View style={styles.avatar}>
        <Text style={styles.avatarText}>
          {displayName.charAt(0).toUpperCase()}
        </Text>
      </View>
      <Text style={styles.nameText}>{displayName}</Text>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: "row",
    alignItems: "center",
    paddingVertical: 12,
    paddingHorizontal: 16,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: "#E5E5E5",
  },
  pressed: {
    backgroundColor: "#F3F4F6",
  },
  avatar: {
    width: 44,
    height: 44,
    borderRadius: 22,
    backgroundColor: colors.ORANGE_150,
    justifyContent: "center",
    alignItems: "center",
    marginRight: 16,
  },
  avatarText: {
    fontSize: 18,
    fontWeight: "600",
    color: colors.GRAY_700,
  },
  nameText: {
    fontSize: 16,
    fontWeight: "500",
    color: "#1F2937",
  },
});
