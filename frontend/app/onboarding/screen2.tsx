import { View, Text, TouchableOpacity } from "react-native";
import { router } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import Programming from '../../image/programming.svg';

export default function Screen2() {
  return (
    <View className="flex-1 justify-center items-center bg-white">
      <Programming width={280} height={230} style={{ marginTop: 92 }} />
      <Text className="text-2xl font-bold mb-10 text-center">
        At anytime{"\n"}
        Your Study space awaits.
      </Text>

      <Text className="text-gray-500 text-lg font-semibold text-center mb-10">
        Reserve a room instantly, whether it’s{"\n"}
        day or night — your learning never{"\n"}
        stops.</Text>

      <TouchableOpacity
        activeOpacity={0.8}
        onPress={() => router.push("/onboarding/screen3")}
      >
        <View className="w-[100px] h-[100px] justify-center items-center relative">

          <View
            style={{
              position: "absolute",
              width: 90,
              height: 90,
              borderRadius: 50,
              borderWidth: 4,
              borderColor: "#FDBA29",
              borderBottomColor: "#FFF1B1",
              transform: [{ rotate: "135deg" }],
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
