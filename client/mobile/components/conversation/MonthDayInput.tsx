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

type MonthDayItem = {
  month: number;
  day: number;
  weekdayLabel: string;
  monthLabel: string;
};

function buildMonthDayItems(year: number): MonthDayItem[] {
  const items: MonthDayItem[] = [];

  if (year === new Date().getFullYear() + 1) {
    for (let month = 1; month <= 12; month++) {
      const lastDay = new Date(year, month, 0).getDate();

      for (let day = 1; day <= lastDay; day++) {
        const date = new Date(year, month - 1, day);
        items.push({
          month,
          day,
          weekdayLabel: date.toLocaleDateString("en-US", { weekday: "short" }),
          monthLabel: date.toLocaleDateString("en-US", { month: "short" }),
        });
      }
    }
    return items;
  }
  let month = new Date().getMonth() + 1;
  let lastDay = new Date(year, month, 0).getDate();
  for (let day = new Date().getDate(); day <= lastDay; day++) {
    const date = new Date(year, month - 1, day);
    items.push({
      month,
      day,
      weekdayLabel: date.toLocaleDateString("en-US", { weekday: "short" }),
      monthLabel: date.toLocaleDateString("en-US", { month: "short" }),
    });
  }
  for (month; month <= 12; month++) {
    lastDay = new Date(year, month, 0).getDate();
    for (let day = 1; day <= lastDay; day++) {
      const date = new Date(year, month - 1, day);
      items.push({
        month,
        day,
        weekdayLabel: date.toLocaleDateString("en-US", { weekday: "short" }),
        monthLabel: date.toLocaleDateString("en-US", { month: "short" }),
      });
    }
  }
  return items;
}

export default function MonthDayInput() {
  const { control, getValues } = useFormContext();
  const [modalVisible, setModalVisible] = useState(false);

  const openModal = () => {
    setModalVisible(true);
  };

  const closeModal = () => {
    setModalVisible(false);
  };

  return (
    <Controller
      name="monthDay"
      control={control}
      render={({ field: { onChange, value } }) => {
        const year = Number(getValues("year")) || new Date().getFullYear();
        const allMonthDays = buildMonthDayItems(year);
        const [selectedMonth, selectedDay] = String(value ?? "").split(".");
        const selected = allMonthDays.find(
          (item) =>
            String(item.month) === selectedMonth &&
            String(item.day) === selectedDay,
        );
        const display = selected
          ? `${selected.weekdayLabel} ${selected.day} ${selected.monthLabel}`
          : "Select date";

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
                  {allMonthDays.map((item, index) => (
                    <Pressable
                      key={index}
                      style={styles.row}
                      onPress={() => {
                        onChange(`${item.month}.${item.day}`);
                        closeModal();
                      }}
                    >
                      <Text style={styles.dateText} numberOfLines={1}>
                        {item.weekdayLabel} {item.day} {item.monthLabel}
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
    paddingVertical: 12,
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
    justifyContent: "center",
    gap: 12,
  },
  dateText: {
    flex: 1,
    fontSize: 16,
    color: colors.BLACK,
  },
  monthText: {
    fontSize: 14,
    color: colors.GRAY_700,
  },
});
