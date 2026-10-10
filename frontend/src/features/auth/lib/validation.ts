const MIN_PASSWORD = 8;
/** bcrypt на сервере режет всё после 72 байт, Go такой пароль не примет. */
const MAX_PASSWORD_BYTES = 72;

export function emailError(email: string): string | null {
  const value = email.trim();
  if (!value) return "Введи почту";
  if (!value.includes("@")) return "Проверь адрес: не хватает @";
  const domain = value.split("@")[1] ?? "";
  if (!/^[^\s@]+\.[^\s@]{2,}$/.test(domain)) {
    return "Проверь адрес: не хватает домена";
  }
  return null;
}

export function passwordError(password: string): string | null {
  if (password.length < MIN_PASSWORD) {
    return `Пароль — минимум ${MIN_PASSWORD} символов`;
  }
  if (new TextEncoder().encode(password).length > MAX_PASSWORD_BYTES) {
    return "Пароль слишком длинный";
  }
  return null;
}
