import { request } from "@/shared/api";
import type { Word } from "../model/types";

export const wordApi = {
  list: (categoryId?: number) =>
    request<Word[]>(
      categoryId ? `/words?category_id=${categoryId}` : "/words",
    ),
};
