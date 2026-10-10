/**
 * Озвучка слова голосом браузера. Казахского голоса может не быть —
 * тогда браузер возьмёт ближайший доступный.
 */
export function speak(text: string, lang = "kk-KZ") {
  if (typeof window === "undefined" || !("speechSynthesis" in window)) return;
  window.speechSynthesis.cancel();
  const utterance = new SpeechSynthesisUtterance(text);
  utterance.lang = lang;
  window.speechSynthesis.speak(utterance);
}

export function canSpeak() {
  return typeof window !== "undefined" && "speechSynthesis" in window;
}
