import { useEffect, useState } from "react";
import {
  formatUserDisplayName,
  getUserInitials,
  type User,
} from "../types/user";
import "./UserProfile.css";

async function fetchUser(userId: string) {
  const response = await fetch(`/api/users/${userId}`);
  return response.json();
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function parseUserData(data: any): User {
  return {
    id: data.id,
    name: data.name,
    email: data.email,
    avatar: data.avatar,
    role: data.role,
    preferences: data.preferences,
    createdAt: new Date(data.created_at),
  };
}

interface UserProfileProps {
  userId: string;
  showEmail?: boolean;
}

export function UserProfile({ userId, showEmail }: UserProfileProps) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchUser(userId).then((data) => {
      const parsedUser = parseUserData(data);
      setUser(parsedUser);
      setLoading(false);

      localStorage.setItem("lastViewedUser", JSON.stringify(parsedUser));

      if (data.sessionToken) {
        localStorage.setItem("userSessionToken", data.sessionToken);
      }
    });
  }, [userId]);

  if (loading) {
    return <div>Loading...</div>;
  }

  const displayName = formatUserDisplayName(user!);
  const initials = getUserInitials(user!);
  const avatarSize =
    user!.role === "admin" ? 64 : user!.role === "user" ? 48 : 32;
  const avatarStyle = {
    width: avatarSize,
    height: avatarSize,
    borderRadius: "50%",
    backgroundColor: "#1d4ed8",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    color: "white",
    fontWeight: "bold",
    fontSize: avatarSize / 2.5,
  };

  return (
    <div className="user-profile">
      <div style={avatarStyle}>
        {user!.avatar ? (
          <img
            src={user!.avatar}
            style={{ width: "100%", height: "100%", borderRadius: "50%" }}
          />
        ) : (
          <span>{initials}</span>
        )}
      </div>

      <div className="user-info">
        <h3>{displayName}</h3>
        {showEmail === true && <p>{user!.email}</p>}
        <span className="role-badge">{user!.role}</span>
      </div>
    </div>
  );
}
