import { View, Text, TouchableOpacity } from "react-native";
import { router } from "expo-router";

export default function WelcomeScreen() {
  return (
    <View className="flex-1 bg-white justify-center items-center px-6 relative">

      <Text className="text-[16px] text-[#666] mb-10 mt-[80px]">
        Have a better sharing experience
      </Text>

      {/* BUTTON GROUP FIXED AT BOTTOM */}
      <View className="absolute bottom-10 w-full items-center">

        <TouchableOpacity
          className="w-[95%] py-4 rounded-xl mb-4 bg-[#FDBA29] items-center"
          onPress={() => router.push("/(auth)/signup")}
        >
          <Text className="text-white text-[18px] font-semibold">
            Create an account
          </Text>
        </TouchableOpacity>

        <TouchableOpacity
          className="w-[95%] py-4 rounded-xl bg-white border items-center"
          style={{ borderColor: "#FDBA29", borderWidth: 2 }}
          onPress={() => router.push("/(auth)/login")}
        >
          <Text className="text-[18px] font-semibold" style={{ color: "#FDBA29" }}>
            Log in
          </Text>
        </TouchableOpacity>

      </View>

    </View>
  );
}
