import React from "react";
import {
  Pressable,
  PressableProps,
  StyleProp,
  StyleSheet,
  Text,
  ViewStyle,
} from "react-native";
import { colors } from "@/constants";

interface CustomButtonProps extends PressableProps {
  label: string;
  size?: "medium" | "large";
  variant?: "filled" | "standard" | "outlined";
  style?: StyleProp<ViewStyle>;
}

export default function CustomButton({
  label,
  size = "large",
  variant = "filled",
  style = null,
  ...props
}: CustomButtonProps) {
  return (
    <Pressable
      {...props}
      style={({ pressed }) => [
        styles.container,
        styles[size],
        styles[variant],
        props.disabled && styles.disabled,
        pressed && styles.pressed,
        style,
      ]}
    >
      <Text style={styles[`${variant}Text`]}>{label}</Text>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  container: {
    borderRadius: 8,
    justifyContent: "center",
    alignItems: "center",
  },
  large: {
    width: "100%",
    height: 44,
  },
  medium: {
    height: 38,
    alignSelf: "center",
    paddingHorizontal: 12,
  },
  filled: {
    backgroundColor: colors.BLACK,
  },
  standard: {},
  outlined: {
    backgroundColor: colors.WHITE,
    borderWidth: 1,
    borderColor: colors.BLACK,
  },
  pressed: {
    opacity: 0.8,
  },
  disabled: {
    backgroundColor: colors.GRAY_300,
  },
  standardText: {
    fontSize: 17,
    fontWeight: "bold",
    color: colors.BLACK,
  },
  filledText: {
    fontSize: 17,
    fontWeight: "bold",
    color: colors.WHITE,
  },
  outlinedText: {
    fontSize: 17,
    fontWeight: "bold",
    color: colors.BLACK,
  },
});
