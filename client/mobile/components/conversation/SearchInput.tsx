import { StyleSheet, TextInput, TextInputProps, View } from "react-native";
import { colors } from "@/constants";
import Ionicons from "@expo/vector-icons/Ionicons";

interface SearchInputProps extends TextInputProps {
  onSubmit?: () => void;
  submitKeyword: string;
  keyword: string;
  onCancel: () => void;
}

export default function SearchInput({
  onSubmit,
  submitKeyword,
  keyword,
  onCancel,
  ...props
}: SearchInputProps) {
  return (
    <View style={styles.container}>
      <TextInput
        style={styles.input}
        autoCapitalize="none"
        placeholderTextColor={colors.GRAY_500}
        returnKeyType="search"
        onSubmitEditing={onSubmit}
        {...props}
      />
      {submitKeyword === keyword ? (
        <Ionicons name="close" size={20} onPress={() => onCancel()} />
      ) : (
        <Ionicons
          name="search"
          size={20}
          onPress={props.onPress ?? onSubmit}
          color={colors.GRAY_500}
        />
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    gap: 8,
    flexDirection: "row",
    alignItems: "center",
    paddingVertical: 8,
    height: 42,
    paddingRight: 13,
    backgroundColor: colors.WHITE,
    borderRadius: 40,
  },
  input: {
    flex: 1,
    fontSize: 17,
    paddingVertical: 0,
    paddingLeft: 20,
    color: colors.BLACK,
    fontFamily: "SpaceMono",
  },
});
