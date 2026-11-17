/** Первая буква заглавная, остальное как есть: «мен шай» → «Мен шай». */
export function capitalize(text: string) {
  return text.charAt(0).toLocaleUpperCase("kk") + text.slice(1);
}
