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

type LengthItem = {
  length: number;
};

function buildLengthItems(): LengthItem[] {
  const items: LengthItem[] = [];
  for (let length = 20; length <= 180; length += 1) {
    items.push({ length });
  }
  return items;
}

export default function LengthInput() {
  const { control } = useFormContext();
  const [modalVisible, setModalVisible] = useState(false);
  const [allLengths] = useState<LengthItem[]>(() => buildLengthItems());

  const openModal = () => {
    setModalVisible(true);
  };

  const closeModal = () => {
    setModalVisible(false);
  };

  return (
    <Controller
      name="length"
      control={control}
      rules={{
        validate: (data: string) => {
          if (data.trim().length === 0) {
            return "length is required";
          }

          const length = Number(data);
          if (Number.isNaN(length) || length < 20 || length > 180) {
            return "length should be a number between 20 and 180";
          }
        },
      }}
      render={({ field: { onChange, value }, fieldState: { error } }) => {
        const selected = allLengths.find(
          (item) => String(item.length) === String(value),
        );
        const display = selected?.length ?? 20;

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
                  {allLengths.map((item, idx) => (
                    <Pressable
                      key={idx}
                      style={styles.row}
                      onPress={() => {
                        onChange(String(item.length));
                        closeModal();
                      }}
                    >
                      <Text style={styles.valueText} numberOfLines={1}>
                        {item.length}
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
