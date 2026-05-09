import { Pressable, StyleSheet, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { router } from "expo-router";
import { useActionSheet } from "@expo/react-native-action-sheet";
import { useAuth } from "@/hooks/useAuth";
import { colors } from "@/constants";

export default function AccountScreen() {
  const { showActionSheetWithOptions } = useActionSheet();
  const { deleteAccountMutation } = useAuth();
  const handleDeleteAccount = () => {
    showActionSheetWithOptions(
      {
        title: "This is irreversible",
        options: [`Delete your account`, "Cancel"],
        cancelButtonIndex: 1,
        destructiveButtonIndex: 0,
      },
      (selectedIndex?: number) => {
        switch (selectedIndex) {
          case 0:
            deleteAccountMutation.mutate(undefined, {
              onSuccess: async () => {
                router.replace("/");
              },
            });
            break;
          case 1:
            break;
        }
      },
    );
  };

  return (
    <SafeAreaView style={styles.safeArea}>
        <View style={styles.menuContainer}>
          <Pressable 
            onPress={() => router.push("/profile/name")}
            style={({ pressed }) => [
              styles.menuItem,
              pressed && styles.menuItemPressed,
            ]}
          >
            <Text style={styles.menuText}>Change your profile</Text>
          </Pressable>
        </View>

        <View style={styles.footer}>
          <Pressable 
            onPress={handleDeleteAccount}
            style={({ pressed }) => [pressed && styles.deletePressed]}
          >
            <Text style={styles.deleteText}>Delete your account</Text>
          </Pressable>
        </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: {
    flex: 1,
    backgroundColor: colors.SAND_110,
  },
  menuContainer: {
    gap: 10,
  },
  menuItem: {
    padding: 20,
    flex: 1,
    borderWidth: StyleSheet.hairlineWidth,
    borderBottomColor: colors.BLACK,
  },
  menuItemPressed: {
    opacity: 0.6,
  },
  menuText: {
    fontSize: 18,
    color: colors.BLACK,
  },
  footer: {
    position: "absolute",
    bottom: 40,
    left: 0,
    right: 0,
    alignItems: "center",
  },
  deleteText: {
    color: colors.GRAY_400,
    fontSize: 12,
  },
  deletePressed: {
    opacity: 0.6,
  },
});
