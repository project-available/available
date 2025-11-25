import { View, Text, TouchableOpacity } from "react-native";
import { router } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import Workspace from '../../image/workspace.svg';

export default function Screen1() {
  return (
    <View className="flex-1 items-center">

    <Workspace width={280} height={230} style={{ marginTop: 48 }} />

      <Text className="text-2xl font-bold mb-10">Anywhere you are </Text>
      <Text className="text-2xl font-bold mb-10">Study your way.</Text>

      <Text className="text-gray-500 text-lg font-semibold">Book your favorite study room from anywhere — focus and learn without limits.</Text>

      <TouchableOpacity
        activeOpacity={0.8}
        onPress={() => router.push("/onboarding/screen2")}
      >
        <View className="w-[100px] h-[100px] justify-center items-center relative">

          <View
            style={{
              position: "absolute",
              width: 100,
              height: 100,
              borderRadius: 50,
              borderWidth: 4,
              borderColor: "transparent",
              borderRightColor: "#FBBF24",
              borderTopColor: "#FBBF24",
              transform: [{ rotate: "-45deg" }],
            }}
          />

          <View className="w-[75px] h-[75px] rounded-full bg-[#FDBA29] justify-center items-center">
            <Ionicons name="arrow-forward" size={28} color="white" />
          </View>
        </View>
      </TouchableOpacity>
    </View>
  );
}
