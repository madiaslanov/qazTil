import { queryOptions } from "@tanstack/react-query";

import { request } from "@/shared/api";
import type { Word } from "../model/types";

export const wordQueries = {
  list: (categoryId: number | null) =>
    queryOptions({
      queryKey: ["words", categoryId ?? "all"],
      queryFn: () =>
        request<Word[]>(
          categoryId ? `/words?category_id=${categoryId}` : "/words",
        ),
    }),
};
