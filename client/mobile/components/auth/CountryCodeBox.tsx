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
import { CountryCode, getCountryCallingCode } from "libphonenumber-js";
import { colors } from "@/constants";
import countries from "i18n-iso-countries";

type CountryItem = {
  cca2: string;
  name: string;
  label: string;
};

function buildCountryItems(): CountryItem[] {
  countries.registerLocale(require("i18n-iso-countries/langs/en.json"));

  const names = countries.getNames("en", { select: "official" });
  delete names["AQ"];
  delete names["BV"];
  delete names["GS"];
  delete names["HM"];
  delete names["TF"];
  delete names["UM"];
  delete names["PN"];

  const items: CountryItem[] = [];

  for (const [cca2, name] of Object.entries<string>(names)) {
    const callingCode = getCountryCallingCode(cca2 as CountryCode);
    const label = `+${callingCode}`;
    items.push({ cca2, name, label });
  }

  items.sort((a, b) => a.name.localeCompare(b.name));
  return items;
}

export default function CountryCodeBox() {
  const { control } = useFormContext();
  const [modalVisible, setModalVisible] = useState(false);
  const [allCountries] = useState<CountryItem[]>(() => buildCountryItems());

  const openModal = () => {
    setModalVisible(true);
  };

  const closeModal = () => {
    setModalVisible(false);
  };

  return (
    <Controller
      name="countryCode"
      control={control}
      render={({ field: { onChange, value } }) => {
        const selected = allCountries.find((c) => c.cca2 === value);
        const display = selected?.label ?? "+82";
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
                  {allCountries.map((c, index) => (
                    <Pressable
                      key={index}
                      style={styles.row}
                      onPress={() => {
                        onChange(c.cca2);
                        closeModal();
                      }}
                    >
                      <Text style={styles.countryName} numberOfLines={1}>
                        {c.name}
                      </Text>
                      <Text style={styles.countryCode} numberOfLines={1}>
                        {c.label}
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
  countryName: {
    flex: 1,
    fontSize: 16,
    color: colors.BLACK,
  },
  countryCode: {
    fontSize: 14,
    color: colors.GRAY_700,
  },
});
