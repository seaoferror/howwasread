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

  for (let year = 1900; year <= 2100; year++) {
    items.push({ year });
  }

  return items;
}

export default function BirthYearBox() {
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
      name="birthYear"
      control={control}
      render={({ field: { onChange, value } }) => {
        const selected = allYears.find((item) => item.year === value);
        const display = selected?.year ?? 2000;
        return (
          <>
            <Pressable onPress={openModal} style={styles.box}>
              <Text style={styles.boxText} numberOfLines={1}>
                {display}
              </Text>
            </Pressable>

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
                        onChange(item.year);
                        closeModal();
                      }}
                    >
                      <Text style={styles.year} numberOfLines={1}>
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
  box: {
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
  boxText: {
    color: colors.BLACK,
    fontSize: 16,
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
  year: {
    flex: 1,
    fontSize: 16,
    color: colors.BLACK,
  },
});
