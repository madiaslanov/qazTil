export type User = {
  id: number;
  email: string;
  display_name: string;
  created_at: string;
};

/** Ответ /auth/register и /auth/login. */
export type Session = {
  token: string;
  /** ISO-время, после которого Go отклонит токен. */
  expires_at: string;
  user: User;
};

export type LoginPayload = {
  email: string;
  password: string;
};

export type RegisterPayload = LoginPayload & {
  display_name: string;
};
