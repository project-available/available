import { View, Text, TouchableOpacity } from "react-native";
import { router } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import Workspace from "../../image/workspace.svg";
import AsyncStorage from "@react-native-async-storage/async-storage";

export default function Screen1() {
  const skipOnboarding = async () => {
    await AsyncStorage.setItem("hasSeenOnboarding", "true");
    router.replace("/(tabs)");
  };

  return (
    <View className="flex-1 justify-center items-center bg-white">
      <Workspace width={280} height={230} style={{ marginTop: 92 }} />

      <TouchableOpacity
        onPress={skipOnboarding}
        style={{
          position: "absolute",
          top: 50,
          right: 20,
          padding: 10,
        }}
      >
        <Text className="text-[#FDBA29] text-lg font-semibold">Skip</Text>
      </TouchableOpacity>

      <Text className="text-2xl font-bold mb-10 text-center">
        Anywhere you are {"\n"}
        Study your way.
      </Text>

      <Text className="text-gray-500 text-base font-semibold text-center mb-10">
        Book your favorite study room from{"\n"}
        anywhere — focus and learn without{"\n"}
        limits.
      </Text>

      <TouchableOpacity
        activeOpacity={0.8}
        onPress={() => router.push("/onboarding/screen2")}
      >
        <View className="w-[100px] h-[100px] justify-center items-center relative">
          <View
            style={{
              position: "absolute",
              width: 90,
              height: 90,
              borderRadius: 50,
              borderWidth: 4,
              borderColor: "#FFF1B1",
              borderRightColor: "#FDBA29",
              transform: [{ rotate: "-45deg" }],
            }}
          />

          <View className="w-[75px] h-[75px] rounded-full bg-[#FDBA29] justify-center items-center">
            <Ionicons name="arrow-forward" size={28} color="black" />
          </View>
        </View>
      </TouchableOpacity>
    </View>
  );
}
