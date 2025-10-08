import { Tabs } from "expo-router";
import { View, Text, TouchableOpacity, Platform } from "react-native";
import type { BottomTabBarProps } from "@react-navigation/bottom-tabs";
import Ionicons from "@expo/vector-icons/Ionicons";

function CustomTabBar({ state, descriptors, navigation }: BottomTabBarProps) {
  return (
    <View
      className="flex-row justify-between items-center bg-white mx-4 rounded-full shadow-lg"
      style={{
        marginBottom: Platform.OS === "ios" ? 30 : 20,
        height: 64,
        padding: 8,
      }}
    >
      {state.routes.map((route, index) => {
        const { options } = descriptors[route.key];
        const label = options.title ?? route.name;
        const isFocused = state.index === index;

        const onPress = () => {
          const event = navigation.emit({
            type: "tabPress",
            target: route.key,
          });
          if (!isFocused && !event.defaultPrevented) {
            navigation.navigate(route.name);
          }
        };

        const iconName =
          route.name === "index"
            ? "home"
            : route.name === "room"
            ? "grid-outline"
            : route.name === "notification"
            ? "notifications-outline"
            : "person-circle-outline";

        return (
          <TouchableOpacity
            key={route.key}
            onPress={onPress}
            className={`flex-row items-center justify-center rounded-full gap-1 ${
              isFocused ? "bg-[#FDBA29] px-4 h-12 w-[97px]" : "bg-[#FCF9F4] h-[48px] w-[48px]"
            }`}
          >
            <Ionicons
              name={iconName}
              size={24} // icon luôn 24x24
              color={isFocused ? "#2E2F30" : "gray"}
            />
            {isFocused && (
              <Text className="text-xs leading-4 font-semibold text-[#2E2F30]">
                {label}
              </Text>
            )}
          </TouchableOpacity>

        );
      })}
    </View>
  );
}

export default function TabLayout() {
  return (
    <Tabs
      screenOptions={{
        headerShown: false,
      }}
      tabBar={(props) => <CustomTabBar {...props} />}
    >
      <Tabs.Screen name="index" options={{ title: "Home" }} />
      <Tabs.Screen name="room" options={{ title: "Room" }} />
      <Tabs.Screen name="notification" options={{ title: "Notification" }} />
      <Tabs.Screen name="profile" options={{ title: "Profile" }} />
    </Tabs>
  );
}
