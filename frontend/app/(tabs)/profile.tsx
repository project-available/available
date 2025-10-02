import { View, Text, Image } from 'react-native';

export default function Profile() {
  return (
    <View className="flex-1 items-center justify-center">
      <Image
        source={{ uri: 'https://i.pravatar.cc/150' }}
        className="w-30 h-30 rounded-full mb-3"
      />
      <Text className="text-lg font-semibold">Profile screen</Text>
    </View>
  );
}
