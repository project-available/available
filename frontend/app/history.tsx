import { View, Text, TouchableOpacity } from "react-native";
import Ionicons from "@expo/vector-icons/Ionicons";
import { useRouter } from "expo-router";

export default function History() {
  const router = useRouter();

  return (
    <View className="flex-1 bg-[#f8f5f2]">
      {/* Custom Header */}
      <View className="flex-row items-center mt-[64px] ml-[24px] gap-3">
        <TouchableOpacity onPress={() => router.back()}>
          <Ionicons name="chevron-back" size={24} color="#000" />
        </TouchableOpacity>

        <Text className="text-[18px] leading-[24px] font-semibold">
          History
        </Text>
      </View>
    </View>
  );
}
