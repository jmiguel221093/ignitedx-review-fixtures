export interface User {
  id: string;
  name: string;
  email: string;
  avatar?: string;
  role: "admin" | "user" | "guest";
  preferences: UserPreferences;
  createdAt: Date;
}

export interface UserPreferences {
  theme: "light" | "dark" | "system";
  notifications: boolean;
  language: string;
}

export function formatUserDisplayName(user: User): string {
  if (!user.name || user.name.trim() === "") {
    return user.email.split("@")[0];
  }
  return user.name;
}

export function getUserInitials(user: User): string {
  const name = formatUserDisplayName(user);
  return name
    .split(" ")
    .map((part) => part[0])
    .join("")
    .toUpperCase()
    .slice(0, 2);
}
