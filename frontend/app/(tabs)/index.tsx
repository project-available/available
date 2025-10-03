import { View, Text, TouchableOpacity } from "react-native";
import { useNavigation } from "@react-navigation/native";
import Ionicons from "@expo/vector-icons/Ionicons";


export default function HomeTab() {
  const navigation = useNavigation();

  return (
    <View className="flex-1 bg-[#f8f5f2] p-4">
      {/* Header */}
    <View className="mt-12 flex-row items-center justify-between bg-white rounded-full px-6 py-2 h-16">
      {/* Logo avatar */}
      <View className="w-[52px] h-[52px] rounded-full bg-[#f8f5f2] items-center justify-center">
        <Ionicons name="person-circle-outline" size={52} color="#000" />
      </View>
      {/* Text */}
      <View className="ml-3 flex-1">
        <Text className="text-[12px] leading-4 font-inter">Hello!!!</Text>
        <Text className="text-[16px] leading-6 font-semibold">Nhi</Text>
      </View>

      {/* Notification icon */}
      <TouchableOpacity className="ml-3 w-10 h-10 rounded-full bg-[#f8f5f2] items-center justify-center">
        <Ionicons name="notifications-outline" size={24} color="black" />
      </TouchableOpacity>
    </View>
  

      {/* Title */}
      <Text className="text-xl font-bold mb-4">Services</Text>

      {/* Services Grid */}
      <View className="flex-row flex-wrap justify-between gap-4">
        <TouchableOpacity
          className="bg-blue-500 rounded-xl w-[48%] h-32 justify-center items-center"
          onPress={() => navigation.navigate("room")}
        >
          <Text className="text-white font-bold">Book</Text>
          <Text className="text-white">Making your booking</Text>
        </TouchableOpacity>

        <TouchableOpacity
          className="bg-white rounded-xl w-[48%] h-32 justify-center items-center"
          onPress={() => navigation.navigate("room")}
        >
          <Text className="font-bold">Cancel</Text>
          <Text>Cancel your booking</Text>
        </TouchableOpacity>

        <TouchableOpacity
          className="bg-white rounded-xl w-[48%] h-32 justify-center items-center"
          onPress={() => navigation.navigate("notification")}
        >
          <Text className="font-bold">View History</Text>
          <Text>View booking history</Text>
        </TouchableOpacity>

        <TouchableOpacity
          className="bg-white rounded-xl w-[48%] h-32 justify-center items-center"
          onPress={() => navigation.navigate("profile")}
        >
          <Text className="font-bold">Available</Text>
          <Text>View available rooms</Text>
        </TouchableOpacity>
      </View>
    </View>
  );
}
