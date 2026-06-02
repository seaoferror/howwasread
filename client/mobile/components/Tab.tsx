import { colors } from "@/constants";
import React, { ReactNode } from "react";
import { Pressable, StyleSheet, Text } from "react-native";

interface TabProps {
  isActive: boolean;
  onPress?: () => void;
  children: ReactNode;
}

export default function Tab({ isActive, onPress, children }: TabProps) {
  return (
    <Pressable
      style={[styles.container, isActive && styles.activeContainer]}
      onPress={onPress}
    >
      <Text style={[styles.text, isActive && styles.activeText]}>
        {children}
      </Text>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  container: {
    height: 38,
    flex: 1,
    paddingVertical: 8,
    justifyContent: "center",
    alignItems: "center",
    borderBottomColor: colors.WHITE,
    borderBottomWidth: 2,
  },
  activeContainer: {
    borderBottomColor: colors.BLACK,
    borderBottomWidth: 2,
  },
  text: {
    fontSize: 14,
    color: colors.GRAY_500,
  },
  activeText: {
    color: colors.BLACK,
    fontWeight: "700",
  },
});

