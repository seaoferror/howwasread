import React, { useState } from "react";
import {
  Modal,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from "react-native";
import { Controller, useFormContext } from "react-hook-form";
import { colors } from "@/constants";

type YearItem = {
  year: number;
};

function buildYearItems(): YearItem[] {
  const items: YearItem[] = [];
  const currentYear = new Date().getFullYear();
  items.push({ year: currentYear });
  items.push({ year: currentYear + 1 });

  return items;
}

export default function YearInput() {
  const { control } = useFormContext();
  const [modalVisible, setModalVisible] = useState(false);
  const [allYears] = useState<YearItem[]>(() => buildYearItems());

  const openModal = () => {
    setModalVisible(true);
  };

  const closeModal = () => {
    setModalVisible(false);
  };

  return (
    <Controller
      name="year"
      control={control}
      rules={{
        validate: (data: string) => {
          if (String(data ?? "").trim().length === 0) {
            return "year is required";
          }
        },
      }}
      render={({ field: { onChange, value }, fieldState: { error } }) => {
        const selected = allYears.find(
          (item) => String(item.year) === String(value),
        );
        const display = selected?.year ?? new Date().getFullYear();
        return (
          <>
            <Pressable
              onPress={openModal}
              style={[styles.box, Boolean(error) && styles.boxError]}
            >
              <Text style={styles.boxText} numberOfLines={1}>
                {display}
              </Text>
            </Pressable>
            {Boolean(error?.message) && (
              <Text style={styles.error}>{error?.message}</Text>
            )}

            <Modal
              visible={modalVisible}
              transparent
              animationType="none"
              onRequestClose={closeModal}
            >
              <Pressable style={styles.backdrop} onPress={closeModal} />

              <View style={styles.picker}>
                <View style={styles.handle} />

                <ScrollView>
                  {allYears.map((item, index) => (
                    <Pressable
                      key={index}
                      style={styles.row}
                      onPress={() => {
                        onChange(String(item.year));
                        closeModal();
                      }}
                    >
                      <Text style={styles.valueText} numberOfLines={1}>
                        {item.year}
                      </Text>
                    </Pressable>
                  ))}
                </ScrollView>
              </View>
            </Modal>
          </>
        );
      }}
    />
  );
}

const styles = StyleSheet.create({
  label: {
    fontSize: 12,
    color: colors.GRAY_700,
    marginBottom: 5,
  },
  box: {
    flex: 1,
    borderWidth: 1,
    borderColor: colors.GRAY_200,
    borderRadius: 10,
    paddingHorizontal: 12,
    paddingVertical: 12,
    minWidth: 92,
    alignItems: "center",
    justifyContent: "center",
    backgroundColor: colors.WHITE,
  },
  boxError: {
    backgroundColor: colors.RED_100,
  },
  boxText: {
    color: colors.BLACK,
    fontSize: 16,
  },
  error: {
    fontSize: 12,
    marginTop: 5,
    color: colors.RED_500,
  },
  backdrop: {
    ...StyleSheet.absoluteFill,
    backgroundColor: "rgba(0,0,0,0.35)",
  },
  picker: {
    position: "absolute",
    left: 0,
    right: 0,
    bottom: 0,
    height: "50%",
    backgroundColor: colors.WHITE,
    padding: 20,
    borderTopLeftRadius: 16,
    borderTopRightRadius: 16,
  },
  handle: {
    alignSelf: "center",
    width: 44,
    height: 5,
    borderRadius: 999,
    backgroundColor: colors.GRAY_200,
    marginBottom: 12,
  },
  row: {
    paddingVertical: 14,
    borderBottomWidth: StyleSheet.hairlineWidth,
    borderBottomColor: colors.GRAY_200,
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    gap: 12,
  },
  valueText: {
    flex: 1,
    fontSize: 16,
    color: colors.BLACK,
  },
});
