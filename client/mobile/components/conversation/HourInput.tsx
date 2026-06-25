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

type HourItem = {
  hour: number;
};

function buildHourItems(): HourItem[] {
  const items: HourItem[] = [];

  for (let hour = 23; hour >= 0; hour--) {
    items.push({ hour });
  }

  return items;
}

export default function HourInput() {
  const { control } = useFormContext();
  const [modalVisible, setModalVisible] = useState(false);
  const [allHours] = useState<HourItem[]>(() => buildHourItems());

  const openModal = () => {
    setModalVisible(true);
  };

  const closeModal = () => {
    setModalVisible(false);
  };

  return (
    <Controller
      name="hour"
      control={control}
      rules={{
        validate: (data: string) => {
          if (String(data ?? "").trim().length === 0) {
            return "hour is required";
          }
        },
      }}
      render={({ field: { onChange, value }, fieldState: { error } }) => {
        const selected = allHours.find(
          (item) => String(item.hour) === String(value),
        );
        const display = selected
          ? String(selected.hour).padStart(2, "0")
          : "20";
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
                  {allHours.map((item, index) => (
                    <Pressable
                      key={index}
                      style={styles.row}
                      onPress={() => {
                        onChange(String(item.hour));
                        closeModal();
                      }}
                    >
                      <Text style={styles.valueText} numberOfLines={1}>
                        {`At ${String(item.hour)}`}
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
    paddingVertical: 12,
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
